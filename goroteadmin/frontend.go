package goroteadmin

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var buildFS embed.FS

// BuildHTTPFS returns the embedded frontend build files as an http.FileSystem
func BuildHTTPFS() http.FileSystem {
	sub, err := fs.Sub(buildFS, "dist")
	if err != nil {
		panic(fmt.Errorf("erro ao criar sub FS: %w", err))
	}

	// Teste rápido pra ver se o index.html existe
	_, err = sub.Open("index.html")
	if err != nil {
		panic(fmt.Errorf("index.html não encontrado no embed: %w", err))
	}

	return http.FS(sub)
}
