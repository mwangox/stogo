// Package stogo defines StooClient that carries methods for connecting and interacting
// with StooKV server. The methods wrap around the gRPC operations to StooKV.
package stogo

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/mwangox/stogo/config"
	"github.com/mwangox/stogo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// StooClient holds stoo client and the associated configurations.
type StooClient struct {
	Config *config.StooConfig
	client proto.KVServiceClient
}

// ErrDefaultNamespaceAndProfileMustBeDefined thrown by *default methods when called while default
// namespace and profile are not defined.
var ErrDefaultNamespaceAndProfileMustBeDefined = errors.New("default namespace and profile must be set to use this method")

// NewStoreClient constructs stoo client from given configurations.
//
// Minimum configurations usage example:
//
//	stooConfig := config.NewStooConfig("localhost:50051", 20*time.Second)
//
// Extended configurations usage example:
//
//	  	stooConfig := config.NewStooConfig("localhost:50051", 20*time.Second).
//				WithDefaultNamespace("my-app").
//			    WithDefaultProfile("prod").
//	  		WithUseTls(true).
//	 		WithTls(&config.TLS{
//				SkipTlsVerification: false,
//				CaCertPath:          "/stookv/ca_cert.pem",
//				ServerNameOverride:  "stookv.example.com",
//			})
//
//		client := stogo.NewStoreClient(stooConfig)
func NewStoreClient(cfg *config.StooConfig) *StooClient {
	var options []grpc.DialOption
	if cfg.GetUseTls() {
		if !cfg.GetTls().SkipTlsVerification {
			creds, err := credentials.NewClientTLSFromFile(cfg.GetTls().CaCertPath, cfg.GetTls().ServerNameOverride)
			if err != nil {
				log.Fatalf("Failed to read CA cert: %v", err)
			}
			options = append(options, grpc.WithTransportCredentials(creds))
		} else {
			options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
		}
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(cfg.GetEndpoint(), options...)
	if err != nil {
		log.Fatalf("Failed to establish connection to stooKV: %v", err)
	}

	client := proto.NewKVServiceClient(conn)
	return &StooClient{
		Config: cfg,
		client: client,
	}
}

// Get gets a value stored using namespace, profile and key.
//
//	 Usage example:
//		   data, err := client.Get("my-app", "prod", "database.username")
//		   if err != nil {
//		     log.Fatalf("Error reading key from server %v", err)
//		   }
//		   log.Printf("Result: %v", data)
func (c *StooClient) Get(namespace, profile, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GetReadTimeout())
	defer cancel()

	res, err := c.client.GetService(ctx, &proto.GetRequest{
		Namespace: namespace,
		Profile:   profile,
		Key:       key,
	})
	return res.GetData(), err
}

// Set sets a key to a namespace and profile.
//
// Usage example:
//
//	   res, err := client.Set("my-app", "prod", "database.username", "lauryn.hill")
//		  if err != nil {
//		      log.Fatalf("Error in setting value %v", err)
//		  }
//		  log.Printf("Set result: %v", res)
func (c *StooClient) Set(namespace, profile, key, value string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GetReadTimeout())
	defer cancel()
	res, err := c.client.SetKeyService(ctx, &proto.SetKeyRequest{
		Namespace: namespace,
		Profile:   profile,
		Key:       key,
		Value:     value,
	})
	return res.GetData(), err
}

// SetSecret sets a key to a namespace and profile in an encrypted format.
// Usage example:
//
//	   res, err := client.SetSecret("my-app", "prod", "database.password", "the-scrore@1996")
//		  if err != nil {
//		      log.Fatalf("Error in setting secret value %v", err)
//		  }
//		  log.Printf("SetSecret result: %v", res)
func (c *StooClient) SetSecret(namespace, profile, key, value string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GetReadTimeout())
	defer cancel()
	res, err := c.client.SetSecretKeyService(ctx, &proto.SetKeyRequest{
		Namespace: namespace,
		Profile:   profile,
		Key:       key,
		Value:     value,
	})
	return res.GetData(), err
}

// Delete removes a key from a given namespace and profile
//
// Usage example:
//
//	   res, err := client.Delete("my-app", "prod", "database.password")
//	   if err != nil {
//		    log.Fatalf("Error deleting a key %v", err)
//	   }
//	   log.Printf("delete result: %v", res)
func (c *StooClient) Delete(namespace, profile, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GetReadTimeout())
	defer cancel()
	res, err := c.client.DeleteKeyService(ctx, &proto.DeleteKeyRequest{
		Namespace: namespace,
		Profile:   profile,
		Key:       key,
	})
	return res.GetData(), err
}

// GetAllByNamespaceAndProfile gets all keys from a given namespace and profile.
//
// Usage example:
//
//	  all, err := client.GetAllByNamespaceAndProfile("my-app", "prod")
//	  if err != nil {
//		   log.Fatalf("Error reading all keys from server %v", err)
//	  }
//	  log.Printf("all keys values : %v", all)
func (c *StooClient) GetAllByNamespaceAndProfile(namespace, profile string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GetReadTimeout())
	defer cancel()
	res, err := c.client.GetServiceByNamespaceAndProfile(ctx, &proto.GetByNamespaceAndProfileRequest{
		Namespace: namespace,
		Profile:   profile,
	})
	return res.GetData(), err
}

// GetDefault gets a value for a key in a given default namespace and profile.
func (c *StooClient) GetDefault(key string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()
	if err := validateDefaultNamespaceAndProfile(defaultNamespace, defaultProfile); err != nil {
		return "", err
	}
	return c.Get(defaultNamespace, defaultProfile, key)
}

// SetDefault sets value for a key in a given default namespace and profile.
func (c *StooClient) SetDefault(key, value string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()
	if err := validateDefaultNamespaceAndProfile(defaultNamespace, defaultProfile); err != nil {
		return "", err
	}
	return c.Set(defaultNamespace, defaultProfile, key, value)
}

// SetSecretDefault sets secret value for a key in a given default namespace and profile.
func (c *StooClient) SetSecretDefault(key, value string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()
	if err := validateDefaultNamespaceAndProfile(defaultNamespace, defaultProfile); err != nil {
		return "", err
	}
	return c.SetSecret(defaultNamespace, defaultProfile, key, value)
}

// DeleteDefault removes a key from a given default namespace and profile.
func (c *StooClient) DeleteDefault(key string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()
	if err := validateDefaultNamespaceAndProfile(defaultNamespace, defaultProfile); err != nil {
		return "", err
	}
	return c.Delete(defaultNamespace, defaultProfile, key)
}

// GetAllByDefaultNamespaceAndProfile gets all key value pairs from a given default namespace and profile.
func (c *StooClient) GetAllByDefaultNamespaceAndProfile() (map[string]string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()
	if err := validateDefaultNamespaceAndProfile(defaultNamespace, defaultProfile); err != nil {
		return nil, err
	}
	return c.GetAllByNamespaceAndProfile(defaultNamespace, defaultProfile)

}

// validateDefaultNamespaceAndProfile checks if all defaultNamespace and defaultProfile are being set.
func validateDefaultNamespaceAndProfile(defaultNamespace, defaultProfile string) error {
	if defaultNamespace != "" && defaultProfile != "" {
		return nil
	}
	return ErrDefaultNamespaceAndProfileMustBeDefined
}
