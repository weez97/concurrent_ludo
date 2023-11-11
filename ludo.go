package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numPlayers  = 4
	finalTile   = 56
	obstaclePct = 10
)

var (
	obstacles   []int
	currentTurn int
	pieces      [][]int
	gameOver    bool
	round       int
	mutex       sync.Mutex
	wg          sync.WaitGroup
)

func main() {
	rand.Seed(time.Now().UnixNano())

	obstacles = placeRandomObstacles(obstaclePct, finalTile)

	initializePlayers()

	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go playTurn(i)
	}

	wg.Wait()
}

func initializePlayers() {
	pieces = make([][]int, numPlayers)
	for i := range pieces {
		pieces[i] = make([]int, 4)
		for j := range pieces[i] {
			pieces[i][j] = -1 // Inicializar todas las piezas en -1
		}
	}
}

func hasWon(playerPieces []int) bool {
	for _, piece := range playerPieces {
		if piece != finalTile {
			return false
		}
	}
	return true
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

func passTurn() {
	currentTurn = (currentTurn + 1) % numPlayers
}

func throwDice() (int, int, int) {
	return rand.Intn(6) + 1, rand.Intn(6) + 1, rand.Intn(2)
}

func placeRandomObstacles(obstacleCount, finalTile int) []int {
	obstacles := make([]int, 0)

	for len(obstacles) < obstacleCount {
		position := rand.Intn(finalTile)
		if position != 0 && position != finalTile && !contains(obstacles, position) {
			obstacles = append(obstacles, position)
		}
	}

	return obstacles
}

func playTurn(playerID int) {
	defer wg.Done()

	for {
		mutex.Lock()
		if gameOver {
			mutex.Unlock()
			break
		}

		if currentTurn == playerID {
			fmt.Printf("Player %d turn\n", playerID+1)
			x, y, operation := throwDice()
			doubles := isDoubles(x, y)

			if doubles {
				fmt.Printf("Got doubles! First two dice: [%d %d]\n", x, y)

				if canMove(pieces[playerID]) {
					movePiece(x, y, operation, pieces[playerID], obstacles)
					fmt.Printf("Dice rolled: [%d %d %d]\n", x, y, operation)
					fmt.Printf("Player Pieces: %+v\n", pieces[playerID])
				} else {
					freePiece(pieces[playerID])
					movePiece(0, 0, 0, pieces[playerID], obstacles)
				}
			} else if canMove(pieces[playerID]) {
				movePiece(x, y, operation, pieces[playerID], obstacles)
				fmt.Printf("Dice rolled: [%d %d %d]\n", x, y, operation)
				fmt.Printf("Player Pieces: %+v\n", pieces[playerID])
			}

			if hasWon(pieces[playerID]) {
				gameOver = true
				doubles = false
				fmt.Printf("\nPlayer %d won!\n", playerID+1)
			}

			for doubles && !gameOver {
				x, y, operation = throwDice()
				movePiece(x, y, operation, pieces[playerID], obstacles)
				fmt.Printf("Dice rolled: [%d %d %d]\n", x, y, operation)
				fmt.Printf("Player Pieces: %+v\n", pieces[playerID])
				doubles = isDoubles(x, y)
			}

			passTurn()

			if currentTurn == 0 && !gameOver {
				round++
				fmt.Printf("\n--- Round %d ---\n", round)
				fmt.Printf("Player Positions: %v\n", pieces)
			}
		}

		mutex.Unlock()
		//time.Sleep(100 * time.Millisecond) // Reducir el tiempo de espera entre iteraciones
	}
}

func contains(slice []int, element int) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
