package generator

type JWTAuthenticatorConfig struct {
	JwksUrls       []string `json:"jwks_urls"`
	TrustedIssuers []string `json:"trusted_issuers"`
	RequiredScope  []string `json:"required_scope"`
}
