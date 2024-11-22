.PHONY: gen-docs
gen-docs:
	@~/go/bin/swag init -g ./api/main.go -d cmd,internal --parseDependency && ~/go/bin/swag fmt