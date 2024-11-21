.PHONY: gen-docs
gen-docs:
	@~/go/bin/swag init -g ./api/main.go -d cmd,internal && ~/go/bin/swag fmt