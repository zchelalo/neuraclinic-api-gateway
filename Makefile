ifneq ("$(wildcard .env)", "")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

DOCKER_COMPOSE_FILE = ./.docker/compose.yml
DOCKER_NETWORK = neuraclinic-network
LOCAL_PROTO_CONTRACTS = ../neuraclinic-proto-contracts

setup:
	$(MAKE) create-envs
	$(MAKE) create-network
	$(MAKE) compose-build-detached

create-envs:
	test -f .env || cp .env.example .env

create-network:
	docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || docker network create $(DOCKER_NETWORK)

proto:
ifneq ("$(wildcard $(LOCAL_PROTO_CONTRACTS)/buf.yaml)", "")
	cd $(LOCAL_PROTO_CONTRACTS) && buf generate \
		--template ../neuraclinic-api-gateway/buf.gen.yaml \
		--output ../neuraclinic-api-gateway \
		--path proto/auth/v1/auth.proto \
		--path proto/user/v1/user.proto \
		--path proto/record/v1/patient.proto \
		--path proto/record/v1/appointment.proto \
		--path proto/record/v1/note.proto \
		--path proto/record/v1/attachment.proto \
		--path proto/record/v1/familiogram.proto \
		--path proto/location/v1/location.proto \
		--path proto/file_management/v1/file_management.proto \
		--path proto/shared/v1/shared.proto
else
	buf generate buf.build/zchelalo-labs/neuraclinic-proto-contracts \
		--path auth/v1/auth.proto \
		--path user/v1/user.proto \
		--path record/v1/patient.proto \
		--path record/v1/appointment.proto \
		--path record/v1/note.proto \
		--path record/v1/attachment.proto \
		--path record/v1/familiogram.proto \
		--path location/v1/location.proto \
		--path file_management/v1/file_management.proto \
		--path shared/v1/shared.proto
endif

compose:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up

compose-detached:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

compose-build:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

compose-build-detached:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

compose-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out

build:
	mkdir -p dist
	go build -buildvcs=false -trimpath -o dist/neuraclinic-api-gateway ./cmd

.PHONY: setup create-envs create-network proto compose compose-detached compose-build compose-build-detached compose-down fmt lint test coverage build
