package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	boardSize   = 10
	finalTile   = 56
	obstaclePct = 10
	numPlayers  = 2
)

type Player struct {
	ID     int
	Pieces []int
}

type GameState struct {
	Players     []Player
	Board       [][]string
	Obstacles   []int
	Round       int
	GameOver    bool
	GameStarted bool
	WinnerID    int
	Mutex       sync.Mutex
}

var (
	board     [][]string
	obstacles []int
	gameState GameState
	server    *http.Server
	gameOver  chan struct{}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	board = make([][]string, boardSize)
	for i := range board {
		board[i] = make([]string, boardSize)
	}

	obstacles = placeRandomObstacles(obstaclePct, finalTile)

	gameState = GameState{
		Players:     make([]Player, 0),
		Board:       board,
		Obstacles:   obstacles,
		Round:       0,
		GameOver:    false,
		GameStarted: false,
	}

	gameOver = make(chan struct{})

	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/state", stateHandler)

	server = &http.Server{Addr: ":8080"}

	go func() {
		fmt.Println("Host is running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Error en el servidor HTTP:", err)
		}
	}()

	<-gameOver

	server.Shutdown(nil)
	fmt.Println("Server closed")
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	gameState.Mutex.Lock()
	defer gameState.Mutex.Unlock()

	if gameState.GameStarted {
		http.Error(w, "The game has already started, no more players can join", http.StatusForbidden)
		return
	}

	player := Player{
		ID:     len(gameState.Players) + 1,
		Pieces: []int{-1, -1, -1, -1},
	}

	gameState.Players = append(gameState.Players, player)

	fmt.Printf("Player %d with ID %d has joined\n", len(gameState.Players), player.ID)

	if len(gameState.Players) == numPlayers {
		gameState.GameStarted = true
		go runGame()
	}

	w.Write([]byte(fmt.Sprintf("Joined as Player %d\n", player.ID)))
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	gameState.Mutex.Lock()
	defer gameState.Mutex.Unlock()

	jsonState, err := json.Marshal(gameState)
	if err != nil {
		http.Error(w, "Error al serializar el estado del juego", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonState)
}

func runGame() {
	for {
		gameState.Mutex.Lock()

		if gameState.GameOver {
			gameState.Mutex.Unlock()
			close(gameOver)
			break
		}

		gameState.Round++
		fmt.Printf("\n--- Round %d ---\n", gameState.Round)

		for i := 0; i < numPlayers; i++ {
			if !gameState.GameOver {
				gameState.Mutex.Unlock()
				playTurn(i)
				gameState.Mutex.Lock()
			}
		}

		gameState.Mutex.Unlock()

		time.Sleep(50 * time.Millisecond)
	}
}

func playTurn(playerID int) {
	fmt.Printf("Player %d turn\n", playerID+1)
	x, y, operation := throwDice()
	doubles := isDoubles(x, y)

	if doubles {
		fmt.Printf("Got doubles! First two dice: [%d %d]\n", x, y)

		if canMove(gameState.Players[playerID].Pieces) {
			movePiece(x, y, operation, gameState.Players[playerID].Pieces, gameState.Obstacles)
			fmt.Printf("Dice rolled: [%d %d %d]\n", x, y, operation)
			fmt.Printf("Player %d Pieces: %+v\n", playerID+1, gameState.Players[playerID].Pieces)
		} else {
			freePiece(gameState.Players[playerID].Pieces)
			movePiece(0, 0, 0, gameState.Players[playerID].Pieces, gameState.Obstacles)
		}
	} else if canMove(gameState.Players[playerID].Pieces) {
		movePiece(x, y, operation, gameState.Players[playerID].Pieces, gameState.Obstacles)
		fmt.Printf("Dice rolled: [%d %d %d]\n", x, y, operation)
		fmt.Printf("Player %d Pieces: %+v\n", playerID+1, gameState.Players[playerID].Pieces)
	}

	if hasWon(gameState.Players[playerID].Pieces) {
		gameState.GameOver = true
		gameState.WinnerID = playerID + 1
		doubles = false
		fmt.Printf("\nPlayer %d won!\n", playerID+1)
	}
}

func throwDice() (int, int, int) {
	return rand.Intn(6) + 1, rand.Intn(6) + 1, rand.Intn(2)
}

func movePiece(x, y, operation int, playerPieces, obstacles []int) {
	pieceIndex := -1
	highestValue := -1

	for index, piece := range playerPieces {
		if finalTile > piece && piece > highestValue {
			highestValue = piece
			pieceIndex = index
		}
	}

	if pieceIndex == -1 {
		return
	}

	diceSum := x + y
	if operation == 1 {
		diceSum = x - y
	}

	newPosition := playerPieces[pieceIndex] + diceSum

	if playerPieces[pieceIndex] == -1 && newPosition < 0 {
		fmt.Println("Cannot move before exiting the starting position.")
		return
	}

	if !isValidMove(newPosition) {
		return
	}

	playerPieces[pieceIndex] = newPosition
}
func isValidMove(newPosition int) bool {
	return newPosition <= finalTile && newPosition >= 0
}
func freePiece(playerPieces []int) {
	for index, piece := range playerPieces {
		if piece == -1 {
			playerPieces[index] = 0
			return
		}
	}
}

func isDoubles(x, y int) bool {
	return x == y
}

func canMove(playerPieces []int) bool {
	for _, piece := range playerPieces {
		if piece != -1 && piece != finalTile {
			return true
		}
	}
	return false
}

func hasWon(playerPieces []int) bool {
	for _, piece := range playerPieces {
		if piece != finalTile {
			return false
		}
	}
	return true

}

func placeRandomObstacles(obstacleCount, finalTile int) []int {
	obstacles := make([]int, 0)

	for len(obstacles) < obstacleCount {
		position := rand.Intn(finalTile)
		if position != 0 && position != finalTile && !contains(obstacles, position) {
			obstacles = append(obstacles, position)
			row, col := position/boardSize, position%boardSize
			board[row][col] = "X"
		}
	}

	return obstacles
}

func contains(slice []int, element int) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
