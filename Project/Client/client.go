package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "Project/protofile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config definisce la struttura per il file di configurazione JSON contenente solo il valore di load_balancer_address
type Config struct {
	LoadBalancerAddress string `json:"load_balancer_address"`
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

func main() {

	// Connessione al server di bilanciamento del carico
	//loadBalancerAddress := "localhost:50051"
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	loadBalancerConn, err := grpc.Dial(config.LoadBalancerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to load balancer: %v", err)
	}

	defer loadBalancerConn.Close()
	// Creazione del client per il server di bilanciamento del carico
	loadBalancerClient := pb.NewPrimesClient(loadBalancerConn)
	// Verifica della connessione chiamando un metodo di servizio del server di bilanciamento del carico
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// Input del numero
	for {
		// Input del numero
		var input int32
		fmt.Print("Inserisci un numero (0 per uscire): ")
		fmt.Scan(&input)

		if input == 0 {
			break // Esci dal ciclo se l'utente inserisce 0
		}
		response, err := loadBalancerClient.GetPrimes(ctx, &pb.RequestParams{N: input})
		if err != nil {
			log.Fatalf("RPC call to load balancer failed: %v", err)
		}

		// Stampa dei risultati
		if response.IsPrime {
			fmt.Printf("Il numero %d è un numero primo. Risultato dalla replica %d.\n", input, response.ReplicaId)
		} else {
			fmt.Printf("Il numero %d non è un numero primo. Risultato dalla replica %d.\n", input, response.ReplicaId)
		}
	}
}
