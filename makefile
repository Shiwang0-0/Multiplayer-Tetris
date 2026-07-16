build:
	@go build -o multiplayertetris ./cmd/tetris

run: build
	@./multiplayertetris