dev: tags lint doc build test
all: doc build
doc: man
man: discoteq.1
build: discoteq


bootstrap:
	# support tools
	go get -u github.com/tools/godep
	go get -u github.com/ddollar/forego
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/vet
	go get -u golang.org/x/tools/cmd/godoc
	go get -u golang.org/x/tools/cmd/goimports
	# support services
	go get -u github.com/ctdk/goiardi
	@(which pandoc > /dev/null && echo "pandoc found.") || echo "pandoc not found! You'll need to install this on your own."
	@(which knife > /dev/null && echo "knife found.") || echo "knife not found! You probably want to install it with ChefDK: https://downloads.chef.io/chef-dk/"
	@(which consul > /dev/null && echo "consul found.") || echo "consul not found! You'll need to install this on your own."
	@echo "You're ready to contribute to discoteq!"

lint:
	goimports -l *.go **/*.go
	golint
	go vet

proc:
	forego start

test: t/chef-test.bats t/consul-test.bats discoteq
	t/chef-test.bats
	t/consul-test.bats

tags:
	ctags -R .

discoteq.1: discoteq.1.md
	pandoc -s -t man -f markdown discoteq.1.md > discoteq.1

discoteq: chef/service.go config/config.go consul/service.go discoteq.go
	go build -o discoteq

clean:
	rm -f discoteq
	rm -f discoteq.1
	rm -f tags
