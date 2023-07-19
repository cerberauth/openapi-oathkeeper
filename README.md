# OpenAPI to Ory Oathkeeper rules

<p align="left">
    <a href="https://github.com/cerberauth/openapi-oathkeeper/actions/workflows/ci.yml"><img src="https://github.com/cerberauth/openapi-oathkeeper/actions/workflows/ci.yml/badge.svg?branch=main&event=push" alt="CI Tasks for Ory Hydra"></a>
    <a href="https://codecov.io/gh/cerberauth/openapi-oathkeeper"><img src="https://codecov.io/gh/cerberauth/openapi-oathkeeper/branch/main/graph/badge.svg?token=BD1WPXJDAW"/></a>
    <a href="https://goreportcard.com/report/github.com/cerberauth/openapi-oathkeeper"><img src="https://goreportcard.com/badge/github.com/cerberauth/openapi-oathkeeper" alt="Go Report Card"></a>
    <a href="https://pkg.go.dev/github.com/cerberauth/openapi-oathkeeper"><img src="https://pkg.go.dev/badge/www.github.com/cerberauth/openapi-oathkeeper" alt="PkgGoDev"></a>
</p>

This project aims to automating the generation of OathKeeper rules from an OpenAPI 3 contract and save a lot of time and effort, especially for larger projects with many endpoints or many services. By leveraging the information in the OpenAPI 3 contract, your tool can generate secure and consistent OathKeeper rules that enforce authentication and authorization policies for each API endpoint. This can improve the overall security of the API and ensure that access is granted only to authorized parties. Additionally, this tool can simplify the development process by reducing the amount of manual work required to write and maintain OathKeeper rules.

## Ory Oathkeeper

If you're not yet familiar with Ory Oathkeeper, I highly recommend checking it out as a powerful and flexible Identity & Access Proxy. You can find more information and get started with [Ory Oathkeeper](https://github.com/ory/oathkeeper).

> ORY Oathkeeper is an Identity & Access Proxy (IAP) and Access Control Decision API that authorizes HTTP requests based on sets of Access Rules. The BeyondCorp Model is designed by Google and secures applications in Zero-Trust networks.

> An Identity & Access Proxy is typically deployed in front of (think API Gateway) web-facing applications and is capable of authenticating and optionally authorizing access requests. The Access Control Decision API can be deployed alongside an existing API Gateway or reverse proxy. ORY Oathkeeper's Access Control Decision API works with:

## Get Started

To use this tool, first you have to download the binary from the latest [release](https://github.com/cerberauth/openapi-oathkeeper/releases). Then provide the path to your OpenAPI 3 contract file. Once you have specified these options, the tool will analyze your contract and generate OathKeeper rules that enforce the specified access policies. You can then save these rules to a file to make it read by Oathkeeper.

## Features

The OpenAPI to Oathkeeper CLI supports following security scheme type:
* `openIdConnect`
* `oauth2`
* `http` with scheme `bearer`

As the authenticator rule may require additional information in order to make authorization and authentication properly, additional information can be passed either by OpenAPI 3 Extensions or CLI arguments.

### Available authenticator options

| Name     | Security Schemes                  | OpenAPI Extension Name     |
|----------|-----------------------------------|----------------------------|
| JWKS URI | `oauth2`, `http`                  | `x-authenticator-jwks-uri` |
| Issuer   | `oauth2`, `http`                  | `x-authenticator-issuer`   |
| Audience | `openIdConnect`, `oauth2`, `http` | `x-authenticator-audience` |

### Example

Here's an example of the same OpenAPI contract but in JSON format

```json sample.openapi.json
{
    "openapi": "3.0.0",
    "info": {
        "title": "My API",
        "version": "1.0.0"
    },
    "servers": [
        {
            "url": "https://api.example.com",
            "description": "Production server"
        }
    ],
    "paths": {
        "/users/{id}": {
            "get": {
                "summary": "Get user by ID",
                "operationId": "getUserById",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "description": "The user id. ",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "type": "object",
                                    "properties": {
                                        "id": {
                                            "type": "integer"
                                        },
                                        "email": {
                                            "type": "string"
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "security": [
                    {
                        "openidconnect": [
                            "user:read"
                        ]
                    }
                ]
            },
            "put": {
                "tags": [
                    "user"
                ],
                "summary": "Update user",
                "description": "This can only be done by the logged in user.",
                "operationId": "updateUser",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "description": "user id that need to be updated",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "requestBody": {
                    "description": "Update an existent user in the store",
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/User"
                            }
                        }
                    }
                },
                "responses": {
                    "default": {
                        "description": "successful operation"
                    }
                },
                "security": [
                    {
                        "openidconnect": [
                            "user:write"
                        ]
                    }
                ]
            }
        }
    },
    "components": {
        "schemas": {
            "User": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "integer",
                        "format": "int64",
                        "example": 10
                    },
                    "email": {
                        "type": "string",
                        "example": "john@email.com"
                    }
                }
            }
        },
        "securitySchemes": {
            "openidconnect": {
                "type": "openIdConnect",
                "openIdConnectUrl": "https://project.console.ory.sh/.well-known/openid-configuration",
                "x-authenticator-audience": "https://api.cerberauth.com"
            }
        }
    }
}
```

This contract defines a single endpoint at /users/{id} that returns a user object response. The endpoint is secured with OpenID Connect authentication using the `openidconnect` security scheme. The components section defines the `openidconnect` security scheme, including the URL of the OpenID Connect configuration.

To generate rules using the tool, simply run the command in your terminal with the appropriate arguments.

```shell
openapi-oathkeeper generate -f test/stub/sample.openapi.json
```

Here is a Ory Oathkeeper rules output

```json
[
    {
        "id": "getUserById",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<^(https://api\\.example\\.com)(/users/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/?)$>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "jwks_urls": [
                        "https://console.ory.sh/.well-known/jwks.json"
                    ],
                    "trusted_issuers": [
                        "https://console.ory.sh"
                    ],
                    "required_scope": [
                        "user:read"
                    ],
                    "target_audience": [
                        "https://api.cerberauth.com"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "errors": [
            {
                "handler": "json",
                "config": null
            }
        ],
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "updateUser",
        "version": "",
        "description": "This can only be done by the logged in user.",
        "match": {
            "methods": [
                "PUT"
            ],
            "url": "<^(https://api\\.example\\.com)(/users/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/?)$>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "jwks_urls": [
                        "https://console.ory.sh/.well-known/jwks.json"
                    ],
                    "trusted_issuers": [
                        "https://console.ory.sh"
                    ],
                    "required_scope": [
                        "user:write"
                    ],
                    "target_audience": [
                        "https://api.cerberauth.com"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "errors": [
            {
                "handler": "json",
                "config": null
            }
        ],
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    }
]
```

### Command line documentation

The documentation is available as markdown files in the [docs](./docs/openapi-oathkeeper.md) directory or by running `openapi-oathkeeper help`.

## Roadmap

Please note that this tool is currently in beta stage and there may be limitations and bugs. Improvements and new features should come to make it more powerful and useful for developers. Any feedback or suggestions are greatly appreciated!

You can find the milestones and future enhancements planned for this tool on the project's [GitHub milestones page]((https://github.com/cerberauth/openapi-oathkeeper/milestones)).

## Useful Links

- [ORY Oathkeeper](https://github.com/ory/oathkeeper)
- [OpenAPI 3.x Specification](https://swagger.io/specification/)

## Maintainers

[![Emmanuel Gautier](https://avatars0.githubusercontent.com/u/2765366?s=144)](https://www.emmanuelgautier.com) |
--- |
[Emmanuel Gautier](https://www.emmanuelgautier.com) |

## License

MIT © [CerberAuth](https://www.cerberauth.com)
