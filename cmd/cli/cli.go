package main

import (
	"fmt"
	"os"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
