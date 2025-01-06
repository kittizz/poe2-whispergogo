package main

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const (
	WHISPERGOGO_BASE_URL string = "https://whispergogo-server.xver.cloud"
)

var (
	POE2_PROCESS_NAMES = []string{
		"PathOfExileSteam.exe",
		"PathOfExile_KG.exe",
	}
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg, err := GetConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	gameStatus := new(bool)
	lineStream := make(chan string, 100) // buffer size 100

	pw := NewProcessWatcher(lineStream)
	if err := pw.WatchProcess(ctx, gameStatus); err != nil {
		fmt.Printf("Error watching process: %v\n", err)
		return
	}

	// Create an instance of the app structure
	app := NewApp(cfg, gameStatus, lineStream)

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "WhisperGoGo | Path Of Exile 2",
		MinWidth:  1300,
		MinHeight: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "b6678fc0-9be3-4d93-895d-200f2fe3b1c1",
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
		OnShutdown: func(ctx context.Context) {

			// เพิ่ม timeout สำหรับการรอ goroutine
			done := make(chan struct{})
			go func() {
				app.shutdown(ctx)
				cancel()
				close(lineStream)

				close(done)
			}()

			select {
			case <-done:
				fmt.Println("Shutdown complete")
			case <-time.After(5 * time.Second):
				fmt.Println("Shutdown timeout after 5 seconds")
			}
		},
		OnStartup: app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
