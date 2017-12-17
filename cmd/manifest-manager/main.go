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

	app.UsageText = fmt.Sprintf(
		"manifest-manager [options] <manifest-file>\n\n%s%s%s%s%s",
		"   examples:\n",
		"      manifest-manager -l manifest.txt\n",
		"      manifest-manager -a '/tmp/**/*' manifest.txt\n",
		"      manifest-manager -d '/var/log/*.log' manifest.txt\n",
		"      manifest-manager -c manifest.txt",
	)

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

func GetUsageText(app *cli.App) string {
	return fmt.Sprintf(
		"\n\nUSAGE:\n%s%s", "   ", app.UsageText,
	)
}

func GenerateFlagLists(addParam, deleteParam string, listFlag, cleanFlag bool) ([]string, []bool) {
	_stringparams := []string{addParam, deleteParam}
	stringparams := []string{}
	for _, param := range _stringparams {
		if param == "" {
			continue
		}
		stringparams = append(stringparams, param)
	}

	_boolparams := []bool{listFlag, cleanFlag}
	boolparams := []bool{}
	for _, param := range _boolparams {
		if param == false {
			continue
		}
		boolparams = append(boolparams, param)
	}

	return stringparams, boolparams
}

func ApplyAction(app *cli.App) *cli.App {

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			err := errors.New("path to manifest file required!")
			logger.Fatalln(err.Error())
			return err
		}

		if len(c.Args()) > 1 {
			err := errors.New("too many arguments encountered, try wrapping glob patterns in quotes")
			logger.Fatalln(err.Error(), GetUsageText(app))
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

		stringparams, boolparams := GenerateFlagLists(addParam, deleteParam, listFlag, cleanFlag)

		if len(stringparams)+len(boolparams) > 1 {
			err := errors.New("too many options encountered")
			logger.Fatalln(err.Error(), GetUsageText(app))
			return err
		}

		if len(stringparams)+len(boolparams) == 0 {
			err := errors.New("no option found")
			logger.Fatalln(err.Error(), GetUsageText(app))
			return err
		}

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
