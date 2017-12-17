package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/urfave/cli.v1"

	"github.com/ultrabluewolf/manifest-manager/files"
	manifestpkg "github.com/ultrabluewolf/manifest-manager/manifest"

	mmlogger "github.com/ultrabluewolf/manifest-manager/logger"
)

const appName = "manifest-manager"
const appVersion = "0.0.1"
const appDesc = "a utility to help in the generation and maintenance of manifest files."

var logger = mmlogger.New()

func main() {
	app := SetupApp()
	ApplyFlags(app)
	ApplyAction(app)
	app.Run(os.Args)
}

func SetupApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Usage = appDesc
	return app
}

func ApplyFlags(app *cli.App) *cli.App {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "add, a",
			Usage: "add all files matching the glob to the target manifest",
		},
		cli.StringFlag{
			Name:  "delete, d",
			Usage: "remove all files matching the glob from target manifest",
		},
		cli.BoolFlag{
			Name:  "clean, c",
			Usage: "clean up the file by removing all paths in it that do not exist",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "shows the content of the manifest file",
		},
	}
	return app
}

func ApplyAction(app *cli.App) *cli.App {

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			err := errors.New("path to manifest file required!")
			logger.Fatalln(err.Error())
			return err
		}

		var (
			manifestFilePath = c.Args().First()
			// flags/options
			addParam    = c.String("add")
			deleteParam = c.String("delete")
			listFlag    = c.Bool("list")
			cleanFlag   = c.Bool("clean")
		)

		if !files.Exists(manifestFilePath) {
			if err := manifestpkg.New(manifestFilePath).Save(); err != nil {
				logger.Fatalln("manifest initialization issue -", err.Error())
				return err
			}
		}

		manifest, err := manifestpkg.ParseManifestFile(manifestFilePath)
		if err != nil {
			logger.Fatalln("manifest file parsing issue -", err.Error())
			return err
		}

		logger.Debug("parsed -", manifest.Files)

		if listFlag {
			fmt.Println(strings.Join(manifest.FileList(), "\n"))
			return nil

		} else if addParam != "" {
			if err = manifest.Add(addParam); err != nil {
				logger.Fatalln("manifest add issue -", err.Error())
				return err
			}

		} else if deleteParam != "" {
			if err = manifest.Remove(deleteParam); err != nil {
				logger.Fatalln("manifest delete issue -", err.Error())
				return err
			}

		} else if cleanFlag {
			if err = manifest.Prune(); err != nil {
				logger.Fatalln("manifest clean issue -", err.Error())
				return err
			}
		}

		logger.Debug("modified -", manifest.Files)

		if err = manifest.Save(); err != nil {
			logger.Fatalln("manifest save issue -", err.Error())
			return err
		}

		return nil
	}
	return app
}
