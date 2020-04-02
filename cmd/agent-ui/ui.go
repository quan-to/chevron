package main

import (
	"github.com/asticode/go-astikit"
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
				Label: astikit.StrPtr("File"),
				SubMenu: []*astilectron.MenuItemOptions{
					{Role: astilectron.MenuItemRoleClose},
				},
			},
			{
				Label: astikit.StrPtr("Tools"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Accelerator: astilectron.NewAccelerator("Alt", "CommandOrControl", "I"),
						Role:        astilectron.MenuItemRoleToggleDevTools,
					},
					{
						Label: astikit.StrPtr("Add Private Key"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							_ = w.SendMessage(bootstrap.MessageOut{
								Name: messageLoadPrivateKey,
							})

							return
						},
					},
				},
			},
			{
				Label: astikit.StrPtr("Edit"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label:       astikit.StrPtr("Undo"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "Z"),
						Role:        astikit.StrPtr("undo:"),
					},
					{
						Label:       astikit.StrPtr("Redo"),
						Accelerator: astilectron.NewAccelerator("Shift", "CmdOrCtrl", "Z"),
						Role:        astikit.StrPtr("redo"),
					},
					{
						Type: astikit.StrPtr("separator"),
					},
					{
						Label:       astikit.StrPtr("Cut"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "X"),
						Role:        astikit.StrPtr("cut"),
					},
					{
						Label:       astikit.StrPtr("Copy"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "C"),
						Role:        astikit.StrPtr("copy"),
					},
					{
						Label:       astikit.StrPtr("Paste"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "V"),
						Role:        astikit.StrPtr("paste"),
					},
					{
						Label:       astikit.StrPtr("Select All"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "A"),
						Role:        astikit.StrPtr("selectAll"),
					},
				},
			},
		},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			log.Info("Astilectron backend is ready")
			w = ws[0]
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(900),
				Width:           astikit.IntPtr(1600),
			},
		}},
	}

	if err := bootstrap.Run(bootConfig); err != nil {
		log.Error("%s: %s", err, "running bootstrap failed")
	}
}
