// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

package storage // import "github.com/skycoin/cx/goyacc/fileutil/storage"

import (
	"flag"
	"os"
	"runtime"
	"testing"
)

var (
	devFlag = flag.Bool("dev", false, "enable dev tests")
	goFlag  = flag.Int("go", 1, "GOMAXPROCS")
)

func TestMain(m *testing.M) {
	flag.Parse()
	runtime.GOMAXPROCS(*goFlag)
	os.Exit(m.Run())
}
