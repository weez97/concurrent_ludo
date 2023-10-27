package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	// Player initialization
	players := make([][]int, 4)
	for i := 0; i < 4; i++ {
		players[i] = []int{-1, -1, -1, -1}
	}

	// Map generation
	var gameMap []int
	mSize := 57
	for i := 1; i < mSize; i++ {
		if rand.Float64() < 1.0/5.0 {
			gameMap = append(gameMap, 0)
		} else {
			gameMap = append(gameMap, 1)
		}
	}

	// Game variables initialization
	won := false
	c := 0

	// A channel to coordinate player moves
	moveChannels := make([]chan []int, 4)
	for i := 0; i < 4; i++ {
		moveChannels[i] = make(chan []int)
		defer close(moveChannels[i])
	}

	// WaitGroup to handle our concurrency
	var wg sync.WaitGroup

	// Launch goroutines for each player
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go playGame(i, players[i], gameMap, mSize, moveChannels[i], &wg)
	}

	// Create a WaitGroup to ensure all players finish their turns
	var playersWG sync.WaitGroup

	// Game loop
	for !won {
		c++
		fmt.Println("\nRonda:", c)

		// Instruct each player to take their turn asynchronously
		for i := 0; i < 4; i++ {
			playersWG.Add(1)
			go func(playerID int) {
				moveChannels[playerID] <- []int{playerID} // Send the player ID to initiate the turn
				playersWG.Done()
			}(i)
		}

		// Wait for all players to finish their turns
		playersWG.Wait()

		won = checkWin(players, mSize)
	}

	fmt.Println("\nFinalizado -", players)
}

func playGame(playerID int, player []int, gameMap []int, mSize int, moveChannel chan []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// Receive the player's ID to initiate the turn
		turnData := <-moveChannel
		if turnData[0] != playerID {
			// Ignore messages not meant for this player
			continue
		}

		roll := []int{rand.Intn(6) + 1, []int{-1, 1}[rand.Intn(2)], rand.Intn(6) + 1}

		if roll[0] == roll[2] && contains(player, -1) {
			for idx, piece := range player {
				if piece == -1 {
					player[idx] = 0
					break
				}
			}
		} else {
			move := roll[0] + roll[1]*roll[2]
			for idx, piece := range player {
				if piece < mSize && piece >= -1 { // Ensure piece position is never less than -1
					if piece+move >= mSize-1 {
						player[idx] = mSize
						break
					}
					if piece+move >= 0 && gameMap[piece+move] == 1 {
						player[idx] = piece + move
						break
					}
				}
			}
		}
		fmt.Printf("Player %d: %v\n", playerID+1, player)
	}
}

func contains(arr []int, value int) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}

func checkWin(players [][]int, mSize int) bool {
	for _, player := range players {
		if equal(player, []int{mSize, mSize, mSize, mSize}) {
			return true
		}
	}
	return false
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
