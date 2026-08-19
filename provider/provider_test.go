package provider

import (
	"context"
	"testing"

	"github.com/NX211/traefik-proxmox-provider/internal"
)

func TestProviderConfig(t *testing.T) {
	config := CreateConfig()
	if config.PollInterval != "30s" {
		t.Errorf("Expected default PollInterval to be '30s', got %s", config.PollInterval)
	}
	if config.ApiValidateSSL != "true" {
		t.Errorf("Expected default ApiValidateSSL to be 'true', got %s", config.ApiValidateSSL)
	}
	if config.ApiLogging != "info" {
		t.Errorf("Expected default ApiLogging to be 'info', got %s", config.ApiLogging)
	}
}

func TestProviderNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true, // We expect an error because the domain doesn't exist
		},
		{
			name:    "Nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "Missing poll interval",
			config: &Config{
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Invalid poll interval",
			config: &Config{
				PollInterval:   "invalid",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(context.Background(), tt.config, "test-provider")
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && provider != nil {
				t.Error("Expected provider to be nil when there's an error")
			}
		})
	}
}

func TestProviderValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: false,
		},
		{
			name:    "Nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "Missing poll interval",
			config: &Config{
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Missing endpoint",
			config: &Config{
				PollInterval:   "5s",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Missing token ID",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Missing token",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Valid multi-node config",
			config: &Config{
				PollInterval: "5s",
				Nodes: []NodeConfig{
					{
						Name:        "cluster1",
						ApiEndpoint: "https://pve1.example.com",
						ApiTokenId:  "test@pam!test",
						ApiToken:    "token-1",
					},
					{
						Name:        "cluster2",
						ApiEndpoint: "https://pve2.example.com",
						ApiTokenId:  "test@pam!test2",
						ApiToken:    "token-2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi-node config with invalid node",
			config: &Config{
				PollInterval: "5s",
				Nodes: []NodeConfig{
					{
						Name:        "cluster1",
						ApiEndpoint: "https://pve1.example.com",
						ApiTokenId:  "test@pam!test",
						ApiToken:    "token-1",
					},
					{
						Name:       "cluster2",
						ApiTokenId: "test@pam!test2",
						ApiToken:   "token-2",
						// Missing ApiEndpoint
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProviderParserConfig(t *testing.T) {
	tests := []struct {
		name        string
		apiEndpoint string
		tokenID     string
		token       string
		wantErr     bool
	}{
		{
			name:        "Valid config",
			apiEndpoint: "https://proxmox.example.com",
			tokenID:     "test@pam!test",
			token:       "test-token",
			wantErr:     false,
		},
		{
			name:        "Missing endpoint",
			apiEndpoint: "",
			tokenID:     "test@pam!test",
			token:       "test-token",
			wantErr:     true,
		},
		{
			name:        "Missing token ID",
			apiEndpoint: "https://proxmox.example.com",
			tokenID:     "",
			token:       "test-token",
			wantErr:     true,
		},
		{
			name:        "Missing token",
			apiEndpoint: "https://proxmox.example.com",
			tokenID:     "test@pam!test",
			token:       "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := newParserConfig(tt.apiEndpoint, tt.tokenID, tt.token, "debug", true)
			if (err != nil) != tt.wantErr {
				t.Errorf("newParserConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if config.ApiEndpoint != tt.apiEndpoint {
					t.Errorf("Expected ApiEndpoint to be %s, got %s", tt.apiEndpoint, config.ApiEndpoint)
				}
				if config.TokenId != tt.tokenID {
					t.Errorf("Expected TokenId to be %s, got %s", tt.tokenID, config.TokenId)
				}
				if config.Token != tt.token {
					t.Errorf("Expected Token to be %s, got %s", tt.token, config.Token)
				}
			}
		})
	}
}

func TestProviderService(t *testing.T) {
	config := map[string]string{
		"traefik.enable":                 "true",
		"traefik.http.routers.test.rule": "Host(`test.example.com`)",
	}

	service := internal.NewService(123, "test-service", config)
	if service.ID != 123 {
		t.Errorf("Expected service ID to be 123, got %d", service.ID)
	}
	if service.Name != "test-service" {
		t.Errorf("Expected service name to be 'test-service', got %s", service.Name)
	}
	if len(service.Config) != 2 {
		t.Errorf("Expected service config to have 2 items, got %d", len(service.Config))
	}
	if len(service.IPs) != 0 {
		t.Errorf("Expected service IPs to be empty, got %d items", len(service.IPs))
	}
}

func TestGetServiceURLs(t *testing.T) {
	tests := []struct {
		name         string
		service      internal.Service
		serviceName  string
		nodeName     string
		expectedURLs []string
	}{
		{
			name:        "IP set, default port and scheme",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip": "1.2.3.4",
				},
			},
			expectedURLs: []string{"http://1.2.3.4:80"},
		},
		{
			name:        "IP and scheme set, default port (http)",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "http",
				},
			},
			expectedURLs: []string{"http://1.2.3.4:80"},
		},
		{
			name:        "IP and scheme set, default port (https)",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "https",
				},
			},
			expectedURLs: []string{"https://1.2.3.4:443"},
		},
		{
			name:        "IP, port and scheme set",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "https",
					"traefik.http.services.service.loadbalancer.server.port":   "8080",
				},
			},
			expectedURLs: []string{"https://1.2.3.4:8080"},
		},
		{
			name:        "URL is set",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.url": "http://test.com:1234",
				},
			},
			expectedURLs: []string{"http://test.com:1234"},
		},
		{
			name:        "URL overrides other settings",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.url":    "http://test.com:1234",
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "https",
					"traefik.http.services.service.loadbalancer.server.port":   "8080",
				},
			},
			expectedURLs: []string{"http://test.com:1234"},
		},
		{
			name:        "Multiple discovered IPs fan out sorted",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{},
				IPs: []internal.IP{
					{Address: "10.0.0.2", AddressType: "ipv4"},
					{Address: "10.0.0.1", AddressType: "ipv4"},
				},
			},
			expectedURLs: []string{"http://10.0.0.1:80", "http://10.0.0.2:80"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := getServiceURLs(tt.service, tt.serviceName, tt.nodeName)
			if len(urls) != len(tt.expectedURLs) {
				t.Fatalf("len(urls) = %d, want %d. Got %v", len(urls), len(tt.expectedURLs), urls)
			}
			for i, want := range tt.expectedURLs {
				if urls[i] != want {
					t.Errorf("urls[%d] = %q, want %q", i, urls[i], want)
				}
			}
		})
	}
}

func TestNormalizeNodes(t *testing.T) {
	t.Run("Legacy flat config", func(t *testing.T) {
		config := &Config{
			PollInterval:   "5s",
			ApiEndpoint:    "https://proxmox.example.com",
			ApiTokenId:     "test@pam!test",
			ApiToken:       "test-token",
			ApiLogging:     "debug",
			ApiValidateSSL: "false",
		}
		nodes := normalizeNodes(config)
		if len(nodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(nodes))
		}
		if nodes[0].Name != "default" {
			t.Errorf("Expected name 'default', got %s", nodes[0].Name)
		}
		if nodes[0].ApiEndpoint != "https://proxmox.example.com" {
			t.Errorf("Expected endpoint from flat config, got %s", nodes[0].ApiEndpoint)
		}
		if nodes[0].ApiLogging != "debug" {
			t.Errorf("Expected logging 'debug', got %s", nodes[0].ApiLogging)
		}
	})

	t.Run("Multi-node config", func(t *testing.T) {
		config := &Config{
			PollInterval: "5s",
			Nodes: []NodeConfig{
				{
					Name:        "cluster1",
					ApiEndpoint: "https://pve1.example.com",
					ApiTokenId:  "test@pam!test",
					ApiToken:    "token-1",
				},
				{
					Name:        "cluster2",
					ApiEndpoint: "https://pve2.example.com",
					ApiTokenId:  "test@pam!test2",
					ApiToken:    "token-2",
				},
			},
		}
		nodes := normalizeNodes(config)
		if len(nodes) != 2 {
			t.Fatalf("Expected 2 nodes, got %d", len(nodes))
		}
		if nodes[0].Name != "cluster1" {
			t.Errorf("Expected name 'cluster1', got %s", nodes[0].Name)
		}
		if nodes[1].Name != "cluster2" {
			t.Errorf("Expected name 'cluster2', got %s", nodes[1].Name)
		}
	})

	t.Run("Defaults filled in", func(t *testing.T) {
		config := &Config{
			PollInterval: "5s",
			Nodes: []NodeConfig{
				{
					ApiEndpoint: "https://pve1.example.com",
					ApiTokenId:  "test@pam!test",
					ApiToken:    "token-1",
					// Name, ApiLogging, ApiValidateSSL omitted
				},
			},
		}
		nodes := normalizeNodes(config)
		if nodes[0].Name != "node-0" {
			t.Errorf("Expected auto-generated name 'node-0', got %s", nodes[0].Name)
		}
		if nodes[0].ApiLogging != "info" {
			t.Errorf("Expected default logging 'info', got %s", nodes[0].ApiLogging)
		}
		if nodes[0].ApiValidateSSL != "true" {
			t.Errorf("Expected default SSL validation 'true', got %s", nodes[0].ApiValidateSSL)
		}
	})
}

func TestGenerateConfiguration_MultiCluster(t *testing.T) {
	// Two clusters, each with one node and one service that uses
	// only traefik.enable=true (so the auto-generated defaultID is used).
	clusterMaps := []clusterServiceMap{
		{
			clusterName: "cluster1",
			services: map[string][]internal.Service{
				"node1": {
					{
						ID:   100,
						Name: "webserver",
						IPs:  []internal.IP{{Address: "10.0.0.1", AddressType: "ipv4"}},
						Config: map[string]string{
							"traefik.enable": "true",
						},
					},
				},
			},
		},
		{
			clusterName: "cluster2",
			services: map[string][]internal.Service{
				"node2": {
					{
						ID:   100, // Same VMID as cluster1 — should not collide
						Name: "webserver",
						IPs:  []internal.IP{{Address: "10.0.1.1", AddressType: "ipv4"}},
						Config: map[string]string{
							"traefik.enable": "true",
						},
					},
				},
			},
		},
	}

	cfg := generateConfiguration(clusterMaps, true)

	// With multiCluster=true, auto-generated IDs should be prefixed.
	expectedService1 := "cluster1-webserver-100"
	expectedService2 := "cluster2-webserver-100"

	if _, ok := cfg.HTTP.Services[expectedService1]; !ok {
		t.Errorf("Expected service %q not found. Services: %v", expectedService1, keys(cfg.HTTP.Services))
	}
	if _, ok := cfg.HTTP.Services[expectedService2]; !ok {
		t.Errorf("Expected service %q not found. Services: %v", expectedService2, keys(cfg.HTTP.Services))
	}
	if _, ok := cfg.HTTP.Routers[expectedService1]; !ok {
		t.Errorf("Expected router %q not found. Routers: %v", expectedService1, keys(cfg.HTTP.Routers))
	}
	if _, ok := cfg.HTTP.Routers[expectedService2]; !ok {
		t.Errorf("Expected router %q not found. Routers: %v", expectedService2, keys(cfg.HTTP.Routers))
	}

	// Verify they have different server URLs
	svc1 := cfg.HTTP.Services[expectedService1]
	svc2 := cfg.HTTP.Services[expectedService2]
	if len(svc1.LoadBalancer.Servers) == 0 || len(svc2.LoadBalancer.Servers) == 0 {
		t.Fatal("Expected non-empty servers")
	}
	if svc1.LoadBalancer.Servers[0].URL == svc2.LoadBalancer.Servers[0].URL {
		t.Errorf("Services should have different URLs, both got %s", svc1.LoadBalancer.Servers[0].URL)
	}
}

func TestGenerateConfiguration_SingleCluster_NoPrefixing(t *testing.T) {
	// Single cluster: defaultID should NOT be prefixed.
	clusterMaps := []clusterServiceMap{
		{
			clusterName: "default",
			services: map[string][]internal.Service{
				"node1": {
					{
						ID:   200,
						Name: "myapp",
						IPs:  []internal.IP{{Address: "10.0.0.5", AddressType: "ipv4"}},
						Config: map[string]string{
							"traefik.enable": "true",
						},
					},
				},
			},
		},
	}

	cfg := generateConfiguration(clusterMaps, false)

	expectedID := "myapp-200"
	if _, ok := cfg.HTTP.Services[expectedID]; !ok {
		t.Errorf("Expected service %q not found (should not be prefixed in single-cluster mode). Services: %v",
			expectedID, keys(cfg.HTTP.Services))
	}
}

// keys is a test helper that returns sorted map keys for readable error messages.
func keys[V any](m map[string]V) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func TestGenerateConfiguration_MultipleRouters(t *testing.T) {
	// A single guest defining multiple routers and multiple services should
	// wire each router to the service whose name matches the router, rather than
	// pointing every router at the first service.
	service := internal.Service{
		ID:   100,
		Name: "multi",
		Config: map[string]string{
			"traefik.enable":                                     "true",
			"traefik.http.routers.blog.rule":                     "Host(`blog.example.com`)",
			"traefik.http.routers.shop.rule":                     "Host(`shop.example.com`)",
			"traefik.http.services.blog.loadbalancer.server.url": "http://23.45.67.89:8080",
			"traefik.http.services.shop.loadbalancer.server.url": "http://23.45.67.90:8080",
		},
	}
	clusterMaps := []clusterServiceMap{
		{
			clusterName: "default",
			services:    map[string][]internal.Service{"node1": {service}},
		},
	}

	config := generateConfiguration(clusterMaps, false)

	if len(config.HTTP.Routers) != 2 {
		t.Fatalf("Expected 2 routers, got %d: %v", len(config.HTTP.Routers), config.HTTP.Routers)
	}

	for _, name := range []string{"blog", "shop"} {
		router, exists := config.HTTP.Routers[name]
		if !exists {
			t.Errorf("Expected router %q to exist", name)
			continue
		}
		if router.Service != name {
			t.Errorf("Router %q should target service %q, got %q", name, name, router.Service)
		}
		if router.Priority != 0 {
			t.Errorf("Router %q Priority = %d, want 0 (unset default)", name, router.Priority)
		}
	}
}

func TestMapKeysToSlice_Sorted(t *testing.T) {
	in := map[string]bool{"charlie": true, "alpha": true, "bravo": true, "delta": true}
	got := mapKeysToSlice(in)
	want := []string{"alpha", "bravo", "charlie", "delta"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("got[%d] = %q, want %q (full got = %v)", i, got[i], v, got)
		}
	}
}

func TestIsBoolLabelEnabled(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
		want   bool
	}{
		{"missing key", map[string]string{}, false},
		{"empty value", map[string]string{"traefik.enable": ""}, false},
		{"true literal", map[string]string{"traefik.enable": "true"}, true},
		{"True mixed case", map[string]string{"traefik.enable": "True"}, true},
		{"TRUE upper", map[string]string{"traefik.enable": "TRUE"}, true},
		{"t shorthand", map[string]string{"traefik.enable": "t"}, true},
		{"T shorthand", map[string]string{"traefik.enable": "T"}, true},
		{"1 numeric", map[string]string{"traefik.enable": "1"}, true},
		{"false literal", map[string]string{"traefik.enable": "false"}, false},
		{"0 numeric", map[string]string{"traefik.enable": "0"}, false},
		{"f shorthand", map[string]string{"traefik.enable": "f"}, false},
		// Intentionally NOT supported (Traefik core doesn't accept these):
		{"yes rejected", map[string]string{"traefik.enable": "yes"}, false},
		{"on rejected", map[string]string{"traefik.enable": "on"}, false},
		{"junk rejected", map[string]string{"traefik.enable": "blarg"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBoolLabelEnabled(tt.labels, "traefik.enable")
			if got != tt.want {
				t.Errorf("isBoolLabelEnabled(%v) = %v, want %v", tt.labels, got, tt.want)
			}
		})
	}
}

func TestHandleRouterTLS_ArrayDomains(t *testing.T) {
	tests := []struct {
		name           string
		config         map[string]string
		expectedMain   []string
		expectedSANs   [][]string
		expectNil      bool
	}{
		{
			name: "Array syntax with main and sans",
			config: map[string]string{
				"traefik.http.routers.test.tls.domains[0].main": "example.com",
				"traefik.http.routers.test.tls.domains[0].sans": "*.example.com,www.example.com",
				"traefik.http.routers.test.tls.domains[1].main": "another.com",
			},
			expectedMain: []string{"example.com", "another.com"},
			expectedSANs: [][]string{{"*.example.com", "www.example.com"}, nil},
		},
		{
			name: "Simple domains fallback",
			config: map[string]string{
				"traefik.http.routers.test.tls.domains": "example.com,another.com",
			},
			expectedMain: []string{"example.com", "another.com"},
			expectedSANs: [][]string{nil, nil},
		},
		{
			name:      "No TLS config",
			config:    map[string]string{},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := internal.Service{Config: tt.config}
			tlsConfig := handleRouterTLS(service, "traefik.http.routers.test")

			if tt.expectNil {
				if tlsConfig != nil {
					t.Error("Expected nil TLS config")
				}
				return
			}

			if tlsConfig == nil {
				t.Fatal("Expected non-nil TLS config")
			}

			if len(tlsConfig.Domains) != len(tt.expectedMain) {
				t.Fatalf("Expected %d domains, got %d", len(tt.expectedMain), len(tlsConfig.Domains))
			}

			for i, domain := range tlsConfig.Domains {
				if domain.Main != tt.expectedMain[i] {
					t.Errorf("Domain[%d].Main = %s, want %s", i, domain.Main, tt.expectedMain[i])
				}
				if tt.expectedSANs[i] != nil {
					if len(domain.SANs) != len(tt.expectedSANs[i]) {
						t.Errorf("Domain[%d].SANs length = %d, want %d", i, len(domain.SANs), len(tt.expectedSANs[i]))
					}
				}
			}
		})
	}
}
