package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-github/v78/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/make-go-great/netrc-go"
)

const (
	postFilesPath    = "posts"
	templatePostPath = "templates/post.html"

	generatedPath = "docs"

	extHTML = ".html"

	netrcPath          = "~/.netrc"
	netrcMachineGitHub = "github.com"
)

type templatePostData struct {
	Index string
	Body  string
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
	netrcData, err := netrc.ParseFile(netrcPath)
	if err != nil {
		log.Fatalln("Failed to read file netrc", err)
	}

	var ghAccessToken string
	for _, machine := range netrcData.Machines {
		if machine.Name == netrcMachineGitHub {
			ghAccessToken = machine.Password
			break
		}
	}
	if ghAccessToken == "" {
		log.Fatalln("Empty GitHub token in netrc")
	}

	ghTokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: strings.TrimSpace(ghAccessToken),
		},
	)
	ghHTTPClient := oauth2.NewClient(ctx, ghTokenSrc)
	ghClient := github.NewClient(ghHTTPClient)

	eg := new(errgroup.Group)

	// Generate post files
	for _, postFile := range postFiles {
		if postFile.IsDir() {
			continue
		}

		postFilename := postFile.Name()

		eg.Go(func() error {
			// Prepare post file
			mdFilename := filepath.Join(postFilesPath, postFilename)
			mdFileBytes, err := os.ReadFile(mdFilename)
			if err != nil {
				return fmt.Errorf("os: failed to read file %s: %w", mdFilename, err)
			}

			ghMarkdown, _, err := ghClient.Markdown.Render(ctx, string(mdFileBytes), nil)
			if err != nil {
				return fmt.Errorf("github: failed to markdown: %w", err)
			}

			// Prepare html file
			htmlFilename := strings.TrimSuffix(postFilename, filepath.Ext(postFilename)) + extHTML
			htmlFilepath := filepath.Join(generatedPath, htmlFilename)

			htmlFile, err := os.OpenFile(htmlFilepath, os.O_RDWR|os.O_CREATE, 0o600)
			if err != nil {
				return fmt.Errorf("os: failed to open file %s: %w", htmlFilepath, err)
			}
			defer htmlFile.Close()

			// Ignore index in index file
			indexHTML := `<h2><a href="index.html"><code>~</code></a></h2>`
			if strings.Contains(postFilename, "index") {
				indexHTML = ""
			}

			if err := templatePost.Execute(htmlFile, templatePostData{
				Index: indexHTML,
				Body:  ghMarkdown,
			}); err != nil {
				return fmt.Errorf("template: failed to execute: %w", err)
			}

			return nil
		})

	}

	// Wait for all HTTP fetches to complete.
	if err := eg.Wait(); err != nil {
		log.Fatalln("I am sorry :(", err)
	} else {
		log.Println("Build successfully. Are you happy now?")
	}
}
