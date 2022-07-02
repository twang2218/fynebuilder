package main

import (
	_ "embed"
	"flag"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
	"github.com/twang2218/fynebuilder"
	"github.com/twang2218/fynebuilder/theme"
)

//go:embed assets/icon.png
var icon []byte

var verbose = flag.Bool("v", false, "verbose")
var min_size = fyne.NewSize(800, 600)

func main() {
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	a := app.New()
	a.Settings().SetTheme(&theme.UnicodeTheme{})

	var res = fynebuilder.ResourceDict{
		"icon.png": fyne.NewStaticResource("icon.png", icon),
	}

	a.SetIcon(res["icon.png"])

	//      窗口
	w := a.NewWindow("GUI Builder")

	renderUI := func(c *fyne.Container, xml string) {
		if len(xml) == 0 {
			return
		}
		objs, err := fynebuilder.LoadFromString([]byte(xml), res)
		if err != nil {
			log.Error(err)
		} else if objs.GetTop() != nil {
			c.RemoveAll()
			c.Add(objs.GetTop())
		}
	}

	watcher := fynebuilder.NewWatcher("main.ui", res, func(objs fynebuilder.ObjectDict) {
		if objs.GetTop() == nil {
			return
		}

		e, ok := objs.Get("code").(*widget.Entry)
		if !ok {
			return
		}
		c, ok := objs.Get("visual").(*fyne.Container)
		if !ok {
			return
		}

		renderUI(c, e.Text)
		w.Resize(min_size)
		w.SetContent(objs.GetTop())

		e.OnChanged = func(content string) {
			renderUI(c, content)
			w.Resize(min_size)
		}
	})
	defer watcher.Close()

	w.Resize(min_size)

	w.ShowAndRun()
}
