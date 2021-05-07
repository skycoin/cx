package evolve

import (
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"sync"

	"github.com/skycoin/cx-evolves/cmd/maze"
	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

type EvolveConfig struct {
	MazeBenchmark       bool
	ConstantsBenchmark  bool
	EvensBenchmark      bool
	OddsBenchmark       bool
	PrimesBenchmark     bool
	CompositesBenchmark bool
	RangeBenchmark      bool
	NetworkSimBenchmark bool

	MazeGame       *maze.Game
	MazeHeight     int
	MazeWidth      int
	RandomMazeSize bool

	NumberOfRounds int

	UpperRange int
	LowerRange int

	EpochLength int
	PlotFitness bool
	SaveAST     bool
	UseAntiLog2 bool

	WorkerPortNum    int
	WorkersAvailable int
}

// Used for concurrent output evaluation
var wg = sync.WaitGroup{}

func (pop *Population) Evolve(cfg EvolveConfig) {
	var histoValues []float64
	var averageValues []float64
	var mostFit []float64
	var availPorts []int
	var game maze.Game
	var saveDirectory string

	output := make([]float64, pop.PopulationSize)
	numIter := pop.Iterations
	solProt := pop.FunctionToEvolve
	fnToEvolveName := solProt.Name
	sPrgrm := cxast.SerializeCXProgramV2(pop.Individuals[0], true)

	setEpochLength(&cfg)
	saveDirectory = makeDirectory(&cfg)

	availWorkers := worker.GetAvailableWorkers(cfg.WorkersAvailable)
	availPorts = append(availPorts, availWorkers...)

	// Evolution process.
	for c := 0; c < int(numIter); c++ {
		// Maze Creation if Maze Benchmark
		if cfg.MazeBenchmark {
			generateNewMaze(c, &cfg, &game)
			cfg.MazeGame = &game
		}

		// Selection process.
		pop1Idx, pop2Idx := tournamentSelection(output, 0.5, true)
		dead1Idx, dead2Idx := tournamentSelection(output, 0.5, false)

		pop1MainPkg, err := pop.Individuals[pop1Idx].GetPackage(cxconstants.MAIN_PKG)
		if err != nil {
			panic(err)
		}

		parent1, err := pop1MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		pop2MainPkg, err := pop.Individuals[pop2Idx].GetPackage(cxconstants.MAIN_PKG)
		if err != nil {
			panic(err)
		}

		parent2, err := pop2MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		// Crossover process.
		crossoverFn := pop.getCrossoverFn()
		child1, child2 := crossoverFn(parent1, parent2)

		// Mutation process.
		_ = sPrgrm
		_ = dead1Idx
		_ = dead2Idx
		_ = child1
		_ = child2
		randomMutation(pop, sPrgrm)

		// Replacing individuals in population.
		replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
		replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)

		runtime.GOMAXPROCS(4)
		// Evaluation process.
		for i := range pop.Individuals {
			// Wait until ports are available.
			for len(availPorts) == 0 {
			}

			// Get port number from the first index and
			// remove it from list of available ports.
			portNum := availPorts[0]
			availPorts = RemoveIndex(availPorts, 0)
			cfg.WorkerPortNum = portNum

			wg.Add(1)
			go func(j int, availPorts *[]int, portNum int) {
				defer wg.Done()

				pop.Individuals[j].PrintProgram()
				output[j], err = RunBenchmark(pop.Individuals[j], solProt, &cfg)
				if err != nil {
					output[j] = float64(math.MaxInt32)
				}

				// Append back the worker port number used so that
				// it can be used by another go routine.
				*availPorts = append(*availPorts, portNum)

				fmt.Printf("output of program[%v]:%v\n", j, output[j])
			}(i, &availPorts, portNum)
		}
		wg.Wait()

		var fittestIndex int = 0
		err = UpdateGraphValues(output, &fittestIndex, &histoValues, &mostFit, &averageValues, &cfg, pop.PopulationSize)
		if err != nil {
			panic(err)
		}

		if cfg.SaveAST {
			err := SaveAST(pop.Individuals[fittestIndex], saveDirectory, c)
			if err != nil {
				panic(err)
			}
		}
	}

	if cfg.PlotFitness {
		saveGraphs(averageValues, mostFit, histoValues, saveDirectory)
	}
}

func RunBenchmark(cxprogram *cxast.CXProgram, solProt *cxast.CXFunction, cfg *EvolveConfig) (intOut float64, err error) {
	if cfg.MazeBenchmark {
		intOut, err = mazeMovesEvaluation(cxprogram, solProt, cfg)
		if err != nil {
			return 0, err
		}
	}

	if cfg.ConstantsBenchmark {
		intOut = perByteEvaluation_Constants(cxprogram, solProt, cfg)
	}

	if cfg.EvensBenchmark {
		intOut = perByteEvaluation_Evens(cxprogram, solProt, cfg)
	}

	if cfg.OddsBenchmark {
		intOut = perByteEvaluation_Odds(cxprogram, solProt, cfg)
	}

	if cfg.PrimesBenchmark {
		intOut = perByteEvaluation_Primes(cxprogram, solProt, cfg)
	}

	if cfg.CompositesBenchmark {
		intOut = perByteEvaluation_Composites(cxprogram, solProt, cfg)
	}

	if cfg.RangeBenchmark {
		intOut = perByteEvaluation_Range(cxprogram, solProt, cfg)
	}

	if cfg.NetworkSimBenchmark {
		intOut = perByteEvaluation_NetworkSim(cxprogram, solProt, cfg)
	}
	return intOut, nil
}

func SaveAST(cxprogram *cxast.CXProgram, saveDir string, generationNum int) error {
	// Save best ASTs per generation
	saveASTDirectory := saveDir + "AST/"
	astName := fmt.Sprintf("Generation_%v", generationNum)

	// Save as human-readable string .txt format
	astAsString := []byte(cxast.ToString(cxprogram))
	if err := ioutil.WriteFile(saveASTDirectory+astName+".txt", astAsString, 0644); err != nil {
		return err
	}

	// Save as serialized bytes.
	astInBytes := cxast.SerializeCXProgramV2(cxprogram, false)
	if err := ioutil.WriteFile(saveASTDirectory+astName+"_serialized"+".ast", []byte(astInBytes), 0644); err != nil {
		return err
	}

	return nil
}

func UpdateGraphValues(output []float64, fittestIndex *int, histoValues, mostFit, averageValues *[]float64, cfg *EvolveConfig, popuSize int) error {
	var total float64 = 0
	var fittest float64 = output[0]
	for z := 0; z < len(output); z++ {
		fitness := output[z]
		total = total + fitness

		// Get Best fitness per generation
		if fitness < fittest {
			fittest = fitness
			*fittestIndex = z
		}

		// Add fitness for histogram
		*histoValues = append(*histoValues, float64(fitness))
	}

	ave := total / float64(popuSize)

	if cfg.UseAntiLog2 {
		ave = math.Pow(2, ave)
		fittest = math.Pow(2, fittest)
	}

	// Add average values for Average fitness graph
	*averageValues = append(*averageValues, ave)

	// Add fittest values for Fittest per generation graph
	*mostFit = append(*mostFit, fittest)
	return nil
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
