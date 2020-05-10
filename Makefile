CIFILE=bash build/ci/ci.sh
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

test:
	@$(CIFILE) onetest $(RUN_ARGS)

tests:
	@$(CIFILE) alltests

app_run:
	@$(CIFILE) runapp
