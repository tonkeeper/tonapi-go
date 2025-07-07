.PHONY: generate-client


generate-client:
	ogen -clean -config .ogen.yml -package tonapi -target . api/openapi.yml

