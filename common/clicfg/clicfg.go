package clicfg

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/neo4j/cli/common/clicfg/credentials"
	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var ConfigPrefix string

const (
	DefaultAuraBaseUrl     = "https://api.neo4j.io/v1"
	DefaultAuraBetaBaseUrl = "https://api.neo4j.io/v1beta5"
	DefaultAuraAuthUrl     = "https://api.neo4j.io/oauth/token"
	DefaultAuraBetaEnabled = "false"
)

var ValidOutputValues = [3]string{"default", "json", "table"}

type Config struct {
	Version     string
	Aura        *AuraConfig
	Credentials *credentials.Credentials
}

func NewConfig(fs afero.Fs, version string) (*Config, error) {
	configPath := filepath.Join(ConfigPrefix, "neo4j", "cli")

	Viper := viper.New()

	Viper.SetFs(fs)
	Viper.SetConfigName("config")
	Viper.SetConfigType("json")
	Viper.AddConfigPath(configPath)
	Viper.SetConfigPermissions(0600)

	bindEnvironmentVariables(Viper)
	setDefaultValues(Viper)

	if err := Viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := fs.MkdirAll(configPath, 0755); err != nil {
				return nil, err
			}
			if err = Viper.SafeWriteConfig(); err != nil {
				return nil, err
			}
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	credentials, err := credentials.NewCredentials(fs, ConfigPrefix)
	if err != nil {
		return nil, err
	}

	return &Config{
		Version: version,
		Aura: &AuraConfig{
			fs:    fs,
			viper: Viper, pollingOverride: PollingConfig{
				MaxRetries: 60,
				Interval:   20,
			},
			ValidConfigKeys: []string{"auth-url", "base-url", "default-tenant", "output", "beta-enabled"},
		},
		Credentials: credentials,
	}, nil
}

func bindEnvironmentVariables(Viper *viper.Viper) {
	Viper.BindEnv("aura.base-url", "AURA_BASE_URL")
	Viper.BindEnv("aura.auth-url", "AURA_AUTH_URL")
}

func setDefaultValues(Viper *viper.Viper) {
	Viper.SetDefault("aura.base-url", DefaultAuraBaseUrl)
	Viper.SetDefault("aura.auth-url", DefaultAuraAuthUrl)
	Viper.SetDefault("aura.output", "default")
	Viper.SetDefault("aura.beta-enabled", DefaultAuraBetaEnabled)
}

type AuraConfig struct {
	viper           *viper.Viper
	fs              afero.Fs
	pollingOverride PollingConfig
	ValidConfigKeys []string
}

type PollingConfig struct {
	Interval   int
	MaxRetries int
}

func (config *AuraConfig) IsValidConfigKey(key string) bool {
	return slices.Contains(config.ValidConfigKeys, key)
}

func (config *AuraConfig) Get(key string) interface{} {
	return config.viper.Get(fmt.Sprintf("aura.%s", key))
}

func (config *AuraConfig) Set(key string, value string) error {
	filename := config.viper.ConfigFileUsed()
	data, err := fileutils.ReadFileSafe(config.fs, filename)
	if err != nil {
		return err
	}

	updateConfig, err := sjson.Set(string(data), fmt.Sprintf("aura.%s", key), value)
	if err != nil {
		return err
	}

	updatedAuraBaseUrl := config.auraBaseUrlOnBetaEnabledChange(key, value)
	if updatedAuraBaseUrl != "" {
		intermediateUpdateConfig, err := sjson.Set(string(updateConfig), "aura.base-url", updatedAuraBaseUrl)
		if err != nil {
			return err
		}
		updateConfig = intermediateUpdateConfig
	}

	return fileutils.WriteFile(config.fs, filename, []byte(updateConfig))
}

func (config *AuraConfig) Print(cmd *cobra.Command) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config.viper.Get("aura")); err != nil {
		return err
	}

	return nil
}

func (config *AuraConfig) BaseUrl() string {
	return config.viper.GetString("aura.base-url")
}

func (config *AuraConfig) BindBaseUrl(flag *pflag.Flag) error {
	return config.viper.BindPFlag("aura.base-url", flag)
}

func (config *AuraConfig) AuthUrl() string {
	return config.viper.GetString("aura.auth-url")
}

func (config *AuraConfig) BindAuthUrl(flag *pflag.Flag) error {
	return config.viper.BindPFlag("aura.auth-url", flag)
}

func (config *AuraConfig) Output() string {
	return config.viper.GetString("aura.output")
}

func (config *AuraConfig) BindOutput(flag *pflag.Flag) error {
	return config.viper.BindPFlag("aura.output", flag)
}

func (config *AuraConfig) AuraBetaEnabled() string {
	return config.viper.GetString("aura.beta-enabled")
}

func (config *AuraConfig) DefaultTenant() string {
	return config.viper.GetString("aura.default-tenant")
}

func (config *AuraConfig) PollingConfig() PollingConfig {
	return config.pollingOverride
}

func (config *AuraConfig) SetPollingConfig(maxRetries int, interval int) {
	config.pollingOverride = PollingConfig{
		MaxRetries: maxRetries,
		Interval:   interval,
	}
}

func (config *AuraConfig) auraBaseUrlOnBetaEnabledChange(key string, value string) string {
	if key == "beta-enabled" {
		nextBaseUrl := DefaultAuraBaseUrl
		if value == "true" {
			nextBaseUrl = DefaultAuraBetaBaseUrl
		}
		return nextBaseUrl
	}
	return ""
}
