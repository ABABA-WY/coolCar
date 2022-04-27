set PROTO_PATH=./auth/api
set GO_OUT_PATH=./auth/api/gen/v1

protoc --go_out=paths=source_relative:gen/v1 auth.proto
protoc --go-grpc_out=paths=source_relative:gen/v1 auth.proto
protoc --grpc-gateway_out=paths=source_relative,grpc_api_configuration=auth.yaml:gen/v1 auth.proto
set PBTS_BIN_DIR=../../../wx/miniprogram/node_modules/.bin
set PBTS_OUT_DIR=../../../wx/miniprogram/service/proto_gen/auth
%PBTS_BIN_DIR%/pbts -o %PBTS_OUT_DIR%/auth_pb.d.ts %PBTS_OUT_DIR%/auth_pb.js
%PBTS_BIN_DIR%/pbjs -t static -w es6 auth.proto --no--create --no--encode --no--decode --no--verify --no--delimited -o %PBTS_OUT_DIR%/auth_pb.js





