package stogo

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/stogo/config"
	"github.com/stogo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type StooClient struct {
	Config *config.StooConfig
	client proto.KVServiceClient
}

var ErrDefaultNamespaceAndProfileMustBeDefined = errors.New("default namespace and profile must be set to use this method")

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
	return &StooClient{Config: cfg,
		client: client,
	}
}

// Get value stored using namespace, profile and key
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

// Set a key to a namespace and profile
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

// Delete a key from a given namespace and profile
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

// GetAll get all keys from a given namespace and profile
func (c *StooClient) GetAll(namespace, profile string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GetReadTimeout())
	defer cancel()
	res, err := c.client.GetServiceByNamespaceAndProfile(ctx, &proto.GetByNamespaceAndProfileRequest{
		Namespace: namespace,
		Profile:   profile,
	})
	return res.GetData(), err
}

// GetDefault get a value for a key in a given default namespace and profile
func (c *StooClient) GetDefault(key string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()

	if defaultNamespace != "" && defaultProfile != "" {
		return c.Get(defaultNamespace, defaultProfile, key)
	}
	return "", ErrDefaultNamespaceAndProfileMustBeDefined
}

// SetDefault set value for a key in a given default namespace and profile
func (c *StooClient) SetDefault(key, value string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()

	if defaultNamespace != "" && defaultProfile != "" {
		return c.Set(defaultNamespace, defaultProfile, key, value)
	}
	return "", ErrDefaultNamespaceAndProfileMustBeDefined
}

// DeleteDefault delete a key from a given default namespace and profile
func (c *StooClient) DeleteDefault(key string) (string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()

	if defaultNamespace != "" && defaultProfile != "" {
		return c.Delete(defaultNamespace, defaultProfile, key)
	}
	return "", ErrDefaultNamespaceAndProfileMustBeDefined

}

// GetAllDefault all key value pairs from given default namespace and profile
func (c *StooClient) GetAllDefault() (map[string]string, error) {
	defaultNamespace := c.Config.GetDefaultNamespace()
	defaultProfile := c.Config.GetDefaultProfile()

	if defaultNamespace != "" && defaultProfile != "" {
		return c.GetAll(defaultNamespace, defaultProfile)
	}
	return nil, ErrDefaultNamespaceAndProfileMustBeDefined
}
