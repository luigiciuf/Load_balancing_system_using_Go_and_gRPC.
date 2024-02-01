package main

import (
	pb "Project/protofile"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"os"
	"sync"
)

type ReplicaServer struct {
	pb.UnimplementedPrimesServer
}

// Config definisce la struttura per il file di configurazione JSON contenente solo il valore di replica_address
type Config struct {
	ReplicaAddresses []string `json:"replica_addresses"`
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
func (s *ReplicaServer) GetPrimes(ctx context.Context, request *pb.RequestParams) (*pb.Response, error) {
	n := request.N
	primes := s.generatePrimes(n)

	// Imposta response.IsPrime sulla base dei risultati
	isPrime := s.isPrime(n)

	return &pb.Response{Primes: primes, IsPrime: isPrime}, nil
}

func (s *ReplicaServer) isPrime(num int32) bool {
	if num < 2 {
		//fmt.Printf("%d non è primo perché è inferiore a 2\n", num)
		return false
	}

	// Verifico i divisori da 2 fino alla radice quadrata del numero
	maxDivisor := int32(math.Sqrt(float64(num)))
	for i := int32(2); i <= maxDivisor; i++ {
		if num%i == 0 {
			//fmt.Printf("%d non è primo perché è divisibile per %d\n", num, i)
			return false
		}
	}
	//fmt.Printf("%d è primo\n", num)
	return true
}

func (s *ReplicaServer) generatePrimes(n int32) []int32 {
	var primes []int32

	for i := int32(2); i <= n; i++ {
		if s.isPrime(i) {
			primes = append(primes, i)
		}
	}
	//log.Printf("Generated primes: %v", primes)
	return primes
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	var wg sync.WaitGroup
	for _, port := range config.ReplicaAddresses {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			listen, err := net.Listen("tcp", p)
			if err != nil {
				log.Fatalf("Failed to listen on port %s: %v", p, err)
			}

			srv := grpc.NewServer()
			pb.RegisterPrimesServer(srv, &ReplicaServer{})

			log.Printf("Replica is listening on port %s", p)
			if err := srv.Serve(listen); err != nil {
				log.Printf("Failed to serve on port %s: %v", p, err)
			}
		}(port)
	}

	wg.Wait()
}
