#!/bin/sh
# support tools
go get -u github.com/tools/godep
go get -u github.com/ddollar/forego
go get -u code.google.com/p/mango-doc
# support services
go get -u github.com/ctdk/goiardi
# dependant libs
go get -u github.com/marpaia/chef-golang

nohup goiardi &
pushd stubs
knife upload .



