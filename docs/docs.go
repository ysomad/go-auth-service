// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "post": {
                "description": "Create a new user with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Sign up",
                "operationId": "signup",
                "parameters": [
                    {
                        "description": "To create a new user email and password should be provided",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.createUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.messageResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/v1.validationErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.messageResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Receive user data by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get",
                "operationId": "get",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.messageResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Update user data partially",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Partial update",
                "operationId": "update",
                "parameters": [
                    {
                        "description": "Provide at least one user field to update user data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.partialUpdateRequest"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.messageResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/v1.validationErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}/archive": {
            "patch": {
                "description": "Archive or restore User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Archive or restore User",
                "operationId": "archive",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "To archive or restore a user is_archive should be provided",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.archiveUserRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.messageResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/v1.validationErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.User": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2021-08-31T16:55:18.080768Z"
                },
                "email": {
                    "type": "string",
                    "example": "user@mail.com"
                },
                "firstName": {
                    "type": "string",
                    "example": "Alex"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "isActive": {
                    "type": "boolean",
                    "example": true
                },
                "isArchive": {
                    "type": "boolean",
                    "example": false
                },
                "lastName": {
                    "type": "string",
                    "example": "Malykh"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2021-08-31T16:55:18.080768Z"
                },
                "username": {
                    "type": "string",
                    "example": "username"
                }
            }
        },
        "v1.archiveUserRequest": {
            "type": "object",
            "required": [
                "isArchive"
            ],
            "properties": {
                "isArchive": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "v1.createUserRequest": {
            "type": "object",
            "required": [
                "confirmPassword",
                "email",
                "password"
            ],
            "properties": {
                "confirmPassword": {
                    "type": "string",
                    "example": "secret"
                },
                "email": {
                    "type": "string",
                    "example": "user@mail.com"
                },
                "password": {
                    "type": "string",
                    "example": "secret"
                }
            }
        },
        "v1.messageResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                }
            }
        },
        "v1.partialUpdateRequest": {
            "type": "object",
            "properties": {
                "firstName": {
                    "type": "string",
                    "example": "Alex"
                },
                "lastName": {
                    "type": "string",
                    "example": "Malykh"
                },
                "username": {
                    "type": "string",
                    "example": "username"
                }
            }
        },
        "v1.validationErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "ModelName.FieldName": "validation error message"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "0.0.0.0:8080",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "Golang auth service",
	Description: "REST API authentication and user management service",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
