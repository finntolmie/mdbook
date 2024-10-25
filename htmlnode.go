package main

import (
	"errors"
	"fmt"
	"strings"
)

type HTMLNode interface {
	toHTML() (string, error)
}

type LeafNode struct {
	Tag   string
	Value string
	Props map[string]string
}

func (ln LeafNode) toHTML() (string, error) {
	if ln.Tag == "" {
		return ln.Value, nil
	}
	return fmt.Sprintf("<%s%s>%s</%s>", ln.Tag, propsToHTML(ln.Props), ln.Value, ln.Tag), nil
}

type ParentNode struct {
	Tag      string
	Children []HTMLNode
	Props    map[string]string
}

func (pn ParentNode) toHTML() (string, error) {
	if pn.Tag == "" {
		return "", errors.New("parent nodes must have a tag")
	}
	elements := make([]string, 0, len(pn.Children))
	for _, el := range pn.Children {
		html, _ := el.toHTML()
		elements = append(elements, html)
	}
	return fmt.Sprintf("<%s%s>%s</%s>", pn.Tag, propsToHTML(pn.Props), strings.Join(elements, ""), pn.Tag), nil
}

func propsToHTML(props map[string]string) string {
	if len(props) == 0 {
		return ""
	}
	var sb strings.Builder
	elements := make([]string, 0, len(props))
	for key, val := range props {
		elements = append(elements, fmt.Sprintf("%s=\"%s\"", key, val))
	}
	sb.WriteString(" " + strings.Join(elements, " "))
	return sb.String()
}
