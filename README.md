## DUPLICATE FILES CHECK

Recursively walks through a directory, and it's sub-directories finding duplicate files.

Outputs duplicate file paths to a text file


## Running

`$ go run main.go`

Provide a path to traverse.

> Enter a file path: /home/wiredmartian/downloads


## High Performance

If you have a huge folder and you require more perfomance

`$ git checkout routine`

`$ go run main.go`

This branch takes advantage of <b>go concurrency</b>


## Packaging

1. `$ go build -o ./bin/app`

Builds the binaries into a `bin` folder and names executable as `app`. Folder and executable could be named anything

2. `$ ./bin/app`

Will run the executable and outputs the results into a newly created folder name `outputs`

