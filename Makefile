BIN_NAME="main"
OUT_DIR="out"
SRC="main.go"

build:
	@mkdir -p out
	@go build -v -o $(OUT_DIR)/$(BIN_NAME) $(SRC)

exec:
	@./$(OUT_DIR)

run:
	@make build
	@make exec