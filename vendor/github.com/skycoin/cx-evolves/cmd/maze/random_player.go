package maze

import (
	"fmt"
	"math/rand"
	"time"
)

type AgentInput struct {
	Move                     int
	BlindRunner              bool
	PassMazeData             bool
	WallDistanceInputEnabled bool
	AgentPositionEnabled     bool
}

type Player func(gameMove *GameMove) (AgentInput, error)

var defaultRandomPlayer Player = func(gameMove *GameMove) (AgentInput, error) {
	var defaultRandomPlayerInside func(gameMove *GameMove) (AgentInput, error)

	defaultRandomPlayerInside = func(gameMove *GameMove) (AgentInput, error) {
		var input int
		var options []int
		var agentInput AgentInput
		var CheckPrev bool
		agentInput.PassMazeData = true
		agentInput.AgentPositionEnabled = true
		agentInput.BlindRunner = false
		agentInput.WallDistanceInputEnabled = true
		CheckPrev = true

		if gameMove.ErrorCode != 0 {
			fmt.Printf("%v\n", gameMove.ErrorMsg)
		}

		// if agentInput.WallDistanceInputEnabled {
		// 	// fmt.Printf("Wall Distance To North: %v\n", gameMove.NumberOfSquaresLeftNorth)
		// 	// fmt.Printf("Wall Distance To South: %v\n", gameMove.NumberOfSquaresLeftSouth)
		// 	// fmt.Printf("Wall Distance To East: %v\n", gameMove.NumberOfSquaresLeftEast)
		// 	// fmt.Printf("Wall Distance To West: %v\n", gameMove.NumberOfSquaresLeftWest)
		// }
		if agentInput.BlindRunner {
			options = []int{Up, Down, Left, Right}
		} else {
			// If UP is open and ((check previous and previous move is not from there) or (if dont check previous))
			if gameMove.MazeData.IsWallOpen(gameMove.AgentPosition, Up) && checkPrev(CheckPrev, gameMove.PrevInput, Down) {
				options = append(options, Up)
			}

			// If DOWN is open and ((check previous and previous move is not from there) or (if dont check previous))
			if gameMove.MazeData.IsWallOpen(gameMove.AgentPosition, Down) && checkPrev(CheckPrev, gameMove.PrevInput, Up) {
				options = append(options, Down)
			}

			// If LEFT is open and ((check previous and previous move is not from there) or (if dont check previous))
			if gameMove.MazeData.IsWallOpen(gameMove.AgentPosition, Left) && checkPrev(CheckPrev, gameMove.PrevInput, Right) {
				options = append(options, Left)
			}

			// If Right is open and ((check previous and previous move is not from there) or (if dont check previous))
			if gameMove.MazeData.IsWallOpen(gameMove.AgentPosition, Right) && checkPrev(CheckPrev, gameMove.PrevInput, Left) {
				options = append(options, Right)
			}
		}

		rand.Seed(time.Now().UTC().UnixNano())
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })
		for input == 0 {
			if len(options) != 0 {
				input = options[rand.Int()%len(options)]
			}
			if len(options) == 0 {
				gameMove.PrevInput = 0
				aInput, err := InputCallback(defaultRandomPlayerInside, gameMove)
				if err != nil {
					return agentInput, err
				}
				input = aInput.Move
			}
		}
		agentInput.Move = input
		return agentInput, nil
	}

	return defaultRandomPlayerInside(gameMove)
}

func InputCallback(player func(gameMove *GameMove) (AgentInput, error), move *GameMove) (AgentInput, error) {
	return player(move)
}
