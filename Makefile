CIFILE=bash build/ci/ci.sh
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

env_up:
	@$(CIFILE) envup

env_down:
	@$(CIFILE) envdown

test:
	@$(CIFILE) onetest $(RUN_ARGS)

tests:
	@$(CIFILE) alltests

runapp:
	@$(CIFILE) runapp

ci: env_down env_up tests env_down

localapp:
	go run cmd/arwos-server/main.go

codegen:
	go generate -v ./...