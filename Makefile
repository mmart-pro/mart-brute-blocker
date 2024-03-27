ENV_NAME=./.env.dev
include ${ENV_NAME}
DIR_API=./bin/api
BIN_API=${DIR_API}/mbb-api
API_CFG=${DIR_API}/config.json
API_CFG_TMPL=./configs/api-config.json.tmpl
DIR_CLI=./bin/cli
BIN_CLI=${DIR_CLI}/mbb-cli
CLI_CFG=${DIR_CLI}/config.json
CLI_CFG_TMPL=./configs/cli-config.json.tmpl

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

dirs:
	mkdir -p ${DIR_API}
	mkdir -p ${DIR_CLI}

config: dirs
	env `cat ${ENV_NAME}` envsubst < ${CLI_CFG_TMPL} > ${CLI_CFG};
	env `cat ${ENV_NAME}` envsubst < ${API_CFG_TMPL} > ${API_CFG};

generate:
#	go generate ./...
	rm -Rf internal/grpc/pb
	mkdir -p internal/grpc/pb
	protoc -I ./api --go_out=. --go-grpc_out=. api/mbb/mbb.proto

build: config generate
	go build -v -o $(BIN_API) -ldflags "$(LDFLAGS)" ./cmd/mbb-api
	go build -v -o $(BIN_CLI) -ldflags "$(LDFLAGS)" ./cmd/mbb-cli

version: build
	$(BIN_API) --version
	$(BIN_CLI) --version

migrate:
#	psql -U $(STORAGE_USER) -tc "SELECT 1 FROM pg_database WHERE datname = '$(STORAGE_DB_NAME)'" | grep -q 1 || psql -U $(STORAGE_USER) -c "CREATE DATABASE $(STORAGE_DB_NAME);"
	goose -dir migrations postgres "host=$(STORAGE_HOST) port=$(STORAGE_PORT) user=$(STORAGE_USER) password=$(STORAGE_PASSWORD) dbname=$(STORAGE_DB_NAME) sslmode=disable" up

run: build migrate
	$(BIN_API) --config $(API_CFG)

# build-img:
# 	docker build \
# 		--build-arg=LDFLAGS="$(LDFLAGS)" \
# 		-t $(DOCKER_IMG) \
# 		-f build/Dockerfile .

# run-img: build-img
# 	docker run $(DOCKER_IMG)

test:
	go test -race -count 100 ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.56.2
#	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...


# up: config
# 	env `cat ./.env.${ENV_NAME}` docker-compose -f ./docker/docker-compose.yaml up -d --build

# down: cfg-clean
# 	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml down

# restart: down up


# cfg-clean:
# 	rm -f ${API_CFG}

# integration-test: config
# 	set -e ;\
# 	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml up -d --build;\
# 	test_status_code=0 ;\
# 	docker build -t tests --network host -f tests/Dockerfile . || test_status_code=$$? ;\
# 	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml down --rmi local --volumes --remove-orphans; \
#     rm -f ${API_CFG}; \
# 	exit $$test_status_code ;

# integration-test-cleanup: cfg-clean
# 	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml down --rmi local --volumes --remove-orphans;


.PHONY: dirs config build run version migrate generate test lint #build-img run-img  