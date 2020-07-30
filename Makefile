lint:
	command -v golangci-lint >/dev/null 2>&1 || { go install github.com/golangci/golangci-lint/cmd/golangci-lint; } && \
    golangci-lint --config=./.golangci.yaml run ./...

build:
	go build -o .bin/image-previewer ./cmd/main.go

test:
	go test -race -count 100 ./pkg/...

integration-tests:
	docker-compose -f ./docker-compose-tests.yaml up --build --abort-on-container-exit --exit-code-from integration-tests && \
	docker-compose -f ./docker-compose-tests.yaml down

run:
	docker-compose -f ./docker-compose.yaml up -d

down:
	docker-compose -f ./docker-compose.yaml down