BINARY_NAME=gobreach_server


build:
	go build -o ${BINARY_NAME} ./cmd/server/main.go

run:
	go build -o ${BINARY_NAME} ./cmd/server/main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test ./... -v -short

testCov:
	go test ./... -v -short -coverprofile cover.out && \
	go tool cover -html=cover.out

