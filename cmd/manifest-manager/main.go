package main

import (
	"os"

	"github.com/ultrabluewolf/manifest-manager/cli"
)

const appName = "manifest-manager"
const appVersion = "0.0.1"
const appDesc = "a utility to help in the generation and maintenance of manifest files."

func main() {
	app := cli.New(appName, appVersion, appDesc)

	cli.ApplyFlags(app)
	cli.ApplyAction(app)

	app.Run(os.Args)
}
