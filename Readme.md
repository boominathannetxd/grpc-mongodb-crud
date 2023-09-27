# Basic Crud using MongoDB & gRPC with Buf CLI 

`go mod init grpc-mongodb-crud`

# Install Buf Cli 
`go install github.com/bufbuild/buf/cmd/buf@latest`


# Check buf installed or not 
`buf --version` <!-- check installation --> 

# Initiate Buf 
`buf mod init`  <!--It will create buf.yaml file -->


# Install proto compiler 
`sudo apt install -y protobuf-compiler`
`protoc --version` <!-- check installation --> 

# Install the protocol compiler plugins 
`go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`
`protoc-gen-go --version` <!-- check installation --> 
`protoc-gen-go-grpc --version` <!-- check installation --> 

# Create directory to generate buf files 
`mkdir gen/go` 

# Create proto file in the path api/*.proto

# Generate buf files using crud.proto 
`buf generate` <!-- It will generate two files in gen/go/* -->


































