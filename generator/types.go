package generator

type JWTAuthenticatorConfig struct {
	jwks_urls       []string
	trusted_issuers []string
	required_scope  []string
}
