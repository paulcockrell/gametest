# If you change or add assets you can add entries here to convert the images to byteslices
assets:
	file2byteslice -input=./resources/images/vaxerman.png -output=./resources/images/vaxerman.go -package=images -var=VaxerMan_png
	file2byteslice -input=./resources/images/bullet.png -output=./resources/images/bullet.go -package=images -var=Bullet_png
	file2byteslice -input=./resources/images/runner-left.png -output=./resources/images/runner_left.go -package=images -var=RunnerLeft_png
	file2byteslice -input=./resources/images/runner-right.png -output=./resources/images/runner_right.go -package=images -var=RunnerRight_png
	file2byteslice -input=./resources/images/tiles.png -output=./resources/images/tiles.go -package=images -var=Tiles_png

# Run locally for development
run:
	go run main.go vaxerman.go level.go

# Build WASM for web browser
buildweb:
	GOOS=js GOARCH=wasm go build -o ./build/web/gametest.wasm github.com/paulcockrell/gametest

# Run web version locally
runweb:
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./build/web`)))'