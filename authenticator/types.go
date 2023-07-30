package authenticator

type AuthenticatorType string

const (
	AuthenticatorTypeNoop          AuthenticatorType = "noop"
	AuthenticatorTypeOpenIdConnect AuthenticatorType = "openidconnect"
	AuthenticatorTypeOAuth2        AuthenticatorType = "oauth2"
	AuthenticatorTypeHttp          AuthenticatorType = "http"
)
