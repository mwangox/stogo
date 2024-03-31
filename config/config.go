package config

import "time"

type StooConfig struct {
	endpoint         string
	useTls           bool
	readTimeout      time.Duration
	defaultNamespace string
	defaultProfile   string
	tls              *TLS
}

type TLS struct {
	SkipTlsVerification bool
	CaCertPath          string
	ServerNameOverride  string
}

func NewDefaultStooConfig() *StooConfig {
	return &StooConfig{
		endpoint:    "localhost:50051",
		readTimeout: 10 * time.Second,
	}
}

func NewStooConfig(endpoint string, timeout time.Duration) *StooConfig {
	if endpoint != "" && timeout != 0 {
		return &StooConfig{
			endpoint:    endpoint,
			readTimeout: timeout,
		}
	}
	panic("endpoint and timeout must be provided")
}

func (s *StooConfig) WithUseTls(useTls bool) *StooConfig {
	s.useTls = useTls
	return s
}

func (s *StooConfig) WithDefaultNamespace(defaultNamespace string) *StooConfig {
	s.defaultNamespace = defaultNamespace
	return s
}

func (s *StooConfig) WithDefaultProfile(defaultProfile string) *StooConfig {
	s.defaultProfile = defaultProfile
	return s
}

func (s *StooConfig) WithTls(tls *TLS) *StooConfig {
	if tls != nil {
		s.tls = tls
	}
	return s
}

func (s *StooConfig) GetUseTls() bool {
	return s.useTls
}

func (s *StooConfig) GetDefaultNamespace() string {
	return s.defaultNamespace
}

func (s *StooConfig) GetDefaultProfile() string {
	return s.defaultProfile
}

func (s *StooConfig) GetEndpoint() string {
	return s.endpoint
}

func (s *StooConfig) GetReadTimeout() time.Duration {
	return s.readTimeout
}

func (s *StooConfig) GetTls() *TLS {
	return s.tls
}
