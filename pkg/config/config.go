package config

import (
	"fmt"
	"strings"

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
	v := viper.New()
	v.SetConfigName("config")      // name of config file (without extension)
	v.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("$PWD/config") // call multiple times to add many search paths
	v.AddConfigPath(".")           // optionally look for config in the working directory
	v.AddConfigPath("/config")     // optionally look for config in the working directory
	v.AddConfigPath("./config")    // optionally look for config in the working directory
	err := v.ReadInConfig()        // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Env variables override config file
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &appConfig
}
