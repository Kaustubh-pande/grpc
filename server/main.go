package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/PandeKaustubhS/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Drug) (*pb.Drug, error)
	GetAll() []*pb.Drug
}

type Repository struct {
	mu   sync.RWMutex
	Drug []*pb.Drug
}

func (repo *Repository) Create(Drug *pb.Drug) (*pb.Drug, error) {
	repo.mu.Lock()
	updated := append(repo.Drug, Drug)
	repo.Drug = updated
	repo.mu.Unlock()
	return Drug, nil
}

type service struct {
	repo repository
}

func (s *service) CreateDrug(ctx context.Context, req *pb.Drug) (*pb.Response, error) {

	// Save our consignment
	_, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	return &pb.Response{Created: true}, nil
}
func (repo *Repository) GetAll() []*pb.Drug {
	return repo.Drug
}
func (s *service) GetDrug(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	drug := s.repo.GetAll()
	return &pb.GetResponse{Drug: drug}, nil
}
func main() {

	repo := &Repository{}

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterDrugServiceServer(s, &service{repo})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// func allUsers(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "All Users Endpoint Hit")
// }

// func newUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "New User Endpoint Hit")
// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Delete User Endpoint Hit")
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Update User Endpoint Hit")
// }
// func handleRequests() {
// 	myRouter := mux.NewRouter() //.StrictSlash(true)
// 	myRouter.HandleFunc("/users", allUsers).Methods("GET")
// 	myRouter.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
// 	myRouter.HandleFunc("/user/{name}/{email}", updateUser).Methods("PUT")
// 	myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")
// 	log.Fatal(http.ListenAndServe(":8081", myRouter))
// }

// func main() {
// 	fmt.Println("in main ")

// 	// Handle Subsequent requests
// 	handleRequests()
// }
