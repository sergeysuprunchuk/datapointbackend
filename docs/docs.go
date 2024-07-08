// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/queries/execute": {
            "post": {
                "tags": [
                    "запросы"
                ],
                "parameters": [
                    {
                        "description": "запрос",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Query"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/sources": {
            "get": {
                "tags": [
                    "источники"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Source"
                            }
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "источники"
                ],
                "parameters": [
                    {
                        "description": "источник",
                        "name": "source",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Source"
                        }
                    }
                ],
                "responses": {}
            },
            "patch": {
                "tags": [
                    "источники"
                ],
                "parameters": [
                    {
                        "description": "источник",
                        "name": "source",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Source"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/sources/drivers": {
            "get": {
                "tags": [
                    "источники"
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
        "/sources/{id}": {
            "get": {
                "tags": [
                    "источники"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "идентификатор источника",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Source"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "источники"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "идентификатор источника",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/sources/{id}/functions": {
            "get": {
                "tags": [
                    "источники"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "идентификатор источника",
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
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/sources/{id}/tables": {
            "get": {
                "tags": [
                    "источники"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "идентификатор источника",
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
                                "$ref": "#/definitions/database.Table"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "database.Column": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "database.Condition": {
            "type": "object",
            "properties": {
                "columns": {
                    "description": "предыдущий и текущий столбец таблицы.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.QColumn"
                    }
                },
                "operator": {
                    "type": "string"
                }
            }
        },
        "database.QColumn": {
            "type": "object",
            "properties": {
                "func": {
                    "description": "используется только в select.",
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "payload": {
                    "description": "специальные данные, привязанные к этому столбцу.",
                    "type": "object",
                    "additionalProperties": {}
                },
                "tableKey": {
                    "description": "ключ таблицы, которой принадлежит столбец.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/database.QTableKey"
                        }
                    ]
                },
                "value": {
                    "description": "используется в insert, update, delete и where.",
                    "type": "array",
                    "items": {}
                }
            }
        },
        "database.QTable": {
            "type": "object",
            "properties": {
                "increment": {
                    "description": "приращение имени для создания уникальных псевдонимов.",
                    "type": "integer"
                },
                "name": {
                    "description": "имя таблицы.",
                    "type": "string"
                },
                "next": {
                    "description": "используется только в select.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.QTable"
                    }
                },
                "rule": {
                    "description": "правило объединения с предыдущей таблицей.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/database.Rule"
                        }
                    ]
                }
            }
        },
        "database.QTableKey": {
            "type": "object",
            "properties": {
                "increment": {
                    "description": "приращение имени для создания уникальных псевдонимов.",
                    "type": "integer"
                },
                "name": {
                    "description": "имя таблицы.",
                    "type": "string"
                }
            }
        },
        "database.Rule": {
            "type": "object",
            "properties": {
                "conditions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Condition"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "database.Table": {
            "type": "object",
            "properties": {
                "columns": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Column"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "entity.Query": {
            "type": "object",
            "properties": {
                "columns": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.QColumn"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "orderBy": {
                    "description": "используется только в select.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.QColumn"
                    }
                },
                "sourceId": {
                    "type": "string"
                },
                "table": {
                    "$ref": "#/definitions/database.QTable"
                },
                "type": {
                    "type": "string"
                },
                "where": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.QColumn"
                    }
                }
            }
        },
        "entity.Source": {
            "type": "object",
            "properties": {
                "connected": {
                    "type": "boolean"
                },
                "databaseName": {
                    "type": "string"
                },
                "driver": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
