package fynebuilder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cast"
)

type attributes map[string]string

func (a *attributes) GetString(key string) string {
	return (*a)[key]
}

func (a *attributes) GetInt(key string) int {
	return cast.ToInt((*a)[key])
}

func (a *attributes) GetFloat32(key string) float32 {
	return cast.ToFloat32((*a)[key])
}

func (a *attributes) GetFloat64(key string) float64 {
	return cast.ToFloat64((*a)[key])
}

type ObjectSet map[string]fyne.CanvasObject

func (s ObjectSet) Get(key string) fyne.CanvasObject {
	return s[key]
}

const key_top = "_top_"

func (s ObjectSet) GetTop() fyne.CanvasObject {
	return s[key_top]
}

type objectTag struct {
	Tag        string
	Attributes attributes
	Object     fyne.CanvasObject
}

func Load(file string, res map[string]*fyne.StaticResource) (ObjectSet, error) {
	// open xml file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	log.Tracef("Opened XML %q file", file)
	//	create decoder for XML
	d := xml.NewDecoder(f)

	//	map which contains interested objects
	objs := make(ObjectSet)

	//	stack is used to preserve the parent path
	var stack []objectTag

	//	go through the XML elements
	for {
		//	get next element
		t, err := d.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return objs, err
		}

		switch t := t.(type) {
		case xml.StartElement:
			//	open tag
			var current objectTag
			current.Tag = t.Name.Local
			current.Attributes = get_attributes(t.Attr)
			current.Object = newWidget(current, res)
			log.Tracef("<%s>", current.Tag)
			//	preserve the TOP object in final result
			if len(objs) == 0 {
				objs[key_top] = current.Object
			}
			//	preserve the object with ID in final result
			if id, ok := current.Attributes["id"]; ok {
				objs[id] = current.Object
			}
			//	add current object to the parent container if the parent is a container
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				if c, ok := parent.Object.(*fyne.Container); ok {
					c.Add(current.Object)
					// fmt.Printf("Add %s to %s\n", current.Tag, parent.Tag)
				}
			}
			//	push the current tag into the stack
			stack = append(stack, current)
		case xml.EndElement:
			//	close tag
			current := stack[len(stack)-1]
			if current.Tag != t.Name.Local {
				log.Errorf("tag not match: expected: <%s>, actual: <%s>", current.Tag, t.Name.Local)
			} else {
				log.Tracef("</%s>", current.Tag)
				//	pop the last tag from the stack
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			//	content
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				switch o := current.Object.(type) {
				case *widget.Label:
					//	<Label>xxxxxxx</Label>
					o.SetText(string(t))
					log.Tracef("<Label>.SetText(%s)", t)
				case *canvas.Text:
					o.Text = string(t)
					log.Tracef("<Text>.Text = %s", t)
				}
			}
		}
	}

	//	validate
	if len(stack) > 0 {
		return objs, fmt.Errorf("XML file is corrupted. %v", stack)
	}

	return objs, nil
}

func get_attributes(attrs []xml.Attr) attributes {
	m := make(attributes, len(attrs))
	for _, attr := range attrs {
		m[attr.Name.Local] = attr.Value
	}
	return m
}
