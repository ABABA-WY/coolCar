
protoc --go_out=paths=source_relative:gen/v1 rental.proto
protoc --go-grpc_out=paths=source_relative:gen/v1 rental.proto
protoc --grpc-gateway_out=paths=source_relative,grpc_api_configuration=rental.yaml:gen/v1 rental.proto
set PBTS_BIN_DIR=../../../wx/miniprogram/node_modules/.bin
set PBTS_OUT_DIR=../../../wx/miniprogram/service/proto_gen/rental
%PBTS_BIN_DIR%/pbts -o %PBTS_OUT_DIR%/rental_pb.d.ts %PBTS_OUT_DIR%/rental_pb.js.
%PBTS_BIN_DIR%/pbjs -t static -w es6 rental.proto --no--create --no--encode --no--decode --no--verify --no--delimited --force-number -o %PBTS_OUT_DIR%/rental_pb.js





















