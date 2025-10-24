package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type CertificateManager struct {
	Enabled       bool   `description:"Enable certificate manager integration." json:"enabled,omitempty" export:"true"`
	CertificateId string `description:"ID of the Certificate to fetch." json:"certificateId,omitempty" export:"true"`
}

type RemoteConfig struct {
	Enabled     bool   `description:"Enable remote configuration." json:"enabled,omitempty" export:"true"`
	YdbEndpoint string `description:"YDB database endpoint." json:"ydbEndpoint,omitempty" export:"true"`
}

type Features struct {
	CertificateManager *CertificateManager `description:"Add certificate manager integration." json:"certificateManager,omitempty" export:"true"`
	//TODO: implement remote config when YDB support Coordination Service in Serverless
	//RemoteConfig       *RemoteConfig       `description:"Use YDB as a configuration store." json:"remoteConfig,omitempty" export:"true"`
}

func (f *Features) SetDefaults() {
	if f.CertificateManager == nil {
		f.CertificateManager = &CertificateManager{
			Enabled: false,
		}
	}
	// TODO: implement remote config when YDB support Coordination Service in Serverless
	//if f.RemoteConfig == nil {
	//	f.RemoteConfig = &RemoteConfig{
	//		Enabled: false,
	//	}
	//}
}

type Configuration struct {
	Features                 *Features `description:"Configure features." json:"features,omitempty" export:"true"`
	Port                     int       `description:"Port HTTP configuration server listens on." json:"port,omitempty" export:"true"`
	SAKeyPath                string    `description:"Path to the service account key file." json:"saKeyPath,omitempty" export:"true"`
	CertFetchIntervalSeconds int       `description:"Interval in seconds to fetch certificate content." json:"certFetchIntervalSeconds,omitempty" export:"true"`
	LogLevel                 string    `json:"logLevel,omitempty" export:"true"`
}

func (c *Configuration) SetDefaults() {
	if c.Features == nil {
		c.Features = &Features{}
		c.Features.SetDefaults()
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.CertFetchIntervalSeconds == 0 {
		c.CertFetchIntervalSeconds = 300
	}
}

func NewConfiguration() (*Configuration, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("json")
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	config.SetDefaults()

	port := os.Getenv("PORT")
	if port != "" {
		v, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		config.Port = v
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		config.LogLevel = logLevel
	}

	certID := os.Getenv("CERTIFICATE_ID")
	if certID != "" {
		config.Features.CertificateManager.Enabled = true
		config.Features.CertificateManager.CertificateId = certID
	}
	// TODO: implement remote config when YDB support Coordination Service in Serverless
	//remoteConfig := os.Getenv("REMOTE_CONFIG")
	//if remoteConfig != "" {
	//	config.Features.RemoteConfig.Enabled = true
	//	config.Features.RemoteConfig.YdbEndpoint = remoteConfig
	//}

	return &config, nil
}
