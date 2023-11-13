package config

// DatabaseConfig is the generic structure type for database configuration
type DatabaseConfig struct {
	Type     string `yaml:"type" toml:"Type"`
	Host     string `env:"DB_HOST" yaml:"host" toml:"Host" `
	Port     string `env:"DB_PORT" yaml:"port" toml:"Port"`
	DBName   string `yaml:"db_name" toml:"DBName"`
	User     string `yaml:"user" toml:"User"`
	Password string `env:"DB_PASSWORD" yaml:"password" toml:"Password"`
	SSLMode  string `yaml:"sslmode" toml:"SSLMode"`
	Schema   string `yaml:"schema" toml:"Schema"`

	// Collections is a map of key/values :
	// - keys define the way these collections are called from the code
	// - values are the names of collections on the MongoDB side.
	Collections map[string]string
}

// C returns the collection name related to the given key
func (d *DatabaseConfig) C(key string) string {
	return d.Collections[key]
}
