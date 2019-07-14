package main

import (
	"os"
	"os/user"
	"os/exec"
	"fmt"
	"bytes"
	"time"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
	"bufio"
	"encoding/json"
	// "encoding/hex"

	"runtime"

	"regexp"
	
	"flag"

	"path/filepath"

	"net"
	"net/http"
	
	. "github.com/skycoin/cx/cx"
	. "github.com/skycoin/cx/cxgo/parser"
	. "github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/theherk/viper"

	"github.com/amherag/skycoin/src/cipher"
	// "github.com/amherag/skycoin/src/cipher/encoder"
	"github.com/amherag/skycoin/src/coin"
	"github.com/amherag/skycoin/src/readable"
	"github.com/amherag/skycoin/src/skycoin"
	"github.com/amherag/skycoin/src/fiber"
	"github.com/amherag/skycoin/src/util/logging"
	"github.com/amherag/skycoin/src/visor"
	"github.com/amherag/skycoin/src/api"
	"github.com/amherag/skycoin/src/cli"
	"github.com/amherag/skycoin/src/wallet"

	"errors"
)

const VERSION = "0.7.0"

var (
	log             = logging.MustGetLogger("newcoin")
	apiClient       = &http.Client{Timeout: 10 * time.Second}
	genesisBlockURL = "http://127.0.0.1:%d/api/v1/block?seq=0"
)

var (
	// ErrMissingProjectRoot is returned when the project root parameter is missing
	ErrMissingProjectRoot = errors.New("missing project root")
	// ErrMissingSecretKey is returned when genesis secret is missing when distributing coins
	ErrMissingSecretKey = errors.New("missing genesis secret key")

	genesisSignature = ""
)

func getJSON(url string, target interface{}) error {
	r, err := apiClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func initCXBlockchain (initPrgrm []byte, coinname, seckey string) error {
	var err error
	
	// check that data.db does not exist
	// if it does, delete it
	userHome := UserHome()
	dbPath := filepath.Join(userHome, "."+coinname, "data.db")
	if _, err := os.Stat(dbPath); err == nil {
		log.Infof("deleting %s", dbPath)
		err = os.Remove(dbPath)
		if err != nil {
			return err
		}
	}

	if seckey == "" {
		return ErrMissingSecretKey
	}

	genesisSecKey, err := cipher.SecKeyFromHex(seckey)
	if err != nil {
		return err
	}

	configDir := os.Getenv("GOPATH") + "/src/github.com/skycoin/cx/"
	configFile := "fiber.toml"
	configFilepath := filepath.Join(configDir, configFile)
	// check that the config file exists
	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		return err
	}

	projectRoot := os.Getenv("GOPATH") + "/src/github.com/skycoin/cx"
	if projectRoot == "" {
		return ErrMissingProjectRoot
	}
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return err
	}

	coinFile := filepath.Join(projectRoot, fmt.Sprintf("cmd/%[1]s/%[1]s.go", coinname))
	if _, err := os.Stat(coinFile); os.IsNotExist(err) {
		return err
	}

	// get fiber params
	params, err := fiber.NewConfig(configFile, configDir)

	cmd := exec.Command("go", "run", filepath.Join(projectRoot, fmt.Sprintf("cmd/%[1]s/%[1]s.go", coinname)), "-block-publisher=true", fmt.Sprintf("-blockchain-secret-key=%s", seckey),
		"-disable-incoming", "-max-out-msg-len=134217929")

	var genesisSig string
	var genesisBlock readable.Block

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	// fetch gensisSig and gensisBlock
	go func() {
		defer cmd.Process.Kill()

		genesisSigRegex, err := regexp.Compile(`Genesis block signature=([0-9a-zA-Z]+)`)
		if err != nil {
			log.Error("error in regexp for genesis block signature")
			log.Error(err)
			return
		}

		scanner := bufio.NewScanner(stdoutIn)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {

			m := scanner.Text()
			log.Info("Scanner: " + m)
			if genesisSigRegex.MatchString(m) {
				genesisSigSubString := genesisSigRegex.FindStringSubmatch(m)
				genesisSig = genesisSigSubString[1]

				// get genesis block
				err = getJSON(fmt.Sprintf(genesisBlockURL, params.Node.WebInterfacePort), &genesisBlock)
			
				return
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrIn)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			log.Error(scanner.Text())
		}
	}()

	cmd.Wait()

	// check that we were able to get genesisSig and genesisUxID

	if genesisSig != "" && len(genesisBlock.Body.Transactions) != 0 {
		genesisSignature = genesisSig
		log.Infof("genesis sig: %s", genesisSig)

		// -- create new skycoin daemon to inject distribution transaction -- //
		if err != nil {
			log.Error("error getting fiber parameters")
			return err
		}

		// get node config
		params.Node.DataDirectory = fmt.Sprintf("$HOME/.%s", coinname)
		nodeConfig := skycoin.NewNodeConfig("", params.Node)

		// create a new fiber coin instance
		newcoin := skycoin.NewCoin(
			skycoin.Config{
				Node: nodeConfig,
			},
			log,
		)

		// parse config values
		newcoin.ParseConfig()

		// dconf := newcoin.ConfigureDaemon()
		vconf := newcoin.ConfigureVisor()

		userHome := UserHome()
		dbPath := filepath.Join(userHome, "."+coinname, "data.db")

		// log.Infof("opening visor db: %s", dconf.Visor.DBPath)
		log.Infof("opening visor db: %s", dbPath)
		db, err := visor.OpenDB(dbPath, false)
		if err != nil {
			log.Error("Error opening DB")
			return err
		}
		defer db.Close()

		vs, err := visor.New(vconf, db, nil)
		if err != nil {
			log.Error("Error with NewVisor")
			return err
		}

		headSeq, _, err := vs.HeadBkSeq()
		if err != nil {
			log.Error("Error with HeadBkSeq")
			return err
		} else if headSeq == 0 {
			if len(genesisBlock.Body.Transactions) != 0 {
				var tx coin.Transaction

				UxID := genesisBlock.Body.Transactions[0].Out[0].Hash
				output := cipher.MustSHA256FromHex(UxID)
				tx.PushInput(output)

				addr := cipher.MustDecodeBase58Address("TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX")
				tx.PushOutput(addr, uint64(1e10), 10000, initPrgrm)

				seckeys := make([]cipher.SecKey, 1)
				seckey := genesisSecKey.Hex()
				seckeys[0] = cipher.MustSecKeyFromHex(seckey)
				tx.SignInputs(seckeys)

				tx.UpdateHeader()
				err = tx.Verify()

				if err != nil {
					log.Panic(err)
				}

				_, _, _, err := vs.InjectUserTransaction(tx)
				if err != nil {
					panic(err)
				}
			} else {
				log.Error("ERROR: len genesis block was zero")
			}
		} else {
			log.Error("ERROR: headSeq not zero")
		}
	} else {
		log.Error("error getting genesis block")
	}
	return err
}

func runNode(mode string, options cxCmdFlags) *exec.Cmd {
	switch mode {
	case "publisher":
		return exec.Command("cxcoin", "-enable-all-api-sets",
		"-block-publisher=true",
			"-localhost-only",
			"-disable-default-peers",
			"-custom-peers-file=localhost-peers.txt",
			"-download-peerlist=false",
			"-launch-browser=false",
			fmt.Sprintf("-blockchain-secret-key=%s", options.secKey),
			fmt.Sprintf("-genesis-address=%s", options.genesisAddress),
			fmt.Sprintf("-genesis-signature=%s", options.genesisSignature),
			fmt.Sprintf("-blockchain-public-key=%s", options.pubKey),
			"-max-txn-size-unconfirmed=134217728",
		 	"-max-txn-size-create-block=134217728",
			"-max-block-size=134217728",
			"-max-in-msg-len=134217929",
			"-max-out-msg-len=134217929", // I don't know why this value, but the logger stated a value >= than this is needed
		)
	case "peer":
	return exec.Command("cxcoin", "-enable-all-api-sets",
		"-localhost-only",
		"-disable-default-peers",
		"-custom-peers-file=localhost-peers.txt",
		"-download-peerlist=false",
		"-launch-browser=false",
		fmt.Sprintf("-genesis-address=%s", options.genesisAddress),
		fmt.Sprintf("-genesis-signature=%s", options.genesisSignature),
		fmt.Sprintf("-blockchain-public-key=%s", options.pubKey),
		// "-web-interface-port=$(expr $2 + 420)",
		fmt.Sprintf("-web-interface-port=%d", options.port + 420),
		fmt.Sprintf("-port=%d", options.port),
		fmt.Sprintf("-data-dir=/tmp/%d", options.port),
		"-max-txn-size-unconfirmed=134217728",
		"-max-txn-size-create-block=134217728",
		"-max-block-size=134217728",
		"-max-in-msg-len=134217929",
		"-max-out-msg-len=134217929", // I don't know why this value, but the logger stated a value >= than this is needed
	)
	default:
		return nil
	}
}

func main() {
	checkCXPathSet()

	runtime.LockOSThread()
	runtime.GOMAXPROCS(2)

	options := defaultCmdFlags()
	compileOutput := "o"

	registerFlags(&options)
	flag.Parse()

	// Propagate some options out to other packages
	DebugLexer = options.debugLexer   // in package parser

	if options.publisherMode || options.peerMode {
		var cmd *exec.Cmd
		if options.publisherMode {
			cmd = runNode("publisher", options)
		} else if options.peerMode {
			cmd = runNode("peer", options)
		}

		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		cmd.Start()

		go func() {
			defer cmd.Process.Kill()

			scanner := bufio.NewScanner(stdoutIn)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				m := scanner.Text()
				log.Info("Scanner: " + m)
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderrIn)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				log.Error(scanner.Text())
			}
		}()
	}

	// Generate a CX chain address.
	if options.genAddress {
		// Create a random seed to create a temporary wallet.
		seed := cli.MakeAlphanumericSeed()
		wltOpts := wallet.Options{
			Label: "cxcoin",
			Seed: seed,
		}

		// Generate temporary wallet.
		wlt, err := cli.GenerateWallet(wallet.NewWalletFilename(), wltOpts, 1)
		if err != nil {
			panic(err)
		}

		rw := wallet.NewReadableWallet(wlt)

		output, err := json.MarshalIndent(rw, "", "    ")
		if err != nil {
			panic(err)
		}

		// Print all the wallet data.
		fmt.Println(string(output))

		return
	}
	
	if options.walletMode {
		if options.walletSeed == "" {
			fmt.Println("creating a wallet requires a seed provided with --wallet-seed")
			return
		}
		if options.walletId == "" {
			// Although there is a default walletId.
			// This error should only occur if the user intentionally provides an empty id.
			fmt.Println("creating a wallet requires an id provided with --wallet-id")
			return
		}

		wltOpts := wallet.Options{
			Label: "cxcoin",
			Seed: options.walletSeed,
		}
		
		wlt, err := cli.GenerateWallet(options.walletId, wltOpts, 1)
		if err != nil {
			panic(err)
		}
		// To Do: This needs to be changed or any CX chains will constantly be destroyed after each reboot.
		err = wlt.Save("/tmp/6001/wallets/")
		if err != nil {
			panic(err)
		}

		wltJSON, err := json.MarshalIndent(wlt.Meta, "", "\t")
		if err != nil {
			panic(err)
		}

		// Printing JSON with wallet information
		fmt.Println(string(wltJSON))
		
		return
	}
	
	if options.printHelp {
		printHelp()
		return
	}
	if options.printVersion {
		printVersion()
		return
	}
	if options.newProject {
		initNewProject()
		return
	}
	if options.compileMode {
		options.baseOutput = true
	}
	if options.initialHeap != "" {
		INIT_HEAP_SIZE = parseMemoryString(options.initialHeap)
	}
	if options.maxHeap != "" {
		MAX_HEAP_SIZE = parseMemoryString(options.maxHeap)
		if MAX_HEAP_SIZE < INIT_HEAP_SIZE {
			// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
			INIT_HEAP_SIZE = MAX_HEAP_SIZE
		}
	}
	if options.stackSize != "" {
		STACK_SIZE = parseMemoryString(options.stackSize)
		DataOffset = STACK_SIZE
	}
	if options.minHeapFreeRatio != float64(0) {
		MIN_HEAP_FREE_RATIO = float32(options.minHeapFreeRatio)
	}
	if options.maxHeapFreeRatio != float64(0) {
		MAX_HEAP_FREE_RATIO = float32(options.maxHeapFreeRatio)
	}

	// options, file pointers, filenames
	cxArgs, sourceCode, fileNames := parseArgsForCX(flag.Args())

	PRGRM = MakeProgram()

	if options.webMode {
		ServiceMode()
		return
	}

	if options.ideMode {
		IdeServiceMode()
		ServiceMode()
		return
	}

	if options.webPersistentMode {
		go ServiceMode()
		PersistentServiceMode()
		return
	}

	if options.compileMode && options.replMode {
		fmt.Println("Error: Options --compile and --repl are mutually exclusive.")
		return
	}

	if options.tokenizeMode {
		var r *os.File
		var w *os.File
		var err error

		if len(fileNames) == 0 {
			r = os.Stdin
		} else {
			sourceFilename := fileNames[0]
			if len(fileNames) > 1 {
				fmt.Fprintln(os.Stderr, "Multiple source files detected. Ignoring all except", sourceFilename)
			}
			r, err = os.Open(sourceFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading:", sourceFilename, err)
				return
			}
			defer r.Close()
		}

		if options.compileOutput == "" {
			w = os.Stdout
		} else {
			tokenFilename := options.compileOutput
			w, err = os. Create(tokenFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing:", tokenFilename, err)
				return
			}
			defer w.Close()
		}

		Tokenize(r, w)
		return
	}

	// var bcPrgrm *CXProgram
	var sPrgrm []byte
	if options.transactionMode || options.broadcastMode {
		resp, err := http.Get("http://127.0.0.1:6420/api/v1/programState?addrs=TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &sPrgrm); err != nil {
			fmt.Println(string(body))
			return
		}

		memOff := GetSerializedMemoryOffset(sPrgrm)
		stackSize := GetSerializedStackSize(sPrgrm)
		// sPrgrm with Stack and Heap
		sPrgrmSH := sPrgrm[:memOff]
		// Appending new stack
		sPrgrmSH = append(sPrgrmSH, make([]byte, stackSize)...)
		// Appending data segment
		sPrgrmSH = append(sPrgrmSH, sPrgrm[memOff:]...)
		// Appending new heap
		sPrgrmSH = append(sPrgrmSH, make([]byte, INIT_HEAP_SIZE)...)
		
		// Deserialize(sPrgrm).RunCompiled(0, cxArgs)
		// bcPrgrm = Deserialize(sPrgrm)
		// prevMemSize := len(PRGRM.Memory)
		PRGRM = Deserialize(sPrgrmSH)
		// // Adding 
		// PRGRM.Memory = append(make([]byte, STACK_SIZE), PRGRM.Memory...)
		// PRGRM.Memory
		// PRGRM.Memory = append(PRGRM.Memory, make([]byte, prevMemSize - len(PRGRM.Memory))...)
		DataOffset = PRGRM.HeapStartsAt
	}

	cxgo0.PRGRM0 = PRGRM

	// setting project's working directory
	if !options.replMode && len(sourceCode) > 0 {
		cxgo0.PRGRM0.Path = getWorkingDirectory(sourceCode[0].Name())
	}

	// Copy the contents of the file pointers containing the CX source
	// code into sourceCodeCopy
	sourceCodeCopy := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		sourceCodeCopy[i] = string(tmp.Bytes())
	}

	var prePkg *CXPackage
	parseErrors := 0
	
	// We need to traverse the elements by hierarchy first add all the
	// packages and structs at the same time then add globals, as these
	// can be of a custom type (and it could be imported) the signatures
	// of functions and methods are added in the cxgo0.y pass
	if len(sourceCode) > 0 {
		
		reMultiCommentOpen  := regexp.MustCompile(`/\*`)
		reMultiCommentClose := regexp.MustCompile(`\*/`)
		reComment := regexp.MustCompile("//")

		rePkg       := regexp.MustCompile("package")
		rePkgName   := regexp.MustCompile("(^|[\\s])package\\s+([_a-zA-Z][_a-zA-Z0-9]*)")
		reStrct     := regexp.MustCompile("type")
		reStrctName := regexp.MustCompile("(^|[\\s])type\\s+([_a-zA-Z][_a-zA-Z0-9]*)?\\s")

		reGlbl     := regexp.MustCompile("var")
		reGlblName := regexp.MustCompile("(^|[\\s])var\\s([_a-zA-Z][_a-zA-Z0-9]*)")

		reBodyOpen  := regexp.MustCompile("{")
		reBodyClose := regexp.MustCompile("}")

		// 1. Identify all the packages and structs
		for ix, source := range sourceCodeCopy {
			filename := fileNames[ix]

			reader := strings.NewReader(source)
			scanner := bufio.NewScanner(reader)
			var commentedCode bool
			var lineno = 0
			for scanner.Scan() {
				line := scanner.Bytes()
				lineno++

				// Identify whether we are in a comment or not.
				commentLoc := reComment.FindIndex(line)
				multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
				multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)
				if commentedCode && multiCommentCloseLoc != nil {
					commentedCode = false
				}
				if commentedCode {
					continue
				}
				if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc == nil {
					commentedCode = true
					continue
				}

				// At this point we know that we are *not* in a comment

				// 1a. Identify all the packages
				if loc := rePkg.FindIndex(line); loc != nil {
					if (commentLoc != nil && commentLoc[0] < loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
						// then it's commented out
						continue
					}

					if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
						if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match) - 1]); err != nil {
							// then it hasn't been added
							newPkg := MakePackage(match[len(match) - 1])
							cxgo0.PRGRM0.AddPackage(newPkg)
							prePkg = newPkg
						} else {
							prePkg = pkg
						}
					}
				}

				// 1b. Identify all the structs
				if loc := reStrct.FindIndex(line); loc != nil {
					if (commentLoc != nil && commentLoc[0] < loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
						// then it's commented out
						continue
					}

					if match := reStrctName.FindStringSubmatch(string(line)); match != nil {
						if prePkg == nil {
							println(CompilationError(filename, lineno),
							        "No package defined")
						} else if _, err := cxgo0.PRGRM0.GetStruct(match[len(match) - 1], prePkg.Name); err != nil {
							// then it hasn't been added
							strct := MakeStruct(match[len(match) - 1])
							prePkg.AddStruct(strct)
						}
					}
				}
			}
		} // for range sourceCodeCopy

		// 2. Identify all global variables
		//    We also identify packages again, so we know to what
		//    package we're going to add the variable declaration to.
		for _, source := range sourceCodeCopy {
			// inBlock needs to be 0 to guarantee that we're in the global scope
			var inBlock int
			var commentedCode bool

			scanner := bufio.NewScanner(strings.NewReader(source))
			for scanner.Scan() {
				line := scanner.Bytes()
				
				// we need to ignore function bodies
				// it'll also ignore struct declaration's bodies, but this doesn't matter
				commentLoc := reComment.FindIndex(line)
				
				multiCommentOpenLoc  := reMultiCommentOpen.FindIndex(line)
				multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)

				if commentedCode && multiCommentCloseLoc != nil {
					commentedCode = false
				}

				if commentedCode {
					continue
				}

				if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc == nil {
					commentedCode = true
					// continue
				}

				// we search for packages at the same time, so we can know to what package to add the global
				if loc := rePkg.FindIndex(line); loc != nil {
					if (commentLoc != nil && commentLoc[0] < loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
						// then it's commented out
						continue
					}
					
					if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
						if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match) - 1]); err != nil {
							// then it hasn't been added
							prePkg = MakePackage(match[len(match) - 1])
							cxgo0.PRGRM0.AddPackage(prePkg)
						} else {
							prePkg = pkg
						}
					}
				}
				
				if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
					for _, loc := range locs {
						if !(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
							// then it's outside of a */, e.g. `*/ }`
							if (commentLoc == nil && multiCommentOpenLoc == nil && multiCommentCloseLoc == nil) ||
								(commentLoc != nil && commentLoc[0] > loc[0]) ||
								(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] > loc[0]) ||
								(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] < loc[0]) {
								// then we have an uncommented opening bracket
								inBlock++
							}
						}
					}
				}

				if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
					for _, loc := range locs {
						if !(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
							if (commentLoc == nil && multiCommentOpenLoc == nil && multiCommentCloseLoc == nil) ||
								(commentLoc != nil && commentLoc[0] > loc[0]) ||
								(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] > loc[0]) ||
								(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] < loc[0]) {
								// then we have an uncommented closing bracket
								inBlock--
							}
						}
					}
				}

				// we could have this situation: {var local i32}
				// but we don't care about this, as the later passes will throw an error as it's invalid syntax
				
				if loc := rePkg.FindIndex(line); loc != nil {
					if (commentLoc != nil && commentLoc[0] < loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
						// then it's commented out
						continue
					}
					
					if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
						if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match) - 1]); err != nil {
							// it should be already present
							panic(err)
						} else {
							prePkg = pkg
						}
					}
				}

				// finally, if we read a "var" and we're in global scope, we add the global without any type
				// the type will be determined later on
				if loc := reGlbl.FindIndex(line); loc != nil {
					if (commentLoc != nil && commentLoc[0] < loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) || inBlock != 0 {
						// then it's commented out or inside a block
						continue
					}

					if match := reGlblName.FindStringSubmatch(string(line)); match != nil {
						if _, err := prePkg.GetGlobal(match[len(match) - 1]); err != nil {
							// then it hasn't been added
							arg := MakeArgument(match[len(match) - 1], "", 0)
							arg.Offset = -1
							arg.Package = prePkg
							prePkg.AddGlobal(arg)
						}
					}
				}
			}
		}

		// cxgo0.Parse(allSC)
		for i, source := range sourceCodeCopy {
			source = source + "\n"
			if len(fileNames) > 0 {
				cxgo0.CurrentFileName = fileNames[i]
			}
			parseErrors += cxgo0.Parse(source)
		}
	}

	PRGRM = cxgo0.PRGRM0
	if FoundCompileErrors || parseErrors > 0 {
		os.Exit(CX_COMPILATION_ERROR)
	}

	if osPkg, err := PRGRM.GetPackage(OS_PKG); err == nil {

		arg0 := MakeArgument(OS_ARGS, "", -1).AddType(TypeNames[TYPE_UNDEFINED])
		arg0.Package = osPkg

		arg1 := MakeArgument(OS_ARGS, "", -1).AddType(TypeNames[TYPE_STR])
		arg1 = DeclarationSpecifiers(arg1, 0, DECL_BASIC)
		arg1 = DeclarationSpecifiers(arg1, 0, DECL_SLICE)


		DeclareGlobalInPackage(osPkg, arg0, arg1, nil, false)
	}

	// The last pass of parsing that generates the actual output.
	for i, source := range sourceCodeCopy {
		source = source + "\n"
		LineNo = 1
		b := bytes.NewBufferString(source)
		if len(fileNames) > 0 {
			CurrentFile = fileNames[i]
		}
		parseErrors += Parse(NewLexer(b))
	}

	if FoundCompileErrors || parseErrors > 0 {
		os.Exit(CX_COMPILATION_ERROR)
	}

	if len(sourceCode) == 0 {
		mod := MakePackage(MAIN_PKG)
		PRGRM.AddPackage(mod)
		fn := MakeFunction(MAIN_FUNC, CurrentFile, LineNo)
		mod.AddFunction(fn)

		ReplTargetFn = MAIN_FUNC
	} else {
		if _, err := PRGRM.GetFunction(MAIN_FUNC, MAIN_PKG); err == nil {
			ReplTargetFn = MAIN_FUNC
		} else {
			// then it's a library, not an app
		}
	}

	// adding *init function that initializes all the global variables
	addInitFunction(PRGRM)
	
	LineNo = 0

	if FoundCompileErrors {
		os.Exit(CX_COMPILATION_ERROR)
	}

	if options.replMode || len(sourceCode) == 0 {
		repl()
	} else if !CompileMode && !BaseOutput && len(sourceCode) > 0 {
		if options.blockchainMode {
			// Initializing the CX chain.
			err := PRGRM.RunCompiled(0, cxArgs)
			if err != nil {
				panic(err)
			}
			
			PRGRM.RemovePackage(MAIN_FUNC)
			
			s := Serialize(PRGRM, 1)
			s = ExtractBlockchainProgram(s, s)

			configDir := os.Getenv("GOPATH") + "/src/github.com/skycoin/cx/"
			configFile := "fiber"
			
			cmd := exec.Command("go", "install", "./cmd/newcoin/...")
			cmd.Start()
			cmd.Wait()
			
			cmd = exec.Command("newcoin", "createcoin",
				fmt.Sprintf("--coin=%s", options.programName),
				fmt.Sprintf("--template-dir=%s%s", os.Getenv("GOPATH"), "/src/github.com/skycoin/cx/template"),
				"--config-file=" + configFile + ".toml",
				"--config-dir=" + configDir,
			)
			cmd.Start()
			cmd.Wait()
			
			cmd = exec.Command("go", "install", "./cmd/cxcoin/...")
			cmd.Start()
			cmd.Wait()

			err = initCXBlockchain(s, options.programName, options.secKey)
			if err != nil {
				panic(err)
			}
			fmt.Println("\ngenesis signature:", genesisSignature)

			viper.SetConfigName(configFile) // name of config file (without extension)
			viper.AddConfigPath(".")      // optionally look for config in the working directory
			err = viper.ReadInConfig()   // Find and read the config file
			if err != nil { // Handle errors reading the config file
				panic(err)
			}

			viper.Set("node.genesis_signature_str", genesisSignature)
			viper.WriteConfig()
		
			cmd = exec.Command("newcoin", "createcoin",
				fmt.Sprintf("--coin=%s", options.programName),
				fmt.Sprintf("--template-dir=%s%s", os.Getenv("GOPATH"), "/src/github.com/skycoin/cx/template"),
				"--config-file=" + configFile + ".toml",
				"--config-dir=" + configDir,
			)
			cmd.Start()
			cmd.Wait()
			cmd = exec.Command("go", "install", "./cmd/cxcoin/...")
			cmd.Start()
			cmd.Wait()
		} else if options.broadcastMode {
			// Setting the CX runtime to run `PRGRM`.
			PRGRM.SelectProgram()

			// Resetting memory
			dsPrgrm := Deserialize(sPrgrm)
			PRGRM.StackPointer = 0
			PRGRM.HeapPointer = dsPrgrm.HeapPointer
			
			s := Serialize(PRGRM, 1)
			txnCode := ExtractTransactionProgram(sPrgrm, s)

			// All these HTTP requests need to be dropped in favor of calls to calls to functions
			// from the `cli` or `api` Skycoin packages
			addr := fmt.Sprintf("http://127.0.0.1:%d", options.port + 420)
			skycoinClient := api.NewClient(addr)
			csrfToken, err := skycoinClient.CSRF()
			if err != nil {
				panic(err)
			}

			url := fmt.Sprintf("http://127.0.0.1:%d/api/v1/wallet/transaction", options.port + 420)

			var dataMap map[string]interface{}
			dataMap = make(map[string]interface{}, 0)
			dataMap["mainExprs"] = txnCode
			dataMap["hours_selection"] = map[string]string{"type": "manual"}
			// dataMap["wallet_id"] = map[string]string{"id": options.walletId}
			dataMap["wallet_id"] = string(options.walletId)
			dataMap["to"] = []interface{}{map[string]string{"address": "2PBcLADETphmqWV7sujRZdh3UcabssgKAEB", "coins": "1", "hours": "0"}}

			
			jsonStr, err := json.Marshal(dataMap)
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			req.Header.Set("X-CSRF-Token", csrfToken)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			var respBody map[string]interface{}
			if err := json.Unmarshal(body, &respBody); err != nil {
				// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
				fmt.Println(string(body))
				return
			}

			url = fmt.Sprintf("http://127.0.0.1:%d/api/v1/injectTransaction", options.port + 420)
			dataMap = make(map[string]interface{}, 0)
			dataMap["rawtx"] = respBody["encoded_transaction"]

			jsonStr, err = json.Marshal(dataMap)
			if err != nil {
				panic(err)
			}

			req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			req.Header.Set("X-CSRF-Token", csrfToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err = client.Do(req)
			if err != nil {
				panic(err)
			}

			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
		} else {
			// Normal run of a CX program.
			err := PRGRM.RunCompiled(0, cxArgs)
			if err != nil {
				panic(err)
			}
			
			if AssertFailed() {
				os.Exit(CX_ASSERT)
			}
		}
		
		if BaseOutput {
			//PRGRM.Compile(true)
		}
		if CompileMode {
			baseFilename := fmt.Sprintf("%s.go", compileOutput)
			build := exec.Command("go", "build", baseFilename)
			build.Run()
			removeBase := exec.Command("rm", baseFilename)
			removeBase.Run()
		}
	}
}


// Used for the -heap-initial, -heap-max and -stack-size flags.
// This function parses, for example, "1M" to 1048576 (the corresponding number of bytes)
// Possible suffixes are: G or g (gigabytes), M or m (megabytes), K or k (kilobytes)
func parseMemoryString (s string) int {
	suffix := s[len(s) - 1]
	_, notSuffix := strconv.ParseFloat(string(suffix), 64)

	if notSuffix == nil {
		// then we don't have a suffix
		num, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			// malformed size
			return -1
		}

		return int(num)
	} else {
		// then we have a suffix
		num, err := strconv.ParseFloat(s[:len(s)-1], 64)

		if err != nil {
			// malformed size
			return -1
		}

		// The user can use suffixes to give as input gigabytes, megabytes or kilobytes.
		switch suffix {
		case 'G', 'g':
			return int(num * 1073741824)
		case 'M', 'm':
			return int(num * 1048576)
		case 'K', 'k':
			return int(num * 1024)
		default:
			return -1
		}
	}
}

// var PRGRM *CXProgram

func unsafeEval (code string) (out string) {
	var lexer *Lexer
	defer func() {
		if r := recover(); r != nil {
			out = fmt.Sprintf("%v", r)
			lexer.Stop()
		}  
	}()
	
	// storing strings sent to standard output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	LineNo = 0

	PRGRM = MakeProgram()
	cxgo0.PRGRM0 = PRGRM

	cxgo0.Parse(code)

	PRGRM = cxgo0.PRGRM0

	lexer = NewLexer(bytes.NewBufferString(code))
	Parse(lexer)
	//yyParse(lexer)

	addInitFunction(PRGRM)

	if err := PRGRM.RunCompiled(0, nil); err != nil {
		PRGRM = MakeProgram()
		return fmt.Sprintf("%s", err)
	}

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	
	w.Close()
	os.Stdout = old // restoring the real stdout
	out = <-outC

	PRGRM = MakeProgram()
	return out
}

func Eval (code string) string {
	runtime.GOMAXPROCS(2)
	ch := make(chan string, 1)

	var result string
	
	go func() {
		result = unsafeEval(code)
		ch <- result
	}()

	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

	select {
	case <-ch:
		return result
	case <-timer.C:
		PRGRM = MakeProgram()
		return "Timed out."
	}
}

type SourceCode struct {
	Code string
}

func ServiceMode () {
	host := ":5336"

	mux := http.NewServeMux()
	
	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var b []byte
		var err error
		if b, err = ioutil.ReadAll(r.Body); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		
		var source SourceCode
		if err := json.Unmarshal(b, &source); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := r.ParseForm(); err == nil {
			fmt.Fprintf(w, "%s", Eval(source.Code + "\n"))
		}
	})

	if listener, err := net.Listen("tcp", host); err == nil {
		fmt.Println("Starting CX web service on http://127.0.0.1:5336/")
		http.Serve(listener, mux)
	}
}

func IdeServiceMode() {
	// Leaps's host address
	ideHost := "localhost:5335"

	// Working directory for ide
	sharedPath := fmt.Sprintf("%s/src/github.com/skycoin/cx", os.Getenv("GOPATH"))

	// Start Leaps
	// cmd = `leaps -address localhost:5335 $GOPATH/src/skycoin/cx`
	cmnd := exec.Command("leaps", "-address", ideHost, sharedPath)

	// Just leave start command
	cmnd.Start()
}

func PersistentServiceMode() {
	fmt.Println("Start persistent for service mode!")

	fi := bufio.NewReader(os.Stdin)

	for {
		var inp string
		var ok bool
		
		printPrompt()

		if inp, ok = readline(fi); ok {
			if isJSON(inp) {
				var err error
				client := &http.Client{}
                body := bytes.NewBufferString(inp)
				req, err := http.NewRequest("GET", "http://127.0.0.1:5336/eval", body)
				if err != nil {
					fmt.Println(err)
					return
				}

				if resp, err := client.Do(req); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(resp.Status)
				}
			}
		}
	}
}

func getWorkingDirectory (file string) string {
	file = filepath.FromSlash(file)
	var c int = len(file) - 1
	for ; c > 0; c-- {
		if file[c - 1] == os.PathSeparator {
			break
		}
	}
	
	return file[:c]
}

func printPrompt () {
	if ReplTargetMod != "" {
		fmt.Println(fmt.Sprintf(":package %s ...", ReplTargetMod))
		fmt.Printf("* ")
	} else if ReplTargetFn != "" {
		fmt.Println(fmt.Sprintf(":func %s {...", ReplTargetFn))
		fmt.Printf("\t* ")
	} else if ReplTargetStrct != "" {
		fmt.Println(fmt.Sprintf(":struct %s {...", ReplTargetStrct))
		fmt.Printf("\t* ")
	} else {
		fmt.Printf("* ")
	}
}

func repl () {
	fmt.Println("CX", VERSION)
	fmt.Println("More information about CX is available at http://cx.skycoin.net/ and https://github.com/skycoin/cx/")

	InREPL = true

	// fi := bufio.NewReader(os.NewFile(0, "stdin"))
	fi := bufio.NewReader(os.Stdin)
	// scanner := bufio.NewScanner(os.Stdin)

	for {
		var inp string
		var ok bool

		printPrompt()
		
		if inp, ok = readline(fi); ok {
			if ReplTargetFn != "" {
				inp = fmt.Sprintf(":func %s {\n%s\n}\n", ReplTargetFn, inp)
			}
			if ReplTargetMod != "" {
				inp = fmt.Sprintf("%s", inp)
			}
			if ReplTargetStrct != "" {
				inp = fmt.Sprintf(":struct %s {\n%s\n}\n", ReplTargetStrct, inp)
			}

			b := bytes.NewBufferString(inp)

			Parse(NewLexer(b))
			//yyParse(NewLexer(b))
		} else {
			if ReplTargetFn != "" {
				ReplTargetFn = ""
				continue
			}

			if ReplTargetStrct != "" {
				ReplTargetStrct = ""
				continue
			}

			if ReplTargetMod != "" {
				ReplTargetMod = ""
				continue
			}

			fmt.Printf("\nBye!\n")
			break
		}
	}
}

func initNewProject () {
	var name string
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Name of the project: ")
	name, _ = reader.ReadString('\n')

	fmt.Printf("Creating project %s%s/", SRCPATH, name)

	os.MkdirAll(fmt.Sprintf("%s%s", SRCPATH, name[:len(name) - 1]), 0751)
}

func checkCXPathSet () {
	if os.Getenv("CXPATH") == "" {
		usr, err := user.Current()
		if err != nil { 
			panic(err)
		}

		CXPATH = usr.HomeDir + "/cx/"
		BINPATH = CXPATH + "bin/"
		PKGPATH = CXPATH + "pkg/"
		SRCPATH = CXPATH + "src/"
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	COREPATH = filepath.Dir(ex)
}

func addInitFunction (PRGRM *CXProgram) {
	if main, err := PRGRM.GetPackage(MAIN_PKG); err == nil {
		initFn := MakeFunction(SYS_INIT_FUNC, CurrentFile, LineNo)
		main.AddFunction(initFn)

		FunctionDeclaration(initFn, nil, nil, SysInitExprs)
		PRGRM.SelectFunction(MAIN_FUNC)
	} else {
		panic(err)
	}
}

// ----------------------------------------------------------------
//                     Utility functions

func readline (fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')

	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)

	for _, ch := range s {
		if ch == rune(4) {
			err = io.EOF
			break
		}
	}

	if err != nil {
		return "", false
	}
	
	return s, true
}

func isJSON(str string) bool {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(str), &js)
	return err == nil
}
