package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
	roll := []int{0, 0, 0}
	won := false
	c := 0

	// Game loop
	for !won {
		c++
		fmt.Println("\nRonda:", c)

		for _, player := range players {
			roll[0] = rand.Intn(6) + 1
			roll[1] = []int{-1, 1}[rand.Intn(2)]
			roll[2] = rand.Intn(6) + 1

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
					if piece < mSize && piece != -1 {
						if piece+move >= mSize-1 {
							player[idx] = mSize
							break
						}
						if gameMap[piece+move] == 1 {
							player[idx] = piece + move
							break
						}
					}
				}
			}
			fmt.Println(player)
		}

		won = checkWin(players, mSize)
	}

	fmt.Println("\nFinalizado -", players)
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
