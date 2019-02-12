package scriptoids

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"os"
	"path"
)

const (
	PackageDefinitionFilename = "scriptoid.hcl"
)

// An Environment represents Scriptoids' working environment, consisting of a "bin" directory where executables are
// linked and a "pkg" directory where package sources are stored.
type Environment struct {
	BinDirectory     string
	PackageDirectory string
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetInstalledPackageByName gets a package within this Environment based on its name.
//
// GetInstalledPackageByName looks for a directory with the given name in the environment's PackageDirectory. An empty
// package is returned along with its respective error if one could not be found for any reason.
//
// GetInstalledPackageByName will return a package no matter what if the necessary file (specified in
// PackageDefinitionFilename) exists-- this includes if any essential fields (i.e. "name" or "entrypoint") are blank.
// IsPackageValid can be used to check the package state.
func (e *Environment) GetInstalledPackageByName(name string) (Package, error) {
	packageDefinition := path.Join(e.PackageDirectory, name, PackageDefinitionFilename)
	_, err := os.Stat(packageDefinition)

	if err != nil {
		return Package{}, err
	}

	p := Package{}
	fileText, err := ioutil.ReadFile(packageDefinition)

	if err != nil {
		return Package{}, err
	}

	err = hcl.Decode(&p, string(fileText))

	if err != nil {
		return Package{}, err
	}

	return p, nil
}

// IsPackageValid determines whether or not an existing Package is valid. A package is considered valid if its name is
// not blank and its entry point exists.
func (e *Environment) IsPackageValid(p Package) bool {
	return p.EntryPoint != "" && p.Name != "" && fileExists(path.Join(e.PackageDirectory, p.Name, p.EntryPoint))
}

// IsPackageLinked determines whether or not a given package has been linked to this Environment's BinDirectory.
//
// IsPackageLinked may return true even if the given package is not installed to this Environment's PackageDirectory.
// Any file with the given package's name in the bin directory will count towards a package being linked.
func (e *Environment) IsPackageLinked(p Package) bool {
	return fileExists(path.Join(e.BinDirectory, p.Name))
}

// LinkPackage creates a symbolic link from an installed package's entry point to this Environment's BinDirectory.
//
// LinkPackage requires that the given package is installed, valid, and unlinked.
func (e *Environment) LinkPackage(p Package) error {
	if !e.IsPackageValid(p) {
		return errors.New(fmt.Sprintf("package is in an invalid state: %s", p.Name))
	}

	if e.IsPackageLinked(p) {
		return errors.New(fmt.Sprintf("package is already linked: %s", p.Name))
	}

	entryPoint := path.Join(e.PackageDirectory, p.Name, p.EntryPoint)
	err := os.Symlink(entryPoint, path.Join(e.BinDirectory, p.Name))

	return err
}

// UnlinkPackage removes a link from an installed package's entry point to this Environment's BinDirectory.
//
// UnlinkPackage requires that the given package is already linked.
//
// UnlinkPackage will remove any file with the given package's name from the bin directory.
func (e *Environment) UnlinkPackage(p Package) error {
	if !e.IsPackageLinked(p) {
		return errors.New(fmt.Sprintf("package is not linked: %s", p.Name))
	}

	return os.Remove(path.Join(e.BinDirectory, p.Name))
}
