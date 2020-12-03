package postaco

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/avrebarra/postaco/bindata"
	"github.com/avrebarra/postaco/docbuilder"
	"github.com/avrebarra/postaco/pkg/fsx"
	mkdon "github.com/avrebarra/postaco/pkg/markdown"
	"github.com/avrebarra/postaco/pkg/writerlog"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	// Dependency string `validate:"required"`
}

type PostacoStruct struct {
	config Config
}

func New(cfg Config) Postaco {
	if err := validator.New().Struct(cfg); err != nil {
		panic(err)
	}
	return &PostacoStruct{config: cfg}
}

func (e *PostacoStruct) BuildDir(ctx context.Context, config ConfigBuildDir) (err error) {
	logger := writerlog.WriterLog{Writer: config.Logger}

	logger.Log(fmt.Sprintf("building postaco documents from '%s' to '%s'...", config.SourcePath, config.OutputPath))

	var outpath string
	outpath, err = filepath.Abs(config.OutputPath)
	if err != nil {
		return
	}

	// find supported documents
	filepaths := []string{}
	err = filepath.Walk(config.SourcePath, func(path string, info os.FileInfo, errwalk error) (errops error) {
		if errwalk != nil {
			errops = errwalk
			return
		}

		// check if extension are supported
		for _, ext := range typext {
			if fsx.HasExtension(info.Name(), ext) {
				logger.Log(fmt.Sprintf("detected document: %s", path))
				filepaths = append(filepaths, path)
			}
		}

		return nil
	})
	if err != nil {
		return
	}

	logger.Log(fmt.Sprintf("building documents..."))

	// build postaco doc file
	var doc = Documentation{
		Name:     config.DocTitle,
		Document: []DocumentIndex{},
	}

	for _, fpath := range filepaths {
		relpath, errops := fsx.PathAsRelative(fpath, config.SourcePath)
		if errops != nil {
			err = fmt.Errorf("failed converting path to relative: %s: %w", fpath, errops)
			return
		}

		switch true {
		case fsx.HasExtension(fpath, typext[DoctypePostman]):
			reldestsrc := filepath.Join("documents", relpath)
			reldestmd := fsx.RemoveExtension(reldestsrc, typext[DoctypePostman]) + ".md" // change extension to md

			destsrc := filepath.Join(config.OutputPath, reldestsrc)
			destmd := filepath.Join(config.OutputPath, reldestmd)

			// build document
			logger.Log(fmt.Sprintf("building postman document: %s", destmd))
			err = config.PostmanDocBuilder.Build(context.Background(), fpath, destmd, true)
			if err != nil {
				err = fmt.Errorf("postman document build failed: %w", err)
				return
			}

			// copy source docfile
			err = fsx.CopyFile(fpath, destsrc, 1000)
			if err != nil {
				err = fmt.Errorf("cannot copy source docfile: %w", err)
				return
			}

			// process index
			doccontent, rerr := docbuilder.ReadFile(destmd)
			if rerr != nil {
				err = fmt.Errorf("failed reading file for index: %w", rerr)
				return
			}
			headers := mkdon.ExtractHeaders(doccontent)

			firstlv1 := DocumentIndex{}
			lastlv2 := DocumentIndex{}
			lastlv3 := DocumentIndex{}

			for _, h := range headers {
				// TODO: beware of error here
				tmp := strings.SplitN(h, " ", 2)

				level := len(tmp[0])
				title := tmp[1]
				tag := mkdon.ExtractHeaderTag(title)

				tagclean := strings.ReplaceAll(strings.ReplaceAll(tag, "}", ""), "{", "")
				title = strings.ReplaceAll(title, tag, "")

				switch level {
				case 1:
					if firstlv1.Kind != "" {
						break
					}
					firstlv1 = DocumentIndex{
						Kind:                    "document/postman",
						Title:                   title,
						DocumentRelDir:          filepath.Dir(relpath),
						DocumentRelPathSource:   reldestsrc,
						DocumentRelPathMarkdown: reldestmd,
						Indexables:              []string{title},
						Rel:                     tagclean,
					}
					doc.Document = append(doc.Document, firstlv1)

				case 2:
					lastlv2 = DocumentIndex{
						Kind:                    "bookmark/postman/dir",
						Title:                   title,
						DocumentRelDir:          firstlv1.DocumentRelPathMarkdown,
						DocumentRelPathSource:   "-",
						DocumentRelPathMarkdown: reldestmd,
						Indexables:              []string{title},
						Rel:                     tagclean,
					}
					doc.Document = append(doc.Document, lastlv2)

				case 3:
					lastlv3 = DocumentIndex{
						Kind:                    "bookmark/postman/request",
						Title:                   title,
						DocumentRelDir:          firstlv1.DocumentRelPathMarkdown,
						DocumentRelPathSource:   "-",
						DocumentRelPathMarkdown: reldestmd,
						Indexables:              []string{title},
						Rel:                     tagclean,
					}
					doc.Document = append(doc.Document, lastlv3)
				}
			}
		}
	}

	// write index
	dest := filepath.Join(outpath, "postaco.doc.json")
	logger.Log(fmt.Sprintf("writing index index: %s", dest))
	indexcontent, err := json.MarshalIndent(doc, "", "\t")
	if err != nil {
		return
	}
	errwrite := docbuilder.WriteFile(dest, indexcontent, true)
	if errwrite != nil {
		err = fmt.Errorf("write asset fail: %w", errwrite)
		return
	}

	// write webapp assets
	logger.Log(fmt.Sprintf("writing web assets..."))
	for _, assetname := range bindata.AssetNames() {
		content, errread := bindata.Asset(assetname)
		if err != nil {
			err = fmt.Errorf("reading web asset fail: %w", errread)
			return
		}

		dest := filepath.Join(outpath, assetname)
		logger.Log(fmt.Sprintf("writing: %s", dest))
		errwrite := docbuilder.WriteFile(dest, content, true)
		if errwrite != nil {
			err = fmt.Errorf("write asset fail: %w", errwrite)
			return
		}
	}

	logger.Log(fmt.Sprintf("process done!"))

	return
}
