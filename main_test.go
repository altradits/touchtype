package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postJSON(t *testing.T, path, body string) Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	recorder := httptest.NewRecorder()
	if path == "/api/keypress" {
		handleKeyPress(recorder, req)
	} else {
		startGame(recorder, req)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("%s returned status %d", path, recorder.Code)
	}
	var response Response
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode %s response: %v", path, err)
	}
	return response
}

func TestGameCompletesEveryLessonInOrder(t *testing.T) {
	started := postJSON(t, "/api/start", "")
	if !started.Success || started.State.CurrentText == "" {
		t.Fatal("game did not start with a lesson")
	}

	for _, lesson := range started.State.LessonLines {
		for _, char := range lesson {
			response := postJSON(t, "/api/keypress", `{"key":`+mustJSON(t, string(char))+`}`)
			if !response.Success {
				t.Fatalf("keypress %q was rejected: %s", char, response.Message)
			}
		}
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()
	if !gameState.Completed || gameState.IsActive {
		t.Fatal("game should be completed and inactive after the final lesson")
	}
	if gameState.Position != len(gameState.CurrentText) {
		t.Fatalf("final position = %d, want %d", gameState.Position, len(gameState.CurrentText))
	}
}

func TestWrongKeyDoesNotAdvanceAndNextPromptIsExactKey(t *testing.T) {
	postJSON(t, "/api/start", "")
	wrong := postJSON(t, "/api/keypress", `{"key":"x"}`)
	if wrong.Correct || wrong.State.Position != 0 || wrong.State.Mistakes != 1 || wrong.State.Combo != 0 {
		t.Fatalf("wrong key changed state incorrectly: %+v", wrong.State)
	}
	correct := postJSON(t, "/api/keypress", `{"key":"a"}`)
	if !correct.Correct || correct.State.Position != 1 {
		t.Fatalf("correct key was not acknowledged: %+v", correct)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()
	if len(gameState.FallingLetters) != 1 || gameState.FallingLetters[0].Char != "s" {
		t.Fatalf("prompt = %+v, want one prompt for s", gameState.FallingLetters)
	}
}

func TestStartRejectsGet(t *testing.T) {
	recorder := httptest.NewRecorder()
	startGame(recorder, httptest.NewRequest(http.MethodGet, "/api/start", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /api/start status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}

func TestHealthCheck(t *testing.T) {
	recorder := httptest.NewRecorder()
	healthCheck(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusOK || recorder.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected health response: status=%d body=%q", recorder.Code, recorder.Body.String())
	}
}

func TestDownloadManifestUsesConfiguredBaseURL(t *testing.T) {
	t.Setenv("DOWNLOAD_BASE_URL", "https://downloads.example.com/touchtype/")
	recorder := httptest.NewRecorder()
	getDownloads(recorder, httptest.NewRequest(http.MethodGet, "/api/downloads", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("download manifest status = %d", recorder.Code)
	}
	var manifest struct {
		Configured bool                `json:"configured"`
		Downloads  map[string]Download `json:"downloads"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&manifest); err != nil {
		t.Fatal(err)
	}
	if !manifest.Configured || manifest.Downloads["windows-amd64"].URL != "https://downloads.example.com/touchtype/touchtype-master-windows-amd64.exe" {
		t.Fatalf("unexpected download manifest: %+v", manifest)
	}
}

func TestMissRespawnsTheSamePromptAndAppliesPenalty(t *testing.T) {
	postJSON(t, "/api/start", "")
	gameState.mu.Lock()
	defer gameState.mu.Unlock()
	gameState.Combo = 4
	gameState.Score = 10
	gameState.FallingLetters[0].Y = 240
	advanceFallingLetters()

	if gameState.Mistakes != 1 || gameState.Combo != 0 || gameState.Score != 5 {
		t.Fatalf("miss did not apply the expected penalty: %+v", gameState)
	}
	if len(gameState.FallingLetters) != 1 || !gameState.FallingLetters[0].Active || gameState.FallingLetters[0].Char != "a" {
		t.Fatalf("miss did not respawn the active prompt: %+v", gameState.FallingLetters)
	}
}

func mustJSON(t *testing.T, value string) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
