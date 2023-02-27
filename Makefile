# The binary to build, base on bin
BIN=search-engine-worker

vendor:
	go mod tidy

worker: vendor
	go build -o build/${BIN} cmd/worker/*.go

start_worker: worker
	./build/${BIN} --config-source=file --config-file=env/development.yml # --logger-enable-debug=true

clean:
	if [ -f ${BIN} ] ; then rm ${BIN} ; fi

lint-prepare:
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install build run stop vendor lint-prepare lint
