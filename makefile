all: wasm-basic-triangle

wasm-basic-triangle: output.wasm main.go
	go build -o wasm-basic-triangle main.go

output.wasm: bundle.go
	GOOS=js GOARCH=wasm go build -o bundle.wasm bundle.go

run: output.wasm wasm-basic-triangle
	./wasm-basic-triangle

clean:
	rm -f wasm-basic-triangle *.wasm
