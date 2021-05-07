package evolve

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func makeDirectory(cfg *EvolveConfig) string {
	var dir string

	if cfg.PlotFitness || cfg.SaveAST {
		if cfg.MazeBenchmark {
			// Unixtime-Maze-2x2
			mazeSize := fmt.Sprintf("%vx%v", cfg.MazeWidth, cfg.MazeHeight)
			if cfg.RandomMazeSize {
				mazeSize = "random"
			}

			dir = fmt.Sprintf("./Results/%v-%v-%v/", time.Now().Unix(), "Maze", mazeSize)
		}

		if cfg.ConstantsBenchmark {
			// Unixtime-Constants
			dir = fmt.Sprintf("./Results/%v-Constants/", time.Now().Unix())
		}

		if cfg.EvensBenchmark {
			// Unixtime-Evens
			dir = fmt.Sprintf("./Results/%v-Evens/", time.Now().Unix())
		}

		if cfg.OddsBenchmark {
			// Unixtime-Odds
			dir = fmt.Sprintf("./Results/%v-Odds/", time.Now().Unix())
		}

		if cfg.PrimesBenchmark {
			// Unixtime-Primes
			dir = fmt.Sprintf("./Results/%v-Primes/", time.Now().Unix())
		}

		if cfg.CompositesBenchmark {
			// Unixtime-Composites
			dir = fmt.Sprintf("./Results/%v-Composites/", time.Now().Unix())
		}

		if cfg.RangeBenchmark {
			// Unixtime-Range
			dir = fmt.Sprintf("./Results/%v-Range/", time.Now().Unix())
		}

		if cfg.NetworkSimBenchmark {
			// Unixtime-NetworkSim
			dir = fmt.Sprintf("./Results/%v-NetworkSim/", time.Now().Unix())
		}

		// create directory
		_ = os.Mkdir(dir, 0700)

		if cfg.SaveAST {
			_ = os.Mkdir(dir+"AST/", 0700)
		}
	}
	return dir
}

func setEpochLength(cfg *EvolveConfig) {
	if cfg.EpochLength == 0 {
		cfg.EpochLength = 1
	}
}

func toByteArray(i int32) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(i))
	return arr
}
