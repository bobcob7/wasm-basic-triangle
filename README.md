# Hello WASM

This is a simple project demoing how to using the new WebAssembly feature in Go 1.11.
Because the WASM feature still isn't on an official release of Go this project is bound to change with it.
Also you're gunna need to install Go 1.11beta3 to get this working.

# Installing Go from source

This is how I installed Go 1.11beta3 on my machine. Be sure to have Go installed before you do this.
1. Set where you want your GOROOT to be: `/usr/local`
1. Download the source: `git clone -b go1.11beta3 --single-branch https://go.googlesource.com/go`
2. Build Go: `cd src; ./make.bash`

# Building
Make sure that you have Go 1.11 installed.
1. Make sure that your GOROOT is set correctly in the makefile. If you followed the installation instructions above,
then you don't have to mess with this.
2. Build with the makefile: `make`

# Running
To run you can use the makefile with `make run` or you can run it directly with `./helloWasm`

# References:
- https://www.youtube.com/watch?v=Xp53ln1nFuk
- https://blog.owulveryck.info/2018/06/08/some-notes-about-the-upcoming-webassembly-support-in-go.html
- https://brianketelsen.com/web-assembly-and-go-a-look-to-the-future/
