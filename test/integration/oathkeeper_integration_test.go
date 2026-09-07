//go:build integration

// Package integration exercises the real openapi-oathkeeper CLI together
// with a real, unmodified Ory Oathkeeper gateway (run via Docker), the same
// image pinned by contrib/quickstart.yml. It is meant to be run explicitly,
// e.g. `go test -tags=integration ./test/integration/...`, since it needs
// Docker and shells out to `docker run` with `--network host` (Linux only,
// which matches the ubuntu-latest GitHub Actions runner).
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const oathkeeperImage = "oryd/oathkeeper:v0.40"

// repoRoot returns the module root, computed from this file's own path so
// the test can be run from any working directory.
func repoRoot(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)

	return filepath.Join(filepath.Dir(thisFile), "..", "..")
}

func requireDocker(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker is not available, skipping integration test")
	}
}

// startMockUpstream starts a plain HTTP server that echoes the requested
// path, standing in for the real API behind the Oathkeeper gateway. It also
// serves an (empty) JWKS document so the jwt authenticator can be
// configured without depending on any external service.
func startMockUpstream(t *testing.T) (addr string, closeFn func()) {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"keys":[]}`))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"upstream_received_path": r.URL.Path})
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := &http.Server{Handler: mux}
	go func() {
		_ = srv.Serve(ln)
	}()

	return ln.Addr().String(), func() {
		_ = srv.Close()
	}
}

// buildCLI builds the real, unmodified openapi-oathkeeper binary from the
// module under test.
func buildCLI(t *testing.T, root string) string {
	t.Helper()

	binPath := filepath.Join(t.TempDir(), "openapi-oathkeeper")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = root

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to build CLI: %s", out)

	return binPath
}

// generateRules runs the real CLI's `generate` command against
// test/integration/openapi.yaml and returns the path to the resulting
// access rules file.
func generateRules(t *testing.T, cliPath string, root string, upstreamAddr string) string {
	t.Helper()

	configPath := filepath.Join(t.TempDir(), "generator-config.yaml")
	generatorConfig := fmt.Sprintf(`prefix: integration
server_urls:
  - "http://127.0.0.1:4455"
upstream:
  url: "http://%s"
authenticators:
  bearerAuth:
    handler: "jwt"
    config:
      jwks_urls:
        - "http://%s/jwks.json"
`, upstreamAddr, upstreamAddr)
	require.NoError(t, os.WriteFile(configPath, []byte(generatorConfig), 0600))

	rulesPath := filepath.Join(t.TempDir(), "access-rules.json")
	cmd := exec.Command(
		cliPath, "--sqa-opt-out", "generate",
		"-f", filepath.Join(root, "test", "integration", "openapi.yaml"),
		"-c", configPath,
		"-o", rulesPath,
	)

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to generate rules: %s", out)

	rules, err := os.ReadFile(rulesPath) //nolint:gosec
	require.NoError(t, err)

	// The CLI writes the rules file with mode 0600; make it world-readable
	// so the (non-root) user inside the Oathkeeper container can read it
	// through the bind mount.
	require.NoError(t, os.Chmod(rulesPath, 0644)) //nolint:gosec

	// Regression check for the path-segment containment fix: string path
	// parameters must not compile to the slash-crossing `.+` token (see
	// generator/match_rule.go).
	require.NotContains(t, string(rules), "<.+>",
		"generated rule contains the slash-crossing <.+> token")
	require.Contains(t, string(rules), "<[^/]+>",
		"generated rule is expected to contain the segment-anchored <[^/]+> token")

	return rulesPath
}

// startOathkeeper loads rulesPath into a real, unmodified Ory Oathkeeper
// gateway running in Docker with `--network host`, so it shares the host's
// network namespace: it can reach the mock upstream bound to 127.0.0.1, and
// the test can dial its proxy port directly. This is Linux-only, matching
// the ubuntu-latest CI runner.
func startOathkeeper(t *testing.T, root string, rulesPath string) {
	t.Helper()
	requireDocker(t)

	containerName := fmt.Sprintf("openapi-oathkeeper-integration-%d", os.Getpid())
	oathkeeperConfigPath := filepath.Join(root, "test", "integration", "oathkeeper.yml")

	runArgs := []string{
		"run", "--rm", "-d",
		"--name", containerName,
		"--network", "host",
		"-v", rulesPath + ":/etc/config/oathkeeper/access-rules.json:ro",
		"-v", oathkeeperConfigPath + ":/etc/config/oathkeeper/oathkeeper.yml:ro",
		oathkeeperImage,
		"serve", "proxy", "-c", "/etc/config/oathkeeper/oathkeeper.yml",
	}

	out, err := exec.Command("docker", runArgs...).CombinedOutput()
	require.NoError(t, err, "failed to start oathkeeper: %s", out)

	t.Cleanup(func() {
		logs, _ := exec.Command("docker", "logs", containerName).CombinedOutput()
		t.Logf("oathkeeper logs:\n%s", logs)

		if out, err := exec.Command("docker", "rm", "-f", containerName).CombinedOutput(); err != nil {
			t.Logf("failed to remove oathkeeper container: %s: %s", err, out)
		}
	})

	waitForPort(t, "127.0.0.1:4455", 30*time.Second)
}

func waitForPort(t *testing.T, addr string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(300 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for %s to become available", addr)
}

func get(t *testing.T, url string, headers map[string]string) (int, string) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	require.NoError(t, err)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(body)
}

// TestGeneratedRulesAgainstRealOathkeeper builds the real CLI, generates
// access rules from an OpenAPI document, loads them into a real Ory
// Oathkeeper gateway, and exercises it with real HTTP requests end-to-end.
//
// It doubles as a regression test for the path-parameter regex fix: a rule
// generated for a single path segment (/public/{slug}) must not match
// deeper, unrelated subpaths sharing the same prefix (CWE-284, overbroad
// `<.+>` token previously emitted by generator/match_rule.go).
func TestGeneratedRulesAgainstRealOathkeeper(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	requireDocker(t)

	root := repoRoot(t)
	cliPath := buildCLI(t, root)

	upstreamAddr, closeUpstream := startMockUpstream(t)
	defer closeUpstream()

	rulesPath := generateRules(t, cliPath, root, upstreamAddr)
	startOathkeeper(t, root, rulesPath)

	t.Run("intended single path segment is proxied", func(t *testing.T) {
		status, body := get(t, "http://127.0.0.1:4455/public/alice", nil)
		require.Equal(t, http.StatusOK, status, body)
		require.Contains(t, body, `"/public/alice"`)
	})

	t.Run("deeper subpath sharing the same prefix is rejected", func(t *testing.T) {
		// Before the fix, the generated <.+> token let this request match
		// the /public/{slug} rule and get proxied unauthenticated, even
		// though it targets a completely different, deeper path.
		status, _ := get(t, "http://127.0.0.1:4455/public/alice/admin/delete-all-users", nil)
		require.Equal(t, http.StatusNotFound, status)
	})

	t.Run("secured route rejects unauthenticated requests", func(t *testing.T) {
		status, _ := get(t, "http://127.0.0.1:4455/secure/42", nil)
		require.Equal(t, http.StatusUnauthorized, status)
	})

	t.Run("secured route rejects requests with an invalid token", func(t *testing.T) {
		status, _ := get(t, "http://127.0.0.1:4455/secure/42", map[string]string{
			"Authorization": "Bearer not-a-real-token",
		})
		require.NotEqual(t, http.StatusOK, status)
	})
}
