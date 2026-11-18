package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type CertificateManager struct {
	Enabled        bool     `description:"Enable certificate manager integration." json:"enabled,omitempty" export:"true"`
	CertificateId  string   `description:"ID of the Certificate to fetch." json:"certificateId,omitempty" export:"true"`
	CertificateIds []string `description:"IDs of the Certificates to fetch." json:"certificateIds,omitempty" export:"true"`
}

// IDs returns CertificateId followed by CertificateIds, trimmed and deduplicated.
func (c *CertificateManager) IDs() []string {
	if c == nil {
		return nil
	}
	seen := map[string]struct{}{}
	var ids []string
	for _, id := range append([]string{c.CertificateId}, c.CertificateIds...) {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
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
	}
	c.Features.SetDefaults()
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

func envInt(name string, dst *int) error {
	s := os.Getenv(name)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	*dst = v
	return nil
}

func applyEnv(config *Configuration) error {
	if err := envInt("PORT", &config.Port); err != nil {
		return err
	}
	if err := envInt("CERT_FETCH_INTERVAL_SECONDS", &config.CertFetchIntervalSeconds); err != nil {
		return err
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		config.LogLevel = logLevel
	}

	saKeyPath := os.Getenv("SA_KEY_PATH")
	if saKeyPath != "" {
		config.SAKeyPath = saKeyPath
	}

	cm := config.Features.CertificateManager
	certID := os.Getenv("CERTIFICATE_ID")
	if certID != "" {
		cm.Enabled = true
		cm.CertificateId = certID
	}
	certIDs := os.Getenv("CERTIFICATE_IDS")
	if certIDs != "" {
		cm.Enabled = true
		cm.CertificateIds = strings.Split(certIDs, ",")
	}
	// TODO: implement remote config when YDB support Coordination Service in Serverless
	//remoteConfig := os.Getenv("REMOTE_CONFIG")
	//if remoteConfig != "" {
	//	config.Features.RemoteConfig.Enabled = true
	//	config.Features.RemoteConfig.YdbEndpoint = remoteConfig
	//}
	return nil
}

func NewConfiguration() (*Configuration, error) {
	return load(".")
}

func load(dir string) (*Configuration, error) {
	v := viper.New()
	v.SetConfigName("config") // name of config file (without extension)
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	config.SetDefaults()

	if err := applyEnv(&config); err != nil {
		return nil, err
	}
	if config.CertFetchIntervalSeconds <= 0 {
		return nil, fmt.Errorf("certFetchIntervalSeconds must be positive, got %d", config.CertFetchIntervalSeconds)
	}

	return &config, nil
}
