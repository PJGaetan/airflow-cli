VERSION=2.5.3


deps:
	go mod vendor -v

lint:
	golangci-lint run


test: build-with-coverage 
	@rm -fr .coverdata
	@mkdir -p .coverdata
	@go test -v ./...
	@go tool covdata percent -i=.coverdata
	# $(MAKE) airflow.server-down > /dev/null

# airflow.server-down:
# 	@docker compose down
# 	@rm docker-compose.yaml
#
#
# airflow.server-up:
# 	@curl --verbose https://airflow.apache.org/docs/apache-airflow/$(VERSION)/docker-compose.yaml > docker-compose.yaml
# 	@docker compose up --force-recreate -d --wait

check-coverage: test 
	@go tool covdata textfmt -i=.coverdata -o profile.txt
	@go tool cover -html=profile.txt

build:
	@go build -o airflow-cli airflow-cli/main.go

build-with-coverage:
	@go build -cover -o airflow-cli-coverage airflow-cli/main.go

.DEFAULT_GOAL := build
