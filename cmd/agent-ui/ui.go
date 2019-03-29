package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

func Run() {
	bootConfig := bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug: *debug,
		MenuOptions: []*astilectron.MenuItemOptions{
			{
				Label: astilectron.PtrStr("File"),
				SubMenu: []*astilectron.MenuItemOptions{
					{Role: astilectron.MenuItemRoleClose},
				},
			},
			{
				Label: astilectron.PtrStr("Tools"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Accelerator: astilectron.NewAccelerator("Alt", "CommandOrControl", "I"),
						Role:        astilectron.MenuItemRoleToggleDevTools,
					},
				},
			},
		},
		OnWait: func(a *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			log.Info("Astilectron backend is ready")
			w = ws[0]
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(1600),
				Width:           astilectron.PtrInt(900),
			},
		}},
	}

	if err := bootstrap.Run(bootConfig); err != nil {
		log.Error("%s: %s", err, "running bootstrap failed")
	}
}
