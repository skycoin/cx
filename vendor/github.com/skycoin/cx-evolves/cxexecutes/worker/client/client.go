package client

import (
	"time"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	cxast "github.com/skycoin/cx/cx/ast"
)

type CallWorkerConfig struct {
	Program   *cxast.CXProgram
	Input     []byte
	OutputArg *cxast.CXArgument
}

func CallWorker(cWorker CallWorkerConfig, workerAddr string, result *worker.Result) {
	defer erpc.SetLoggerLevel("INFO")()
	cli := erpc.NewPeer(erpc.PeerConfig{RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()
	cli.SetTLSConfig(erpc.GenerateTLSConfigForClient())
	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(workerAddr)
	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}

	defer sess.Close()

	args := &worker.Args{
		Program:      cxast.SerializeCXProgramV2(cWorker.Program, false),
		Inputs:       cWorker.Input,
		OutputOffset: cWorker.OutputArg.Offset,
		OutputSize:   cWorker.OutputArg.TotalSize,
	}

	stat = sess.Call(
		worker.RunProgram,
		args,
		&result,
	).Status()

	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}

}

// Push push handler
type Push struct {
	erpc.PushCtx
}

// Push handles '/push/status' message
func (p *Push) Status(arg *string) *erpc.Status {
	erpc.Printf("%s", *arg)
	return nil
}
