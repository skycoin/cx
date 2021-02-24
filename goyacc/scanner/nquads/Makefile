# Copyright (c) 2014 The scanner Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: all clean editor todo cover nuke

grep=--include=*.go *.l

all: editor
	go vet || true
	golint .
	go install
	make todo

cover:
	t=$(shell tempfile) ; go test -coverprofile $$t && go tool cover -html $$t && unlink $$t

editor: scanner.go
	go fmt
	go test -i
	go test
	go build

scanner.go: scanner.l
	golex -o $@ $<
	go fmt $@

todo:
	@grep -n ^[[:space:]]*_[[:space:]]*=[[:space:]][[:alpha:]][[:alnum:]]* $(grep) || true
	@grep -n TODO $(grep) || true
	@grep -n BUG $(grep) || true
	@grep -n println $(grep) || true

clean:
	go clean

nuke: clean
	go clean -i
