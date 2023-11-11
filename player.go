package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
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

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(playerID int) {
			defer wg.Done()

			joinResponse, err := http.Get("http://localhost:8080/join")
			if err != nil {
				fmt.Println("Error al unirse al juego:", err)
				return
			}

			defer joinResponse.Body.Close()
			body, err := ioutil.ReadAll(joinResponse.Body)
			if err != nil {
				fmt.Println("Error al leer la respuesta del servidor:", err)
				return
			}

			fmt.Printf("%s\n", string(body))

			for {
				stateResponse, err := http.Get("http://localhost:8080/state")
				if err != nil {
					fmt.Println("Error al obtener el estado del juego:", err)
					return
				}

				defer stateResponse.Body.Close()
				stateBody, err := ioutil.ReadAll(stateResponse.Body)
				if err != nil {
					fmt.Println("Error al leer la respuesta del servidor:", err)
					return
				}

				var gameState GameState
				err = json.Unmarshal(stateBody, &gameState)
				if err != nil {
					fmt.Println("Error al decodificar el estado del juego:", err)
					return
				}

				if gameState.GameStarted {
					break
				}

				fmt.Printf("Esperando a que comience el juego...\n")
				time.Sleep(2 * time.Second)
			}

			processGame(playerID)
		}(i)
	}

	wg.Wait()
}

func processGame(playerID int) {
	for {
		stateResponse, err := http.Get("http://localhost:8080/state")
		if err != nil {
			fmt.Println("Error al obtener el estado del juego:", err)
			return
		}

		defer stateResponse.Body.Close()
		stateBody, err := ioutil.ReadAll(stateResponse.Body)
		if err != nil {
			fmt.Println("Error al leer la respuesta del servidor:", err)
			return
		}

		var gameState GameState
		err = json.Unmarshal(stateBody, &gameState)
		if err != nil {
			fmt.Println("Error al decodificar el estado del juego:", err)
			return
		}

		if gameState.GameOver {
			fmt.Printf("Game finished\n")
			if gameState.WinnerID != 0 {
				fmt.Printf("Player %d won!\n", gameState.WinnerID)
			}
			return
		}

		if !gameState.GameStarted {
			fmt.Printf("Esperando a que comience el juego...\n")
			time.Sleep(2 * time.Second)
			continue
		}

	}
}
