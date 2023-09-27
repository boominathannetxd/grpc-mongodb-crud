create:
	protoc --go_out=./gen/go/ --go-grpc_out=./gen/go/ crud.proto
clean:
	rm -rf ./gen/go/*.go
