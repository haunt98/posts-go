package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

const (
	postFilesPath    = "posts"
	templatePostPath = "templates/post.html"

	generatedPath = "docs"

	extHTML = ".html"
)

type templatePostData struct {
	Body string
}

func main() {
	ctx := context.Background()

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

	// Prepare GitHub
	ghAccessTokenBytes, err := os.ReadFile(".github_access_token")
	if err != nil {
		log.Fatalln("Failed to read file", ".github_access_token", err)
	}

	ghTokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: string(ghAccessTokenBytes),
		},
	)
	ghHTTPClient := oauth2.NewClient(ctx, ghTokenSrc)
	ghClient := github.NewClient(ghHTTPClient)

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

		ghMarkdown, _, err := ghClient.Markdown(ctx, string(mdFileBytes), &github.MarkdownOptions{
			Mode: "markdown",
		})
		if err != nil {
			log.Fatalln("Failed to GitHub markdown", err)
		}

		// Prepare html file
		htmlFilename := strings.TrimSuffix(postFile.Name(), filepath.Ext(postFile.Name())) + extHTML
		htmlFilepath := filepath.Join(generatedPath, htmlFilename)

		htmlFile, err := os.OpenFile(htmlFilepath, os.O_RDWR|os.O_CREATE, 0o600)
		if err != nil {
			log.Fatalln("Failed to open file", htmlFilepath, err)
		}

		if err := templatePost.Execute(htmlFile, templatePostData{
			Body: ghMarkdown,
		}); err != nil {
			log.Fatalln("Failed to execute html template", err)
		}

		htmlFile.Close()
	}
}
