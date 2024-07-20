// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/luma/generations": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Submit luma generate video task",
                "parameters": [
                    {
                        "description": "submit generate video",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.GenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "generate result",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.VideoTask"
                            }
                        }
                    }
                }
            }
        },
        "/luma/generations/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get luma generate video task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "page limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "video tasks",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.VideoTask"
                            }
                        }
                    }
                }
            }
        },
        "/luma/generations/:task_id/extend": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Submit luma extend generate video task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "extend task id",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "submit generate video",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.GenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "generate result",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.VideoTask"
                            }
                        }
                    }
                }
            }
        },
        "/luma/generations/file_upload": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload image to luma",
                "parameters": [
                    {
                        "description": "Upload image params",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.UploadReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "upload result",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.FileUploadResult"
                            }
                        }
                    }
                }
            }
        },
        "/luma/generations/{task_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get luma generate video task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "fetch single task by id",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "video single task",
                        "schema": {
                            "$ref": "#/definitions/main.VideoTask"
                        }
                    }
                }
            }
        },
        "/luma/generations/{task_id}/download_video_url": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get video url without watermark",
                "parameters": [
                    {
                        "type": "string",
                        "description": "fetch by id",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "url",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/luma/subscription/usage": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get current user subscription usage",
                "responses": {
                    "200": {
                        "description": "subscription info",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/luma/users/me": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get current user info",
                "responses": {
                    "200": {
                        "description": "user info",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.FileUploadResult": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "presigned_url": {
                    "type": "string"
                },
                "public_url": {
                    "type": "string"
                }
            }
        },
        "main.GenRequest": {
            "type": "object",
            "properties": {
                "aspect_ratio": {
                    "description": "option",
                    "type": "string"
                },
                "expand_prompt": {
                    "description": "option",
                    "type": "boolean"
                },
                "image_end_url": {
                    "description": "option, uploaded refer image url",
                    "type": "string"
                },
                "image_url": {
                    "description": "option, uploaded refer image url",
                    "type": "string"
                },
                "user_prompt": {
                    "description": "option",
                    "type": "string"
                }
            }
        },
        "main.UploadReq": {
            "type": "object",
            "properties": {
                "url": {
                    "description": "support public url \u0026 base64",
                    "type": "string"
                }
            }
        },
        "main.Video": {
            "type": "object",
            "properties": {
                "height": {
                    "type": "integer"
                },
                "thumbnail": {},
                "url": {
                    "type": "string"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "main.VideoTask": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "estimate_wait_seconds": {},
                "id": {
                    "type": "string"
                },
                "liked": {},
                "prompt": {
                    "type": "string"
                },
                "state": {
                    "description": "\"pending\", \"processing\", \"completed\"",
                    "type": "string"
                },
                "video": {
                    "$ref": "#/definitions/main.Video"
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
