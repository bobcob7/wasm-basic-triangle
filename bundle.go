package main

import (
	"fmt"
	"reflect"
	"runtime"
	"syscall/js"
	"unsafe"
)

// Shamelessly stolen from: https://github.com/golang/go/issues/32402
// Thanks hajimehoshi
func sliceToByteSlice(s interface{}) []byte {
	switch s := s.(type) {
	case []int8:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		return *(*[]byte)(unsafe.Pointer(h))
	case []int16:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 2
		h.Cap *= 2
		return *(*[]byte)(unsafe.Pointer(h))
	case []int32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []int64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint8:
		return s
	case []uint16:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 2
		h.Cap *= 2
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	case []float32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []float64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at sliceToBytesSlice: %T", s))
	}
}
func SliceToTypedArray(s interface{}) js.Value {
	switch s := s.(type) {
	case []int8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int8Array").New(buf, a.Get("byteOffset"), a.Get("byteLength"))
	case []int16:
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int16Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []int32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []uint8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, s)
		runtime.KeepAlive(s)
		return a
	case []uint16:
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Uint16Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []uint32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Uint32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Float32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float64:
		a := js.Global().Get("Uint8Array").New(len(s) * 8)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Float64Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/8)
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at SliceToTypedArray: %T", s))
	}
}

var (
	width   int
	height  int
	gl      js.Value
	glTypes GLTypes
)

type GLTypes struct {
	staticDraw         js.Value
	arrayBuffer        js.Value
	elementArrayBuffer js.Value
	vertexShader       js.Value
	fragmentShader     js.Value
	float              js.Value
	depthTest          js.Value
	colorBufferBit     js.Value
	triangles          js.Value
	unsignedShort      js.Value
}

func (types *GLTypes) New() {
	types.staticDraw = gl.Get("STATIC_DRAW")
	types.arrayBuffer = gl.Get("ARRAY_BUFFER")
	types.elementArrayBuffer = gl.Get("ELEMENT_ARRAY_BUFFER")
	types.vertexShader = gl.Get("VERTEX_SHADER")
	types.fragmentShader = gl.Get("FRAGMENT_SHADER")
	types.float = gl.Get("FLOAT")
	types.depthTest = gl.Get("DEPTH_TEST")
	types.colorBufferBit = gl.Get("COLOR_BUFFER_BIT")
	types.triangles = gl.Get("TRIANGLES")
	types.unsignedShort = gl.Get("UNSIGNED_SHORT")
}

func main() {

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "gocanvas")
	width = doc.Get("body").Get("clientWidth").Int()
	height = doc.Get("body").Get("clientHeight").Int()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)

	gl = canvasEl.Call("getContext", "webgl")
	if gl.IsUndefined() {
		gl = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if gl.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	glTypes.New()

	//// VERTEX BUFFER ////
	var verticesNative = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	var vertices = SliceToTypedArray(verticesNative)
	// Create buffer
	vertexBuffer := gl.Call("createBuffer", glTypes.arrayBuffer)
	// Bind to buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, vertexBuffer)

	// Pass data to buffer
	gl.Call("bufferData", glTypes.arrayBuffer, vertices, glTypes.staticDraw)

	// Unbind buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, nil)

	//// INDEX BUFFER ////
	var indicesNative = []uint32{
		2, 1, 0,
	}
	var indices = SliceToTypedArray(indicesNative)

	// Create buffer
	indexBuffer := gl.Call("createBuffer", glTypes.elementArrayBuffer)

	// Bind to buffer
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, indexBuffer)

	// Pass data to buffer
	gl.Call("bufferData", glTypes.elementArrayBuffer, indices, glTypes.staticDraw)

	// Unbind buffer
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, nil)

	//// Shaders ////

	// Vertex shader source code
	vertCode := `
	attribute vec3 coordinates;
		
	void main(void) {
		gl_Position = vec4(coordinates, 1.0);
	}`

	// Create a vertex shader object
	vertShader := gl.Call("createShader", glTypes.vertexShader)

	// Attach vertex shader source code
	gl.Call("shaderSource", vertShader, vertCode)

	// Compile the vertex shader
	gl.Call("compileShader", vertShader)

	//fragment shader source code
	fragCode := `
	void main(void) {
		gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
	}`

	// Create fragment shader object
	fragShader := gl.Call("createShader", glTypes.fragmentShader)

	// Attach fragment shader source code
	gl.Call("shaderSource", fragShader, fragCode)

	// Compile the fragmentt shader
	gl.Call("compileShader", fragShader)

	// Create a shader program object to store
	// the combined shader program
	shaderProgram := gl.Call("createProgram")

	// Attach a vertex shader
	gl.Call("attachShader", shaderProgram, vertShader)

	// Attach a fragment shader
	gl.Call("attachShader", shaderProgram, fragShader)

	// Link both the programs
	gl.Call("linkProgram", shaderProgram)

	// Use the combined shader program object
	gl.Call("useProgram", shaderProgram)

	//// Associating shaders to buffer objects ////

	// Bind vertex buffer object
	gl.Call("bindBuffer", glTypes.arrayBuffer, vertexBuffer)

	// Bind index buffer object
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, indexBuffer)

	// Get the attribute location
	coord := gl.Call("getAttribLocation", shaderProgram, "coordinates")

	// Point an attribute to the currently bound VBO
	gl.Call("vertexAttribPointer", coord, 3, glTypes.float, false, 0, 0)

	// Enable the attribute
	gl.Call("enableVertexAttribArray", coord)

	//// Drawing the triangle ////

	// Clear the canvas
	gl.Call("clearColor", 0.5, 0.5, 0.5, 0.9)
	gl.Call("clear", glTypes.colorBufferBit)

	// Enable the depth test
	gl.Call("enable", glTypes.depthTest)

	// Set the view port
	gl.Call("viewport", 0, 0, width, height)

	// Draw the triangle
	gl.Call("drawElements", glTypes.triangles, len(indicesNative), glTypes.unsignedShort, 0)
}
