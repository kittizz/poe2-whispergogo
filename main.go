package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type Ntfy string

const (
	NTFY_BASE_URL      Ntfy = "https://ntfy.sh"
	NTFY_PREFIX_TOPICS Ntfy = "wpgogo"

	POE2_PROCESS_NAME = "PathOfExileSteam.exe" // Process name to search for

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgChan := make(chan string, 100) // buffer size 100

	pw := NewProcessWatcher(msgChan)
	if err := pw.WatchProcess(ctx); err != nil {
		fmt.Printf("Error watching process: %v\n", err)
		return
	}
	go func() {
		for logLine := range msgChan {
			msg, err := ParseChatMessage(logLine)
			if err != nil {
				fmt.Printf("Error parsing message: %v\n", err)
				continue
			}

			fmt.Printf("%v> Type: %s, User: %s, Content: %s\n",
				msg.Timestamp,
				msg.MessageType, msg.Username, msg.Content)
		}
	}()

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
		fmt.Println("Error:", err.Error())
	}
}
