# Ludo modificado
El Ludo modificado representa una versión mejorada y ajustada del conocido juego de mesa Ludo. En esta variante, los participantes compiten para conducir a sus personajes a lo largo de un laberinto peligroso lleno de obstáculos y desafíos. Cada jugador se enfrenta al desafío de guiar a sus personajes a través del laberinto con el objetivo de llegar a la meta antes que sus oponentes.
## Planteamiento
### Inicialización del juego y mapa
* Se crea la cantidad de jugadores ingresada con 4 piezas cada uno. El estado inicial de cada pieza se establece en '-1'.
* Se genera el mapa con un 20% de posibilidades de que una de las casillas sea un obstáculo. Cada casilla libre se representa por un 1 y cada casilla con obstáculos con un 0.
### Canales
Creación de canales para coordinar los movimientos de los jugadores de forma asincrónica.
### Goroutines de los jugadores
* Simulan a cada jugador.
* Ejecuta una función que simula el turno de un jugaor en específico.
* Se emplea un sync.WaitGroup denominado "wg" con el propósito de asegurar que todas las goroutines de los jugadores concluyan antes de que el programa termine su ejecución.
### Turno de jugadores
* La función que va a simular los turnos de un jugador se va a llamar "playGame". Esta va a recibir el ID del jugador, sus piezas, el mapa de juego, un canal de movimiento y un sync.WaitGroup.
* Se recibe el ID del jugador desde el canal de movimiento para iniciar su turno.
* Si salen dados iguales y hay piezas en la posición de inicio (-1), el jugador mueve una pieza desde el inicio a la posición 0.
* Se imprime el estado actual de las piezas del jugador para mostrar su progreso.
* Para garantizar la finalización de los turnos se crea un sync.WaitGroup llamado 'playersWG'.
![image](https://github.com/weez97/concurrent_ludo/assets/63934328/dbc3487e-dc6b-4eda-a401-6f3bd30295f5)

## Video explicativo del código
https://youtu.be/8daiG1vD7dI 
