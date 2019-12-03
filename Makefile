OUTFILE = 3700kvstore
PROJECT_GOFILES = ./*

build: path vendor
	go build -o $(OUTFILE)

build_linux: path
	GOOS=linux GOARCH=amd64 go build -o $(OUTFILE)

path:
	export GOPATH=

vendor:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

clean:
	rm $(OUTFILE) $(OUTFILE).tar.gz

bundle:
	tar -czvf $(OUTFILE).tar.gz $(PROJECT_GOFILES)

copy:
	scp -r $(OUTFILE).tar.gz reedda@login.ccs.neu.edu:/home/reedda/cs3700/raft

publish: vendor build_linux bundle copy