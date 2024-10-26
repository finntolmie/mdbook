package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToBlocks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "Basic h1",
			input: `This is **bolded** paragraph

This is another paragraph with *italic* text and ` + "`code`" + ` here
This is the same paragraph on a new line

* This is a list
* with items
`,
			expected: []string{
				"This is **bolded** paragraph",
				"This is another paragraph with *italic* text and `code` here\nThis is the same paragraph on a new line",
				"* This is a list\n* with items",
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, markdownToBlocks(tc.input))
	}
}

func TestBlockToBlockType(t *testing.T) {
	testCases := []struct {
		input    string
		expected BlockType
	}{
		{
			input:    "# Heading 1",
			expected: HEADING,
		},
		{
			input:    "###### Heading 6",
			expected: HEADING,
		},
		{
			input:    "####### Invalid Heading",
			expected: PARAGRAPH,
		},
		{
			input:    "#No space after hash",
			expected: PARAGRAPH,
		},
		{
			input:    "```\ncode block\n```",
			expected: CODEBLOCK,
		},
		{
			input:    "```single line code block```",
			expected: CODEBLOCK,
		},
		{
			input:    "```\nmultiline\ncode block\n```",
			expected: CODEBLOCK,
		},
		{
			input:    "```\nUnclosed code block",
			expected: PARAGRAPH,
		},
		{
			input:    "`Not enough backticks`",
			expected: PARAGRAPH,
		},
		{
			input:    "> Quote line 1\n> Quote line 2",
			expected: QUOTE,
		},
		{
			input:    "> Single line quote",
			expected: QUOTE,
		},
		{
			input:    "> Valid line\nInvalid line",
			expected: PARAGRAPH,
		},
		{
			input:    "No leading > in any line",
			expected: PARAGRAPH,
		},
		{
			input:    "* Item 1\n* Item 2",
			expected: UNORDERED_LIST,
		},
		{
			input:    "- Item 1\n- Item 2",
			expected: UNORDERED_LIST,
		},
		{
			input:    "* Mixed list\n- Works too",
			expected: UNORDERED_LIST,
		},
		{
			input:    "* Valid line\nInvalid line",
			expected: PARAGRAPH,
		},
		{
			input:    "This is not a list",
			expected: PARAGRAPH,
		},
		{
			input:    "1. First item\n2. Second item\n3. Third item",
			expected: ORDERED_LIST,
		},
		{
			input:    "1. Single item",
			expected: ORDERED_LIST,
		},
		{
			input:    "1. First item\n2. Second item\n4. Skips number",
			expected: PARAGRAPH,
		},
		{
			input:    "1. First item\n3. Skips second",
			expected: PARAGRAPH,
		},
		{
			input:    "This is a simple paragraph.",
			expected: PARAGRAPH,
		},
		{
			input:    "No special characters here.",
			expected: PARAGRAPH,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, blockToBlockType(tc.input))
	}
}

func TestMarkdownToHTMLNode(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input: `
This is **bolded** paragraph
text in a p
tag here

`,
			expected: "<div><p>This is <b>bolded</b> paragraph text in a p tag here</p></div>",
		},
		{
			input: `
This is **bolded** paragraph
text in a p
tag here

This is another paragraph with *italic* text and ` + "`code`" + ` here

`,
			expected: "<div><p>This is <b>bolded</b> paragraph text in a p tag here</p><p>This is another paragraph with <i>italic</i> text and <code>code</code> here</p></div>",
		},
		{
			input: `
- This is a list
- with items
- and *more* items

1. This is an ` + "`ordered`" + ` list
2. with items
3. and more items

`,
			expected: "<div><ul><li>This is a list</li><li>with items</li><li>and <i>more</i> items</li></ul><ol><li>This is an <code>ordered</code> list</li><li>with items</li><li>and more items</li></ol></div>",
		},
		{
			input: `
# this is an h1

this is paragraph text

## this is an h2
`,
			expected: "<div><h1>this is an h1</h1><p>this is paragraph text</p><h2>this is an h2</h2></div>",
		},
		{
			input: `
> This is a
> blockquote block

this is paragraph text

`,
			expected: "<div><blockquote>This is a blockquote block</blockquote><p>this is paragraph text</p></div>",
		},
	}

	for _, tc := range testCases {
		node, err := MarkdownToHTMLNode(tc.input)
		assert.Nil(t, err)
		html, err := node.toHTML()
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, html)
	}
}
