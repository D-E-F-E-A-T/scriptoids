package main

import (
	scli "github.com/dhsavell/scriptoids/pkg/cli"
	"github.com/dhsavell/scriptoids/pkg/environment"
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
)

func main() {
	app := cli.NewApp()
	app.Name = "scriptoids"
	app.Usage = "a package manager for small utilities"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "bindir",
			Usage:  "package bin directory (should be in your PATH)",
			EnvVar: "SCRIPTOIDS_BIN",
			Value:  path.Join(os.Getenv("HOME"), ".scriptoids", "bin"),
		},
		cli.StringFlag{
			Name:   "pkgdir",
			Usage:  "package install directory (does not need to be in your PATH)",
			EnvVar: "SCRIPTOIDS_PKG",
			Value:  path.Join(os.Getenv("HOME"), ".scriptoids", "pkg"),
		},
		cli.BoolFlag{
			Name:  "no-color",
			Usage: "if specified, no colored output will be displayed",
		},
		cli.BoolFlag{
			Name:  "no-symbols",
			Usage: `if specified, labels like "Success" will be displayed instead of symbols like check marks`,
		},
	}

	var display *scli.Display
	var env *environment.Environment

	app.Before = func(c *cli.Context) error {
		display = &scli.Display{
			NoColor:   c.GlobalBool("no-color"),
			NoSymbols: c.GlobalBool("no-symbols"),
		}

		env = environment.NewEnvironment(c.GlobalString("bindir"), c.GlobalString("pkgdir"))

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "links packages, enabling them in your PATH",

			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					display.Failure("No packages were specified.")
				} else {
					scli.LinkPackages(display, env, c.Args())
				}

				return nil
			},
		},
		{
			Name:    "unlink",
			Aliases: []string{"u"},
			Usage:   "unlinks packages, removing them from your PATH",

			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					display.Failure("No packages were specified.")
				} else {
					scli.UnlinkPackages(display, env, c.Args())
				}

				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "lists all installed packages",

			Action: func(c *cli.Context) error {
				scli.ListPackages(display, env)
				return nil
			},
		},
		{
			Name:    "init",
			Aliases: []string{"new"},
			Usage:   "initializes a new scriptoid.hcl file",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Usage: "output filename",
					Value: "./scriptoid.hcl",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "scriptoid name",

				},
				cli.StringFlag{
					Name:  "version",
					Usage: "scriptoid version",
					Value: "0.0.0",
				},
				cli.StringFlag{
					Name:  "desc",
					Usage: "scriptoid description",
				},
				cli.StringFlag{
					Name:  "entry",
					Usage: "scriptoid entry point",
				},
			},

			Action: func(c *cli.Context) error {
				scli.InitPackage(
					display,
					c.String("path"),
					c.String("name"),
					c.String("version"),
					c.String("desc"),
					c.String("entry"),
				)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
