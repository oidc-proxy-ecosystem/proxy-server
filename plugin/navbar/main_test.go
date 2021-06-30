package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCSS(t *testing.T) {
	expected := `<html><head><title>Test</title><link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet"/></head><body><p>Test</p></body></html>`
	html := `<html><head><title>Test</title></head><body><p>Test</p></body></html>`
	// html = setCSS([]byte(html))
	assert.Equal(t, expected, html)
}

func TestSetNavBar(t *testing.T) {
	// expected := `<html><head><title>Test</title><link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet"/></head><body><p>Test</p></body></html>`
	html := `<html><head><title>Test</title></head><body><p>Test</p></body></html>`
	// html = setNav([]byte(html))
	println(html)
	// assert.Equal(t, expected, html)
}
