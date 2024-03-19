package config

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	lock        sync.Mutex
	logType     string        // required internally by config
	logLevel    string        // required internally by config
	packs       []int         // will have get/set due to mutex
	ServerPort  int           // free to access by server, only required in setup
	HttpTimeout time.Duration // free to access by server, only required in setup
}

func New() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.Wrap(err, "failed to find config file")
		}

		return nil, errors.Wrap(err, "failed to parse config")
	}

	httpTimeoutDuration, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("httpTimeout")))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http request timeout duration")
	}

	conf := &Config{
		logLevel:    viper.GetString("logLevel"),
		logType:     viper.GetString("logType"),
		packs:       viper.GetIntSlice("packs"),
		ServerPort:  viper.GetInt("serverPort"),
		HttpTimeout: httpTimeoutDuration,
	}

	if err := conf.initLogger(); err != nil {
		return nil, errors.Wrap(err, "failed to init logger")
	}

	logrus.WithFields(logrus.Fields{
		"logLevel":    conf.logLevel,
		"logType":     conf.logType,
		"packs":       conf.packs,
		"serverPort":  conf.ServerPort,
		"httpTimeout": conf.HttpTimeout,
	}).Info("parsed config")

	return conf, nil
}

func (c *Config) initLogger() error {
	level, err := logrus.ParseLevel(c.logLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	switch c.logType {
	// if we want to use ELK stack or similar this makes the log parsing easier
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		return fmt.Errorf("unrecognized log type: %s", c.logType)
	}

	logrus.SetOutput(os.Stdout)

	return nil
}

// SetPacks takes a slice of ints, that represent our package sizes. It locks the config to overwrite the current set.
func (c *Config) SetPacks(packs []int) []int {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.packs = packs
	sort.Ints(c.packs)

	return c.packs
}

// GetPacks returns the current set of packs.
func (c *Config) GetPacks() []int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.packs
}
