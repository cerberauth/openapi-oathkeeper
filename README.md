# Ory Oathkeeper rules from OpenAPI


[![Join Discord](https://img.shields.io/discord/1242773130137833493?label=Discord&style=for-the-badge)](https://www.cerberauth.com/community)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/cerberauth/openapi-oathkeeper/ci.yml?branch=main&label=core%20build&style=for-the-badge)](https://github.com/cerberauth/openapi-oathkeeper/actions/workflows/ci.yml)
![Codecov](https://img.shields.io/codecov/c/gh/cerberauth/openapi-oathkeeper?token=BD1WPXJDAW&style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/cerberauth/openapi-oathkeeper?style=for-the-badge)](https://goreportcard.com/report/github.com/cerberauth/openapi-oathkeeper)
[![GoDoc reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/cerberauth/openapi-oathkeeper)

This CLI generates OathKeeper rules that enforce authentication and authorization policies for each API endpoint from an OpenAPI file.

This project automate the generation of Oathkeeper rules from an OpenAPI contract and save a lot of time especially for larger projects with many endpoints or many services by using the existing documentation provided in an OpenAPI contract. This can improve the overall security of the API and ensure that access is granted only to authorized parties. Additionally, this tool can simplify the development process by reducing the amount of manual work required to write and maintain OathKeeper rules.

## Ory Oathkeeper

If you're not yet familiar with Ory Oathkeeper, Oathkeeper is an Identity & Access Proxy (IAP) and Access Control Decision API that authorizes HTTP requests based on sets of Access Rules. You can find more information and get started with [Ory Oathkeeper](https://github.com/ory/oathkeeper).

> An Identity & Access Proxy is typically deployed in front of (think API Gateway or Service mesh) web-facing applications and is capable of authenticating and optionally authorizing access requests. The Access Control Decision API can be deployed alongside an existing API Gateway or reverse proxy.

## Installation

Below are the instructions to install on Linux, Windows, MacOS, and Docker. You can choose the installation method that best suits your needs and environment.

If none of the installation methods below work for you, you can also download the binary from the latest [release](https://github.com/cerberauth/openapi-oathkeeper/releases).

### Linux (Snap)

Install using [Snap](https://snapcraft.io/openapi-oathkeeper).

```sh
sudo snap install vulnapi
```

### MacOS (Homebrew)

Install using Homebrew.

```sh
brew tap cerberauth/openapi-oathkeeper https://github.com/cerberauth/openapi-oathkeeper
brew install $(brew --repository cerberauth/openapi-oathkeeper)/openapi-oathkeeper.rb
```

## Get Started

Provides the path to your OpenAPI contract file.

```sh
./openapi-oathkeeper generate -f ./openapi.json
```

Once you have specified these options, the tool will analyze your contract and generate OathKeeper rules that enforce the specified access policies. You can then save these rules to a file to make it read by Oathkeeper.

Here is an example Oathkeeper rules output from the [Petstore OpenAPI](./test/stub/petstore.openapi.json)

<details>
    <summary>Oathkeeper rules output</summary>

```json
[
    {
        "id": "addPet",
        "version": "",
        "description": "Add a new pet to the store",
        "match": {
            "methods": [
                "POST"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "createUser",
        "version": "",
        "description": "This can only be done by the logged in user.",
        "match": {
            "methods": [
                "POST"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "createUsersWithListInput",
        "version": "",
        "description": "Creates list of users with given input array",
        "match": {
            "methods": [
                "POST"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user/createWithList"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "deleteOrder",
        "version": "",
        "description": "For valid response try integer IDs with value < 1000. Anything above 1000 or nonintegers will generate API errors",
        "match": {
            "methods": [
                "DELETE"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/store/order/<\\d+>"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "deletePet",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "DELETE"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet/<\\d+>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "deleteUser",
        "version": "",
        "description": "This can only be done by the logged in user.",
        "match": {
            "methods": [
                "DELETE"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user/<.+>"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "findPetsByStatus",
        "version": "",
        "description": "Multiple status values can be provided with comma separated strings",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet/findByStatus"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "findPetsByTags",
        "version": "",
        "description": "Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet/findByTags"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "getInventory",
        "version": "",
        "description": "Returns a map of status codes to quantities",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/store/inventory"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "getOrderById",
        "version": "",
        "description": "For valid response try integer IDs with value <= 5 or > 10. Other values will generate exceptions.",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/store/order/<\\d+>"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "getPetById",
        "version": "",
        "description": "Returns a single pet",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet/<\\d+>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "getUserByName",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user/<.+>"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "loginUser",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user/login"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "logoutUser",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user/logout"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "placeOrder",
        "version": "",
        "description": "Place a new order in the store",
        "match": {
            "methods": [
                "POST"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/store/order"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "updatePet",
        "version": "",
        "description": "Update an existing pet by Id",
        "match": {
            "methods": [
                "PUT"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "updatePetWithForm",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "POST"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet/<\\d+>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
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
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/user/<.+>"
        },
        "authenticators": [
            {
                "handler": "noop",
                "config": null
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    },
    {
        "id": "uploadFile",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "POST"
            ],
            "url": "<(https://cerberauth\\.com/api/v3|http://swagger\\.io/api/v3)>/pet/<\\d+>/uploadImage"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "required_scope": [
                        "write:pets",
                        "read:pets"
                    ]
                }
            }
        ],
        "authorizer": {
            "handler": "allow",
            "config": null
        },
        "mutators": null,
        "errors": null,
        "upstream": {
            "preserve_host": false,
            "strip_path": "",
            "url": ""
        }
    }
]
```
</details>

## Configuration

As the authenticator rule may require additional information in order to make authorization and authentication working properly, additional information can be passed either by OpenAPI Extensions or configuration file.

### Configuration File

The recommended approach involves using dedicated configuration files for your Oathkeeper rules. These configuration files provide a more flexible and user-friendly way of managing your security settings.

Every Oathkeeper rule property can be configured this way. Here are the available properties:

| Field          | Type                                                                               | Key              |
|----------------|------------------------------------------------------------------------------------|------------------|
| Prefix         | string                                                                             | "prefix"         |
| ServerUrls     | []string                                                                           | "server_urls"    |
| Upstream       | [Upstream](https://www.ory.sh/docs/oathkeeper/api-access-rules#access-rule-format) | "upstream"       |
| Authenticators | Map of [Authenticators](https://www.ory.sh/docs/oathkeeper/pipeline/authn)         | "authenticators" |
| Authorizer     | [Authorization Handler](https://www.ory.sh/docs/oathkeeper/pipeline/authz)         | "authorizer"     |
| Mutators       | Array [of Mutator Handlers](https://www.ory.sh/docs/oathkeeper/pipeline/mutator)   | "mutators"       |
| Errors         | Array of [Error Handlers](https://www.ory.sh/docs/oathkeeper/pipeline/error)       | "errors"         |

Below is an example of a configuration file in YAML format:

```yaml
prefix: cerberauth

server_urls:
  - https://www.cerberauth.com/api
  - https://api.cerberauth.com/api

authenticators:
  openidconnect:
    handler: "jwt"
    config:
      target_audience:
      - https://api.cerberauth.com
```

In order to generate rules using the CLI, simply run the command in your terminal with the appropriate arguments.

```shell
./openapi-oathkeeper generate -c ./test/config/sample.yaml -f ./test/stub/sample.openapi.json
```

<details>
  <summary>Oathkeeper rules output</summary>

```json
[
    {
        "id": "cerberauth:getUserById",
        "version": "",
        "description": "",
        "match": {
            "methods": [
                "GET"
            ],
            "url": "<^(https://www\\.cerberauth\\.com/api|https://api\\.cerberauth\\.com/api)(/users/(?:[[:alnum:]]?\\x2D?=?\\??&?_?)+/?)$>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "jwks_urls": [
                        "https://console.ory.sh/.well-known/jwks.json"
                    ],
                    "required_scope": [
                        "user:read"
                    ],
                    "target_audience": [
                        "https://api.cerberauth.com"
                    ],
                    "trusted_issuers": [
                        "https://console.ory.sh"
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
        "id": "cerberauth:updateUser",
        "version": "",
        "description": "This can only be done by the logged in user.",
        "match": {
            "methods": [
                "PUT"
            ],
            "url": "<^(https://www\\.cerberauth\\.com/api|https://api\\.cerberauth\\.com/api)(/users/(?:[[:alnum:]]?\\x2D?=?\\??&?_?)+/?)$>"
        },
        "authenticators": [
            {
                "handler": "jwt",
                "config": {
                    "jwks_urls": [
                        "https://console.ory.sh/.well-known/jwks.json"
                    ],
                    "required_scope": [
                        "user:write"
                    ],
                    "target_audience": [
                        "https://api.cerberauth.com"
                    ],
                    "trusted_issuers": [
                        "https://console.ory.sh"
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
</details>

### OpenAPI Extension

OpenAPI Extensions serve as an extension mechanism for the OpenAPI Specification (OAS). When using OpenAPI-Oathkeeper with OpenAPI Extensions, you can embed Oathkeeper-specific rules directly within your API documentation. This integration can be beneficial when you desire a unified source of truth for both API specifications and security rules.

Here the available configurations:

| Name     | Security Schemes                  | OpenAPI Extension Name     |
|----------|-----------------------------------|----------------------------|
| JWKS URI | `oauth2`, `http`                  | `x-authenticator-jwks-uri` |
| Issuer   | `oauth2`, `http`                  | `x-authenticator-issuer`   |
| Audience | `openIdConnect`, `oauth2`, `http` | `x-authenticator-audience` |

### Example

Here's an example of the same OpenAPI contract but in JSON format

<details>
  <summary>OpenAPI example using OpenAPI Extensions</summary>

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
                "openIdConnectUrl": "https://project.console.ory.sh/.well-known/openid-configuration"
            }
        }
    }
}
```
</details>

### Command line documentation

The documentation is available as markdown files in the [docs](./docs/openapi-oathkeeper.md) directory or by running `openapi-oathkeeper help`.

## Telemetry

The scanner collects anonymous usage data to help improve the tool. This data includes the number of scans performed, number of detected vulnerabilities, and the severity of vulnerabilities. No sensitive information is collected. You can opt-out of telemetry by passing the `--sqa-opt-out` flag.

## Roadmap

Please note that this tool is currently in beta stage and there may be limitations and bugs. Improvements and new features should come to make it more powerful and useful for developers. Any feedback or suggestions are greatly appreciated!

You can find the milestones and future enhancements planned for this tool on the project's [GitHub milestones page]((https://github.com/cerberauth/openapi-oathkeeper/milestones)).

## Useful Links

- [ORY Oathkeeper](https://github.com/ory/oathkeeper)
- [OpenAPI 3.x Specification](https://swagger.io/specification/)

## License

MIT © [CerberAuth](https://www.cerberauth.com)
