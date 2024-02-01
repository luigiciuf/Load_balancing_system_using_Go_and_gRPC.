package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "Project/protofile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoadBalancerServer struct {
	pb.UnimplementedPrimesServer
	replicaClients []*grpc.ClientConn
	currentIndex   int
	mu             sync.Mutex
}

// Config definisce la struttura per il file di configurazione JSON.
type Config struct {
	LoadBalancerAddress string   `json:"load_balancer_address"`
	ReplicaAddresses    []string `json:"replica_addresses"`
}

// loadConfig carica la configurazione da un file JSON.
func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

// NewLoadBalancerServer crea una nuova istanza del server di load balancing
func NewLoadBalancerServer(replicaAddresses []string) (*LoadBalancerServer, error) {
	var replicaClients []*grpc.ClientConn
	for _, addr := range replicaAddresses {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to replica %s: %v", addr, err)
		}
		replicaClients = append(replicaClients, conn)
	}

	return &LoadBalancerServer{replicaClients: replicaClients}, nil
}

// GetPrimes è l'implementazione del servizio RPC per ottenere numeri primi dal server di bilanciamento del carico.
func (s *LoadBalancerServer) GetPrimes(ctx context.Context, request *pb.RequestParams) (*pb.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	replicaClient := s.replicaClients[s.currentIndex]
	s.currentIndex = (s.currentIndex + 1) % len(s.replicaClients)

	client := pb.NewPrimesClient(replicaClient)
	replicaResponse, err := client.GetPrimes(ctx, request)
	if err != nil {
		log.Printf("Error while calling replica: %v", err)
		return nil, err
	}
	//print("########is prime", replicaResponse.IsPrime)
	return &pb.Response{
		Primes:    replicaResponse.Primes,
		IsPrime:   replicaResponse.IsPrime,
		ReplicaId: int32(s.currentIndex + 1),
	}, nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	loadBalancer, err := NewLoadBalancerServer(config.ReplicaAddresses)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}

	listen, err := net.Listen("tcp", config.LoadBalancerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterPrimesServer(srv, loadBalancer)

	log.Println("Load Balancer is listening on port" + config.LoadBalancerAddress)
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
