SHELL 		:= /bin/bash
PACKAGE		:= $(shell go list)
APP_NAME	?= $(shell basename $(PACKAGE))
DATE		?= $(shell date +%FT%T%z)
VERSION		?= $(shell git describe --tags --exact-match 2> /dev/null || \
					git symbolic-ref -q --short HEAD 2> /dev/null || \
					cat $(CURDIR)/.version 2> /dev/null || echo v0)

GOPATH 	= $(shell echo $${GOPATH} | awk -F : '{print $$NF}')
BASE 	= $(GOPATH)/src/$(PACKAGE)
PKGS 	= $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./...))
BIN 	= $(GOPATH)/bin

BUILD_FLAGS  = -tags release \
			   -ldflags '-X $(PACKAGE)/cmd.Version=$(VERSION) -X $(PACKAGE)/cmd.BuildDate=$(DATE)' \
			   -o bin/$(APP_NAME)

# APP
DB_TYPE				 =  mysql
DB_MIGRATION_PATH	 =  data/db/migration
DB_SEED_PATH	 	 =  data/db/seed
DB_SETUP_PATH	 	 =  data/db/setup

DB_MIGRATION_HOST	 	 ?= 172.17.0.2
DB_MIGRATION_USERNAME	 ?= root
DB_MIGRATION_PASSWORD	 ?= root
DB_SCHEMA 				 ?= svc-article
DB_BASE_URI				 ?= mysql://$(DB_MIGRATION_USERNAME):$(DB_MIGRATION_PASSWORD)@tcp($(DB_MIGRATION_HOST))
DB_URI				 	 ?= $(DB_BASE_URI)/$(DB_SCHEMA)

GOMIGRATE = $(BIN)/migrate
$(GOMIGRATE): | $(BASE) ; $(info building migrate...)
	go get -u -d github.com/mattes/migrate/cli github.com/go-sql-driver/mysql
	go build -tags '$(DB_TYPE)' -o $(GOPATH)/bin/migrate github.com/mattes/migrate/cli

# Tools
GO      	= go
GOFMT   	= gofmt
TIMEOUT 	= 180

.PHONY: status compile

status:
	@echo "PACKAGE: $(PACKAGE)"
	@echo "APP_NAME: $(APP_NAME)"
	@echo "DOCKER_TAG: 	$(DOCKER_TAG)"
	@echo "DATE: $(DATE)"
	@echo "VERSION: $(VERSION)"
	@echo "GOPATH: $(GOPATH)"
	@echo "BASE: $(BASE)"
	@echo "PKGS: $(PKGS)"
	@echo "GO: $(GO)"
	@echo "GOFMT: $(GOFMT)"
	@echo "TIMEOUT: $(TIMEOUT)"
	@echo "Build_Flags": $(BUILD_FLAGS)

compile:
	$(GO) build $(BUILD_FLAGS) $(TARGET_FLAGS) main.go

.PHONY: docker

docker.build:		## Build docker image
	docker build --build-arg PACKAGE=$(PACKAGE) -t $(PACKAGE) .

docker.run:
	docker run --name rest-article -p 8080:8080 -d rest-article:latest

DB_DIRECTORY = $(BASE)/data/database

docker.database :
	docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7

#######################
# DB related things
#######################

.PHONY: db.*

DB_MIGRATE 	= @migrate -database '$(DB_URI)' -path $(DB_MIGRATION_PATH)
DB_SEED 	= @migrate -database '$(DB_URI)' -path $(DB_SEED_PATH)
DB_SETUP 	= @migrate -database '$(DB_BASE_URI)/mysql' -path $(DB_SETUP_PATH)

db.status:
	@echo "Database Docker Port $(DB_DOCKER_PORT)"
	@echo "Database Migration Port $(DB_MIGRATION_PORT)"
	@echo "Database Base URL $(DB_BASE_URI)"

db.req: $(GOMIGRATE)

# setup
db.setup: db.setup.up ## Alias for db.setup.up
db.setup.up: db.req ## Run setup "up" tasks.
	$(DB_SETUP) up $(n)
db.setup.down: db.req ## Run setup "down" tasks.
	$(DB_SETUP) down $(n)

# seed
db.seed: db.seed.up ## Alias for db.seed.up
db.seed.up: db.req ## Run seed "up" tasks.
	$(DB_SEED) up $(n)
db.seed.down: db.req ## Run seed "down" tasks.
	$(DB_SEED) down $(n)

# migrate
db.migrate: db.migrate.up ## Alias for db.migrate.up
db.migrate.up: db.req ## Run migration "up" tasks.
	$(DB_MIGRATE) up $(n)
db.migrate.down: db.req ## Run migration "down" tasks.
	$(DB_MIGRATE) down $(n)