package fynebuilder

import (
	"encoding/xml"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TagObject struct {
	Tag    string
	Object fyne.CanvasObject
}

func Load(file string, res map[string]*fyne.StaticResource) fyne.CanvasObject {
	// Open our xmlFile
	xmlFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}

	// log.Printf("Successfully Opened %s", file)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	dec := xml.NewDecoder(xmlFile)

	var stack []TagObject
	var top fyne.CanvasObject
	dict := make(map[string]fyne.CanvasObject)

	var current fyne.CanvasObject
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("xmlselect: %v", err)
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			//	开始
			// fmt.Printf("[Start] %+v\n", tok)

			attrs := xml_to_attrs(tok.Attr)
			current = NewWidget(tok.Name.Local, attrs, res)
			if top == nil {
				top = current
			}
			if id, ok := attrs["id"]; ok {
				dict[id] = current
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				if c, ok := parent.Object.(*fyne.Container); ok {
					c.Add(current)
					// fmt.Printf("Add %s to %s\n", tok.Name.Local, parent.Tag)
				}
			}
			stack = append(stack, TagObject{Tag: tok.Name.Local, Object: current})
		case xml.EndElement:
			//	结束
			stack = stack[:len(stack)-1]
			// fmt.Printf("[Stop]  %+v\n", tok)
		case xml.CharData:
			//	内容
			if len(stack) > 1 {
				item := stack[len(stack)-1]
				switch item.Tag {
				case "Label":
					if label, ok := current.(*widget.Label); ok {
						label.SetText(string(tok))
						// fmt.Printf("<%s> SetText(): %q\n", item.Tag, tok)
					}
					// fmt.Printf("<%s>: %q\n", item.Tag, tok)
				}
			}
		}
	}

	// for id, c := range dict {
	// 	fmt.Printf("[%s]: \t %+v\n", id, c)
	// }

	return top
}

func xml_to_attrs(attrs []xml.Attr) map[string]string {
	m := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		m[attr.Name.Local] = attr.Value
	}
	return m
}
