package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var cfg Config

var DefaultConfigDir = "$XDG_CONFIG_HOME/cobble"
var DefaultConfigPath = filepath.Join(DefaultConfigDir, "config.yaml")
var DefaultScaffoldsDir = filepath.Join(DefaultConfigDir, "scaffolds")

var ConfigFileFlag = "config-file"
var ScaffoldsDirFlag = "scaffolds-dir"
var WorkingDirFlag = "working-dir"

func ScaffoldsDir() (string, error) {
	return filepath.Abs(os.ExpandEnv(cfg.ScaffoldsDir))
}

func WorkingDir() (string, error) {
	return filepath.Abs(os.ExpandEnv(cfg.WorkingDir))
}

type Config struct {
	ScaffoldsDir string `mapstructure:"scaffolds-dir"` // ScaffoldsDirFlag
	WorkingDir   string `mapstructure:"working-dir"`   // WorkingDirFlag
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() error {
	config := viper.GetString(ConfigFileFlag)

	if config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config)
	} else {
		configDir, err := filepath.Abs(os.ExpandEnv(DefaultConfigDir))
		if err != nil {
			return err
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetDefault(ScaffoldsDirFlag, DefaultScaffoldsDir)
	viper.SetDefault(WorkingDirFlag, ".")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	return viper.Unmarshal(&cfg)
}
