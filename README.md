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

To use this tool, you need to provide the path to your OpenAPI 3 contract file, as well as some configuration options. Once you have specified these options, the tool will analyze your contract and generate OathKeeper rules that enforce the specified access policies. You can then save these rules to a file.

### Disclaimer

Please note that this tool only generates Oathkeeper rules based on OpenAPI 3 contracts that use OpenId Connect for authentication and authorization. It may not be suitable for other use cases or contract formats. Additionally, the generated rules should be reviewed and tested thoroughly before being used in a production environment.

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
                        "required": true,
                        "schema": {
                            "type": "integer"
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
                                        "name": {
                                            "type": "string"
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
                            "read:user",
                            "write:user"
                        ]
                    }
                ]
            }
        }
    },
    "components": {
        "securitySchemes": {
            "openidconnect": {
                "type": "openIdConnect",
                "openIdConnectUrl": "https://project.console.ory.sh/.well-known/openid-configuration"
            }
        }
    }
}
```

This contract defines a single endpoint at /users/{id} that returns a user object response. The endpoint is secured with OpenID Connect authentication using the `openidconnect` security scheme. The components section defines the `openidconnect` security scheme, including the URL of the OpenID Connect configuration.

To generate rules using the tool, simply run the command in your terminal with the appropriate arguments.

```shell
openapi-oathkeeper generate -f test/stub/sample.openapi.json --allowed-audiences "openidconnect=https://api.cerberauth.com/"
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
            "url": "<^(https://api\\.example\\.com)(/users/(.+)/?)$>"
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
                        "read:user",
                        "write:user"
                    ],
                    "target_audience": [
                        "https://api.cerberauth.com/"
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

Please note that this tool is currently in alpha stage and there may be limitations and bugs. Improvements and new features should come to make it more powerful and useful for developers. Any feedback or suggestions are greatly appreciated!

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
