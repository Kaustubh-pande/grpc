package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"context"

	pb "github.com/PandeKaustubhS/proto"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "drug.json"
)

func parseFile(file string) (*pb.Drug, error) {
	var drug *pb.Drug
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &drug)
	return drug, err
}
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewDrugServiceClient(conn)
	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	drug, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateDrug(context.Background(), drug)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetDrug(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Drug {
		log.Println(v)
	}
}
