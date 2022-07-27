// Package internal contains the core application components.
package internal

// Executor is the main plugin method to execute plugins logic.
type Executor interface {
	Execute() error
}

// ConfigConfigurator is an interface for WithConfig.
type ConfigConfigurator interface {
	WithConfig(config *Config)
}
