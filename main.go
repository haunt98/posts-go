package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

const (
	postsPath     = "posts"
	headHTMLPath  = "custom/head.html"
	generatedPath = "docs"
	htmlExt       = ".html"
)

func main() {
	// Cleanup generated path
	if err := os.RemoveAll(generatedPath); err != nil {
		log.Fatalln("Failed to remove all", generatedPath, err)
	}

	if err := os.MkdirAll(generatedPath, 0o777); err != nil {
		log.Fatalln("Failed to mkdir all", generatedPath)
	}

	// Read needed files
	headHTML, err := os.ReadFile(headHTMLPath)
	if err != nil {
		log.Fatalln("Failed to read file", headHTML)
	}

	files, err := os.ReadDir(postsPath)
	if err != nil {
		log.Fatalln("Failed to read dir", postsPath)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Generate HTML
		filePath := filepath.Join(postsPath, file.Name())
		md, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln("Failed to read file", filePath)
		}

		htmlFlags := html.CommonFlags |
			html.CompletePage |
			html.TOC |
			html.LazyLoadImages
		htmlRendererOtps := html.RendererOptions{
			Title: file.Name(),
			Head:  headHTML,
			Flags: htmlFlags,
		}

		htmlRenderer := html.NewRenderer(htmlRendererOtps)
		generatedHTML := markdown.ToHTML(md, nil, htmlRenderer)

		generatedFileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + htmlExt
		generatedFilePath := filepath.Join(generatedPath, generatedFileName)

		if err := os.WriteFile(generatedFilePath, generatedHTML, 0o666); err != nil {
			log.Fatalln("Failed to write file", generatedFilePath, err)
		}
	}
}
