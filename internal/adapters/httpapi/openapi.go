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
					"enum":        []any{"noop", "sleep", "download", "spawn", "wait"},
				},
				"Subscription": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":                    map[string]any{"type": "string"},
						"baseUrl":               map[string]any{"type": "string"},
						"label":                 map[string]any{"type": "string"},
						"player":                map[string]any{"type": "string", "description": "auto ou Player N"},
						"lastScheduledEpisode":  map[string]any{"type": "integer", "minimum": 0},
						"lastDownloadedEpisode": map[string]any{"type": "integer", "minimum": 0},
						"lastAvailableEpisode":  map[string]any{"type": "integer", "minimum": 0},
						"nextCheckAt":           map[string]any{"type": "string", "format": "date-time"},
						"lastCheckedAt":         map[string]any{"type": "string", "format": "date-time"},
						"createdAt":             map[string]any{"type": "string", "format": "date-time"},
						"updatedAt":             map[string]any{"type": "string", "format": "date-time"},
					},
					"required":             []any{"id", "baseUrl", "label", "player", "lastScheduledEpisode", "lastDownloadedEpisode", "lastAvailableEpisode", "nextCheckAt", "createdAt", "updatedAt"},
					"additionalProperties": false,
				},
				"SubscriptionList": map[string]any{
					"type":  "array",
					"items": map[string]any{"$ref": "#/components/schemas/Subscription"},
				},
				"CreateSubscriptionRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"baseUrl": map[string]any{"type": "string", "example": "https://anime-sama.si/catalogue/xxx/saison1/vostfr/"},
						"label":   map[string]any{"type": "string", "example": "Solo Leveling"},
						"player":  map[string]any{"type": "string", "example": "auto"},
					},
					"required":             []any{"baseUrl", "label"},
					"additionalProperties": false,
				},
				"SyncSubscriptionResult": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"subscription":        map[string]any{"$ref": "#/components/schemas/Subscription"},
						"selectedPlayer":      map[string]any{"type": "string"},
						"maxAvailableEpisode": map[string]any{"type": "integer"},
						"enqueuedEpisodes":    map[string]any{"type": "array", "items": map[string]any{"type": "integer"}},
						"enqueuedJobIds":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
						"message":             map[string]any{"type": "string"},
					},
					"required":             []any{"subscription", "selectedPlayer", "maxAvailableEpisode", "enqueuedEpisodes", "enqueuedJobIds"},
					"additionalProperties": false,
				},
				"SyncAllSubscriptionsResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"results": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/SyncSubscriptionResult"}},
						"errors": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"id":    map[string]any{"type": "string"},
									"error": map[string]any{"type": "string"},
								},
								"required":             []any{"id", "error"},
								"additionalProperties": false,
							},
						},
					},
					"required":             []any{"results", "errors"},
					"additionalProperties": false,
				},
				"AniListViewer": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":   map[string]any{"type": "integer"},
						"name": map[string]any{"type": "string"},
					},
					"required":             []any{"id", "name"},
					"additionalProperties": false,
				},
				"AniListAiringScheduleEntry": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":       map[string]any{"type": "integer"},
						"airingAt": map[string]any{"type": "integer", "description": "Unix timestamp (seconds)"},
						"episode":  map[string]any{"type": "integer"},
						"media": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"id": map[string]any{"type": "integer"},
								"title": map[string]any{
									"type": "object",
									"properties": map[string]any{
										"romaji":  map[string]any{"type": "string"},
										"english": map[string]any{"type": "string"},
										"native":  map[string]any{"type": "string"},
									},
									"additionalProperties": false,
								},
							},
							"required":             []any{"id", "title"},
							"additionalProperties": false,
						},
					},
					"required":             []any{"id", "airingAt", "episode", "media"},
					"additionalProperties": false,
				},
				"AniListAiringScheduleList": map[string]any{
					"type":  "array",
					"items": map[string]any{"$ref": "#/components/schemas/AniListAiringScheduleEntry"},
				},
				"AniListWatchlistEntry": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"status":   map[string]any{"type": "string"},
						"progress": map[string]any{"type": "integer"},
						"media": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"id":       map[string]any{"type": "integer"},
								"synonyms": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
								"title": map[string]any{
									"type": "object",
									"properties": map[string]any{
										"romaji":  map[string]any{"type": "string"},
										"english": map[string]any{"type": "string"},
										"native":  map[string]any{"type": "string"},
									},
									"additionalProperties": false,
								},
							},
							"required":             []any{"id", "title"},
							"additionalProperties": false,
						},
					},
					"required":             []any{"status", "progress", "media"},
					"additionalProperties": false,
				},
				"AniListWatchlist": map[string]any{
					"type":  "array",
					"items": map[string]any{"$ref": "#/components/schemas/AniListWatchlistEntry"},
				},
				"AniListImportPreviewRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"statuses":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "example": []any{"CURRENT", "PLANNING"}},
						"season":        map[string]any{"type": "integer", "minimum": 1, "example": 1},
						"lang":          map[string]any{"type": "string", "example": "vostfr"},
						"maxCandidates": map[string]any{"type": "integer", "minimum": 1, "maximum": 10, "example": 3},
					},
					"additionalProperties": false,
				},
				"AniListImportPreviewResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"items": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"anilistMediaId": map[string]any{"type": "integer"},
									"title":          map[string]any{"type": "string"},
									"titles":         map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}},
									"synonyms":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
									"candidates": map[string]any{
										"type": "array",
										"items": map[string]any{
											"type": "object",
											"properties": map[string]any{
												"catalogueUrl": map[string]any{"type": "string"},
												"baseUrl":      map[string]any{"type": "string"},
												"slug":         map[string]any{"type": "string"},
												"matchedTitle": map[string]any{"type": "string"},
												"score":        map[string]any{"type": "number", "format": "double"},
											},
											"additionalProperties": false,
										},
									},
								},
								"additionalProperties": false,
							},
						},
					},
					"required":             []any{"items"},
					"additionalProperties": false,
				},
				"AniListImportConfirmRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"items": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"baseUrl": map[string]any{"type": "string"},
									"label":   map[string]any{"type": "string"},
									"player":  map[string]any{"type": "string"},
								},
								"required":             []any{"baseUrl", "label"},
								"additionalProperties": false,
							},
						},
					},
					"required":             []any{"items"},
					"additionalProperties": false,
				},
				"AniListImportConfirmResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"created": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Subscription"}},
						"errors": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"baseUrl": map[string]any{"type": "string"},
									"error":   map[string]any{"type": "string"},
								},
								"additionalProperties": false,
							},
						},
					},
					"required":             []any{"created", "errors"},
					"additionalProperties": false,
				},
				"AniListImportAutoRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"statuses":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "example": []any{"CURRENT", "PLANNING"}},
						"season":        map[string]any{"type": "integer", "minimum": 1, "example": 1},
						"lang":          map[string]any{"type": "string", "example": "vostfr"},
						"maxCandidates": map[string]any{"type": "integer", "minimum": 1, "maximum": 10, "example": 3},
						"minScore":      map[string]any{"type": "number", "format": "double", "minimum": 0, "maximum": 1, "example": 0.95},
					},
					"additionalProperties": false,
				},
				"AniListImportAutoResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"created": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Subscription"}},
						"skipped": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"anilistMediaId": map[string]any{"type": "integer"},
									"title":          map[string]any{"type": "string"},
									"reason":         map[string]any{"type": "string"},
									"baseUrl":        map[string]any{"type": "string"},
									"topScore":       map[string]any{"type": "number", "format": "double"},
								},
								"additionalProperties": false,
							},
						},
						"errors": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"baseUrl": map[string]any{"type": "string"},
									"error":   map[string]any{"type": "string"},
								},
								"additionalProperties": false,
							},
						},
					},
					"required":             []any{"created", "skipped", "errors"},
					"additionalProperties": false,
				},
					"AnimeSamaResolveRequest": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"titles":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "example": []any{"Sousou no Frieren", "Frieren: Beyond Journey's End"}},
							"season":        map[string]any{"type": "integer", "minimum": 1, "example": 1},
							"lang":          map[string]any{"type": "string", "example": "vostfr"},
							"maxCandidates": map[string]any{"type": "integer", "minimum": 1, "maximum": 10, "example": 5},
						},
						"additionalProperties": false,
					},
					"AnimeSamaResolvedCandidate": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"catalogueUrl": map[string]any{"type": "string"},
							"baseUrl":      map[string]any{"type": "string"},
							"slug":         map[string]any{"type": "string"},
							"matchedTitle": map[string]any{"type": "string"},
							"score":        map[string]any{"type": "number", "format": "double"},
						},
						"additionalProperties": false,
					},
					"AnimeSamaResolveResponse": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"candidates": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/AnimeSamaResolvedCandidate"}},
						},
						"required":             []any{"candidates"},
						"additionalProperties": false,
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
						"anilistToken":           map[string]any{"type": "string"},
					},
					"additionalProperties": false,
				},
				"DownloadJobResult": map[string]any{
					"type":        "object",
					"description": "Résultat d'un job download.",
					"properties": map[string]any{
						"url":         map[string]any{"type": "string", "description": "URL réellement utilisée pour télécharger (après résolution)."},
						"sourceUrl":   map[string]any{"type": "string", "description": "URL initiale (avant résolution), utilisée comme Referer quand pertinent."},
						"resolvedUrl": map[string]any{"type": "string", "description": "Alias de url (pour debug)."},
						"mode":        map[string]any{"type": "string", "enum": []any{"http", "ffmpeg"}, "description": "Méthode utilisée pour le téléchargement."},
						"usedReferer": map[string]any{"type": "boolean", "description": "Indique si un header Referer a été envoyé."},
						"usedOrigin":  map[string]any{"type": "boolean", "description": "Indique si un header Origin a été envoyé."},
						"path":        map[string]any{"type": "string", "description": "Chemin final du fichier sur disque."},
						"bytes":       map[string]any{"type": "integer", "minimum": 0, "description": "Nombre d'octets écrits (peut être 0 en mode ffmpeg)."},
						"contentType": map[string]any{"type": "string"},
					},
					"required":             []any{"url", "path", "bytes"},
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
							"description": "Résultat du job (si applicable).",
							"anyOf": []any{
								map[string]any{"$ref": "#/components/schemas/DownloadJobResult"},
								map[string]any{"type": "object", "additionalProperties": true},
							},
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
						map[string]any{"$ref": "#/components/schemas/CreateSpawnJobRequest"},
						map[string]any{"$ref": "#/components/schemas/CreateWaitJobRequest"},
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
				"CreateSpawnJobRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string", "enum": []any{"spawn"}},
						"params": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"jobs": map[string]any{
									"type":        "array",
									"description": "Liste de jobs enfants à créer.",
									"items": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"type":   map[string]any{"$ref": "#/components/schemas/JobType"},
											"params": map[string]any{"type": "object", "additionalProperties": true},
										},
										"required":             []any{"type"},
										"additionalProperties": false,
									},
								},
							},
							"required":             []any{"jobs"},
							"additionalProperties": false,
						},
					},
					"required":             []any{"type", "params"},
					"additionalProperties": false,
				},
				"CreateWaitJobRequest": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string", "enum": []any{"wait"}},
						"params": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"jobIds": map[string]any{
									"type":        "array",
									"description": "Liste de jobs à attendre.",
									"items":       map[string]any{"type": "string"},
								},
								"failOnFailed": map[string]any{"type": "boolean", "description": "Si true, échoue dès qu'un enfant est failed/canceled.", "example": true},
								"timeoutMs":    map[string]any{"type": "integer", "description": "Timeout côté executor en ms.", "example": 600000},
								"pollMs":       map[string]any{"type": "integer", "description": "Intervalle de polling en ms.", "example": 250},
							},
							"required":             []any{"jobIds"},
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
			"/api/v1/subscriptions": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/SubscriptionList"),
						"500": jsonErr,
					},
				},
				"post": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/CreateSubscriptionRequest"},
							},
						},
					},
					"responses": map[string]any{
						"201": jsonOK("#/components/schemas/Subscription"),
						"400": jsonErr,
						"500": jsonErr,
					},
				},
			},
			"/api/v1/subscriptions/{id}": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/Subscription"),
						"404": jsonErr,
						"500": jsonErr,
					},
				},
				"put": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"409":    jsonErr,
								"schema": map[string]any{"$ref": "#/components/schemas/Subscription"},
							},
						},
					},
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/Subscription"),
						"400": jsonErr,
						"404": jsonErr,
						"500": jsonErr,
					},
				},
				"delete": map[string]any{
					"responses": map[string]any{
						"204": map[string]any{"description": "No Content"},
						"404": jsonErr,
						"500": jsonErr,
					},
				},
			},
			"/api/v1/subscriptions/{id}/sync": map[string]any{
				"post": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/SyncSubscriptionResult"),
						"400": jsonErr,
						"404": jsonErr,
						"500": jsonErr,
					},
				},
			},
			"/api/v1/subscriptions/sync-all": map[string]any{
				"post": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/SyncAllSubscriptionsResponse"),
						"500": jsonErr,
					},
				},
			},
			"/api/v1/anilist/viewer": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AniListViewer"),
						"400": jsonErr,
						"502": jsonErr,
					},
				},
			},
			"/api/v1/anilist/airing": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AniListAiringScheduleList"),
						"502": jsonErr,
					},
				},
			},
			"/api/v1/anilist/watchlist": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AniListWatchlist"),
						"400": jsonErr,
						"502": jsonErr,
					},
				},
			},
			"/api/v1/animesama/resolve": map[string]any{
				"post": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/AnimeSamaResolveRequest"},
							},
						},
					},
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AnimeSamaResolveResponse"),
						"400": jsonErr,
						"501": jsonErr,
						"502": jsonErr,
					},
				},
			},
			"/api/v1/import/anilist/preview": map[string]any{
				"post": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/AniListImportPreviewRequest"},
							},
						},
					},
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AniListImportPreviewResponse"),
						"400": jsonErr,
						"502": jsonErr,
					},
				},
			},
			"/api/v1/import/anilist/confirm": map[string]any{
				"post": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/AniListImportConfirmRequest"},
							},
						},
					},
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AniListImportConfirmResponse"),
						"400": jsonErr,
					},
				},
			},
			"/api/v1/import/anilist/auto": map[string]any{
				"post": map[string]any{
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/AniListImportAutoRequest"},
							},
						},
					},
					"responses": map[string]any{
						"200": jsonOK("#/components/schemas/AniListImportAutoResponse"),
						"400": jsonErr,
						"502": jsonErr,
					},
				},
			},
		},
	}

	httpjson.Write(w, http.StatusOK, spec)
}
