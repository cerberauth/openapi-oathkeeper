package main

import "github.com/cerberauth/openapi-oathkeeper/cmd"

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
