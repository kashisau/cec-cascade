BINARY_NAME=cec-cascade

build:
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}.arm64 main.go