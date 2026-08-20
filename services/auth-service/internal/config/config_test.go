package config

import "testing"

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		port        string
		databaseURL string
		environment string
		wantPort    string
		wantEnv     string
		wantErr     bool
	}{
		{
			name:        "loads configuration from environment",
			port:        "9090",
			databaseURL: "postgres://test:test@localhost:5432/test",
			environment: "testing",
			wantPort:    "9090",
			wantEnv:     "testing",
		},
		{
			name:        "uses default port",
			databaseURL: "postgres://test:test@localhost:5432/test",
			environment: "testing",
			wantPort:    "8081",
			wantEnv:     "testing",
		},
		{
			name:        "uses default environment",
			port:        "9090",
			databaseURL: "postgres://test:test@localhost:5432/test",
			wantPort:    "9090",
			wantEnv:     "development",
		},
		{
			name:        "requires database URL",
			port:        "9090",
			environment: "testing",
			wantErr:     true,
		},
		{
			name:        "rejects invalid environment",
			port:        "9090",
			databaseURL: "postgres://test:test@localhost:5432/test",
			environment: "banana",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AUTH_SERVICE_PORT", tt.port)
			t.Setenv("DATABASE_URL", tt.databaseURL)
			t.Setenv("APP_ENV", tt.environment)

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Fatal("Load() expected an error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Load() returned unexpected error: %v", err)
			}

			if cfg.AppPort != tt.wantPort {
				t.Errorf("AppPort = %q, want %q", cfg.AppPort, tt.wantPort)
			}

			if cfg.Environment != tt.wantEnv {
				t.Errorf("Environment = %q, want %q", cfg.Environment, tt.wantEnv)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: Config{
				AppPort:     "8081",
				DatabaseURL: "postgres://test:test@localhost:5432/test",
				Environment: "development",
			},
		},
		{
			name: "missing app port",
			config: Config{
				DatabaseURL: "postgres://test:test@localhost:5432/test",
				Environment: "development",
			},
			wantErr: true,
		},
		{
			name: "missing database URL",
			config: Config{
				AppPort:     "8081",
				Environment: "development",
			},
			wantErr: true,
		},
		{
			name: "invalid environment",
			config: Config{
				AppPort:     "8081",
				DatabaseURL: "postgres://test:test@localhost:5432/test",
				Environment: "banana",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Fatal("Validate() expected an error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Validate() returned unexpected error: %v", err)
			}
		})
	}
}
