package maze

import (
	"fmt"
	"time"

	"gonum.org/v1/plot/plotter"
)

type Game struct {
	maze          *Maze
	PlotHistogram bool
}

type GameMove struct {
	MoveCount                int
	ErrorCode                int    // Error code success - 0
	ErrorMsg                 string // Error message success - empty
	GoalPosition             *Point
	AgentPosition            *Point // Disabled value is nil/null
	NumberOfSquaresLeftNorth int    // Disabled value is -1
	NumberOfSquaresLeftSouth int    // Disabled value is -1
	NumberOfSquaresLeftEast  int    // Disabled value is -1
	NumberOfSquaresLeftWest  int    // Disabled value is -1
	MazeData                 *Maze  // Disabled value is nil/null
	PrevInput                int
}

func (game *Game) Init(w, h int) {
	game.maze = StartMaze(w, h)
}

func (game *Game) MazeGame(numberOfRuns int, player func(gameMove *GameMove) (AgentInput, error)) (int, error) {
	var moves int
	var reachTime time.Duration
	var stopTime time.Duration
	var histoValues plotter.Values
	var gameMove *GameMove
	var movedToAWall bool

	maxMoves := 100 * game.maze.Width * game.maze.Height
	if player == nil {
		player = defaultRandomPlayer
	}

	for run := 0; run < numberOfRuns; run++ {
		// Reset game data
		gameMove = &GameMove{
			MoveCount:                0,
			ErrorCode:                0,
			ErrorMsg:                 "",
			GoalPosition:             nil,
			AgentPosition:            nil,
			NumberOfSquaresLeftNorth: -1,
			NumberOfSquaresLeftSouth: -1,
			NumberOfSquaresLeftEast:  -1,
			NumberOfSquaresLeftWest:  -1,
			MazeData:                 nil,
		}

		reachedGoal := false
		moves = 0
		goalPos := game.maze.Goal
		currPos := game.maze.Start
		gameMove.AgentPosition = currPos
		gameMove.MazeData = game.maze
		startTime := time.Now()

		for !reachedGoal {
			agentInput, err := InputCallback(player, gameMove)
			if err != nil {
				return 0, err
			}
			movedToAWall = false
			gameMove.ErrorCode = 0
			gameMove.ErrorMsg = ""

			switch agentInput.Move {
			case Up:
				next := currPos.Advance(Up)
				if game.maze.Contains(next) && game.maze.IsWallOpen(currPos, Up) {
					currPos = next
				} else {
					movedToAWall = true
				}
			case Down:
				next := currPos.Advance(Down)
				if game.maze.Contains(next) && game.maze.IsWallOpen(currPos, Down) {
					currPos = next
				} else {
					movedToAWall = true
				}
			case Left:
				next := currPos.Advance(Left)
				if game.maze.Contains(next) && game.maze.IsWallOpen(currPos, Left) {
					currPos = next
				} else {
					movedToAWall = true
				}
			case Right:
				next := currPos.Advance(Right)
				if game.maze.Contains(next) && game.maze.IsWallOpen(currPos, Right) {
					currPos = next
				} else {
					movedToAWall = true
				}
			}
			moves = moves + 1
			gameMove.MoveCount++
			if agentInput.AgentPositionEnabled {
				gameMove.AgentPosition = currPos
			}

			if agentInput.PassMazeData {
				gameMove.MazeData = game.maze
			}

			if movedToAWall {
				gameMove.ErrorCode = 1
				gameMove.ErrorMsg = "You moved to a wall."
			}

			if agentInput.WallDistanceInputEnabled {
				gameMove.NumberOfSquaresLeftNorth = game.CountSquaresBeforeWall(*currPos, Up)
				gameMove.NumberOfSquaresLeftSouth = game.CountSquaresBeforeWall(*currPos, Down)
				gameMove.NumberOfSquaresLeftWest = game.CountSquaresBeforeWall(*currPos, Left)
				gameMove.NumberOfSquaresLeftEast = game.CountSquaresBeforeWall(*currPos, Right)
			}

			if *currPos == *goalPos {
				reachedGoal = true
				reachTime = time.Since(startTime)
			}

			if moves > maxMoves {
				stopTime = time.Since(startTime)
				break
			}

			gameMove.PrevInput = agentInput.Move
		}

		if reachedGoal {
			fmt.Println("-------------------------------------------")
			fmt.Printf("Run Number: %v\n", run+1)
			fmt.Printf("You have reached the main goal! congrats!\n")
			fmt.Printf("Time it took to reach the goal: %v\n", reachTime)
			fmt.Printf("Moves it took you to reach goal: %v\n", moves)
		} else {
			fmt.Println("-------------------------------------------")
			fmt.Printf("Run Number: %v\n", run+1)
			fmt.Printf("Sorry, you have not reached the goal\n")
			fmt.Printf("Moves reached more than %v which is 100 times greater than maze size\n", maxMoves)
			fmt.Printf("Moves you took: %v\n", moves)
			fmt.Printf("Time you took: %v\n", stopTime)
		}
		if game.PlotHistogram {
			histoValues = append(histoValues, float64(moves))
		}

	}
	if game.PlotHistogram {
		histogramPlot(histoValues, "Number Of Moves The Random Player Took", histoSaveDirectory+"MovesHistogram-"+time.Now().Format(time.RFC3339)+".png")
	}

	return moves, nil
}

func (game *Game) CountSquaresBeforeWall(currPos Point, direction int) int {
	var count int
	var wallFound bool
	for !wallFound {
		next := currPos.Advance(direction)
		if game.maze.Contains(next) && game.maze.IsWallOpen(&currPos, direction) {
			currPos = *next
			count++
		} else {
			wallFound = true
		}
	}

	return count
}

func checkPrev(checkPrev bool, prevInput int, pos int) bool {
	return (checkPrev && prevInput != pos) || !checkPrev
}
