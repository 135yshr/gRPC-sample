.PHONY: all test build

GO			= GO111MODULE=on go
GO_RUN		= $(GO) run
GO_BUILD	= $(GO) build
GO_TEST		= $(GO) test
GO_CLEAN	= $(GO) clean
GO_TOOL		= $(GO) tool
GO_VET		= $(GO) vet
GO_FMT		= $(GO) fmt
GOLINT		= golint

DOCKER		= docker
DOCKER_COMPOSE	= $(DOCKER) compose
DC_BUILD	= $(DOCKER_COMPOSE) build
DC_UP		= $(DOCKER_COMPOSE) up
DC_DOWN		= $(DOCKER_COMPOSE) down
DC_STOP		= $(DOCKER_COMPOSE) stop
DC_PS		= $(DOCKER_COMPOSE) ps
DC_LOGS		= $(DOCKER_COMPOSE) logs

CLIENT_SRC	= cmd/client/main.go
SERVER_SRC	= cmd/server/main.go
CLIENT_BIN	= client
SERVER_BIN	= server
OUT_DIR		= out/
BIN_DIR		= $(OUT_DIR)bin/
COVER_DIR	= $(OUT_DIR)cover/
COVER_FILE	= $(COVER_DIR)cover.out
COVER_HTML	= $(COVER_DIR)cover.html

VERSION?=0.0.0
DOCKER_REGISTRY?=

all: clean test build

## Set up
initialize:
	@cp githooks/* .git/hooks/
	@chmod +x .git/hooks/*

## Build
run:
	$(GO_RUN) $(MAIN_GO)
build: build-client build-server docker-build
build-client:
	$(GO_BUILD) -o $(BIN_DIR)$(CLIENT_BIN) $(CLIENT_SRC)
build-server:
	$(GO_BUILD) -o $(BIN_DIR)$(SERVER_BIN) $(SERVER_SRC)

clean:
	$(GO_CLEAN)
	@rm -rf $(BIN_DIR)$(BIN_NAME)

## Test
test:
	@mkdir -p $(COVER_DIR)
	GO111MODULE=on $(GO_TEST) -cover ./... -coverprofile=$(COVER_FILE)
	GO111MODULE=on $(GO_TOOL) cover -html=$(COVER_FILE) -o $(COVER_HTML)

## Lint
lint:
	$(GOLINT) -set_exit_status ./...
vet:
	$(GO_VET) ./...
fmt:
	$(GO_FMT) ./...

## Docker
docker-build:
	$(DC_BUILD) --no-cache --force-rm
up:
	$(DC_UP) -d
stop:
	$(DC_STOP)
down:
	$(DC_DOWN) --remove-orphans
restart:
	@make down
	@make up
destroy:
	$(DC_DOWN) --rmi all --volumes --remove-orphans
destroy-volumes:
	$(DC_DOWN) --volumes --remove-orphans
ps:
	$(DC_PS)
logs:
	$(DC_LOGS)
logs-watch:
	$(DC_LOGS) --follow
logs-app:
	$(DC_LOGS) app
logs-app-watch:
	$(DC_LOGS) --follow app
