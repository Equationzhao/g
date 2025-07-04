package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/internal/util"
	"gopkg.in/yaml.v3"
)

const (
	NoConfig          = "-no-config"
	DefaultConfigFile = "g.yaml"
)

func GetUserConfigDir() (string, error) {
	err := InitConfigDir.Do(
		func() error {
			home, err := os.UserConfigDir()
			if err != nil {
				return err
			}
			Dir = filepath.Join(home, "g")
			err = os.MkdirAll(Dir, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return "", err
	}
	return Dir, nil
}

var (
	InitConfigDir util.Once
	Dir           = ""
)

// READ config
// g.yaml
// Args:
//     - args
//     - ...

type Config struct {
	Args            []string  `yaml:"Args"`
	CustomTreeStyle TreeStyle `yaml:"CustomTreeStyle"`
	ThemeLocation   string    `yaml:"Theme"`
	Order           []string  `yaml:"Order"`
}

type TreeStyle struct {
	Child     string `yaml:"Child"`
	LastChild string `yaml:"LastChild"`
	Mid       string `yaml:"Mid"`
	Empty     string `yaml:"Empty"`
}

func (t TreeStyle) IsEmpty() bool {
	return t.Empty == "" && t.Child == "" && t.LastChild == "" && t.Mid == ""
}

func (t TreeStyle) IsEnabled() bool {
	return !t.IsEmpty()
}

type ErrReadConfig struct {
	error
	Location string
}

func (e ErrReadConfig) Error() string {
	if e.Location != "" {
		return fmt.Sprintf("failed to load configuration at %s: %s", e.Location, e.error.Error())
	}
	return fmt.Sprintf("failed to load configuration: %s", e.error.Error())
}

var Default = Config{
	Args: make([]string, 0),
}

var emptyConfig = Config{}

// GetOrder returns the column order from the loaded config
func GetOrder() []string {
	return Default.Order
}

func Load() (*Config, error) {
	Dir, err := GetUserConfigDir()
	if err != nil {
		return nil, err
	}

	location := filepath.Join(Dir, DefaultConfigFile)
	content, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	// parse yaml
	configErr := yaml.Unmarshal(content, &Default)
	if configErr != nil {
		return nil, ErrReadConfig{error: configErr, Location: location}
	}

	for i, v := range Default.Args {
		if v == NoConfig {
			Default = emptyConfig
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
