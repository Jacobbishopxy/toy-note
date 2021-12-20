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
        "contact": {
            "name": "Jacob Bishop",
            "url": "https://github.com/Jacobbishopxy",
            "email": "jacobbishopxy@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/delete-tag/{id}": {
            "delete": {
                "description": "delete a tag by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tag"
                ],
                "summary": "delete a tag by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "tag ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Tag"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/get-tags": {
            "get": {
                "description": "get all tags",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tag"
                ],
                "summary": "get all tags",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Tag"
                            }
                        }
                    }
                }
            }
        },
        "/save-tag": {
            "post": {
                "description": "create a new tag or update an existing tag, based on whether the tag ID is provided",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tag"
                ],
                "summary": "create/update a tag",
                "parameters": [
                    {
                        "description": "tag data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Tag"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Tag"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Affiliate": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "filename": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "object_id": {
                    "type": "string"
                },
                "post_refer": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entity.Post": {
            "type": "object",
            "properties": {
                "affiliates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Affiliate"
                    }
                },
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "subtitle": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Tag"
                    }
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entity.Tag": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Post"
                    }
                },
                "updated_at": {
                    "type": "string"
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
	Host:        "localhost:8080",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Toy-note API",
	Description: "A simple toy-note API",
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
	swag.Register("swagger", &s{})
}