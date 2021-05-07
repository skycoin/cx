package evolve

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/skycoin/cx-evolves/cmd/maze"
	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	workerclient "github.com/skycoin/cx-evolves/cxexecutes/worker/client"
	cxast "github.com/skycoin/cx/cx/ast"
)

// Evaluate Program as the Maze Player
func mazeMovesEvaluation(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, cfg *EvolveConfig) (float64, error) {
	player := func(gameMove *maze.GameMove) (maze.AgentInput, error) {
		agentInput := maze.AgentInput{
			PassMazeData:             true,
			WallDistanceInputEnabled: true,
			AgentPositionEnabled:     true,
		}
		options := []int{maze.Up, maze.Down, maze.Left, maze.Right}

		move, err := perByteEvaluationMaze(ind, solPrototype, MazeEncodeParam(gameMove), cfg.WorkerPortNum)
		if err != nil {
			return maze.AgentInput{}, err
		}
		input := options[int(move)%len(options)]
		agentInput.Move = input

		return agentInput, nil
	}

	moves, err := cfg.MazeGame.MazeGame(1, player)
	if err != nil {
		return 0, err
	}
	return float64(moves), nil
}

// perByteEvaluation for evolve with maze, 13 i32 input, 1 i32 output
func perByteEvaluationMaze(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, inputs [][]byte, portNumber int) (int, error) {
	var move int
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)
	// `inpsOff` helps us keep track of what byte in `inps` we can write to.
	inpsOff := 0

	for c := 0; c < len(inputs); c++ {
		// The size of the input.
		inpSize := solPrototype.Inputs[c].TotalSize
		// The bytes representing the input.
		inp := inputs[c]

		// Copying the input `b`ytes.
		for b := 0; b < len(inp); b++ {
			inps[inpsOff+b] = inp[b]
		}

		// Updating offset.
		inpsOff += inpSize
	}

	var result worker.Result
	workerAddr := fmt.Sprintf(":%v", portNumber)
	workerclient.CallWorker(
		workerclient.CallWorkerConfig{
			Program:   ind,
			Input:     inps,
			OutputArg: solPrototype.Outputs[0],
		},
		workerAddr,
		&result,
	)
	move = int(binary.BigEndian.Uint32(result.Output))

	cxast.PROGRAM = tmp
	return move, nil
}

func MazeEncodeParam(param *maze.GameMove) [][]byte {
	paramCount := 13
	inputs := make([][]byte, paramCount)

	inputs[0] = []byte(fmt.Sprint(int32(param.MoveCount)))
	inputs[1] = []byte(fmt.Sprint(int32(param.ErrorCode)))
	inputs[2] = []byte(param.ErrorMsg)
	inputs[3] = []byte(fmt.Sprint(int32(param.MazeData.Goal.X)))
	inputs[4] = []byte(fmt.Sprint(int32(param.MazeData.Goal.Y)))
	inputs[5] = []byte(fmt.Sprint(int32(param.AgentPosition.X)))
	inputs[6] = []byte(fmt.Sprint(int32(param.AgentPosition.Y)))
	inputs[7] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftNorth)))
	inputs[8] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftSouth)))
	inputs[9] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftEast)))
	inputs[10] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftWest)))
	inputs[11] = []byte(fmt.Sprint(int32(param.MazeData.Height)))
	inputs[12] = []byte(fmt.Sprint(int32(param.MazeData.Width)))

	return inputs

	// Other implementation
	// WIll still se if theres difference
	// if param.MoveCount != 0 {
	// 	inputs[0] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[0], uint32(param.MoveCount))
	// }

	// if param.ErrorCode != 0 {
	// 	inputs[1] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[1], uint32(param.ErrorCode))
	// }

	// if param.ErrorMsg != "" {
	// 	inputs[2] = make([]byte, 8)
	// 	errMsg, _ := strconv.ParseUint(param.ErrorMsg, 10, 32)
	// 	binary.BigEndian.PutUint32(inputs[2], uint32(errMsg))
	// }

	// if param.MazeData.Goal.X != 0 {
	// 	inputs[3] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[3], uint32(param.MazeData.Goal.X))

	// }
	// if param.MazeData.Goal.Y != 0 {
	// 	inputs[4] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[4], uint32(param.MazeData.Goal.Y))
	// }

	// if param.AgentPosition.X != 0 {
	// 	inputs[5] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[5], uint32(param.AgentPosition.X))

	// }
	// if param.AgentPosition.Y != 0 {
	// 	inputs[6] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[6], uint32(param.AgentPosition.Y))
	// }

	// if param.NumberOfSquaresLeftNorth != 0 {
	// 	inputs[7] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[7], uint32(param.NumberOfSquaresLeftNorth))
	// }
	// if param.NumberOfSquaresLeftSouth != 0 {
	// 	inputs[8] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[8], uint32(param.NumberOfSquaresLeftSouth))
	// }

	// if param.NumberOfSquaresLeftEast != 0 {
	// 	inputs[9] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[9], uint32(param.NumberOfSquaresLeftEast))
	// }

	// if param.NumberOfSquaresLeftWest != 0 {
	// 	inputs[10] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[10], uint32(param.NumberOfSquaresLeftWest))
	// }

	// if param.MazeData.Height != 0 {
	// 	inputs[11] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[11], uint32(param.MazeData.Height))
	// }

	// if param.MazeData.Width != 0 {
	// 	inputs[12] = make([]byte, 8)
	// 	binary.BigEndian.PutUint32(inputs[12], uint32(param.MazeData.Width))
	// }

}

func generateNewMaze(generationCount int, cfg *EvolveConfig, game *maze.Game) {
	if generationCount%cfg.EpochLength == 0 || generationCount == 0 {
		if cfg.RandomMazeSize {
			setRandomMazeSize(cfg)
		}

		game.Init(cfg.MazeWidth, cfg.MazeHeight)
	}
}

func setRandomMazeSize(cfg *EvolveConfig) {
	rand.Seed(time.Now().Unix())
	randOptions := []int{2, 3, 4, 5, 6, 7, 8}
	size := randOptions[rand.Int()%len(randOptions)]
	cfg.MazeWidth = size
	cfg.MazeHeight = size
}
