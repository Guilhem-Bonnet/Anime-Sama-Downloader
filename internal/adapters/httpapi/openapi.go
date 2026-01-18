package httpapi

import (
	"net/http"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
)

// handleOpenAPI renvoie une spec OpenAPI minimale (placeholder) pour cadrer l'API.
// Elle sera enrichie au fil des jalons.
func (s *Server) handleOpenAPI(w http.ResponseWriter, r *http.Request) {
	jsonOK := func(schemaRef string) map[string]any {
		return map[string]any{
			"description": "OK",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": map[string]any{"$ref": schemaRef},
				},
			},
		}
	}

	jsonErr := map[string]any{
		"description": "Error",
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{"$ref": "#/components/schemas/Error"},
			},
		},
	}

	spec := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   "ASD API",
			"version": "v1",
		},
		"components": map[string]any{
			"schemas": map[string]any{
				"OpenAPIDocument": map[string]any{
					"type":                 "object",
					"additionalProperties": true,
				},
				"JobType": map[string]any{
					"type":        "string",
					"description": "Type de job (extensible).",
					"enum":        []any{"noop", "sleep", "download"},
				},
				"Error": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"error": map[string]any{"type": "string"},
					},
					"required": []any{"error"},
				},
				"Settings": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"destination":            map[string]any{"type": "string"},
						"outputNamingMode":       map[string]any{"type": "string", "enum": []any{"legacy", "media-server"}},
						"separateLang":           map[string]any{"type": "boolean"},
						"maxWorkers":             map[string]any{"type": "integer", "minimum": 1},
						"maxConcurrentDownloads": map[string]any{"type": "integer", "minimum": 1},
						"jellyfinUrl":            map[string]any{"type": "string"},
						"jellyfinApiKey":         map[string]any{"type": "string"},
						"plexUrl":                map[string]any{"type": "string"},
						"plexToken":              map[string]any{"type": "string"},
						"plexSectionId":          map[string]any{"type": "string"},
					},
					"additionalProperties": false,
				},
				"Job": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":        map[string]any{"type": "string"},
						"type":      map[string]any{"$ref": "#/components/schemas/JobType"},
						"state":     map[string]any{"type": "string", "enum": []any{"queued", "running", "muxing", "completed", "failed", "canceled"}},
						"progress":  map[string]any{"type": "number", "format": "double"},
						"createdAt": map[string]any{"type": "string", "format": "date-time"},
						"updatedAt": map[string]any{"type": "string", "format": "date-time"},
						"params": map[string]any{
							"type":                 "object",
							"description":          "Paramètres du job (dépend du type).",
							"additionalProperties": true,
						},
						"result": map[string]any{
							"type":                 "object",
							"description":          "Résultat du job (si applicable).",
							"additionalProperties": true,
						},
						"errorCode": map[string]any{"type": "string"},
						"error":     map[string]any{"type": "string"},
					},
					"required":             []any{"id", "type", "state", "progress", "createdAt", "updatedAt"},
					"additionalProperties": false,
				},
				"JobList": map[string]any{
					"type":  "array",
					"items": map[string]any{"$ref": "#/components/schemas/Job"},
				},
				"CreateJobRequest": map[string]any{
					"oneOf": []any{
						map[string]any{"$ref": "#/components/schemas/CreateNoopJobRequest"},
						map[string]any{"$ref": "#/components/schemas/CreateSleepJobRequest"},
						map[string]any{"$ref": "#/components/schemas/CreateDownloadJobRequest"},
					},
					"description": "Requête de création d'un job. Les params dépendent du type.",
				},
				"CreateNoopJobRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string", "enum": []any{"noop"}},
						"params": map[string]any{
							"type":                 "object",
							"additionalProperties": true,
							"example":              map[string]any{},
						},
					},
					"required":             []any{"type"},
					"additionalProperties": false,
				},
				"CreateSleepJobRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string", "enum": []any{"sleep"}},
						"params": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"duration":   map[string]any{"type": "string", "description": "Durée Go (ex: 250ms, 2s, 1m)", "example": "2s"},
								"durationMs": map[string]any{"type": "integer", "description": "Durée en ms", "example": 2000},
								"seconds":    map[string]any{"type": "integer", "description": "Durée en secondes", "example": 2},
							},
							"additionalProperties": false,
						},
					},
					"required":             []any{"type"},
					"additionalProperties": false,
				},
				"CreateDownloadJobRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string", "enum": []any{"download"}},
						"params": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"url":      map[string]any{"type": "string", "description": "URL source (http/https)", "example": "https://example.com/video.bin"},
								"filename": map[string]any{"type": "string", "description": "Nom de fichier de sortie (dans settings.destination)", "example": "episode-01.mp4"},
								"path":     map[string]any{"type": "string", "description": "Chemin relatif de sortie (dans settings.destination)", "example": "series/season-1/episode-01.mp4"},
							},
							"required":             []any{"url"},
							"additionalProperties": false,
						},
					},
					"required":             []any{"type", "params"},
					"additionalProperties": false,
				},
			},
		},
		"paths": map[string]any{
			"/api/v1/health": map[string]any{
				"get": map[string]any{"responses": map[string]any{"200": map[string]any{"description": "OK"}}},
			},
			"/api/v1/version": map[string]any{
				"get": map[string]any{"responses": map[string]any{"200": map[string]any{"description": "OK"}}},
			},
			"/api/v1/openapi.json": map[string]any{
				"get": map[string]any{"responses": map[string]any{"200": jsonOK("#/components/schemas/OpenAPIDocument")}},
			},
			"/api/v1/events": map[string]any{
				"get": map[string]any{"responses": map[string]any{"200": map[string]any{"description": "SSE"}}},
			},
			"/api/v1/jobs": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/JobList"),
						"500": jsonErr,
					},
				},
				"post": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/CreateJobRequest"},
							},
						},
					},
					"responses": map[string]any{
						"201": jsonOK("#/components/schemas/Job"),
						"400": jsonErr,
						"500": jsonErr,
					},
				},
			},
			"/api/v1/jobs/{id}": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/Job"),
						"404": jsonErr,
						"500": jsonErr,
					},
				},
			},
			"/api/v1/jobs/{id}/cancel": map[string]any{
				"post": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/Job"),
						"404": jsonErr,
						"500": jsonErr,
					},
				},
			},
			"/api/v1/settings": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/Settings"),
						"500": jsonErr,
					},
				},
				"put": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Settings"},
							},
						},
					},
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/Settings"),
						"400": jsonErr,
						"500": jsonErr,
					},
				},
			},
		},
	}

	httpjson.Write(w, http.StatusOK, spec)
}
