# stogo

`stogo` is the Go library for `stookv`. It has abstracted the low level communication details with
`stookv`.

[![Go Reference](https://pkg.go.dev/badge/golang.org)](https://pkg.go.dev/golang.org/x/text)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](MIT-LICENSE)

## Install

```bash
go get github.com/mwangox/stogo
```

## Usage

```go
package main

import (
	"github.com/mwangox/stogo"
	"github.com/mwangox/stogo/config"
	"log"
	"time"
)

func main() {
	// Create stoo configurations.
	// Extended configurations:
	//stooConfig := config.NewStooConfig("localhost:50051", 20*time.Second).
	//	WithDefaultNamespace("my-app").
	//	WithDefaultProfile("prod").
	//	WithUseTls(true).
	//	WithTls(&config.TLS{
	//		SkipTlsVerification: false,
	//		CaCertPath:          "/stookv/conf/ca_cert.pem",
	//		ServerNameOverride:  "stookv.hostname.com",
	//	})
	// Minimal configurations:
	// Or stooConfig := config.NewDefaultStooConfig()
	stooConfig := config.NewStooConfig("localhost:50051", 20*time.Second)
	
	// Create stoo client.
	client := stogo.NewStoreClient(stooConfig)

	// Set value to a key.
	res, err := client.Set("my-app", "prod", "database.username", "lauryn.hill")
	if err != nil {
		log.Fatalf("Error in setting value %v", err)
	}
	log.Printf("Set result: %v", res)
	
	// Get value from a key.
	data, err := client.Get("my-app", "prod", "database.username")
	if err != nil {
		log.Fatalf("Error reading key from server %v", err)
	}
	log.Printf("Result: %v", data)
	
	// Get all key value pairs from a given namespace and profile
	all, err := client.GetAllByNamespaceAndProfile("my-app", "prod")
	if err != nil {
		log.Fatalf("Error reading all keys from server %v", err)
	}
	log.Printf("all keys values : %v", all)

	// Delete a key
	result, err := client.Delete("my-app", "prod", "database.password")
	if err != nil {
		log.Fatalf("Error deleting a key %v", err)
	}
	log.Printf("delete result: %v", result)
}
```

## License

The project is licensed under [MIT license](./MIT-LICENSE).

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in `stogo` by you, shall be licensed as MIT, without any additional
terms or conditions.