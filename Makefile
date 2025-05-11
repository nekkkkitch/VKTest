buildsppb: 
	protoc --proto_path=pkg/grpc/proto/subpubservice --go_out=pkg/grpc/pb/subpubservice --go-grpc_out=pkg/grpc/pb/subpubservice pkg/grpc/proto/subpubservice/*.proto