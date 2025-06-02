build_windows:
	GOOS=windows GOARCH=amd64 go build -o build/BallGame.exe cmd/game.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o build/BallGame.out cmd/game.go

run: build
	./build/BallGame.out

clean:
	rm -rf build/*