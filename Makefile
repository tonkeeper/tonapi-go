OGEN := go run github.com/ogen-go/ogen/cmd/ogen
SPEC := api/openapi.yml

.PHONY: patch-openapi generate-client verify-generated test

patch-openapi:
	./scripts/patch-openapi.sh $(SPEC)

generate-client: patch-openapi
	$(OGEN) -clean -config .ogen.yml -package tonapi -target . $(SPEC)

# CI guard, regenerates in place so run it on a clean tree.
verify-generated: generate-client
	git diff --exit-code -- $(SPEC) '*_gen.go'

test:
	go test ./...
