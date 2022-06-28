package fynebuilder

import (
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewWidget(tag string, attrs map[string]string, res map[string]*fyne.StaticResource) fyne.CanvasObject {
	switch tag {
	case "Max":
		return container.NewMax()
	case "Center":
		return container.NewCenter()
	case "VBox":
		return container.NewVBox()
	case "HBox":
		return container.NewHBox()
	case "Image":
		sw, sh, src := attrs["width"], attrs["height"], attrs["src"]
		width, err := strconv.Atoi(sw)
		if err != nil {
			log.Println(err)
		}
		height, err := strconv.Atoi(sh)
		if err != nil {
			log.Println(err)
		}
		var imgRes *fyne.StaticResource
		if strings.HasPrefix(src, "embed:") {
			src := strings.TrimPrefix(src, "embed:")
			if r, ok := res[src]; ok {
				imgRes = r
			}
		}
		img := canvas.NewImageFromResource(imgRes)
		img.SetMinSize(fyne.NewSize(float32(width), float32(height)))
		return img
	case "Label":
		return widget.NewLabel("")
	default:
		return container.NewWithoutLayout()
	}
}
