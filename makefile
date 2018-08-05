GOROOT=/usr/local/go

all: helloWasm

helloWasm: output.wasm helloWasm.go
	PATH=$(GOROOT)/misc/wasm:$(GOROOT)/bin:$PATH go build -o helloWasm helloWasm.go

output.wasm: helloWasm.go
	PATH=$(GOROOT)/misc/wasm:$(GOROOT)/bin:$PATH GOOS=js GOARCH=wasm go build -o bundle.wasm bundle.go

run: output.wasm helloWasm
	./helloWasm

clean:
	rm -f helloWasm *.wasm
