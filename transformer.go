package xmler

import (
	"bytes"
	"container/list"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

type Transformer struct {
	elements Elements
}

func NewTransformer() Transformer {
	return Transformer{
		elements: NewElements(),
	}
}

func (self *Transformer) Parse(filePath string) Elements {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	d := xml.NewDecoder(bytes.NewReader(data))

	var lastParent *list.Element

	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}

		switch token := t.(type) {
		case xml.StartElement:
			e := &Element{
				Type: token.Name.Local,
			}

			for _, a := range token.Attr {
				e.addAttr(a.Name.Local, a.Value)
			}

			if lastParent != nil {
				e.Parent = lastParent.Value.(*Element)
			}
			lastParent = self.elements.PushBack(e)

		case xml.EndElement:
			lastParent = lastParent.Prev()

		case xml.CharData:
			// adds tag content from <tag>conent</tag> to element.Value
			value := strings.TrimSpace(string(token[:]))
			if value != "" {
				e := self.elements.Back().Value.(*Element)
				e.Value = value
			}
		}

	}

	return self.elements
}

func (self *Transformer) Transform(xmlFilePath, tplFilePath string) bytes.Buffer {

	elements := self.Parse(xmlFilePath)

	tpl, err := template.ParseFiles(tplFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, elements.Slice()); err != nil {
		log.Fatalf("generating code: %v", err)
	}

	return buf
}
