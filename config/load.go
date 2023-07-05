package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/util"
	"gopkg.in/yaml.v2"
)

const NoConfig = "-no-config"
const DefaultConfigFile = "g.yaml"

func GetUserConfigDir() (string, error) {
	err := InitConfigDir.Do(func() error {
		home, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		ConfigDir = filepath.Join(home, "g")
		err = os.MkdirAll(ConfigDir, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return ConfigDir, nil
}

var InitConfigDir util.Once
var ConfigDir = ""

// READ config
// g.yaml
// Args:
//     - args
//     - ...

type Config struct {
	Args []string `yaml:"Args"`
}

type ErrReadConfig struct {
	error
}

func Load() (*Config, error) {

	Dir, err := GetUserConfigDir()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filepath.Join(Dir, DefaultConfigFile))
	if err != nil {
		return nil, err
	}
	Default := Config{
		Args: make([]string, 0),
	}
	// parse yaml
	configErr := yaml.Unmarshal(content, &Default)
	if configErr != nil {
		return nil, ErrReadConfig{error: configErr}
	}

	for i, v := range Default.Args {
		if v == NoConfig {
			return nil, nil
		}
		// if not prefixed with '-', add '-'
		if !strings.HasPrefix(v, "-") {
			if len(v) == 1 {
				Default.Args[i] = "-" + v
			} else {
				Default.Args[i] = "--" + v
			}
		}
	}

	return &Default, nil
}
