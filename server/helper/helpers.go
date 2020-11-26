package helper

import (
	"strings"

	"github.com/avrebarra/postaco/server/config"
)

func PrepResponse(in config.ServerResponse) config.ServerResponse {
	// send empty object instead of nil
	if in.Data == nil {
		in.Data = map[string]interface{}{}
	}

	// change case message to lower
	in.Message = strings.ToLower(in.Message)

	return in
}
