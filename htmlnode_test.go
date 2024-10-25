package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeafNode(t *testing.T) {
	testCases := []struct {
		name     string
		input    LeafNode
		expected string
	}{
		{
			name: "Basic h1",
			input: LeafNode{
				Tag:   "h1",
				Value: "BIG HEADER",
			},
			expected: `<h1>BIG HEADER</h1>`,
		},
		{
			name: "Basic a with one prop",
			input: LeafNode{
				Tag:   "a",
				Value: "link to something",
				Props: map[string]string{
					"class": "big",
				},
			},
			expected: `<a class="big">link to something</a>`,
		},
	}

	for _, tc := range testCases {
		html, err := tc.input.toHTML()
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, html)
	}
}

func TestParentNode(t *testing.T) {
	testCases := []struct {
		name     string
		input    ParentNode
		expected string
	}{
		{
			name: "Basic div",
			input: ParentNode{
				Tag: "div",
			},
			expected: `<div></div>`,
		},
		{
			name: "Basic div with leaf node",
			input: ParentNode{
				Tag: "div",
				Children: []HTMLNode{
					LeafNode{
						Value: "something in the div",
					},
				},
			},
			expected: `<div>something in the div</div>`,
		},
		{
			name: "Nested div with leaf nodes",
			input: ParentNode{
				Tag: "div",
				Children: []HTMLNode{
					ParentNode{
						Tag: "div",
						Children: []HTMLNode{
							LeafNode{
								Tag:   "a",
								Value: "link",
							},
							LeafNode{
								Value: "text",
							},
						},
					},
				},
			},
			expected: `<div><div><a>link</a>text</div></div>`,
		},
	}

	for _, tc := range testCases {
		html, err := tc.input.toHTML()
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, html)
	}
}
