package fynebuilder

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
)

type ResourceDict map[string]fyne.Resource

type ObjectDict map[string]fyne.CanvasObject

func (s ObjectDict) Get(key string) fyne.CanvasObject {
	return s[key]
}

const key_top = "_top_"

func (s ObjectDict) GetTop() fyne.CanvasObject {
	return s[key_top]
}

type ObjectTag struct {
	Tag        string
	Attributes Attributes
	Object     fyne.CanvasObject
}

func LoadFromFile(file string, res ResourceDict) (ObjectDict, error) {
	// open xml file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return LoadFromString(content, res)
}

func LoadFromString(content []byte, res ResourceDict) (ObjectDict, error) {
	//	validate xml
	if err := xml.Unmarshal(content, new(interface{})); err != nil {
		return nil, err
	}

	//	create decoder for XML
	d := xml.NewDecoder(bytes.NewReader(content))

	//	map which contains interested objects
	objs := make(ObjectDict)

	//	stack is used to preserve the parent path
	var stack []ObjectTag

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
			var current = ObjectTag{
				Tag:        t.Name.Local,
				Attributes: NewAttributesFromXML(t.Attr),
			}
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
			//	attach current object to the parent container if the parent is a container
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				var grand ObjectTag
				if len(stack) > 1 {
					grand = stack[len(stack)-2]
				}
				attach(current, parent, grand, res)
			}
			//	push the current tag into the stack
			stack = append(stack, current)
		case xml.EndElement:
			//	close tag
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				if current.Tag != t.Name.Local {
					log.Errorf("tag not match: expected: <%s>, actual: <%s>", current.Tag, t.Name.Local)
				} else {
					log.Tracef("</%s>", current.Tag)
					//	pop the last tag from the stack
					stack = stack[:len(stack)-1]
				}
			}
		case xml.CharData:
			//	set content
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				setContent(current, string(t))
			}
		}
	}

	//	validate
	if len(stack) > 0 {
		return objs, fmt.Errorf("XML file is corrupted. %v", stack)
	}

	return objs, nil
}
