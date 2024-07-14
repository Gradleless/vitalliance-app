package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/build
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "Vitalliance",
		Width:         640,
		Height:        480,
		DisableResize: false,
		Windows: &windows.Options{
			Theme:             windows.Dark,
			DisableWindowIcon: false,
		},
		Linux: &linux.Options{
			Icon: icon,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Vitalliance",
				Icon:    icon,
				Message: "Made with love by @gradleless",
			},
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 28, G: 35, B: 42, A: 0},
		OnStartup:        app.startup,

		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
