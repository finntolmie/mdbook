package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func splitNodesDelimiter(oldNodes []TextNode, delim string, textType TextType) ([]TextNode, error) {
	res := make([]TextNode, 0)
	for _, node := range oldNodes {
		if node.Type != TEXT {
			res = append(res, node)
			continue
		}
		splitNodes := strings.Split(node.Text, delim)
		if len(splitNodes)%2 == 0 {
			return nil, errors.New("invalid markdown, formatted section not closed")
		}
		for i, splitNode := range splitNodes {
			if splitNode == "" {
				continue
			}
			newNode := TextNode{
				Text: splitNode,
			}
			if i%2 == 0 {
				newNode.Type = TEXT
			} else {
				newNode.Type = textType
			}
			res = append(res, newNode)
		}
	}
	return res, nil
}

type MDImage struct {
	Alt string
	URL string
}

func extractMarkdownImages(text string) []MDImage {
	re := regexp.MustCompile(`!\[(?P<text>.*?)\]\((?P<url>.*?)\)`)
	matches := re.FindAllStringSubmatch(text, -1)
	images := make([]MDImage, len(matches))
	for i, match := range matches {
		images[i] = MDImage{Alt: match[1], URL: match[2]}
	}
	return images
}

type MDLink struct {
	Alt string
	URL string
}

func extractMarkdownLinks(text string) []MDLink {
	re := regexp.MustCompile(`\[(?P<alt>.*?)\]\((?P<url>.*?)\)`)
	matches := re.FindAllStringSubmatch(text, -1)
	images := make([]MDLink, len(matches))
	for i, match := range matches {
		images[i] = MDLink{Alt: match[1], URL: match[2]}
	}
	return images
}

func splitNodesImage(oldNodes []TextNode) ([]TextNode, error) {
	res := make([]TextNode, 0)
	for _, node := range oldNodes {
		if node.Type != TEXT {
			res = append(res, node)
			continue
		}
		original := node.Text
		imageMatches := extractMarkdownImages(original)
		if len(imageMatches) == 0 {
			res = append(res, node)
			continue
		}
		for _, image := range imageMatches {
			originalSplit := strings.Split(original, fmt.Sprintf("![%s](%s)", image.Alt, image.URL))
			if len(originalSplit) != 2 {
				return nil, errors.New("invalid markdown, image section not closed")
			}
			if originalSplit[0] != "" {
				res = append(res, TextNode{
					Text: originalSplit[0],
					Type: TEXT,
				})
			}
			res = append(res, TextNode{
				Text: image.Alt,
				Type: IMAGE,
				URL:  image.URL,
			})
			original = originalSplit[1]
		}
		if original != "" {
			res = append(res, TextNode{
				Text: original,
				Type: TEXT,
			})
		}
	}
	return res, nil
}

func splitNodesLink(oldNodes []TextNode) ([]TextNode, error) {
	res := make([]TextNode, 0)
	for _, node := range oldNodes {
		if node.Type != TEXT {
			res = append(res, node)
			continue
		}
		original := node.Text
		linkMatches := extractMarkdownLinks(original)
		if len(linkMatches) == 0 {
			res = append(res, node)
			continue
		}
		for _, link := range linkMatches {
			originalSplit := strings.Split(original, fmt.Sprintf("[%s](%s)", link.Alt, link.URL))
			if len(originalSplit) != 2 {
				return nil, errors.New("invalid markdown, link section not closed")
			}
			if originalSplit[0] != "" {
				res = append(res, TextNode{
					Text: originalSplit[0],
					Type: TEXT,
				})
			}
			res = append(res, TextNode{
				Text: link.Alt,
				Type: LINK,
				URL:  link.URL,
			})
			original = originalSplit[1]
		}
		if original != "" {
			res = append(res, TextNode{
				Text: original,
				Type: TEXT,
			})
		}
	}
	return res, nil
}

func textToTextNodes(text string) ([]TextNode, error) {
	nodes := []TextNode{{
		Text: text,
		Type: TEXT,
	}}
	boldNodes, err := splitNodesDelimiter(nodes, "**", BOLD)
	if err != nil {
		return nil, err
	}

	italicNodes, err := splitNodesDelimiter(boldNodes, "*", ITALIC)
	if err != nil {
		return nil, err
	}
	codeNodes, err := splitNodesDelimiter(italicNodes, "`", CODE)
	if err != nil {
		return nil, err
	}
	imageNodes, err := splitNodesImage(codeNodes)
	if err != nil {
		return nil, err
	}
	linkNodes, err := splitNodesLink(imageNodes)
	if err != nil {
		return nil, err
	}
	return linkNodes, nil
}
