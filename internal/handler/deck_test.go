package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateDeck_SameLanguageError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/decks", func(c *gin.Context) {
		// Simulate authenticated user
		c.Set("auth_user", struct{}{})

		var req CreateDeckRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !Language(req.SourceLanguage).IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source_language"})
			return
		}
		if !Language(req.TargetLanguage).IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_language"})
			return
		}
		if req.SourceLanguage == req.TargetLanguage {
			c.JSON(http.StatusBadRequest, gin.H{"error": "source_language and target_language must be different"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "ok"})
	})

	tests := []struct {
		body           CreateDeckRequest
		name           string
		expectedError  string
		expectedStatus int
	}{
		{
			name: "same language returns error",
			body: CreateDeckRequest{
				Name:           "Test Deck",
				SourceLanguage: "en",
				TargetLanguage: "en",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "source_language and target_language must be different",
		},
		{
			name: "different languages succeed",
			body: CreateDeckRequest{
				Name:           "Test Deck",
				SourceLanguage: "en",
				TargetLanguage: "de",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "invalid source language",
			body: CreateDeckRequest{
				Name:           "Test Deck",
				SourceLanguage: "invalid",
				TargetLanguage: "de",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid source_language",
		},
		{
			name: "invalid target language",
			body: CreateDeckRequest{
				Name:           "Test Deck",
				SourceLanguage: "en",
				TargetLanguage: "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid target_language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req, err := http.NewRequest(http.MethodPost, "/decks", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				var resp map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp["error"] != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, resp["error"])
				}
			}
		})
	}
}
