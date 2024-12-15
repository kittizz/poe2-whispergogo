package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

const CONFIG_PATH = "whispergogo.json"

var (
	config *Config
	once   sync.Once
)

type Config struct {
	v          *viper.Viper
	mu         sync.RWMutex
	NtfyTopics string    `mapstructure:"ntfy.topics"`
	Keywords   []Keyword `mapstructure:"keywords"`
}
type Keyword struct {
	Keyword string `mapstructure:"keyword"`
	Enable  bool   `mapstructure:"enable"`
}

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			v:  viper.New(),
			mu: sync.RWMutex{},
		}
		config.autoInit()
	})
	return config
}

func (c *Config) autoInit() {
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		c.setDefaults()
		if err := c.Save(); err != nil {
			panic(fmt.Sprintf("Failed to save default config: %v", err))
		}
	}

	c.v.SetConfigType("json")
	c.v.SetConfigFile(CONFIG_PATH)

	if err := c.v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Failed to read config: %v", err))
	}

	if err := c.v.Unmarshal(c); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config: %v", err))
	}
}

func (c *Config) setDefaults() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Failed to get hostname: %v", err))
	}
	c.NtfyTopics = hostname

}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.v.MergeConfigMap(map[string]interface{}{
		"ntfy.topics": c.NtfyTopics,
		"keywords":    c.Keywords,
	}); err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	return c.v.WriteConfigAs(CONFIG_PATH)
}

func (c *Config) Update(updateFn func(*Config)) error {

	updateFn(c)
	return c.Save()
}
