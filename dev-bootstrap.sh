#!/bin/sh

# fixes build/test problems
if [ ! -d $GOPATH/src/github.com/discoteq/discoteq.go ]; then
  mkdir -p $GOPATH/src/github.com/discoteq
  ln -s ./ $GOPATH/src/github.com/discoteq/discoteq.go
fi

make dev-bootstrap
