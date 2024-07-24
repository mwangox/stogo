// Package config defines configurations required in establishing connection
// to stookv instance.
package config

import (
	"log"
	"time"
)

// StooConfig holds data to be used during interactions with StooKV using StooClient.
type StooConfig struct {
	// endpoint grpc endpoint, http or https.
	endpoint string
	// useTls flag that tells if StooKV has enabled https on not.
	useTls bool
	// readTimeout max duration of time for a client to wait for a response.
	readTimeout time.Duration
	// defaultNamespace default namespace to be used by *default methods.
	defaultNamespace string
	// defaultProfile default profile to be used by *default methods.
	defaultProfile string
	// tls holds data to be used during TLS handshake.
	tls *TLS
}

// TLS holds data to be used during TLS handshake.
type TLS struct {
	// SkipTlsVerification tells the client to either skip the verification process or not.
	SkipTlsVerification bool
	// CaCertPath CA certificate to be used for StooKV server verification during handshake only if SkipTlsVerification is false
	// which is the default behaviour.
	CaCertPath string
	// ServerNameOverride StooKV server hostname to be used during TLS hostname verification.
	ServerNameOverride string
}

// DefaultTimeout default timeout to be used if not specified.
const DefaultTimeout = 10 * time.Second

// NewDefaultStooConfig creates StooConfig from default settings.
func NewDefaultStooConfig() *StooConfig {
	return &StooConfig{
		endpoint:    "localhost:50051",
		readTimeout: 10 * time.Second,
	}
}

// NewStooConfig creates a new StooConfig, stops if endpoint is empty.
func NewStooConfig(endpoint string, timeout time.Duration) *StooConfig {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	if endpoint != "" {
		return &StooConfig{
			endpoint:    endpoint,
			readTimeout: timeout,
		}
	}
	log.Fatal("endpoint must be defined")
	return nil
}

// WithUseTls sets useTls.
func (s *StooConfig) WithUseTls(useTls bool) *StooConfig {
	s.useTls = useTls
	return s
}

// WithDefaultNamespace sets defaultNamespace.
func (s *StooConfig) WithDefaultNamespace(defaultNamespace string) *StooConfig {
	s.defaultNamespace = defaultNamespace
	return s
}

// WithDefaultProfile sets defaultProfile.
func (s *StooConfig) WithDefaultProfile(defaultProfile string) *StooConfig {
	s.defaultProfile = defaultProfile
	return s
}

// WithTls sets tls.
func (s *StooConfig) WithTls(tls *TLS) *StooConfig {
	if tls != nil {
		s.tls = tls
	}
	return s
}

// GetUseTls returns useTls.
func (s *StooConfig) GetUseTls() bool {
	return s.useTls
}

// GetDefaultNamespace returns defaultNamespace.
func (s *StooConfig) GetDefaultNamespace() string {
	return s.defaultNamespace
}

// GetDefaultProfile returns defaultProfile.
func (s *StooConfig) GetDefaultProfile() string {
	return s.defaultProfile
}

// GetEndpoint returns endpoint.
func (s *StooConfig) GetEndpoint() string {
	return s.endpoint
}

// GetReadTimeout returns readTimeout.
func (s *StooConfig) GetReadTimeout() time.Duration {
	return s.readTimeout
}

// GetTls returns tls.
func (s *StooConfig) GetTls() *TLS {
	return s.tls
}
