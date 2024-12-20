package configs

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
	Name     string `mapstructure:"Name"`
}

type ServerConfig struct {
	Port int
	Env  string
}

func Load() (*Config, error) {
	v := viper.New()
	if isDev := os.Getenv("APP_ENV"); isDev == "development" {
		v.SetConfigName("configs-dev")
	} else {
		v.SetConfigName("config")
	}

	v.SetConfigType("yaml")

	v.AddConfigPath("./configs")
	v.AutomaticEnv()

	v.SetEnvPrefix("APP") // prefix for environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, k := range v.AllKeys() {
		val := v.GetString(k)
		if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
			// Extract environment variable name
			envVar := strings.TrimSuffix(strings.TrimPrefix(val, "${"), "}")

			// Check for default value
			defaultValue := ""
			if strings.Contains(envVar, ":") {
				parts := strings.SplitN(envVar, ":", 2)
				envVar = parts[0]
				defaultValue = parts[1]
			}

			// Get environment variable value
			envVal := os.Getenv(envVar)

			if envVal == "" {
				envVal = defaultValue
			}
			// Set the processed value
			v.Set(k, envVal)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
