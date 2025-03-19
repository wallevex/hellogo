@echo off
cd %~dp0%/../

protoc -I . -I ./third_party --go_out=paths=source_relative:. ^
	api/chat/chat.proto

protoc -I . -I ./third_party --go-grpc_out . --go-grpc_opt paths=source_relative ^
	api/chat/chat.proto

protoc -I . -I ./third_party --grpc-gateway_out . ^
	--grpc-gateway_opt logtostderr=true ^
	--grpc-gateway_opt paths=source_relative ^
	api/chat/chat.proto

protoc -I . -I ./third_party --validate_out=paths=source_relative,lang=go:. ^
	api/chat/chat.proto

protoc -I . -I ./third_party --openapiv2_out . ^
	--openapiv2_opt logtostderr=true ^
	--openapiv2_opt json_names_for_fields=true ^
	--openapiv2_opt openapi_naming_strategy=simple ^
	api/chat/chat.proto

cd %~dp0%