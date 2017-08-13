#! /bin/sh
nex cx.nex && goyacc -ex -o cx.go cx.y && go build cx.go cx.nn.go && time ./cx $@
