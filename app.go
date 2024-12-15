package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	config *Config
}

// NewApp creates a new App application struct
func NewApp(config *Config) *App {
	return &App{
		config: config,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) onSecondInstanceLaunch(options options.SecondInstanceData) {
	a.BringToFront()
}
func (a *App) BringToFront() {
	runtime.WindowShow(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}

// GetKeywords returns list of keywords from config
func (a *App) GetKeywords() []Keyword {
	return a.config.Keywords
}

// SetKeywords sets list of keywords in config
func (a *App) SetKeywords(keywords []Keyword) {
	a.config.Keywords = keywords
	a.config.Save()
}

// GetNtfyTopics returns ntfy topics from config
func (a *App) GetNtfyTopics() string {
	return a.config.NtfyTopics
}

// SetNtfyTopics sets ntfy topics in config
func (a *App) SetNtfyTopics(topics string) {
	a.config.NtfyTopics = topics
	a.config.Save()
}

func (a *App) ResetNtfyTopics() {
	a.config.setDefaults()
	a.config.Save()
}

func (a *App) GetDeviceName() string {
	return getDeviceName()
}
