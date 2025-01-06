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
	v              *viper.Viper
	mu             sync.RWMutex
	AlertStatus    bool       `mapstructure:"alert_status"`
	DarkTheme      bool       `mapstructure:"dark_theme"`
	Keywords       []Keyword  `mapstructure:"keywords"`
	ChatFilters    []ChatType `mapstructure:"chat_filters"`
	AlertType      AlertType  `mapstructure:"alert_type"`
	TelegramChatID string     `mapstructure:"telegram_chat_id"`
}

type AlertType string

const (
	AlertTypeBoth       AlertType = "both"
	AlertTypeChatFilter AlertType = "chat_filter"
	AlertTypeKeyword    AlertType = "keyword"
)

type Keyword struct {
	Keyword string `mapstructure:"keyword"`
	Enable  bool   `mapstructure:"enable"`
}

func GetConfig() (*Config, error) {
	var initErr error
	once.Do(func() {
		config = &Config{
			v:  viper.New(),
			mu: sync.RWMutex{},
		}
		if err := config.autoInit(); err != nil {
			initErr = err
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return config, nil
}

func (c *Config) autoInit() error {
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		c.setDefaults()
		if err := c.Save(); err != nil {
			return fmt.Errorf("failed to save default config: %v", err)
		}
	}

	c.v.SetConfigType("json")
	c.v.SetConfigFile(CONFIG_PATH)

	if err := c.v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	if err := c.v.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return c.Validate()
}

func (c *Config) Validate() error {

	switch c.AlertType {
	case AlertTypeBoth, AlertTypeChatFilter, AlertTypeKeyword:
	default:
		return fmt.Errorf("invalid alert type: %s", c.AlertType)
	}

	return nil
}

func (c *Config) setDefaults() {
	c.defaultKeywords()
	c.defaultChatFilter()
	c.defaultAlertType()
}

func (c *Config) defaultKeywords() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Keywords = []Keyword{
		{Keyword: "สวัสดี เราต้องการซื้อ", Enable: true},
		// {Keyword: "(แท็บ", Enable: true},
		// {Keyword: "ตำแหน่ง: ซ้าย", Enable: true},
		{Keyword: "Hi, I would like to buy your", Enable: true},
		{Keyword: "Здравствуйте, хочу купить у вас", Enable: true},
		{Keyword: "안녕하세요", Enable: true},
		// {Keyword: `(stash tab "`, Enable: true},
		// {Keyword: `position: left`, Enable: true},
	}
}

func (c *Config) defaultChatFilter() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ChatFilters = []ChatType{
		ChatTypeWhisper,
	}
}

func (c *Config) defaultAlertType() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AlertType = AlertTypeBoth
}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.v.MergeConfigMap(map[string]interface{}{
		"keywords":         c.Keywords,
		"alert_status":     c.AlertStatus,
		"dark_theme":       c.DarkTheme,
		"chat_filters":     c.ChatFilters,
		"alert_type":       c.AlertType,
		"telegram_chat_id": c.TelegramChatID,
	}); err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	return c.v.WriteConfigAs(CONFIG_PATH)
}

func (c *Config) Update(updateFn func(*Config)) error {
	updateFn(c)
	return c.Save()
}
