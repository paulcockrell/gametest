assets:
	file2byteslice -input=./resources/images/bullet.png -output=./resources/images/bullet.go -package=images -var=Bullet_png
	file2byteslice -input=./resources/images/runner-left.png -output=./resources/images/runner_left.go -package=images -var=RunnerLeft_png
	file2byteslice -input=./resources/images/runner-right.png -output=./resources/images/runner_right.go -package=images -var=RunnerRight_png
	file2byteslice -input=./resources/images/tiles.png -output=./resources/images/tiles.go -package=images -var=Tiles_png
run:
	go run main.go runner.go level.go