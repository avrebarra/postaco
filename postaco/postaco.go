package postaco

import (
	"context"
	"io"

	"github.com/avrebarra/postaco/docbuilder"
)

const (
	DoctypePostman = "postman"
)

var (
	typext = map[string]string{
		DoctypePostman: ".postman_collection.json",
	}
)

type Postaco interface {
	BuildDir(ctx context.Context, cfg ConfigBuildDir) (err error)
}

type Documentation struct {
	Name     string
	Document []DocumentIndex `json:"documents"`
}
type DocumentIndex struct {
	Kind       string   `json:"kind"`
	Title      string   `json:"title"`
	Indexables []string `json:"indexables"`

	DocumentRelDir          string `json:"document_dir"`
	DocumentRelPathMarkdown string `json:"document_path_markdown"`
	DocumentRelPathSource   string `json:"document_path_source"`

	Rel string `json:"rel"`
}

type ConfigBuildDir struct {
	SourcePath string `validate:"required,ne="`
	OutputPath string `validate:"required,ne="`

	Logger            io.Writer             `validate:"required"`
	PostmanDocBuilder docbuilder.DocBuilder `validate:"required"`

	DocTitle string `validate:"required,ne="`
}
