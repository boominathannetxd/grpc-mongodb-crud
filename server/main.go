package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	crud "grpc-mongodb-crud/gen/go"
)

const (
	mongoURI       = "mongodb://localhost:27017"
	databaseName   = "mydb"
	collectionName = "crud"
)

type server struct {
	crud.UnimplementedCrudServiceServer
}

func (s *server) Create(ctx context.Context, req *crud.CreateRequest) (*crud.CreateResponse, error) {
	client, err := connectMongo()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(databaseName).Collection(collectionName)
	res, err := collection.InsertOne(ctx, bson.M{"name": req.Name, "age": req.Age})
	if err != nil {
		return nil, err
	}
	objectID, _ := res.InsertedID.(primitive.ObjectID)
	return &crud.CreateResponse{Id: objectID.Hex()}, nil

}

func (s *server) Read(ctx context.Context, req *crud.ReadRequest) (*crud.ReadResponse, error) {
	client, err := connectMongo()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(databaseName).Collection(collectionName)
	var result bson.M
	if err := collection.FindOne(ctx, bson.M{"_id": req.Id}).Decode(&result); err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}

	name, _ := result["name"].(string)
	age, _ := result["age"].(string)

	return &crud.ReadResponse{Name: name, Age: age}, nil
}

func (s *server) Update(ctx context.Context, req *crud.UpdateRequest) (*crud.UpdateResponse, error) {
	client, err := connectMongo()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(databaseName).Collection(collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": req.Id}, bson.M{"$set": bson.M{"name": req.Name, "age": req.Age}})
	if err != nil {
		return nil, err
	}

	return &crud.UpdateResponse{Success: true}, nil
}

func (s *server) Delete(ctx context.Context, req *crud.DeleteRequest) (*crud.DeleteResponse, error) {
	client, err := connectMongo()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(databaseName).Collection(collectionName)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": req.Id})
	if err != nil {
		return nil, err
	}

	return &crud.DeleteResponse{Success: true}, nil
}

func connectMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	crud.RegisterCrudServiceServer(s, &server{})
	reflection.Register(s)
	fmt.Println("Server is listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
