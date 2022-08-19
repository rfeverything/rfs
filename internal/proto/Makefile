all: genproto

genproto:
	protoc meta_server.proto --go_out=./meta_server --go-grpc_out=./meta_server --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc rfs.proto --go_out=./rfs --go-grpc_out=./rfs --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc volume_server.proto --go_out=./volume_server --go-grpc_out=./volume_server --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative