// Package docs berisi spesifikasi Swagger/OpenAPI Delcom Blog API.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "Delcom Blog API",
    "description": "REST API Delcom Blog multi-tier + AI Generator terintegrasi Search Engine. Respons standar: {status, message, data?}",
    "version": "1.0.0"
  },
  "basePath": "/",
  "paths": {
    "/api/public/categories": {
      "get": {
        "tags": [
          "Public"
        ],
        "summary": "Daftar kategori artikel",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        }
      }
    },
    "/api/public/posts": {
      "get": {
        "tags": [
          "Public"
        ],
        "summary": "Daftar artikel terbit (q, category)",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "query",
            "name": "q",
            "type": "string"
          },
          {
            "in": "query",
            "name": "category",
            "type": "string"
          }
        ]
      }
    },
    "/api/public/posts/{slug}": {
      "get": {
        "tags": [
          "Public"
        ],
        "summary": "Detail artikel publik berdasarkan slug",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "path",
            "name": "slug",
            "required": true,
            "type": "string"
          }
        ]
      }
    },
    "/api/auth/register": {
      "post": {
        "tags": [
          "Auth"
        ],
        "summary": "Registrasi penulis baru",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "name": {
                  "type": "string"
                },
                "username": {
                  "type": "string"
                },
                "email": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                }
              },
              "required": [
                "username",
                "email",
                "password"
              ]
            }
          }
        ]
      }
    },
    "/api/auth/login": {
      "post": {
        "tags": [
          "Auth"
        ],
        "summary": "Login dan dapatkan token JWT",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "username": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                }
              },
              "required": [
                "username",
                "password"
              ]
            }
          }
        ]
      }
    },
    "/api/auth/logout": {
      "post": {
        "tags": [
          "Auth"
        ],
        "summary": "Logout",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/auth/profile": {
      "get": {
        "tags": [
          "Auth"
        ],
        "summary": "Profil pengguna login",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      },
      "put": {
        "tags": [
          "Auth"
        ],
        "summary": "Perbarui profil",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "name": {
                  "type": "string"
                },
                "email": {
                  "type": "string"
                },
                "bio": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                }
              }
            }
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/posts": {
      "get": {
        "tags": [
          "Posts"
        ],
        "summary": "Artikel milik penulis",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      },
      "post": {
        "tags": [
          "Posts"
        ],
        "summary": "Buat artikel",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "title": {
                  "type": "string"
                },
                "category": {
                  "type": "string"
                },
                "excerpt": {
                  "type": "string"
                },
                "content": {
                  "type": "string"
                },
                "cover_image_path": {
                  "type": "string"
                },
                "is_published": {
                  "type": "boolean"
                }
              },
              "required": [
                "title",
                "content"
              ]
            }
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/posts/stats": {
      "get": {
        "tags": [
          "Posts"
        ],
        "summary": "Statistik artikel penulis",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/posts/ai-generate": {
      "post": {
        "tags": [
          "AI"
        ],
        "summary": "Generate draf artikel AI + Search Engine",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "topic": {
                  "type": "string"
                },
                "category": {
                  "type": "string"
                },
                "tone": {
                  "type": "string"
                },
                "include_web_search": {
                  "type": "boolean"
                }
              },
              "required": [
                "topic"
              ]
            }
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/posts/{id}": {
      "get": {
        "tags": [
          "Posts"
        ],
        "summary": "Detail artikel milik penulis",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      },
      "put": {
        "tags": [
          "Posts"
        ],
        "summary": "Perbarui artikel",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "title": {
                  "type": "string"
                },
                "category": {
                  "type": "string"
                },
                "excerpt": {
                  "type": "string"
                },
                "content": {
                  "type": "string"
                },
                "cover_image_path": {
                  "type": "string"
                },
                "is_published": {
                  "type": "boolean"
                }
              }
            }
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      },
      "delete": {
        "tags": [
          "Posts"
        ],
        "summary": "Hapus artikel",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/posts/{id}/toggle-publish": {
      "patch": {
        "tags": [
          "Posts"
        ],
        "summary": "Toggle status Draf/Terbit",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/admin/users": {
      "get": {
        "tags": [
          "Admin"
        ],
        "summary": "Seluruh pengguna + total artikel",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/admin/posts": {
      "get": {
        "tags": [
          "Admin"
        ],
        "summary": "Seluruh artikel sistem",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/admin/users/{id}/posts": {
      "get": {
        "tags": [
          "Admin"
        ],
        "summary": "Artikel milik pengguna tertentu",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "success"
          },
          "201": {
            "description": "created"
          },
          "400": {
            "description": "fail"
          },
          "401": {
            "description": "unauthorized"
          },
          "403": {
            "description": "forbidden"
          },
          "404": {
            "description": "not found"
          }
        },
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    }
  },
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "Format: Bearer <token>"
    }
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Delcom Blog API",
	Description:      "REST API Delcom Blog",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
