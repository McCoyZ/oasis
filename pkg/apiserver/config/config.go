package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"zmc.io/oasis/pkg/simple/client/k8s"
	"zmc.io/oasis/pkg/simple/client/servicemesh"
)

const (
	// DefaultConfigurationName is the default name of configuration
	defaultConfigurationName = "oasis.yaml"

	// DefaultConfigurationPath the default location of the configuration file
	defaultConfigurationPath = "hack"
)

type Config struct {
	KubernetesOptions  *k8s.KubernetesOptions `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty" mapstructure:"kubernetes"`
	ServiceMeshOptions *servicemesh.Options   `json:"servicemesh,omitempty" yaml:"servicemesh,omitempty" mapstructure:"servicemesh"`
}

func New() *Config {
	return &Config{
		KubernetesOptions:  k8s.NewKubernetesOptions(),
		ServiceMeshOptions: servicemesh.NewServiceMeshOptions(),
	}
}

func TryLoadFromDisk() (*Config, error) {
	viper.SetConfigName(defaultConfigurationName)
	viper.AddConfigPath(defaultConfigurationPath)

	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(err)
			return nil, err
		} else {
			return nil, fmt.Errorf("error parsing configuration file %s", err)
		}
	}

	conf := New()

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// convertToMap simply converts config to map[string]bool
// to hide sensitive information
func (conf *Config) ToMap() map[string]bool {
	conf.stripEmptyOptions()
	result := make(map[string]bool, 0)

	if conf == nil {
		return result
	}

	c := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < c.NumField(); i++ {
		name := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if strings.HasPrefix(name, "-") {
			continue
		}

		if c.Field(i).IsNil() {
			result[name] = false
		} else {
			result[name] = true
		}
	}

	return result
}

func (conf *Config) stripEmptyOptions() {
	if conf.ServiceMeshOptions != nil && conf.ServiceMeshOptions.IstioPilotHost == "" {
		conf.ServiceMeshOptions = nil
	}
}
