# Image byteslice generation

## Prerequisites

Install Go package [file2byteslice](https://github.com/hajimehoshi/file2byteslice)

## Generate

Example:
```
file2byteslice -input=runner-left.png -output=runner_left.go -package=images -var=RunnerLeft_png
```