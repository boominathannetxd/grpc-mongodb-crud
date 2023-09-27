package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"

	crud "grpc-mongodb-crud/gen/go"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	c := crud.NewCrudServiceClient(conn)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: client <command> [options]")
		return
	}

	cmd := args[1]
	switch cmd {
	case "create":
		if len(args) != 4 {
			fmt.Println("Usage: client create <name> <age>")
			return
		}
		name := args[2]
		age := args[3]
		createResponse, err := c.Create(context.Background(), &crud.CreateRequest{Name: name, Age: age})
		if err != nil {
			log.Fatalf("Create failed: %v", err)
		}
		fmt.Printf("Created with ID: %s\n", createResponse.Id)
	case "read":
		if len(args) != 3 {
			fmt.Println("Usage: client read <id>")
			return
		}
		id := args[2]
		readResponse, err := c.Read(context.Background(), &crud.ReadRequest{Id: id})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				fmt.Println("Not found")
				return
			}
			log.Fatalf("Read failed: %v", err)
		}
		fmt.Printf("Name: %s, Age: %s\n", readResponse.Name, readResponse.Age)
	case "update":
		if len(args) != 5 {
			fmt.Println("Usage: client update <id> <name> <age>")
			return
		}
		id := args[2]
		name := args[3]
		age := args[4]
		_, err := c.Update(context.Background(), &crud.UpdateRequest{Id: id, Name: name, Age: age})
		if err != nil {
			log.Fatalf("Update failed: %v", err)
		}
		fmt.Println("Update successful")
	case "delete":
		if len(args) != 3 {
			fmt.Println("Usage: client delete <id>")
			return
		}
		id := args[2]
		_, err := c.Delete(context.Background(), &crud.DeleteRequest{Id: id})
		if err != nil {
			log.Fatalf("Delete failed: %v", err)
		}
		fmt.Println("Delete successful")
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
}
