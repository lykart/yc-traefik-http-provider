package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

type Config struct {
	CertificateId string `mapstructure:"certificate_id" json:"certificateId,omitempty"`
	Port          int    `mapstructure:"port" json:"port,omitempty"`
	SAKeyPath     string `mapstructure:"sa_key" json:"saKey,omitempty"`
	SAKey         *iamkey.Key
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("json")
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.BindEnv("CERTIFICATE_ID")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("PORT")
	if err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.SAKeyPath != "" {
		log.Printf("SAKey loaded: %s", config.SAKeyPath)
		data, err := os.ReadFile(config.SAKeyPath)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, &config.SAKey); err != nil {
			return nil, err
		}
	}

	log.Printf("config loaded: %+v", config)

	if config.CertificateId == "" {
		return nil, fmt.Errorf("CERTIFICATE_ID is required")
	}

	return &config, nil
}
