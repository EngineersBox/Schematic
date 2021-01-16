BIN_NAME="main"
OUT_DIR="out/$(BIN_NAME)"
SRC="main.go"

build:
	@go build -o $(OUT_DIR) $(SRC)

exec:
	@./$(OUT_DIR)

run:
	@make build
	@make exec