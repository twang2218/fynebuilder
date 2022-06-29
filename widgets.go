package fynebuilder

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/image/colornames"
)

func newWidget(tag objectTag, res map[string]*fyne.StaticResource) fyne.CanvasObject {
	switch tag.Tag {
	case "Max":
		return container.NewMax()
	case "Center":
		return container.NewCenter()
	case "VBox":
		return container.NewVBox()
	case "HBox":
		return container.NewHBox()
	case "Image":
		width, height := tag.Attributes.GetFloat32("width"), tag.Attributes.GetFloat32("height")
		src := tag.Attributes.GetString("src")
		var imgRes *fyne.StaticResource
		if strings.HasPrefix(src, "embed:") {
			src := strings.TrimPrefix(src, "embed:")
			if r, ok := res[src]; ok {
				imgRes = r
			}
		}
		img := canvas.NewImageFromResource(imgRes)
		img.SetMinSize(fyne.NewSize(width, height))
		return img
	case "Label":
		return widget.NewLabel("")
	case "Text":
		var c color.Color
		name := strings.ToLower(tag.Attributes.GetString("color"))
		if item, ok := colornames.Map[name]; ok {
			c = item
		} else {
			c = color.Black
		}
		return canvas.NewText("", c)
	default:
		return container.NewWithoutLayout()
	}
}
