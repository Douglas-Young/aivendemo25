build:
    go build -o bin/aivendemo ./src/cmd/main.go

run:
    go run ./src/cmd/main.go

test:
    go test ./...

clean:
    rm -rf bin/