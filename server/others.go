package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/avrebarra/postaco/server/config"
	"github.com/avrebarra/postaco/server/helper"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var resp config.ServerResponse
	req := c.Request()

	// try to catch for http error
	if httpError, ok := err.(*echo.HTTPError); ok {
		switch httpError.Code {
		case http.StatusNotFound:
			resp = config.HashServerResponse[config.ServerRespErrPathNotFound]
			resp.Message = fmt.Sprintf("%s: %s %s", resp.Message, req.Method, req.RequestURI)

		default:
			resp = config.HashServerResponse[config.ServerRespErrUnexpected]
			resp.Message = fmt.Sprintf("%s: bad http status: %s", resp.Message, httpError.Message)
		}
	}

	// try to catch for server error
	if serr, ok := err.(config.ErrServer); ok {
		// set default
		resp.Status = serr.Code
		resp.Message = serr.Error()

		// find respective responses from list
		for _, rtemplate := range config.HashServerResponse {
			if rtemplate.Status == serr.Code {
				errmsg := ""
				if resp.Message != "" {
					errmsg = ": " + resp.Message
				}
				resp = rtemplate
				resp.Message += errmsg
				break
			}
		}
	}

	// if still empty throw general internal error
	if resp.Status == "" {
		fmt.Printf("[%s] - unexpected error: %s \n", time.Now().Format(time.RFC3339), err.Error())
		resp = config.HashServerResponse[config.ServerRespErrUnexpected]
		resp.Message = fmt.Sprintf("%s: %s", resp.Message, err.Error())
	}

	err = c.JSON(http.StatusOK, helper.PrepResponse(resp))
	if err != nil {
		fmt.Println(err)
		return
	}
}
