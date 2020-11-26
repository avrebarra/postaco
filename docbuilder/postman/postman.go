package postman

import (
	"bytes"
	"context"
	"net/http"

	"github.com/avrebarra/postaco/docbuilder"
	"github.com/avrebarra/postaco/docbuilder/templates"
	"github.com/avrebarra/postaco/pkg/stringcurl"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct{}

type Default struct {
	config Config
}

func New(cfg Config) docbuilder.DocBuilder {
	if err := validator.New().Struct(cfg); err != nil {
		panic(err)
	}
	return &Default{config: cfg}
}

func (e *Default) Build(ctx context.Context, srcpath, distpath string, force bool) (err error) {
	bts, err := docbuilder.ReadFile(srcpath)
	if err != nil {
		return
	}

	co, err := parseCollection(bts)
	if err != nil {
		return
	}

	payload, err := buildTemplatePayload(co)
	if err != nil {
		return
	}

	docstr := templates.MakePostmanSimple(payload)

	err = docbuilder.WriteFile(distpath, []byte(docstr), force)
	if err != nil {
		return
	}

	return
}

func parseCollection(src []byte) (co Collection, err error) {
	co = Collection{}
	err = co.ParseFrom(bytes.NewBuffer(src))
	if err != nil {
		return
	}
	return
}

func buildTemplatePayload(collection Collection) (payload templates.PostmanSimplePayload, err error) {
	payload = templates.PostmanSimplePayload{
		Name: collection.Info.Name,
		RequestDirectories: []struct {
			Name     string
			Requests []templates.PostmanSimpleRequestMarkupDataV1
		}{},
		Description: collection.Info.Description,
	}

	for _, group := range collection.Items {
		reqdir := struct {
			Name     string
			Requests []templates.PostmanSimpleRequestMarkupDataV1
		}{
			Name:     group.Name,
			Requests: []templates.PostmanSimpleRequestMarkupDataV1{},
		}

		for _, request := range group.Items {
			reqmarkupdata := templates.PostmanSimpleRequestMarkupDataV1{
				Directory:          group.Name,
				Name:               request.Name,
				Description:        request.Request.Description,
				HTTPVerb:           request.Request.Method,
				URL:                request.Request.URL.Raw,
				ExampleRequestBody: request.Request.Body.Raw,

				Responses: []templates.PostmanSimpleExampleResponse{},
				CURL:      "", // will be added below

				QueryParams: []templates.PostmanSimpleRequestQueryParam{},
				URLParams:   []templates.PostmanSimpleRequestURLParam{},
			}

			for _, response := range request.Responses {
				reqmarkupdata.Responses = append(reqmarkupdata.Responses, templates.PostmanSimpleExampleResponse{
					Code:         response.Code,
					Name:         response.Name,
					Status:       response.Status,
					ResponseBody: response.Body,
				})
			}

			// add curl string
			curlstr := ""
			mockrequest, _ := buildRequest(request.Request)
			if mockrequest != nil {
				scurl, _ := stringcurl.FromRequest(mockrequest)
				curlstr = string(scurl)
			}
			reqmarkupdata.CURL = curlstr

			// register to list
			reqdir.Requests = append(reqdir.Requests, reqmarkupdata)
		}

		payload.RequestDirectories = append(payload.RequestDirectories, reqdir)
	}

	return
}

func buildRequest(r Request) (req *http.Request, err error) {
	data := bytes.NewBufferString(r.Body.Raw)
	req, err = http.NewRequest(r.Method, r.URL.Raw, data)
	if err != nil {
		return
	}

	for _, he := range r.Headers {
		req.Header.Set(he.Key, he.Value)
	}

	return
}
