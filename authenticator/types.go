package authenticator

type AuthenticatorType string

const (
	AuthenticatorTypeNoop          AuthenticatorType = "noop"
	AuthenticatorTypeOpenIdConnect AuthenticatorType = "openIdConnect"
	AuthenticatorTypeOAuth2        AuthenticatorType = "oauth2"
	AuthenticatorTypeHttp          AuthenticatorType = "http"
)

type JWTAuthenticatorConfig struct {
	JwksUrls       []string `json:"jwks_urls"`
	TrustedIssuers []string `json:"trusted_issuers"`
	RequiredScope  []string `json:"required_scope"`
	TargetAudience []string `json:"target_audience"`
}
