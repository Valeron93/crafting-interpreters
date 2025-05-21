test: build
	./run_all_programs.sh

generate:
	go generate ./ast ./scanner

build:
	mkdir -p ./build/
	go build -o ./build/vl

clean:
	rm -rf build/

fmt:
	go fmt ./parser ./interpreter ./scanner ./cmd/astgen ./stack ./util ./resolver

.PHONY: build