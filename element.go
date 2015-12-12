package xmler

import "container/list"

type Attr struct {
	Name  string
	Value string
}

type Element struct {
	Id     string
	Type   string
	Value  string
	Parent *Element
	Attrs  []Attr
}

func (self *Element) addAttr(name, value string) {
	if name == "id" {
		self.Id = value
	} else {
		self.Attrs = append(
			self.Attrs,
			Attr{
				Name:  name,
				Value: value,
			},
		)
	}
}

func (self *Element) Name() string {
	if self.Id != "" {
		return self.Id
	}

	return self.Type
}

func (self *Element) IdentifierName() string {
	instanceName := self.Name()
	if self.Id == "" && self.Parent != nil {
		return self.Parent.Name() + instanceName
	}

	return instanceName
}

type Elements struct {
	*list.List
}

func NewElements() Elements {
	return Elements{
		list.New(),
	}
}

func (self *Elements) Slice() []*Element {
	var elements []*Element
	for e := self.Front(); e != nil; e = e.Next() {
		elements = append(elements, e.Value.(*Element))
	}
	return elements
}
