// main.go
package main

import (
	"embed"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

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
	Char      string
	X         int
	Y         int
	Speed     float64
	Active    bool
	TargetX   int
	TargetY   int
	StartTime time.Time
	Lane      int
}

type KeyPress struct {
	Key string `json:"key"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	State   *GameState  `json:"state,omitempty"`
	Letter  *FallingLetter `json:"letter,omitempty"`
}

var gameState *GameState

func init() {
	rand.Seed(time.Now().UnixNano())
	
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
	// Serve static files
	http.HandleFunc("/", serveFrontend)
	
	// API endpoints
	http.HandleFunc("/api/start", startGame)
	http.HandleFunc("/api/keypress", handleKeyPress)
	http.HandleFunc("/api/state", getState)
	http.HandleFunc("/api/reset", resetGame)
	http.HandleFunc("/api/falling", getFallingLetters)
	
	// Start background goroutine for falling letters
	go updateFallingLetters()
	
	log.Println("🎮 Touch Typing Master - Server starting on http://localhost:8080")
	log.Println("⌨️  Practice all keyboard letters for Go programming!")
	log.Println("📦  All assets are self-contained - works offline!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Could not load frontend", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
		return
	}
	http.NotFound(w, r)
}

func startGame(w http.ResponseWriter, r *http.Request) {
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
	
	generateFallingLetters()
	
	respondJSON(w, Response{
		Success: true,
		Message: "Game started! Type the falling letters!",
		State:   gameState,
	})
}

func resetGame(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
				generateFallingLetters()
				
				respondJSON(w, Response{
					Success: true,
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
	
	// Add new falling letter if needed
	if len(gameState.FallingLetters) < 25 {
		addFallingLetter()
	}
	
	respondJSON(w, Response{
		Success: true,
		State:   gameState,
	})
}

func getState(w http.ResponseWriter, r *http.Request) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()
	
	respondJSON(w, Response{
		Success: true,
		State:   gameState,
	})
}

func getFallingLetters(w http.ResponseWriter, r *http.Request) {
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

func generateFallingLetters() {
	text := gameState.CurrentText
	if len(text) == 0 {
		return
	}
	
	// Get unique characters from text
	charMap := make(map[string]bool)
	for _, ch := range text {
		charStr := string(ch)
		if charStr != " " {
			charMap[charStr] = true
		}
	}
	
	chars := []string{}
	for ch := range charMap {
		chars = append(chars, ch)
	}
	
	if len(chars) == 0 {
		return
	}
	
	// Create lanes for letters (20 lanes across the screen)
	numLanes := 20
	laneWidth := 800 / numLanes
	
	// Assign each unique character to a lane
	charLanes := make(map[string]int)
	for i, ch := range chars {
		charLanes[ch] = i % numLanes
	}
	
	// Generate 15-20 falling letters
	numLetters := 15 + rand.Intn(6)
	for i := 0; i < numLetters && i < len(chars)*2; i++ {
		char := chars[rand.Intn(len(chars))]
		lane := charLanes[char]
		x := lane*laneWidth + laneWidth/2 - 15
		y := -50 - rand.Intn(200) // Start above screen
		speed := 0.8 + rand.Float64()*0.7
		
		gameState.FallingLetters = append(gameState.FallingLetters, FallingLetter{
			Char:      char,
			X:         x,
			Y:         y,
			Speed:     speed,
			Active:    true,
			TargetX:   x,
			TargetY:   400 + rand.Intn(100),
			Lane:      lane,
			StartTime: time.Now(),
		})
	}
}

func addFallingLetter() {
	charMap := make(map[string]bool)
	for _, ch := range gameState.CurrentText {
		charStr := string(ch)
		if charStr != " " {
			charMap[charStr] = true
		}
	}
	
	chars := []string{}
	for ch := range charMap {
		chars = append(chars, ch)
	}
	
	if len(chars) == 0 {
		return
	}
	
	char := chars[rand.Intn(len(chars))]
	
	// Find a free lane for this character
	numLanes := 20
	laneWidth := 800 / numLanes
	
	// Try to find a lane that doesn't have too many letters
	laneCounts := make([]int, numLanes)
	for _, letter := range gameState.FallingLetters {
		if letter.Active && letter.Lane >= 0 && letter.Lane < numLanes {
			laneCounts[letter.Lane]++
		}
	}
	
	// Find lane with fewest letters
	minLane := 0
	minCount := laneCounts[0]
	for i := 1; i < numLanes; i++ {
		if laneCounts[i] < minCount {
			minCount = laneCounts[i]
			minLane = i
		}
	}
	
	lane := minLane
	x := lane*laneWidth + laneWidth/2 - 15
	y := -50 - rand.Intn(100)
	speed := 0.8 + rand.Float64()*0.7
	
	gameState.FallingLetters = append(gameState.FallingLetters, FallingLetter{
		Char:      char,
		X:         x,
		Y:         y,
		Speed:     speed,
		Active:    true,
		TargetX:   x,
		TargetY:   450 + rand.Intn(50),
		Lane:      lane,
		StartTime: time.Now(),
	})
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
		
		for i := range gameState.FallingLetters {
			if gameState.FallingLetters[i].Active {
				// Move letter down
				gameState.FallingLetters[i].Y += int(gameState.FallingLetters[i].Speed * 3)
				
				// Check if letter passed the target zone
				if gameState.FallingLetters[i].Y > 600 {
					// Letter missed - count as mistake
					gameState.FallingLetters[i].Active = false
					gameState.Mistakes++
					gameState.Combo = 0
					
					// Only if game is active
					if gameState.IsActive && !gameState.Completed {
						// Penalty: reduce score
						gameState.Score = max(0, gameState.Score-5)
					}
				}
			}
		}
		
		// Add new falling letters periodically
		if gameState.IsActive && !gameState.Completed {
			activeCount := 0
			for _, letter := range gameState.FallingLetters {
				if letter.Active {
					activeCount++
				}
			}
			
			// Maintain 15-20 active letters
			if activeCount < 15 && len(gameState.FallingLetters) < 40 {
				addFallingLetter()
			}
		}
		
		gameState.mu.Unlock()
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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}
