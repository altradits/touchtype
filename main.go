// main.go
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

// Version and author information - set during build
var (
	version        = "1.0.0"
	buildTime      = "unknown"
	gitCommit      = "unknown"
	authorName     = "Stanley Chege Thuita"
	authorEmail    = "altradits@gmail.com"
	authorLinkedIn = "https://www.linkedin.com/in/stanmobitech/"
	authorPhone    = "+254707172370"
)

type GameState struct {
	mu             sync.Mutex
	CurrentText    string          `json:"currentText"`
	Position       int             `json:"position"`
	StartTime      time.Time       `json:"startTime"`
	IsActive       bool            `json:"isActive"`
	WPM            float64         `json:"wpm"`
	Accuracy       float64         `json:"accuracy"`
	TotalChars     int             `json:"totalChars"`
	CorrectChars   int             `json:"correctChars"`
	Mistakes       int             `json:"mistakes"`
	LessonLines    []string        `json:"lessonLines"`
	CurrentLine    int             `json:"currentLine"`
	Completed      bool            `json:"completed"`
	FallingLetters []FallingLetter `json:"-"`
	Score          int             `json:"score"`
	Combo          int             `json:"combo"`
	MaxCombo       int             `json:"maxCombo"`
}

type FallingLetter struct {
	Char    string `json:"char"`
	Display string `json:"display"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Active  bool   `json:"active"`
}

type KeyPress struct {
	Key string `json:"key"`
}

type Response struct {
	Success bool           `json:"success"`
	Correct bool           `json:"correct,omitempty"`
	Message string         `json:"message,omitempty"`
	State   *GameState     `json:"state,omitempty"`
	Letter  *FallingLetter `json:"letter,omitempty"`
}

var gameState *GameState

func init() {
	// Progressive lessons covering all keyboard letters
	lessonLines := []string{
		// Home Row Basics
		"asdf fdsa asdf fdsa asdf fdsa asdf fdsa asdf fdsa asdf fdsa",
		"jkl; ;lkj jkl; ;lkj jkl; ;lkj jkl; ;lkj jkl; ;lkj jkl; ;lkj",

		// Home Row Combinations
		"asdf jkl; asdf jkl; asdf jkl; asdf jkl; asdf jkl; asdf jkl;",
		";lkj fdsa ;lkj fdsa ;lkj fdsa ;lkj fdsa ;lkj fdsa ;lkj fdsa",

		// Top Row (QWERTY)
		"qwerty ytrewq qwerty ytrewq qwerty ytrewq qwerty ytrewq",
		"poiuyt yuiop poiuyt yuiop poiuyt yuiop poiuyt yuiop",
		"qwertyuiop poiuytrewq qwertyuiop poiuytrewq qwertyuiop",

		// Bottom Row (ZXCVBNM)
		"zxcvbnm mnbvcxz zxcvbnm mnbvcxz zxcvbnm mnbvcxz zxcvbnm",
		"asdfghjkl zxcvbnm asdfghjkl zxcvbnm asdfghjkl zxcvbnm",

		// All Letters Mixed
		"the quick brown fox jumps over the lazy dog",
		"pack my box with five dozen liquor jugs",
		"how vexingly quick daft zebras jump",
		"the five boxing wizards jump quickly",

		// Go Programming Keywords
		"package main import fmt func var const type",
		"if else for range switch defer go select",
		"chan map interface struct error nil true false",
		"make new append copy delete close recover",
		"go func main println printf print panic",

		// Go Code Snippets
		"fmt.Println hello world",
		"var x int = 42",
		"if err != nil { return err }",
		"for i := 0; i < 10; i++ {",
		"switch value { case 1: default: }",

		// Complex Combinations
		"func main() { fmt.Println(\"Hello, World!\") }",
		"type User struct { Name string Age int }",
		"go run main.go go build go test go mod tidy",
		"make(chan int, 10) <-chan int chan<- int",

		// Final Challenge - Full Keyboard
		"the quick brown fox jumps over the lazy dog 1234567890",
		"!@#$%^&*()_+-=[]{}|;:',.<>?/~`",
		"abcdefghijklmnopqrstuvwxyz 0123456789",
	}

	gameState = &GameState{
		LessonLines:    lessonLines,
		CurrentLine:    0,
		IsActive:       false,
		FallingLetters: []FallingLetter{},
		Score:          0,
		Combo:          0,
		MaxCombo:       0,
	}

	if len(lessonLines) > 0 {
		gameState.CurrentText = lessonLines[0]
	}
}

func main() {
	// Print author and version information
	fmt.Printf("🎮 Touch Typing Master - v%s\n", version)
	fmt.Printf("👤 Author: %s\n", authorName)
	fmt.Printf("📧 Email: %s\n", authorEmail)
	fmt.Printf("🔗 LinkedIn: %s\n", authorLinkedIn)
	fmt.Printf("📱 Phone: %s\n", authorPhone)
	fmt.Printf("🕐 Build Time: %s\n", buildTime)
	fmt.Printf("🔀 Git Commit: %s\n", gitCommit)
	fmt.Println(strings.Repeat("=", 60))

	// Serve static files
	http.HandleFunc("/", serveFrontend)

	// API endpoints
	http.HandleFunc("/api/start", startGame)
	http.HandleFunc("/api/keypress", handleKeyPress)
	http.HandleFunc("/api/state", getState)
	http.HandleFunc("/api/reset", resetGame)
	http.HandleFunc("/api/falling", getFallingLetters)
	http.HandleFunc("/healthz", healthCheck)

	// Start background goroutine for falling letters
	go updateFallingLetters()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port
	log.Printf("🎮 Touch Typing Master - Server starting on http://localhost%s", address)
	log.Println("⌨️  Practice all keyboard letters for Go programming!")
	log.Println("📦  All assets are self-contained - works offline!")
	log.Printf("👤  Developed by %s\n", authorName)
	log.Fatal(http.ListenAndServe(address, nil))
}

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Could not load frontend", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(content)
		return
	}
	http.NotFound(w, r)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	gameState.CurrentLine = 0
	gameState.CurrentText = gameState.LessonLines[0]
	gameState.Position = 0
	gameState.IsActive = true
	gameState.StartTime = time.Now()
	gameState.WPM = 0
	gameState.Accuracy = 100
	gameState.TotalChars = 0
	gameState.CorrectChars = 0
	gameState.Mistakes = 0
	gameState.Completed = false
	gameState.Score = 0
	gameState.Combo = 0
	gameState.MaxCombo = 0
	gameState.FallingLetters = []FallingLetter{}

	syncFallingLetter()

	respondJSON(w, Response{
		Success: true,
		Message: "Game started! Type the falling letters!",
		State:   gameState,
	})
}

func resetGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	gameState.CurrentLine = 0
	gameState.CurrentText = gameState.LessonLines[0]
	gameState.Position = 0
	gameState.IsActive = false
	gameState.WPM = 0
	gameState.Accuracy = 100
	gameState.TotalChars = 0
	gameState.CorrectChars = 0
	gameState.Mistakes = 0
	gameState.Completed = false
	gameState.Score = 0
	gameState.Combo = 0
	gameState.MaxCombo = 0
	gameState.FallingLetters = []FallingLetter{}

	respondJSON(w, Response{
		Success: true,
		Message: "Game reset",
		State:   gameState,
	})
}

func handleKeyPress(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var keyPress KeyPress
	if err := json.NewDecoder(r.Body).Decode(&keyPress); err != nil {
		respondJSON(w, Response{Success: false, Message: "Invalid request"})
		return
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if !gameState.IsActive || gameState.Completed {
		respondJSON(w, Response{Success: false, Message: "Game not active"})
		return
	}

	// Check if the key matches the expected character
	if gameState.Position >= len(gameState.CurrentText) {
		respondJSON(w, Response{Success: false, Message: "Line complete"})
		return
	}

	expected := string(gameState.CurrentText[gameState.Position])
	key := keyPress.Key

	gameState.TotalChars++

	if key == expected {
		gameState.CorrectChars++
		gameState.Position++
		gameState.Combo++
		gameState.Score += 10 * gameState.Combo // Combo multiplier

		if gameState.Combo > gameState.MaxCombo {
			gameState.MaxCombo = gameState.Combo
		}

		// Remove the correct falling letter
		removeFallingLetter(key)

		// Check if line is complete
		if gameState.Position >= len(gameState.CurrentText) {
			gameState.CurrentLine++
			if gameState.CurrentLine < len(gameState.LessonLines) {
				// Move to next line
				gameState.CurrentText = gameState.LessonLines[gameState.CurrentLine]
				gameState.Position = 0
				gameState.FallingLetters = []FallingLetter{}
				syncFallingLetter()

				respondJSON(w, Response{
					Success: true,
					Correct: true,
					Message: "✨ Level Complete! Moving to next lesson!",
					State:   gameState,
				})
				return
			} else {
				// All lessons complete
				gameState.Completed = true
				gameState.IsActive = false
				calculateStats()
				respondJSON(w, Response{
					Success: true,
					Correct: true,
					Message: "🏆 MASTER TYPIST! You've conquered all lessons!",
					State:   gameState,
				})
				return
			}
		}
	} else {
		gameState.Mistakes++
		gameState.Combo = 0 // Reset combo on mistake
	}

	// Update stats
	calculateStats()

	syncFallingLetter()

	respondJSON(w, Response{
		Success: true,
		Correct: key == expected,
		State:   gameState,
	})
}

func getState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	respondJSON(w, Response{
		Success: true,
		State:   gameState,
	})
}

func getFallingLetters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	active := []FallingLetter{}
	for _, letter := range gameState.FallingLetters {
		if letter.Active {
			active = append(active, letter)
		}
	}

	respondJSON(w, map[string]interface{}{
		"success": true,
		"letters": active,
	})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	respondJSON(w, map[string]string{"status": "ok"})
}

// syncFallingLetter shows one useful prompt: the exact key required next.
// A single target keeps the visual game aligned with the lesson rather than
// asking learners to ignore unrelated falling characters.
func syncFallingLetter() {
	gameState.FallingLetters = nil
	if !gameState.IsActive || gameState.Position >= len(gameState.CurrentText) {
		return
	}

	char := string(gameState.CurrentText[gameState.Position])
	display := char
	if char == " " {
		display = "␣"
	}
	gameState.FallingLetters = []FallingLetter{{
		Char: char, Display: display, X: 50, Y: -40, Active: true,
	}}
}

func removeFallingLetter(char string) {
	for i := range gameState.FallingLetters {
		if gameState.FallingLetters[i].Active && gameState.FallingLetters[i].Char == char {
			gameState.FallingLetters[i].Active = false
			break
		}
	}
}

func updateFallingLetters() {
	ticker := time.NewTicker(30 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		gameState.mu.Lock()
		advanceFallingLetters()
		gameState.mu.Unlock()
	}
}

func advanceFallingLetters() {
	if !gameState.IsActive || gameState.Completed {
		return
	}
	for i := range gameState.FallingLetters {
		if !gameState.FallingLetters[i].Active {
			continue
		}
		gameState.FallingLetters[i].Y += 2
		if gameState.FallingLetters[i].Y > 240 {
			gameState.FallingLetters[i].Active = false
			gameState.Mistakes++
			gameState.Combo = 0
			gameState.Score = max(0, gameState.Score-5)
		}
	}
	if len(gameState.FallingLetters) == 0 || !gameState.FallingLetters[0].Active {
		syncFallingLetter()
	}
}

func calculateStats() {
	elapsed := time.Since(gameState.StartTime).Minutes()
	if elapsed > 0 && gameState.TotalChars > 0 {
		gameState.WPM = (float64(gameState.TotalChars) / 5.0) / elapsed
		gameState.Accuracy = float64(gameState.CorrectChars) / float64(gameState.TotalChars) * 100
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	respondJSON(w, Response{Success: false, Message: "Method not allowed"})
}
