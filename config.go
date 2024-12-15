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
	v           *viper.Viper
	mu          sync.RWMutex
	AlertStatus bool `mapstructure:"alert_status"`
	DarkTheme   bool `mapstructure:"dark_theme"`

	NtfyTopics string    `mapstructure:"ntfy_topics"`
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
	c.defaultNtfyTopics()
	c.defaultKeywords()

}
func (c *Config) defaultNtfyTopics() {
	c.NtfyTopics = getDeviceName()
}
func (c *Config) defaultKeywords() {
	c.Keywords = []Keyword{
		{Keyword: "keyword1", Enable: true},
		{Keyword: "keyword2", Enable: true},
		// Add more default keywords as needed
	}
}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.v.MergeConfigMap(map[string]interface{}{
		"ntfy_topics":  c.NtfyTopics,
		"keywords":     c.Keywords,
		"alert_status": c.AlertStatus,
		"dark_theme":   c.DarkTheme,
	}); err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	return c.v.WriteConfigAs(CONFIG_PATH)
}

func (c *Config) Update(updateFn func(*Config)) error {

	updateFn(c)
	return c.Save()
}
