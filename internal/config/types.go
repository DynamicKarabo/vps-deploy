package config

// Service represents a single deployable service defined in deploy.yaml.
type Service struct {
	Name            string   `yaml:"-" json:"-"`
	Host            string   `yaml:"host" json:"host"`
	User            string   `yaml:"user" json:"user"`
	KeyPath         string   `yaml:"key_path" json:"key_path"`
	DeployCommands  []string `yaml:"deploy_commands" json:"deploy_commands"`
	HealthCheckURL  string   `yaml:"health_check_url" json:"health_check_url"`
	RollbackCommand string   `yaml:"rollback_command" json:"rollback_command"`
}

// Config is the top-level deploy.yaml structure.
type Config struct {
	Services map[string]*Service `yaml:"services" json:"services"`
}
