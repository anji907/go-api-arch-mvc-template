IMAGE_TAG = web

generate-code-from-openapi: ## Generate code from openapi
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	oapi-codegen --config=./api/config.yaml ./api/openapi.yaml

external-up: ## Up external containers
	docker-compose up -d mysql swagger-ui

external-down: ## Down external containers
	docker-compose down

mysql-cli: ## Connect to mysql cli
	docker-compose run mysql-cli

run: ## Run app
	export APP_ENV=development
	go run main.go

docker-build: ## Build image
	docker build --tag $(IMAGE_TAG) .

docker-run: ## Run docker
	docker run -p 8080:8080 -i -t $(IMAGE_TAG)

docker-compose-up: docker-build ## Run docker compose up
	docker-compose up -d --wait mysql web swagger-ui

docker-compose-down: ## Run docker compose down
	docker-compose down

unittest: ## Run unittest
	go clean -testcache
	go test -v `go list ./... | grep -v /integration | grep -v /testutils | grep -v /app`
	go test -v -p 1 ./app/...

test-cover: ## Run test cover
	go test -coverprofile=coverage.out `go list ./... | grep -v /integration` && go tool cover -html=coverage.out

integration-test: generate-code-from-openapi docker-build ## Run integration test
	export APP_ENV=integration
	-docker-compose up -d --wait
	-go clean -testcache
	-go test -v `go list ./... | grep /integration`
	+docker-compose down

lint: ## Run lint
	# Install with `brew install golangci-lint` on Mac
	golangci-lint run

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
