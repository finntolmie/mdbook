package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMDImages(t *testing.T) {
	testCases := []struct {
		input    string
		expected []MDImage
	}{
		{
			input: "This is text with an ![image](https://i.imgur.com/zjjcJKZ.png)",
			expected: []MDImage{
				{
					Alt: "image",
					URL: "https://i.imgur.com/zjjcJKZ.png",
				},
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, extractMarkdownImages(tc.input))
	}
}

func TestExtractMDLinks(t *testing.T) {
	testCases := []struct {
		input    string
		expected []MDLink
	}{
		{
			input: "This is text with a [link](https://boot.dev) and [another link](https://blog.boot.dev)",
			expected: []MDLink{
				{
					Alt: "link",
					URL: "https://boot.dev",
				},
				{
					Alt: "another link",
					URL: "https://blog.boot.dev",
				},
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, extractMarkdownLinks(tc.input))
	}
}

func TestTextToTextNodes(t *testing.T) {
	testCases := []struct {
		input    string
		expected []TextNode
	}{
		{
			input: "This is **text** with an *italic* word and a `code block` and an ![obi wan image](https://i.imgur.com/fJRm4Vk.jpeg) and a [link](https://boot.dev)",
			expected: []TextNode{
				{
					"This is ", TEXT, "",
				},
				{
					"text", BOLD, "",
				},
				{
					" with an ", TEXT, "",
				},
				{
					"italic", ITALIC, "",
				},
				{
					" word and a ", TEXT, "",
				},
				{
					"code block", CODE, "",
				},
				{
					" and an ", TEXT, "",
				},
				{
					"obi wan image", IMAGE, "https://i.imgur.com/fJRm4Vk.jpeg",
				},
				{
					" and a ", TEXT, "",
				},
				{
					"link", LINK, "https://boot.dev",
				},
			},
		},
	}

	for _, tc := range testCases {
		textNodes, err := textToTextNodes(tc.input)
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, textNodes)
	}
}
