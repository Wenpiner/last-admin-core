.PHONY: gen-rpc
gen-rpc: # 生成 RPC 的代码
	goctl rpc protoc ./rpc/core.proto --style=go_zero --go_out=./rpc/types --go-grpc_out=./rpc/types --zrpc_out=./rpc -m
	@echo "Generate RPC files successfully"

.PHONY: gen-api
gen-api: # 生成 API 的代码
	goctl api go -api ./api/desc/all.api -dir ./api -style=go_zero  --home=../.goctl-template
	@echo "Generate API files successfully"

.PHONY: swagger
swagger: # 生成 Swagger 文档
	goctl api swagger --api ./api/desc/all.api --dir ./api --filename ./swagger
	@echo "Generate Swagger files successfully"

.PHONY: swagger-serve
swagger-serve: # 启动 Swagger 服务
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	lsof -i:36666 | awk 'NR!=1 {print $2}' | xargs killall -9 || true
	swagger serve -F=swagger --port 36666 ./api/swagger.json

.PHONY: ent-generate
ent-generate: # 生成 Ent 的代码
	cd rpc/ent && go run entgo.io/ent/cmd/ent generate --template glob="./template/*.tmpl" ./schema
	@echo "Generate Ent files successfully"


# 同步Ent到数据库
.PHONY: ent-migrate
ent-migrate: # 同步 Ent 到数据库
	@go run ./scripts/migrate/main.go