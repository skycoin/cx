package maze

import (
	"fmt"
	"math/rand"
	"time"
)

// Maze cell configurations
// The paths of the maze is represented in the binary representation.
const (
	Up = 1 << iota
	Down
	Left
	Right
)

// Directions is the set of all the directions
var Directions = []int{Up, Down, Left, Right}

// The differences in the x-y coordinate
var dy = map[int]int{Up: -1, Down: 1, Left: 0, Right: 0}
var dx = map[int]int{Up: 0, Down: 0, Left: -1, Right: 1}

// Opposite directions
var Opposite = map[int]int{Up: Down, Down: Up, Left: Right, Right: Left}

type Point struct {
	X, Y int
}

// Advance the point forward by the argument direction
func (point *Point) Advance(direction int) *Point {
	return &Point{point.X + dx[direction], point.Y + dy[direction]}
}

type Maze struct {
	Width  int
	Height int
	// Each cell of maze has 4 bits (for whether there is an opening N, opening S, opening W, opening E) on the current cell
	// index=x+(y*width) each cell of maze has 4 bits
	Cells       []int
	CurrentMove int    // starts at zero, increments every move
	Goal        *Point // Goal position random
	Start       *Point // Start Position random
}

// NewMaze creates a new maze.
// The starting point is set randomly.
func NewMaze(height int, width int) *Maze {
	rand.Seed(time.Now().UnixNano())

	cells := make([]int, width*height)
	start := &Point{
		X: (rand.Int() % width),
		Y: (rand.Int() % height),
	}
	return &Maze{
		Width:       width,
		Height:      height,
		Cells:       cells,
		CurrentMove: 0,
		Start:       start,
	}
}

//  Generate generates a maze.
func (maze *Maze) Generate() {
	startTime := time.Now()
	point := maze.Start
	stack := []*Point{maze.Start}
	for len(stack) > 0 {
		for {
			point = maze.Next(point)
			if point == nil {
				break
			}
			stack = append(stack, point)
		}

		if len(stack) > 0 {
			stack = stack[:len(stack)-1] // Pop
			if len(stack) > 0 {
				point = stack[len(stack)-1]
			}
		}
	}
	totalTime := time.Since(startTime)
	fmt.Printf("Total time to Generate Maze: %v\n", totalTime)
}

// Next advances the maze path randomly and returns the new point
func (maze *Maze) Next(point *Point) *Point {
	neighbors := maze.Neighbors(point)
	if len(neighbors) == 0 {
		return nil
	}
	direction := neighbors[rand.Int()%len(neighbors)]
	maze.Cells[maze.getIndex(point.X, point.Y)] |= direction
	next := point.Advance(direction)
	maze.Cells[maze.getIndex(next.X, next.Y)] |= Opposite[direction]

	return next
}

// Contains judges whether the argument point is inside the maze or not
func (maze *Maze) Contains(point *Point) bool {
	return 0 <= point.X && point.X < maze.Width && 0 <= point.Y && point.Y < maze.Height
}

// Neighbors gathers the nearest undecided points
func (maze *Maze) Neighbors(point *Point) (neighbors []int) {
	for _, direction := range Directions {
		next := point.Advance(direction)
		if maze.Contains(next) && maze.Cells[maze.getIndex(next.X, next.Y)] == 0 {
			neighbors = append(neighbors, direction)
		}
	}
	return neighbors
}

// IsWallOpen checks if the bit in b located in pos is true.
func (maze *Maze) IsWallOpen(point *Point, pos int) bool {
	return ((maze.Cells[maze.getIndex(point.X, point.Y)]) & pos) == pos
}

func (maze *Maze) getIndex(x, y int) int {
	return x + (y * maze.Width)
}

func StartMaze(width, height int) *Maze {
	maze := NewMaze(height, width)
	maze.Generate()
	maze.ValidateMaze()
	maze.PrintMaze()
	maze.SetGoalPoint()

	return maze
}
