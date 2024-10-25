package main

import "errors"

type TextType int

const (
	TEXT TextType = iota
	BOLD
	ITALIC
	CODE
	LINK
	IMAGE
)

type TextNode struct {
	Text string
	Type TextType
	URL  string
}

func (tn TextNode) toHTMLNode() (HTMLNode, error) {
	switch tn.Type {
	case TEXT:
		return LeafNode{Value: tn.Text}, nil
	case BOLD:
		return LeafNode{Tag: "b", Value: tn.Text}, nil
	case ITALIC:
		return LeafNode{Tag: "i", Value: tn.Text}, nil
	case CODE:
		return LeafNode{Tag: "code", Value: tn.Text}, nil
	case LINK:
		return LeafNode{
			Tag:   "a",
			Value: tn.URL,
			Props: map[string]string{"href": tn.URL},
		}, nil
	case IMAGE:
		return LeafNode{
			Tag: "img",
			Props: map[string]string{
				"src": tn.URL,
				"alt": tn.Text,
			},
		}, nil
	default:
		return nil, errors.New("unrecognised text type")
	}
}
