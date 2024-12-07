.PHONY: gen-docs
gen-docs:
	@~/go/bin/swag init -g ./api/main.go -d cmd,internal --parseDependency && ~/go/bin/swag fmt

.PHONY: gen-docs-wins
gen-docs-wins:
	@swag init -g ./api/main.go -d cmd,internal --parseDependency && swag fmt