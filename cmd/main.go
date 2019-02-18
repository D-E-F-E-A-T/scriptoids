package main

import (
	"fmt"
	"github.com/dhsavell/scriptoids/pkg/environment"
	. "github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
)

var (
	NoColor   = false
	NoSymbols = false
)

func printInfo(msg string) {
	prefix := "."
	if NoSymbols {
		prefix = "Info:"
	}

	if NoColor {
		fmt.Println(prefix, msg)
	} else {
		fmt.Println(Bold(Black(prefix)), Bold(Black(msg)))
	}
}

func printSuccess(msg string) {
	prefix := "✔"
	if NoSymbols {
		prefix = "Success:"
	}

	if NoColor {
		fmt.Println(prefix, msg)
	} else {
		fmt.Println(Green(prefix), msg)
	}
}

func printFailure(msg string) {
	prefix := "✘"
	if NoSymbols {
		prefix = "Error:"
	}

	if NoColor {
		fmt.Println(prefix, msg)
	} else {
		fmt.Println(Red(prefix), msg)
	}
}

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
			Name:        "no-color",
			Usage:       "if specified, no colored output will be displayed",
			Destination: &NoColor,
		},
		cli.BoolFlag{
			Name:        "no-symbols",
			Usage:       `if specified, labels like "Success" will be displayed instead of symbols like check marks`,
			Destination: &NoSymbols,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "links packages, enabling them in your PATH",

			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					fmt.Println("no packages specified")
				}

				env := environment.NewEnvironment(c.GlobalString("bindir"), c.GlobalString("pkgdir"))
				successes := c.NArg()

				for _, name := range c.Args() {
					pkg, err := env.GetInstalledPackageByName(name)
					if err != nil {
						printFailure(fmt.Sprintf("No such package %s, skipping...", name))
						successes--
						continue
					}

					err = env.LinkPackage(pkg)
					if err != nil {
						printFailure(fmt.Sprintf("Failed to create link for package %s. Is it already linked?", name))
						successes--
						continue
					}

					printInfo(fmt.Sprintf("%s => %s", path.Join(env.PackageDirectory, pkg.EntryPoint),
						path.Join(env.BinDirectory, pkg.Name)))
				}

				if successes > 0 {
					printSuccess(fmt.Sprintf("Linked %d packages", successes))
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
					fmt.Println("no packages specified")
				}

				env := environment.NewEnvironment(c.GlobalString("bindir"), c.GlobalString("pkgdir"))
				successes := c.NArg()

				for _, name := range c.Args() {
					pkg, err := env.GetInstalledPackageByName(name)
					if err != nil {
						printFailure(fmt.Sprintf("No such package %s, skipping...", name))
						successes--
						continue
					}

					err = env.UnlinkPackage(pkg)
					if err != nil {
						printFailure(fmt.Sprintf("Failed to unlink package %s. Was it linked in the first place?", name))
						successes--
						continue
					}

					printInfo(fmt.Sprintf("Removed %s", path.Join(env.BinDirectory, pkg.Name)))
				}

				if successes > 0 {
					printSuccess(fmt.Sprintf("Unlinked %d packages", successes))
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
