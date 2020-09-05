# GameTest
Golang RPG game - For the WEB (using WASM)

## Setup
1. Clone this repo
2. Install dependencies: `go get ./...`
3. Install image process executable: `go install github.com/hajimehoshi/file2byteslice`
4. Install Go command line tool to execute Go code (a webserver in this case): `go install github.com/shurcooL/goexec`

## Run

### Local machine (native)

Will run the game in a window natively on your computer
```
$> make assets && make run
```

### Local web

Build and run the game via a local webserver for you to play in the web. The Go is converted in to WebAssembly (a binary) that can be run in the browser.

```
$> make assets && make buildweb && make runweb
```

Then navigate to `https://localhost:8080` to play the game

