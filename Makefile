all: tags lint doc build test
tags:
	ctags -R .
doc: man
man: discoteq.1
discoteq.1: discoteq.1.md
	pandoc -s -t man -f markdown discoteq.1.md > discoteq.1

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
	@echo "You're ready to contribute to discoteq!"

proc:
	forego start

lint:
	goimports -l **/*.go
	golint
	go vet

discoteq: discoteq.go chef/config/config.go chef/service.go
	go build -o discoteq

build: discoteq

test: discoteq
	t/chef-test.bats

clean:
	rm -f discoteq
	rm -f tags
	rm -f nohup.out
	rm -f discoteq.1

