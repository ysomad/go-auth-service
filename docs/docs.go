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
                "description": "Register a new user with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Sign Up",
                "operationId": "signup",
                "parameters": [
                    {
                        "description": "To register a new user email and password should be provided",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.CreateUserResponse"
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
            "put": {
                "description": "Update user data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update",
                "operationId": "update",
                "parameters": [
                    {
                        "description": "All required fields should be provided",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateUserRequest"
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
        },
        "/users/{id}/state": {
            "patch": {
                "description": "Update user state",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update state",
                "operationId": "state",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "To change user state is_archive should be provided",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateStateUserRequest"
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
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.CreateUserRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "confirm_password": {
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
        "domain.CreateUserResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2021-08-31T16:55:18.080768Z"
                },
                "email": {
                    "type": "string",
                    "example": "user@mail.com"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "domain.UpdateStateUserRequest": {
            "type": "object",
            "required": [
                "is_active"
            ],
            "properties": {
                "is_active": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "domain.UpdateUserRequest": {
            "type": "object",
            "required": [
                "first_name",
                "last_name",
                "username"
            ],
            "properties": {
                "first_name": {
                    "type": "string",
                    "example": "Alex"
                },
                "last_name": {
                    "type": "string",
                    "example": "Malykh"
                },
                "username": {
                    "type": "string",
                    "example": "username"
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
        "v1.validationErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "CreateUserRequest.ConfirmPassword": "must be equal to Password"
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
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Golang auth service",
	Description: "REST API authentication service",
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
