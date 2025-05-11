start:
	docker-compose build --no-cache \
	&& docker-compose up -d
buildsppb: 
	protoc --proto_path=pkg/grpc/proto/subpubservice --go_out=pkg/grpc/pb/subpubservice --go-grpc_out=pkg/grpc/pb/subpubservice pkg/grpc/proto/subpubservice/*.proto
