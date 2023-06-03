all: wasm-basic-triangle

wasm-basic-triangle: bundle.wasm main.go
	go build -o wasm-basic-triangle main.go

bundle.wasm: bundle.go
	GOOS=js GOARCH=wasm go build -o bundle.wasm bundle.go

run: bundle.wasm wasm-basic-triangle
	./wasm-basic-triangle

clean:
	rm -f wasm-basic-triangle *.wasm
