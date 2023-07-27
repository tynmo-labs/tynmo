
.PHONY: protoc
protoc:
	protoc --go_out=. --go-grpc_out=. ./server/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./network/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./txpool/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./syncer/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./consensus/proto/*.proto

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	$(eval LATEST_VERSION = $(shell git describe --tags --abbrev=0))
	$(eval COMMIT_HASH = $(shell git rev-parse HEAD))
	$(eval BRANCH = $(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n'))
	$(eval TIME = $(shell date))
	go build -o bin/tynmo -ldflags="\
    	-X 'tynmo/versioning.Version=$(LATEST_VERSION)' \
		-X 'tynmo/versioning.Commit=$(COMMIT_HASH)'\
		-X 'tynmo/versioning.Branch=$(BRANCH)'\
		-X 'tynmo/versioning.BuildTime=$(TIME)'" \
	main.go
