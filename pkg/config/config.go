package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port       string `mapstructure:"port" yaml:"port"`
	DBHost     string `mapstructure:"db_host" yaml:"db_host"`
	DBPort     string `mapstructure:"db_port" yaml:"db_port"`
	DBUser     string `mapstructure:"db_user" yaml:"db_user"`
	DBPassword string `mapstructure:"db_password" yaml:"db_password"`
	DBName     string `mapstructure:"db_name" yaml:"db_name"`
	SSLMode    string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
}

func Read() *AppConfig {
	// Load .env file first
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	v := viper.New()

	v.SetConfigName("config")      // name of config file (without extension)
	v.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("$PWD/config") // call multiple times to add many search paths
	v.AddConfigPath(".")           // optionally look for config in the working directory
	v.AddConfigPath("/config")     // optionally look for config in the working directory
	v.AddConfigPath("./config")    // optionally look for config in the working directory

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Warning: config file not found: %v\n", err)
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return &config
}
