package main

import (
	"fmt"

	"github.com/avrebarra/postaco/cmd"
)

//go:generate qtc -dir=docbuilder/templates
//go:generate parcel build webapp/index.html --no-source-maps --no-content-hash -d tmp/.buildweb
//go:generate go-bindata -prefix "tmp/.buildweb/" -fs -o ./bindata/bindata.go -pkg bindata -ignore=\\*.map tmp/.buildweb/...

func main() {
	cmd.Initialize()
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
