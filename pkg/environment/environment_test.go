package environment

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var (
	TestRoot = path.Join("..", "..", "res", "test")
	TestBin  = path.Join(TestRoot, "bin")
	TestPkg  = path.Join(TestRoot, "pkg")
)

func MustGetTestEnv() *Environment {
	binPath, err := filepath.Abs(TestBin)

	if err != nil {
		panic("failed to get absolute test bin path")
	}

	pkgPath, err := filepath.Abs(TestPkg)

	if err != nil {
		panic("failed to get absolute test pkg path")
	}

	return &Environment{
		BinDirectory:     binPath,
		PackageDirectory: pkgPath,
	}
}

func MustCleanupTestBin() {
	dir, err := ioutil.ReadDir(TestBin)

	if err != nil {
		panic(err)
	}

	for _, d := range dir {
		if d.Name() != ".gitkeep" {
			err = os.RemoveAll(path.Join(TestBin, d.Name()))
			if err != nil {
				panic(err)
			}
		}
	}
}

func TestEnvironment_GetInstalledPackageByName_Exists(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()

	pkg, err := env.GetInstalledPackageByName("package1")

	req.Nil(err)
	req.Equal(pkg, Package{
		Name:        "package1",
		Description: "Test package #1",
		Version:     "0.0.0",
		EntryPoint:  "main.sh",
	})
}

func TestEnvironment_GetInstalledPackageByName_Missing(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()

	_, err := env.GetInstalledPackageByName("asdf")

	req.NotNil(err)
}

func TestEnvironment_GetInstalledPackageByName_Invalid(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()

	_, err := env.GetInstalledPackageByName("invalid_missing_scriptoid_hcl")
	req.NotNil(err)

	_, err = env.GetInstalledPackageByName("invalid_scriptoid_hcl_contents")
	req.NotNil(err)
}

func TestEnvironment_IsPackageValid_ValidPackage(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()

	pkg, err := env.GetInstalledPackageByName("package1")
	req.Nil(err)
	req.True(env.IsPackageValid(pkg))
}

func TestEnvironment_IsPackageValid_BlankEntryPoint(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()

	pkg, err := env.GetInstalledPackageByName("invalid_empty_entry_point")
	req.Nil(err)
	req.False(env.IsPackageValid(pkg))
}

func TestEnvironment_IsPackageValid_MissingEntryPoint(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()

	pkg, err := env.GetInstalledPackageByName("invalid_missing_entry_point")
	req.Nil(err)
	req.False(env.IsPackageValid(pkg))
}

func TestEnvironment_LinkPackage_Exists(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()
	defer MustCleanupTestBin()

	pkg, err := env.GetInstalledPackageByName("package1")
	req.Nil(err)

	err = env.LinkPackage(pkg)
	req.Nil(err)
	req.FileExists(path.Join(TestBin, "package1"))
}

func TestEnvironment_LinkPackage_AlreadyExists(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()
	defer MustCleanupTestBin()

	pkg, err := env.GetInstalledPackageByName("package1")
	req.Nil(err)

	err = env.LinkPackage(pkg)
	req.Nil(err)
	req.FileExists(path.Join(TestBin, "package1"))

	err = env.LinkPackage(pkg)
	req.NotNil(err)
}

func TestEnvironment_UnlinkPackage_Exists(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()
	defer MustCleanupTestBin()

	pkg, err := env.GetInstalledPackageByName("package1")
	req.Nil(err)

	err = env.LinkPackage(pkg)
	req.Nil(err)

	err = env.UnlinkPackage(pkg)
	req.Nil(err)
	req.False(fileExists(path.Join(TestBin, "package1")))
}

func TestEnvironment_UnlinkPackage_NotLinked(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()
	defer MustCleanupTestBin()

	pkg, err := env.GetInstalledPackageByName("package1")
	req.Nil(err)

	err = env.UnlinkPackage(pkg)
	req.NotNil(err)
}

func TestEnvironment_IsPackageLinked(t *testing.T) {
	req := require.New(t)
	env := MustGetTestEnv()
	defer MustCleanupTestBin()

	pkg, err := env.GetInstalledPackageByName("package1")
	req.Nil(err)

	req.False(env.IsPackageLinked(pkg))

	err = env.LinkPackage(pkg)
	req.Nil(err)

	req.True(env.IsPackageLinked(pkg))
}