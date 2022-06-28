package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/twang2218/fynebuilder"
	"github.com/twang2218/fynebuilder/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/fsnotify/fsnotify"
)

//go:embed assets/icon.png
var icon []byte
var resourceIcon = fyne.NewStaticResource("icon.png", icon)

//go:embed assets/idcard.jpg
var idcard []byte
var resourceJpegIdCard = fyne.NewStaticResource("idcard.jpg", idcard)

//go:embed assets/background.jpg
var background []byte
var resourceJpegBackground = fyne.NewStaticResource("background.jpg", background)

//go:embed assets/qrcode.png
var qrcode []byte
var resourcePngQRcode = fyne.NewStaticResource("qrcode.png", qrcode)

func monitor(file string, f func()) *fsnotify.Watcher {
	//	hot reload
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	watcher.Add(file)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Printf("%s %s\n", event.Name, event.Op)
				if event.Op == fsnotify.Write {
					//	update
					f()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	return watcher
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&theme.UnicodeTheme{})
	a.SetIcon(resourceIcon)
	//      窗口
	w := a.NewWindow("访客身份校验")

	var embedResourcesDict = map[string]*fyne.StaticResource{
		"idcard.jpg":     resourceJpegIdCard,
		"background.jpg": resourceJpegBackground,
		"qrcode.png":     resourcePngQRcode,
	}

	c := fynebuilder.Load("demo.ui", embedResourcesDict)
	w.SetContent(c)

	watcher := monitor("demo.ui", func() {
		t := time.Now()
		c := fynebuilder.Load("demo.ui", embedResourcesDict)
		w.SetContent(c)
		log.Printf("Reloaded %q in %v.", "demo.ui", time.Since(t))
	})

	w.ShowAndRun()
	watcher.Close()
}
