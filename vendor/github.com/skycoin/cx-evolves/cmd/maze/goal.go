package maze

import (
	"fmt"
	"time"
)

// Distance between cells.
// Used for Dijkstra's algo for finding farthest point from goal point.
const (
	distanceBetweenPoints = 1
)

// SetGoalPoint sets maze goal point by using djikstra's algorithm to find farthest
// point from the starting point.
func (maze *Maze) SetGoalPoint() {
	startTime := time.Now()
	var GoalPoint Point
	var maxDistance int
	var pathValues []int
	length := len(maze.Cells)
	// Create an array for each cell
	for i := 0; i < length; i++ {
		pathValues = append(pathValues, 0)
	}

	maze.UpdateValue(&pathValues, *maze.Start, *maze.Start, &GoalPoint, &maxDistance)
	maze.Goal = &GoalPoint
	maze.CurrentMove = maxDistance
	totalTime := time.Since(startTime)

	fmt.Printf("Total time to get furthest point: %v\n", totalTime)
	fmt.Printf("Start Point: %v\n", maze.Start)
	fmt.Printf("Goal Point: %v\n", GoalPoint)
	fmt.Printf("Minimum moves: %v\n", maxDistance)
}

//  UpdateValue is a recursion implementation for dijkstra's algorithm for the maze.
func (maze *Maze) UpdateValue(pathValues *[]int, point, prev Point, goalPoint *Point, maxDistance *int) {
	var newVal int
	// check up
	up := point.Advance(Up)
	if maze.Contains(up) && maze.IsWallOpen(&point, Up) && *up != prev {
		newVal = (*pathValues)[maze.getIndex(point.X, point.Y)] + distanceBetweenPoints
		(*pathValues)[maze.getIndex(up.X, up.Y)] = newVal
		if newVal > *maxDistance {
			*maxDistance = newVal
			*goalPoint = *up
		}
		maze.UpdateValue(pathValues, *up, point, goalPoint, maxDistance)
	}

	// check down
	down := point.Advance(Down)
	if maze.Contains(down) && maze.IsWallOpen(&point, Down) && *down != prev {
		newVal = (*pathValues)[maze.getIndex(point.X, point.Y)] + distanceBetweenPoints
		(*pathValues)[maze.getIndex(down.X, down.Y)] = newVal
		if newVal > *maxDistance {
			*maxDistance = newVal
			*goalPoint = *down
		}
		maze.UpdateValue(pathValues, *down, point, goalPoint, maxDistance)
	}

	// check right
	right := point.Advance(Right)
	if maze.Contains(right) && maze.IsWallOpen(&point, Right) && *right != prev {
		newVal = (*pathValues)[maze.getIndex(point.X, point.Y)] + distanceBetweenPoints
		(*pathValues)[maze.getIndex(right.X, right.Y)] = newVal
		if newVal > *maxDistance {
			*maxDistance = newVal
			*goalPoint = *right
		}
		maze.UpdateValue(pathValues, *right, point, goalPoint, maxDistance)
	}

	// check left
	left := point.Advance(Left)
	if maze.Contains(left) && maze.IsWallOpen(&point, Left) && *left != prev {
		newVal = (*pathValues)[maze.getIndex(point.X, point.Y)] + distanceBetweenPoints
		(*pathValues)[maze.getIndex(left.X, left.Y)] = newVal
		if newVal > *maxDistance {
			*maxDistance = newVal
			*goalPoint = *left
		}
		maze.UpdateValue(pathValues, *left, point, goalPoint, maxDistance)
	}
}
