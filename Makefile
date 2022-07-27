export SHELL:=/bin/bash
export BASE_NAME:=$(shell basename ${PWD})
export IMAGE_BASE_NAME:=kz/$(shell basename ${PWD})
export NETWORK:=${BASE_NAME}-network

default: help

help: ## Prints help for targets with comments
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-50s\033[0m %s\n", $$1, $$2}'
	@echo ""

#######
# Run #
#######

BASE:=\
	-f docker-compose/server.yml \
	-f docker-compose/client.yml

compose:
	@docker-compose ${COMPOSE} \
		-p ${BASE_NAME} \
		up --build --remove-orphans --force-recreate # --abort-on-container-exit

up-base: ## Start the example
	@COMPOSE="${BASE}" make compose

#######
# Dev #
#######

codegen: ## Generate code
	@mkdir -p \
		internal/model
	@protoc --go_out=internal/model --go-grpc_out=internal/model proto/*.proto

###############
# Danger Zone #
###############

reset: ## Cleanup
	@docker stop $(shell docker ps -aq) || true
	@docker system prune || true
	@docker volume rm $(shell docker volume ls -q) || true
	@docker rmi -f ${IMAGE_BASE_NAME}-go:latest || true
