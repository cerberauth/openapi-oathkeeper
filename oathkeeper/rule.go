// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oathkeeper

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type RuleMatch struct {
	// An array of HTTP methods (e.g. GET, POST, PUT, DELETE, ...). When ORY Oathkeeper searches for rules
	// to decide what to do with an incoming request to the proxy server, it compares the HTTP method of the incoming
	// request with the HTTP methods of each rules. If a match is found, the rule is considered a partial match.
	// If the matchesUrl field is satisfied as well, the rule is considered a full match.
	Methods []string `json:"methods"`

	// This field represents the URL pattern this rule matches. When ORY Oathkeeper searches for rules
	// to decide what to do with an incoming request to the proxy server, it compares the full request URL
	// (e.g. https://mydomain.com/api/resource) without query parameters of the incoming
	// request with this field. If a match is found, the rule is considered a partial match.
	// If the matchesMethods field is satisfied as well, the rule is considered a full match.
	//
	// You can use regular expressions or glob patterns in this field to match more than one url.
	// The matching strategy is determined by configuration parameter MatchingStrategy.
	// Regular expressions and glob patterns are encapsulated in brackets < and >.
	// The following regexp example matches all paths of the domain `mydomain.com`: `https://mydomain.com/<.*>`.
	// The glob equivalent of the above regexp example is `https://mydomain.com/<*>`.
	URL string `json:"url"`
}

type RuleHandler struct {
	// Handler identifies the implementation which will be used to handle this specific request. Please read the user
	// guide for a complete list of available handlers.
	Handler string `json:"handler"`

	// Config contains the configuration for the handler. Please read the user
	// guide for a complete list of each handler's available settings.
	Config json.RawMessage `json:"config"`
}

type RuleErrorHandler struct {
	// Handler identifies the implementation which will be used to handle this specific request. Please read the user
	// guide for a complete list of available handlers.
	Handler string `json:"handler"`

	// Config defines additional configuration for the response handler.
	Config json.RawMessage `json:"config"`
}

type OnErrorRequest struct {
	// ContentType defines the content type(s) that should match. Wildcards such as `application/*` are supported.
	ContentType []string `json:"content_type"`

	// Accept defines the accept header that should match. Wildcards such as `application/*` are supported.
	Accept []string `json:"accept"`
}

type URLProvider interface{}

// Rule is a single rule that will get checked on every HTTP request.
type Rule struct {
	// ID is the unique id of the rule. It can be at most 190 characters long, but the layout of the ID is up to you.
	// You will need this ID later on to update or delete the rule.
	ID string `json:"id"`

	// Version represents the access rule version. Should match one of ORY Oathkeepers release versions. Supported since
	// v0.20.0-beta.1+oryOS.14.
	Version string `json:"version"`

	// Description is a human readable description of this rule.
	Description string `json:"description"`

	// Match defines the URL that this rule should match.
	Match URLProvider `json:"match" faker:"urlProvider"`

	// Authenticators is a list of authentication handlers that will try and authenticate the provided credentials.
	// Authenticators are checked iteratively from index 0 to n and if the first authenticator to return a positive
	// result will be the one used.
	//
	// If you want the rule to first check a specific authenticator  before "falling back" to others, have that authenticator
	// as the first item in the array.
	Authenticators []RuleHandler `json:"authenticators"`

	// Authorizer is the authorization handler which will try to authorize the subject (authenticated using an Authenticator)
	// making the request.
	Authorizer RuleHandler `json:"authorizer"`

	// Mutators is a list of mutation handlers that transform the HTTP request. A common use case is generating a new set
	// of credentials (e.g. JWT) which then will be forwarded to the upstream server.
	//
	// Mutations are performed iteratively from index 0 to n and should all succeed in order for the HTTP request to be forwarded.
	Mutators []RuleHandler `json:"mutators"`

	// Errors is a list of error handlers. These will be invoked if any part of the system returns an error. You can
	// configure error matchers to listen on certain errors (e.g. unauthorized) and execute specific logic (e.g. redirect
	// to the login endpoint, return with an XML error, return a json error, ...).
	Errors []RuleErrorHandler `json:"errors"`

	// Upstream is the location of the server where requests matching this rule should be forwarded to.
	Upstream RuleUpstream `json:"upstream"`
}

type RuleUpstream struct {
	// PreserveHost, if false (the default), tells ORY Oathkeeper to set the upstream request's Host header to the
	// hostname of the API's upstream's URL. Setting this flag to true instructs ORY Oathkeeper not to do so.
	PreserveHost bool `json:"preserve_host"`

	// StripPath if set, replaces the provided path prefix when forwarding the requested URL to the upstream URL.
	StripPath string `json:"strip_path"`

	// URL is the URL the request will be proxied to.
	URL string `json:"url"`
}

var _ json.Unmarshaler = new(Rule)

func (r *Rule) UnmarshalJSON(raw []byte) error {
	var rr struct {
		ID             string             `json:"id"`
		Version        string             `json:"version"`
		Description    string             `json:"description"`
		Authenticators []RuleHandler      `json:"authenticators"`
		Authorizer     RuleHandler        `json:"authorizer"`
		Mutators       []RuleHandler      `json:"mutators"`
		Errors         []RuleErrorHandler `json:"errors"`
		Upstream       RuleUpstream       `json:"upstream"`

		RawMatch json.RawMessage `json:"match"`
		Match    URLProvider
	}

	if err := json.Unmarshal(raw, &rr); err != nil {
		return errors.WithStack(err)
	}

	// copy all fields
	r.ID = rr.ID
	r.Version = rr.Version
	r.Description = rr.Description
	r.Match = rr.Match
	r.Authenticators = rr.Authenticators
	r.Authorizer = rr.Authorizer
	r.Mutators = rr.Mutators
	r.Errors = rr.Errors
	r.Upstream = rr.Upstream

	return nil
}

// GetID returns the rule's ID.
func (r *Rule) GetID() string {
	return r.ID
}
