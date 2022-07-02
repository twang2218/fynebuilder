package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/twang2218/fynebuilder"
	"github.com/twang2218/fynebuilder/theme"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/idcard.jpg
var idcard []byte

//go:embed assets/background.jpg
var background []byte

//go:embed assets/qrcode.png
var qrcode []byte

func main() {
	a := app.New()
	a.Settings().SetTheme(&theme.UnicodeTheme{})

	var res = fynebuilder.ResourceDict{
		"icon.png":       fyne.NewStaticResource("icon.png", icon),
		"idcard.jpg":     fyne.NewStaticResource("idcard.jpg", idcard),
		"background.jpg": fyne.NewStaticResource("background.jpg", background),
		"qrcode.png":     fyne.NewStaticResource("qrcode.png", qrcode),
	}

	a.SetIcon(res["icon.png"])

	//      窗口
	w := a.NewWindow("访客身份校验")

	watcher := fynebuilder.NewWatcher("demo.ui", res, func(objs fynebuilder.ObjectDict) {
		w.SetContent(objs.GetTop())
	})
	defer watcher.Close()

	w.ShowAndRun()
}
