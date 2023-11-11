# Ludo modificado
El Ludo modificado representa una versión mejorada y ajustada del conocido juego de mesa Ludo. En esta variante, los participantes compiten para conducir a sus personajes a lo largo de un laberinto peligroso lleno de obstáculos y desafíos. Cada jugador se enfrenta al desafío de guiar a sus personajes a través del laberinto con el objetivo de llegar a la meta antes que sus oponentes. Los obstáculos añaden un elemento de estrategia al juego. La comunicación entre el servidor y los jugadores se realiza a través de solicitudes HTTP.
## Planteamiento
### Inicialización del juego y mapa
* El juego es un juego de tablero simple con un tablero de tamaño "boardSize" (10x10 en este caso).
* El objetivo del juego es que los jugadores muevan sus piezas desde la posición inicial "0" hasta la casilla final "finalTile".
* El servidor HTTP se inicia en "http://localhost:8080".
* El tablero, los obstáculos y el estado del juego se inicializan en la función main.
* El servidor maneja dos rutas: "/join" para que los jugadores se unan y "/state" para obtener el estado actual del juego.
* Los jugadores se unen al juego haciendo una solicitud a "http://localhost:8080/join".
* Cada jugador tiene su propia goroutine que realiza solicitudes periódicas a "http://localhost:8080/state" para obtener el estado del juego.
* Se crea la cantidad de jugadores ingresada con N piezas cada uno. El estado inicial de cada pieza se establece en "-1".
* Se genera el mapa con un 20% de posibilidades de que una de las casillas sea un obstáculo. Cada casilla libre se representa por un "1" y cada casilla con obstáculos con un "0".
### Canales
* Creación de canales para coordinar los movimientos de los jugadores de forma asincrónica.
* Se utiliza un canal "gameOver" para notificar al servidor cuando el juego ha terminado.
### Goroutines de los jugadores
* Cada jugador tiene su propia goroutine que se ejecuta en la función "processGame".
* La goroutine realiza solicitudes periódicas para obtener el estado del juego y espera a que el juego comience si aún no ha empezado.
* Se emplea un "sync.WaitGroup" denominado "wg" con el propósito de asegurar que todas las goroutines de los jugadores concluyan antes de que el programa termine su ejecución.
### Turno de jugadores
* Los jugadores toman turnos en la función "runGame".
* Cada jugador realiza su turno llamando a la función "playTurn" en la goroutine correspondiente.
* Se recibe el ID del jugador desde el canal de movimiento para iniciar su turno.
* Si salen dados iguales y hay piezas en la posición de inicio "-1", el jugador mueve una pieza desde el inicio a la posición "0".
* Se imprime el estado actual de las piezas del jugador para mostrar su progreso.
* Para garantizar la finalización de los turnos se crea un "sync.WaitGroup" llamado "playersWG".
<img width="485" alt="image" src="https://github.com/weez97/concurrent_ludo/assets/38121350/ecb304ff-89cc-4d62-835d-260bda72b10d">
## Video explicativo del código
https://youtu.be/8daiG1vD7dI 
