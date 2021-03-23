# Copyright (c) 2014 The scanner Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: all clean cover editor later nuke todo

grep=--include=*.go --include=*.l

all: editor
	go tool vet -printfuncs "Log:0,Logf:0" *.go || true
	go install
	make todo

clean:
	@go clean

cover:
	t=$(shell tempfile) ; go test -coverprofile $$t && go tool cover -html $$t && unlink $$t

editor: scanner.go
	go fmt
	go test -i
	go test
	go build

later:
	@grep -n $(grep) LATER * || true
	@grep -n $(grep) MAYBE * || true

nuke: clean
	go clean -i

scanner.go: scanner.l
	golex -o $@ $<

todo:
	@grep -nr $(grep) ^[[:space:]]*_[[:space:]]*=[[:space:]][[:alpha:]][[:alnum:]]* * || true
	@grep -nr $(grep) TODO * || true
	@grep -nr $(grep) BUG * || true
	@grep -nr $(grep) println * || true
