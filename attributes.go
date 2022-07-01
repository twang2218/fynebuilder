package fynebuilder

import (
	"encoding/xml"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/image/colornames"
)

type Attributes map[string]string

func NewAttributesFromXML(attrs []xml.Attr) Attributes {
	m := make(Attributes, len(attrs))
	for _, attr := range attrs {
		m[attr.Name.Local] = attr.Value
	}
	return m
}

func (a *Attributes) Contains(key string) bool {
	_, ok := (*a)[key]
	return ok
}

func (a *Attributes) ContainsAny(keys []string) bool {
	for _, k := range keys {
		if _, ok := (*a)[k]; ok {
			return true
		}
	}
	return false
}

func (a *Attributes) ContainsStyle() bool {
	return a.ContainsAny([]string{"bold", "italic", "monospace", "symbol", "tabwidth"})
}

func (a *Attributes) GetString(key string) string {
	return (*a)[key]
}

func (a *Attributes) GetList(key string) []string {
	content := (*a)[key]
	l := strings.Split(content, ",")
	for i, item := range l {
		l[i] = strings.TrimSpace(item)
	}
	return l
}

func (a *Attributes) GetBool(key string) bool {
	return cast.ToBool((*a)[key])
}

func (a *Attributes) GetInt(key string) int {
	return cast.ToInt((*a)[key])
}

func (a *Attributes) GetFloat32(key string) float32 {
	return cast.ToFloat32((*a)[key])
}

func (a *Attributes) GetFloat64(key string) float64 {
	return cast.ToFloat64((*a)[key])
}

func (a *Attributes) GetImage(key string, res ResourceDict) *canvas.Image {
	r := a.GetResource(key, res)
	if r != nil {
		return canvas.NewImageFromResource(r)
	}
	return nil
}

func (a *Attributes) GetResource(key string, res ResourceDict) fyne.Resource {
	src := a.GetString(key)

	if strings.HasPrefix(src, "embed:") {
		if res != nil {
			src := strings.TrimPrefix(src, "embed:")
			if r, ok := res[src]; ok {
				return r
			} else {
				return nil
			}
		}
	} else if strings.HasPrefix(src, "http") {
		r, err := fyne.LoadResourceFromURLString(src)
		if err != nil {
			log.Error(err)
		}
		return r
	} else if len(src) > 0 {
		r, err := fyne.LoadResourceFromPath(src)
		if err != nil {
			log.Error(err)
		}
		return r
	}

	return nil
}

func (a *Attributes) GetColor(key string) color.Color {
	name := strings.ToLower(a.GetString(key))
	if item, ok := colornames.Map[name]; ok {
		return item
	}
	return color.Black
}

func (a *Attributes) GetTextStyle() fyne.TextStyle {
	t := fyne.TextStyle{}
	t.Bold = a.GetBool("bold")
	t.Italic = a.GetBool("italic")
	t.Monospace = a.GetBool("monospace")
	t.Symbol = a.GetBool("symbol")
	t.TabWidth = a.GetInt("tabwidth")
	return t
}

func (a *Attributes) GetTextAlignment() fyne.TextAlign {
	switch strings.ToLower(a.GetString("align")) {
	default:
	case "center":
		return fyne.TextAlignCenter
	case "leading":
		return fyne.TextAlignLeading
	case "trailing":
		return fyne.TextAlignTrailing
	}
	return fyne.TextAlignCenter
}

func (a *Attributes) GetSize() fyne.Size {
	width := a.GetFloat32("width")
	height := a.GetFloat32("height")
	return fyne.NewSize(width, height)
}

func (a *Attributes) GetPosition() fyne.Position {
	x, y := a.GetFloat32("x"), a.GetFloat32("y")
	return fyne.NewPos(x, y)
}
