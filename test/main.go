package main

import (
	"github.com/mwangox/stogo"
	"github.com/mwangox/stogo/config"
	"log"
	"time"
)

func main() {
	stooConfig := config.NewStooConfig("localhost:50051", 20*time.Second).
		WithDefaultNamespace("my-app").
		WithDefaultProfile("prod").
		WithUseTls(true).
		WithTls(&config.TLS{
			SkipTlsVerification: true,
			CaCertPath:          "/opt/systems/apps/stogo/proto/ca_cert.pem",
			ServerNameOverride:  "x.test.example.com",
		})

	client := stogo.NewStoreClient(stooConfig)
	data, err := client.Get("my-app", "prod", "database.username")
	if err != nil {
		log.Fatalf("Error reading key from server %v", err)
	}
	log.Printf("Result: %v, err = %v", data, err)

	res, err := client.Set("my-app", "prod", "database.username", "root")
	if err != nil {
		log.Fatalf("Error in setting value %v", err)
	}
	log.Printf("Set result: %v, err = %v", res, err)

	all, err := client.GetAll("my-app", "prod")
	if err != nil {
		log.Fatalf("Error reading all keys from server %v", err)
	}
	log.Printf("all keys : %v, err = %v", all, err)

	delRes, err := client.Delete("my-app", "prod", "database.password")
	if err != nil {
		log.Fatalf("Error deleting a key %v", err)
	}
	log.Printf("delete result: %v, err = %v", delRes, err)

}
