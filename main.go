package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tdewolff/minify/v2"
	minify_css "github.com/tdewolff/minify/v2/css"
	minify_html "github.com/tdewolff/minify/v2/html"
	"github.com/yuin/goldmark"
	gm_extension "github.com/yuin/goldmark/extension"
	gm_html "github.com/yuin/goldmark/renderer/html"
)

const (
	postFilesPath    = "posts"
	templatePostPath = "templates/post.html"

	templateCSSPath = "templates/styles.css"
	cssFilename     = "styles.css"

	generatedPath = "docs"

	extHTML = ".html"

	mimeTypeHTML = "text/html"
	mimeTypeCSS  = "text/css"
)

type templatePostData struct {
	Body string
}

func main() {
	// Cleanup generated path
	if err := os.RemoveAll(generatedPath); err != nil {
		log.Fatalln("Failed to remove all", generatedPath, err)
	}

	if err := os.MkdirAll(generatedPath, 0o777); err != nil {
		log.Fatalln("Failed to mkdir all", generatedPath)
	}

	// Read post files directory
	postFiles, err := os.ReadDir(postFilesPath)
	if err != nil {
		log.Fatalln("Failed to read dir", postFilesPath)
	}

	// Prepare template
	templatePostBytes, err := os.ReadFile(templatePostPath)
	if err != nil {
		log.Fatalln("Failed to read file", templatePostPath, err)
	}

	templatePost, err := template.New("post").Parse(string(templatePostBytes))
	if err != nil {
		log.Fatalln("Failed to parse template", err)
	}

	// Prepare parse markdown
	gm := goldmark.New(
		goldmark.WithExtensions(
			gm_extension.GFM,
		),
		goldmark.WithRendererOptions(
			gm_html.WithHardWraps(),
		),
	)

	// Prepare minify
	m := minify.New()
	m.AddFunc(mimeTypeHTML, minify_html.Minify)
	m.AddFunc(mimeTypeCSS, minify_css.Minify)

	// Generate post files
	for _, postFile := range postFiles {
		if postFile.IsDir() {
			continue
		}

		// Prepare post file
		mdFilename := filepath.Join(postFilesPath, postFile.Name())
		mdFileBytes, err := os.ReadFile(mdFilename)
		if err != nil {
			log.Fatalln("Failed to read file", mdFilename, err)
		}

		// Prepare html file
		htmlFilename := strings.TrimSuffix(postFile.Name(), filepath.Ext(postFile.Name())) + extHTML
		htmlFilepath := filepath.Join(generatedPath, htmlFilename)

		htmlFile, err := os.OpenFile(htmlFilepath, os.O_RDWR|os.O_CREATE, 0o600)
		if err != nil {
			log.Fatalln("Failed to open file", htmlFilepath, err)
		}

		// Parse markdown
		var markdownBuf bytes.Buffer
		if err := gm.Convert(mdFileBytes, &markdownBuf); err != nil {
			log.Fatalln("Failed to convert markdown", err)
		}

		tmpReader, tmpWriter := io.Pipe()

		// Template
		go func() {
			if err := templatePost.Execute(tmpWriter, templatePostData{
				Body: markdownBuf.String(),
			}); err != nil {
				log.Fatalln("Failed to execute html template", err)
			}
			tmpWriter.Close()
		}()

		// Minify
		if err := m.Minify(mimeTypeHTML, htmlFile, tmpReader); err != nil {
			log.Fatalln("Failed to minify html", err)
		}
		tmpReader.Close()
		htmlFile.Close()
	}

	// Copy css file
	templateCSSFile, err := os.OpenFile(templateCSSPath, os.O_RDONLY, 0o600)
	if err != nil {
		log.Fatalln("Failed to open file", templateCSSPath, err)
	}

	cssFilename := filepath.Join(generatedPath, cssFilename)
	cssFile, err := os.OpenFile(cssFilename, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		log.Fatalln("Failed to open file", cssFilename, err)
	}

	if err := m.Minify(mimeTypeCSS, cssFile, templateCSSFile); err != nil {
		log.Fatalln("Failed to minify css", err)
	}
}
