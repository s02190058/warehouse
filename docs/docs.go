// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Bakanov Artem",
            "url": "https://t.me/s02190058",
            "email": "sklirian@mail.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/warehouses/{id}": {
            "get": {
                "description": "Number of remaining products.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "warehouse"
                ],
                "summary": "OK status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "warehouse id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/warehouse.ProductRemains"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/warehouses/{id}:release": {
            "post": {
                "description": "Release products with the specified codes.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "warehouse"
                ],
                "summary": "OK status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "warehouse id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "product\tcodes\tto\tbe\treleased",
                        "name": "releasedCodes",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/warehouses/{id}:reserve": {
            "post": {
                "description": "Reserves products with the specified codes.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "warehouse"
                ],
                "summary": "OK status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "warehouse id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "product\tcodes\tto\tbe\treserved",
                        "name": "reservedCodes",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Shows that service is available.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "healthcheck"
                ],
                "summary": "OK status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "warehouse.ProductRemains": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "remains": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.o",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "warehouse App",
	Description:      "Warehouse management platform.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
