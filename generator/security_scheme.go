package generator

import (
	"errors"

	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/getkin/kin-openapi/openapi3"
)

type SecurityScheme struct {
	jwksUri         string
	allowedIssuer   string
	allowedAudience string

	authenticator *authenticator.Authenticator
}

var (
	JWKSUriExtensionName  = "x-authenticator-jwks-uri"
	IssuerExtensionName   = "x-authenticator-issuer"
	AudienceExtensionName = "x-authenticator-audience"
)

func getFromExtension(ss *openapi3.SecuritySchemeRef, name string) *string {
	v, e := ss.Value.Extensions[name]
	if !e {
		return nil
	}

	ext := v.(string)
	return &ext
}

func NewAuthenticatorFromSecurityScheme(ss *openapi3.SecuritySchemeRef, jwksUri *string, allowedIssuer *string, allowedAudience *string) (authenticator.Authenticator, error) {
	if jwksUri == nil {
		jwksUri = getFromExtension(ss, JWKSUriExtensionName)
	}

	if allowedIssuer == nil {
		allowedIssuer = getFromExtension(ss, IssuerExtensionName)
	}

	if allowedAudience == nil {
		allowedAudience = getFromExtension(ss, AudienceExtensionName)
	}

	sstype := ss.Value.Type
	switch sstype {
	case string(authenticator.AuthenticatorTypeOpenIdConnect):
		return authenticator.NewAuthenticatorOpenIdConnect(ss, allowedAudience)

	case string(authenticator.AuthenticatorTypeOAuth2):
		if jwksUri == nil || allowedIssuer == nil {
			return nil, errors.New("jwks uri and issuer must be set for oauth2 security scheme")
		}

		return authenticator.NewAuthenticatorOAuth2(ss, *jwksUri, *allowedIssuer, allowedAudience)

	case string(authenticator.AuthenticatorTypeHttp):
		if ss.Value.Scheme != "bearer" {
			return nil, errors.New("http security scheme must be bearer")
		}
		if jwksUri == nil || allowedIssuer == nil {
			return nil, errors.New("jwks uri and issuer must be set for oauth2 security scheme")
		}
		return authenticator.NewAuthenticatorHttpBearer(ss, *jwksUri, *allowedIssuer, *allowedAudience)

	default:
		return nil, errors.New("unknown security scheme")
	}
}

func (s *SecurityScheme) GetJwksUri() (string, error) {
	if s.jwksUri == "" {
		return "", errors.New("no jwksUris found for given security scheme")
	}

	return s.jwksUri, nil
}

func (s *SecurityScheme) GetAllowedIssuer() (string, error) {
	if s.allowedIssuer == "" {
		return "", errors.New("no allowed issuer found for given security scheme")
	}

	return s.allowedIssuer, nil
}

func (s *SecurityScheme) GetAllowedAudience() string {
	return s.allowedAudience
}

func (s *SecurityScheme) GetAuthenticator() *authenticator.Authenticator {
	return s.authenticator
}
