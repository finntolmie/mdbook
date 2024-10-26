package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/finntolmie/gofilib"
)

type BlockType int

const (
	PARAGRAPH BlockType = iota
	HEADING
	CODEBLOCK
	QUOTE
	UNORDERED_LIST
	ORDERED_LIST
)

func markdownToBlocks(markdown string) []string {
	return gofilib.Filter(gofilib.Map(strings.Split(markdown, "\n\n"), strings.TrimSpace), func(s string) bool {
		return s != ""
	})
}

func blockToBlockType(block string) BlockType {
	if match, _ := regexp.MatchString("^#{1,6} [ \\S]*$", block); match {
		return HEADING
	}
	if match, _ := regexp.MatchString("^```([\\s\\S]*?)```$", block); match {
		return CODEBLOCK
	}
	if match, _ := regexp.MatchString("^(>[ \\S]*\n?)+$", block); match {
		return QUOTE
	}
	if match, _ := regexp.MatchString("^([*-] [ \\S]*\n?)+$", block); match {
		return UNORDERED_LIST
	}
	re := regexp.MustCompile(`^(\d+)\. .*?$`)
	lines := strings.Split(block, "\n")
	for i := 1; i <= len(lines); i++ {
		matches := re.FindAllString(lines[i-1], -1)
		if len(matches) > 0 {
			listIdx, _ := strconv.Atoi(matches[0][0:strings.Index(matches[0], ".")])
			if listIdx == i {
				if i == len(lines) {
					return ORDERED_LIST
				}
				continue
			}
		} else {
			break
		}
	}
	return PARAGRAPH
}

func MarkdownToHTMLNode(markdown string) (HTMLNode, error) {
	children, err := gofilib.MapErr(markdownToBlocks(markdown), blockToHTMLNode)
	if err != nil {
		return nil, err
	}
	return ParentNode{
		Tag:      "div",
		Children: children,
	}, nil
}

func blockToHTMLNode(block string) (HTMLNode, error) {
	blockType := blockToBlockType(block)
	switch blockType {
	case PARAGRAPH:
		return paragraphToHTMLNode(block)
	case HEADING:
		return headerToHTMLNode(block)
	case CODEBLOCK:
		return codeToHTMLNode(block)
	case ORDERED_LIST:
		return orderedListToHTMLNode(block)
	case UNORDERED_LIST:
		return unorderedListToHTMLNode(block)
	case QUOTE:
		return quoteToHTMLNode(block)
	}
	return nil, errors.New("unrecognised block type")
}

func textToChildren(text string) ([]HTMLNode, error) {
	children, err := textToTextNodes(text)
	if err != nil {
		return nil, err
	}
	nodes := make([]HTMLNode, 0, len(children))
	for _, child := range children {
		node, err := child.toHTMLNode()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func paragraphToHTMLNode(block string) (HTMLNode, error) {
	children, err := textToChildren(strings.Join(strings.Split(block, "\n"), " "))
	if err != nil {
		return nil, err
	}
	return ParentNode{
		Tag:      "p",
		Children: children,
	}, nil
}

func headerToHTMLNode(block string) (HTMLNode, error) {
	re := regexp.MustCompile(`^(#{1,6})\s([ \S]*)$`)
	matches := re.FindStringSubmatch(block)

	if matches == nil || len(matches) < 3 {
		return nil, errors.New("no header match found")
	}
	headerLevel := len(matches[1])
	headerText := matches[2]
	children, err := textToChildren(headerText)
	if err != nil {
		return nil, err
	}
	return ParentNode{
		Tag:      fmt.Sprintf("h%d", headerLevel),
		Children: children,
	}, nil
}

func codeToHTMLNode(block string) (HTMLNode, error) {
	children, err := textToChildren(block[4 : len(block)-3])
	if err != nil {
		return nil, err
	}
	return ParentNode{
		Tag: "pre",
		Children: []HTMLNode{
			ParentNode{
				Tag:      "code",
				Children: children,
			},
		},
	}, nil
}

func unorderedListToHTMLNode(block string) (HTMLNode, error) {
	items := strings.Split(block, "\n")
	children := make([]HTMLNode, 0, len(items))
	for _, item := range items {
		itemNodes, err := textToChildren(item[2:])
		if err != nil {
			return nil, err
		}
		children = append(children, ParentNode{
			Tag:      "li",
			Children: itemNodes,
		})
	}
	return ParentNode{
		Tag:      "ul",
		Children: children,
	}, nil
}

func orderedListToHTMLNode(block string) (HTMLNode, error) {
	items := strings.Split(block, "\n")
	children := make([]HTMLNode, 0, len(items))
	for _, item := range items {
		itemNodes, err := textToChildren(item[3:])
		if err != nil {
			return nil, err
		}
		children = append(children, ParentNode{
			Tag:      "li",
			Children: itemNodes,
		})
	}
	return ParentNode{
		Tag:      "ol",
		Children: children,
	}, nil
}

func quoteToHTMLNode(block string) (HTMLNode, error) {
	children, err := textToChildren(strings.Join(gofilib.Map(strings.Split(block, "\n"), func(s string) string {
		return s[2:]
	}), " "))
	if err != nil {
		return nil, err
	}
	return ParentNode{
		Tag:      "blockquote",
		Children: children,
	}, nil
}
