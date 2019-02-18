package cli

import (
	"fmt"
	"github.com/dhsavell/scriptoids/pkg/environment"
	"path"
)

func LinkPackages(d *Display, env *environment.Environment, pkgNames []string) {
	var successes int

	for _, name := range pkgNames {
		pkg, err := env.GetInstalledPackageByName(name)
		if err != nil {
			d.Failure("No such package %s, skipping...", name)
			continue
		}

		err = env.LinkPackage(pkg)
		if err != nil {
			d.Failure("Failed to create link for package %s. Is it already linked?", name)
			continue
		}

		d.Info(fmt.Sprintf("%s => %s", path.Join(env.PackageDirectory, pkg.EntryPoint), path.Join(env.BinDirectory, pkg.Name)))
		successes++
	}

	if successes > 0 {
		pkgStr := "pkgStr"
		if successes == 1 {
			pkgStr = "package"
		}

		d.Success("Linked %d %s", successes, pkgStr)
	}
}

func UnlinkPackages(d *Display, env *environment.Environment, pkgNames []string) {
	var successes int

	for _, name := range pkgNames {
		pkg, err := env.GetInstalledPackageByName(name)
		if err != nil {
			d.Failure("No such package %s, skipping...", name)
			continue
		}

		err = env.UnlinkPackage(pkg)
		if err != nil {
			d.Failure("Failed to unlink package %s. Was it linked in the first place?", name)
			continue
		}

		d.Info("Removed %s", path.Join(env.BinDirectory, pkg.Name))
		successes++
	}

	if successes > 0 {
		pkgStr := "pkgStr"
		if successes == 1 {
			pkgStr = "package"
		}

		d.Success("Unlinked %d %s", successes, pkgStr)
	}
}

func ListPackages(d *Display, env *environment.Environment) {
	pkgs, err := env.GetAllInstalledPackages()

	if err != nil {
		d.Failure("Failed to list packages. Does the package directory exist?")
		return
	}

	fmt.Printf("%-15s %-10s %-10s %-10s %s\n", "Name", "Version", "Status", "Linked?", "Description")
	for _, pkg := range pkgs {
		pkgStatus := "OK"

		ok, err := env.IsPackageValid(pkg)
		if err != nil || !ok {
			pkgStatus = "Error"
		}

		linkedStr := "No"
		if env.IsPackageLinked(pkg) {
			linkedStr = "Yes"
		}

		fmt.Printf("%-15s %-10s %-10s %-10s %s\n", pkg.Name, pkg.Version, pkgStatus, linkedStr, pkg.Description)
	}
}
