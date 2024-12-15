package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type Ntfy string

const (
	NTFY_BASE_URL      Ntfy = "https://ntfy.sh"
	NTFY_PREFIX_TOPICS Ntfy = "wpgogo"
)

var (
	AllNtfys = []struct {
		Value  Ntfy
		TSName string
	}{
		{NTFY_BASE_URL, "NTFY_BASE_URL"},
		{NTFY_PREFIX_TOPICS, "NTFY_PREFIX_TOPICS"},
	}
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := GetConfig()

	// Create an instance of the app structure
	app := NewApp(cfg)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "WhisperGoGo | Path Of Exile 2",
		Width:  1300,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DisableResize:    true,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "b6678fc0-9be3-4d93-895d-200f2fe3b1c1",
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
		OnStartup: app.startup,
		Bind: []any{
			app,
		},
		EnumBind: []any{
			AllNtfys,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
