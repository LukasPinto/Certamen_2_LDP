package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"
)

func ctrl_c() {
	fmt.Print("\n\n[!] Saliendo.....\n")
}

// si is_occupied es un numero > 0, signinfica que esta ocupado por alguna serpiente
//  si is_occupied es 0, significa que esta libre
// si is_occupied es -1, siginifica que esta ocupado por una fruta

type Cell struct {
	cord_y      int // coordenada de la ceula en y
	cord_x      int // la coordenada de la celula en x
	is_occupied int // numero que contiene el valor de la snake o la fruta si es que esta ocupada
}

type Board struct {
	columns int
	rows    int
	cells   []Cell
	snakes  []Snake
	fruit   Fruit
}

type Fruit struct {
	cord_x int
	cord_y int
}

type Snake struct {
	head   Cell   // par correspondiente a la cabeza de la serpiente
	body   []Cell // pares correspondiente al cuerpo de la serpiente
	number int    // numero correspondiente a la serpiete ej --> serpiete 1,serpiente 2,etc
}

// funcion crea el tablero con valores vacios
func create_board(board *Board) {
	var aux_cell Cell
	for i := 0; i < board.rows; i++ {
		for j := 0; j < board.columns; j++ {
			aux_cell.cord_x = j
			aux_cell.cord_y = i
			board.cells = append(board.cells, aux_cell)

		}
	}

}

func create_fruit(board *Board) {

	var rand_y int
	var rand_x int
	find := true
	for find {
		rand_y = rand.Intn(board.rows)
		rand_x = rand.Intn(board.columns)

		for _, cell := range board.cells {

			if cell.cord_x == rand_x && cell.cord_y == rand_y && cell.is_occupied == 0 {
				board.fruit.cord_y = cell.cord_y
				board.fruit.cord_x = cell.cord_x
				find = false
			}
		}

	}
	for i := 0; i < board.rows*board.columns; i++ {
		if board.cells[i].cord_x == board.fruit.cord_x && board.cells[i].cord_y == board.fruit.cord_y {
			board.cells[i].is_occupied = -1
		}
	}

}

func create_snake(board *Board, snake_number int) {
	var rand_y int
	var rand_x int
	find := true
	for find {
		var snake Snake
		rand_y = rand.Intn(board.rows)
		rand_x = rand.Intn(board.columns)
		fmt.Printf(" Serpiente : %v,%v\n", rand_y, rand_x)
		for _, cell := range board.cells {

			if cell.cord_x == rand_x && cell.cord_y == rand_y && cell.is_occupied == 0 {

				snake.head.cord_x = cell.cord_x
				snake.head.cord_y = cell.cord_y
				snake.head.is_occupied = 1
				snake.number = snake_number
				board.snakes = append(board.snakes, snake)
				find = false
			}
		}

	}
	//esta funcion simplemente actualiza el valor de la serpiente en el tablero
	for i := 0; i < board.rows*board.columns; i++ {
		for j := 0; j < len(board.snakes); j++ {
			if board.snakes[j].head.cord_x == board.cells[i].cord_x && board.snakes[j].head.cord_y == board.cells[i].cord_y && board.cells[i].is_occupied == 0 {
				board.cells[i].is_occupied = board.snakes[j].number
			}
		}
	}

}

func swap_snake_add(board *Board, snake *Snake, old_snake *Snake) {
	var cell Cell
	cell.cord_x = old_snake.head.cord_x
	cell.cord_y = old_snake.head.cord_y
	cell.is_occupied = old_snake.number

	snake.body = append([]Cell{cell}, snake.body...)
	snake.head.cord_x = board.fruit.cord_x
	snake.head.cord_y = board.fruit.cord_y
	snake.head.is_occupied = 1

}

func swap_snake_move(board *Board, snake *Snake, old_snake *Snake) {
	var cell Cell
	cell.cord_x = old_snake.head.cord_x
	cell.cord_y = old_snake.head.cord_y
	cell.is_occupied = 1

	snake.body = append([]Cell{cell}, snake.body...)
	snake.body = snake.body[:len(snake.body)-1]

}

func snake_eat(board *Board, snake_index int, old_snake *Snake) {

	if board.snakes[snake_index].head.cord_x == board.fruit.cord_x && board.snakes[snake_index].head.cord_y == board.fruit.cord_y {
		swap_snake_add(board, &board.snakes[snake_index], old_snake)

		create_fruit(board)
	} else {
		swap_snake_move(board, &board.snakes[snake_index], old_snake)

	}

}

// verifico si la posicion a la que se va a mover la serpiente esta ocupada o no
func is_safe_move(board *Board, snake_cord_x int, snake_cord_y int) bool {

	for _, cell := range board.cells {
		if snake_cord_x == cell.cord_x && snake_cord_y == cell.cord_y && cell.is_occupied > 0 {
			return false
		}
	}
	return true
}

func move_snake(board *Board, snake_number int) {

	var old_snake Snake
	old_snake = board.snakes[snake_number]
	//  condiciones para hacer que la snake se mueve hacia la fruta
	// tambien se contempla si es que la fruta aparece atras, que la snake se mueva hacia un lado y luego continue
	if board.snakes[snake_number].head.cord_x < board.fruit.cord_x && is_safe_move(board, board.snakes[snake_number].head.cord_x+1, board.snakes[snake_number].head.cord_y) {
		board.snakes[snake_number].head.cord_x++
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_x > board.fruit.cord_x && is_safe_move(board, board.snakes[snake_number].head.cord_x-1, board.snakes[snake_number].head.cord_y) {
		board.snakes[snake_number].head.cord_x--
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_y < board.fruit.cord_y && is_safe_move(board, board.snakes[snake_number].head.cord_x, board.snakes[snake_number].head.cord_y+1) {
		board.snakes[snake_number].head.cord_y++
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_y > board.fruit.cord_y && is_safe_move(board, board.snakes[snake_number].head.cord_x, board.snakes[snake_number].head.cord_y-1) {
		board.snakes[snake_number].head.cord_y--
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_x < board.columns-1 && board.fruit.cord_x == board.snakes[snake_number].head.cord_x && is_safe_move(board, board.snakes[snake_number].head.cord_x+1, board.snakes[snake_number].head.cord_y) {
		board.snakes[snake_number].head.cord_x++
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_x > 0 && board.fruit.cord_x == board.snakes[snake_number].head.cord_x && is_safe_move(board, board.snakes[snake_number].head.cord_x-1, board.snakes[snake_number].head.cord_y) {
		board.snakes[snake_number].head.cord_x--
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_y < board.rows-1 && board.fruit.cord_y == board.snakes[snake_number].head.cord_y && is_safe_move(board, board.snakes[snake_number].head.cord_x, board.snakes[snake_number].head.cord_y+1) {
		board.snakes[snake_number].head.cord_y++
		snake_eat(board, snake_number, &old_snake)
	} else if board.snakes[snake_number].head.cord_y > 0 && board.fruit.cord_y == board.snakes[snake_number].head.cord_y && is_safe_move(board, board.snakes[snake_number].head.cord_x, board.snakes[snake_number].head.cord_y-1) {
		board.snakes[snake_number].head.cord_y--
		snake_eat(board, snake_number, &old_snake)
	} else { // por ultimo si no encuentra un camino, que se mueva para intentar encontrar un espacio o que se termine encerrando
		if board.snakes[snake_number].head.cord_y+1 <= board.rows-1 && is_safe_move(board, board.snakes[snake_number].head.cord_x, board.snakes[snake_number].head.cord_y+1) {
			board.snakes[snake_number].head.cord_y++
			snake_eat(board, snake_number, &old_snake)
		} else if board.snakes[snake_number].head.cord_y-1 >= 0 && is_safe_move(board, board.snakes[snake_number].head.cord_x, board.snakes[snake_number].head.cord_y-1) {
			board.snakes[snake_number].head.cord_y--
			snake_eat(board, snake_number, &old_snake)
		} else if board.snakes[snake_number].head.cord_x+1 <= board.columns-1 && is_safe_move(board, board.snakes[snake_number].head.cord_x+1, board.snakes[snake_number].head.cord_y) {
			board.snakes[snake_number].head.cord_x++
			snake_eat(board, snake_number, &old_snake)
		} else if board.snakes[snake_number].head.cord_x-1 >= 0 && is_safe_move(board, board.snakes[snake_number].head.cord_x-1, board.snakes[snake_number].head.cord_y) {
			board.snakes[snake_number].head.cord_x--
			snake_eat(board, snake_number, &old_snake)
		}

	}
	//actualiza valor de la cabeza de la serpiente
	if len(old_snake.body) > 0 {
		for i := 0; i < board.rows*board.columns; i++ {
			if old_snake.body[len(old_snake.body)-1].cord_x == board.cells[i].cord_x && old_snake.body[len(old_snake.body)-1].cord_y == board.cells[i].cord_y {
				board.cells[i].is_occupied = 0
			}
		}
	} else {
		for i := 0; i < board.rows*board.columns; i++ {
			if old_snake.head.cord_x == board.cells[i].cord_x && old_snake.head.cord_y == board.cells[i].cord_y {
				board.cells[i].is_occupied = 0
			}
		}
	}

	// actualizo el valor de la celda donde estaba anteriormente la serpiente
	for i := 0; i < board.rows*board.columns; i++ {
		for j := 0; j < len(board.snakes); j++ {
			if board.cells[i].cord_x == board.snakes[j].head.cord_x && board.cells[i].cord_y == board.snakes[j].head.cord_y {
				board.cells[i].is_occupied = board.snakes[j].number

			} else {
				for y := 0; y < len(board.snakes[j].body); y++ {
					if board.cells[i].cord_x == board.snakes[j].body[y].cord_x && board.cells[i].cord_y == board.snakes[j].body[y].cord_y {
						board.cells[i].is_occupied = board.snakes[j].number

					}
				}
			}
		}

	}

}

func show_board(board *Board) {
	colors := [6]string{
		"\033[32m", // green
		"\033[33m", //yellow
		"\033[35m", //purple
		"\033[36m", //Cyan
		"\033[34m", //Blue
		"\033[97m", // White

	}
	simbols := [6]string{
		"▤",
		"▥",
		"▦",
	}
	//
	cont := 0
	for i := 0; i < board.rows*board.columns; i++ {
		if cont < board.columns {
			if board.cells[i].is_occupied > 0 {
				aux_color := colors[board.cells[i].is_occupied%6]
				fmt.Printf(aux_color+"%v\033[0m\t", simbols[board.cells[i].is_occupied%3])
			} else if board.cells[i].is_occupied == -1 {
				fmt.Printf("\033[31m●\033[0m\t")
			} else {
				fmt.Printf("\033[37m□\033[0m\t")
			}
			//fmt.Printf("%v\t", board.cells[i].is_occupied)

			cont++
		} else {
			fmt.Println()
			cont = 0
			if board.cells[i].is_occupied > 0 {
				aux_color := colors[board.cells[i].is_occupied%6]
				fmt.Printf(aux_color+"%v\033[0m\t", simbols[board.cells[i].is_occupied%3])
			} else if board.cells[i].is_occupied == -1 {
				fmt.Printf("\033[31m●\033[0m\t")
			} else {
				fmt.Printf("\033[37m□\033[0m\t")
			}
			//fmt.Printf("%v\t", board.cells[i].is_occupied)
			cont++
		}

	}
	fmt.Println("\n\n")

}
func main() {
	rand.Seed(time.Now().UnixNano()) // semilla para el numero aleatorio
	wg := &sync.WaitGroup{}
	//parseo de argumentos
	columns := flag.Int("c", 10, "Número de columnas")
	rows := flag.Int("r", 10, "Número de filas")
	timer := flag.Int("timer", 500, "Numero de temporizado por movimiento")
	snakes := flag.Int("n", 2, "Cantidad de serpientes")
	gen := flag.Bool("gen", false, "Mostrar cada generacion del tablero")
	flag.Parse()
	fmt.Println(*columns, "x", *rows)
	// incializacion rapida de un array de m*n
	board := make(chan Board)
	var aux_board Board
	aux_board.columns = *columns
	aux_board.rows = *rows
	create_board(&aux_board)
	create_fruit(&aux_board)
	// se crean las snakes
	for i := 1; i <= *snakes; i++ {
		create_snake(&aux_board, i)
	}

	//go read_input(wg, board)
	fmt.Print("\033[H\033[2J")
	show_board(&aux_board)
	prev_board := aux_board // variable que guarda el valor anterior del tablero

	//esta funcion anonima se encarga de estar escuchando los cambios del tablero en base al canal de board
	go func() {

		cont := 0 // contador para determinar deadlock
		for {

			select {
			case aux_board = <-board: // recibe los valores que le manda la funcion update_board_by_snake

				if !(*gen) { // si el parametro gen esta activo, imprimo todas las generaciones del juego
					fmt.Printf("\033[0;0H") // retorno de  carro
				}
				time.Sleep(time.Duration(*timer) * time.Millisecond) // para manejar el tiempo de visualizacion de cada generacion

				show_board(&aux_board)
				if cont == 100 { // si el contador de deadlock llega a 100, significa que el tablero no ha cambiado en 100 generaciones
					//por lo tanto esta en deadlock
					fmt.Printf("DEADLOCK!!!\n\nGame Over!!!\n\n")
					syscall.Exit(0)

				}
				if reflect.DeepEqual(aux_board, prev_board) { //compara el tablero con su version anterior
					cont++

				} else if !reflect.DeepEqual(aux_board, prev_board) {
					cont = 0
				}
				prev_board = aux_board

			}

		}

		// printf para limpiar la pantalla
	}()

	for {
		//por muevo cada serpiente a traves de una goroutine y luego espero a que finalicen
		for i := 0; i < *snakes; i++ {
			wg.Add(1) // agrego la gorutina a la cola de espera
			go update_board_by_snake(wg, board, &aux_board, i)

		}
		wg.Wait() // aca espero a que finalicen las gorutinas

	}

}

func update_board_by_snake(wg *sync.WaitGroup, board chan Board, aux_board *Board, number_snake int) {
	//se le pasa el valor actual del tablero
	//se le dice que mueva la snake "number_snake"
	// y luego mande el valor del tablero por el canal de board
	// para que lo actualice en la gorutina que escucha los cambios
	move_snake(aux_board, number_snake)
	board <- *aux_board

	wg.Done()

}

// sirve para atrapar la senar SIGTERM (ctrl + c ) y finalizar el juego
func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ctrl_c()
		os.Exit(1)
	}()

}
