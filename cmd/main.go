package main

import (
	"fmt"
	"github.com/dhsavell/scriptoids/pkg/scriptoids"
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
	}

	app.Commands = []cli.Command{
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "links a package, enabling it in your PATH",

			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					fmt.Println("no packages specified")
				}

				env := scriptoids.NewEnvironment(c.GlobalString("bindir"), c.GlobalString("pkgdir"))

				for _, name := range c.Args() {
					pkg, err := env.GetInstalledPackageByName(name)
					if err != nil {
						return err
					}

					err = env.LinkPackage(pkg)
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			Name:    "unlink",
			Aliases: []string{"u"},
			Usage:   "unlinks a package, removing it from your PATH",

			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					fmt.Println("no packages specified")
				}

				env := scriptoids.NewEnvironment(c.GlobalString("bindir"), c.GlobalString("pkgdir"))

				for _, name := range c.Args() {
					pkg, err := env.GetInstalledPackageByName(name)
					if err != nil {
						return err
					}

					err = env.UnlinkPackage(pkg)
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
