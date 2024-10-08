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
        "/auth/signup": {
            "post": {
                "tags": [
                    "Auth"
                ],
                "summary": "Sign up",
                "parameters": [
                    {
                        "description": "request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app.RegisterUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/app.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/app.RegisterUserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tasks": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Get Tasks",
                "operationId": "GetTasks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "cursor for forward pagination",
                        "name": "cursor",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "maximum number of tasks to return",
                        "name": "per_page",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "completed",
                            "pending"
                        ],
                        "type": "string",
                        "description": "filter by task status",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/app.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/app.CreateTaskResponse"
                                        },
                                        "paging": {
                                            "$ref": "#/definitions/app.PaginationData"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Create Task",
                "operationId": "CreateTasks",
                "parameters": [
                    {
                        "description": "request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app.CreateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/app.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/app.CreateTaskResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tasks/{id}": {
            "delete": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Delete Tasks",
                "operationId": "DeleteTasks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Edit Tasks",
                "operationId": "EditTasks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app.EditTaskResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/app.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/app.EditTaskResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/app.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.CreateTaskRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "app.CreateTaskResponse": {
            "type": "object",
            "properties": {
                "task": {
                    "$ref": "#/definitions/app.Task"
                }
            }
        },
        "app.EditTaskResponse": {
            "type": "object",
            "properties": {
                "task": {
                    "$ref": "#/definitions/app.Task"
                }
            }
        },
        "app.ErrorResponse": {
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
        "app.PaginationData": {
            "type": "object",
            "properties": {
                "item_count": {
                    "type": "integer"
                },
                "next_cursor": {
                    "type": "integer"
                },
                "per_page": {
                    "type": "integer"
                }
            }
        },
        "app.RegisterUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "app.RegisterUserResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/app.User"
                }
            }
        },
        "app.SuccessResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "status": {
                    "type": "string"
                }
            }
        },
        "app.Task": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_completed": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "app.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_name": {
                    "type": "string"
                },
                "update_at": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Task Managment API",
	Description:      "This is a task management api server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
