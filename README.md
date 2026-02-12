# posts-go

[![Go](https://github.com/haunt98/posts-go/actions/workflows/go.yml/badge.svg)](https://github.com/haunt98/posts-go/actions/workflows/go.yml)
[![gitleaks](https://github.com/haunt98/posts-go/actions/workflows/gitleaks.yml/badge.svg)](https://github.com/haunt98/posts-go/actions/workflows/gitleaks.yml)
[![Latest Version](https://img.shields.io/github/v/tag/haunt98/posts-go)](https://github.com/haunt98/posts-go/tags)

Write markdown, convert to html, then publish using Github Pages.

First add GitHub access token in `~/.netrc`.

Steps:

- Write new post in `posts`

- Update index in `posts/index.md`

- Generate HTML by running `just`

## Thanks

- https://github.com/sindresorhus/github-markdown-css
- https://github.com/google/go-github
