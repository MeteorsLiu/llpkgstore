package conan

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/goplus/llpkg/llpkg-tool/pkg/config"
)

func TestConanInstaller(t *testing.T) {
	c := &conanInstaller{
		pkg: config.Package{
			Name:    "cjson",
			Version: "1.7.18",
		},
	}

	if name := c.Name(); name != "conan" {
		t.Errorf("Unexpected name: %s", name)
	}

	tempDir, err := os.MkdirTemp("", "llpkg-tool")
	if err != nil {
		t.Errorf("Unexpected error when creating temp dir: %s", err)
	}

	if err := c.Install(tempDir); err != nil {
		t.Errorf("Install failed: %s", err)
	}

	if err := verify(c, tempDir); err != nil {
		t.Errorf("Verify failed: %s", err)
	}
}

func verify(c *conanInstaller, installDir string) error {
	err := verifyConanInstall(c, installDir)
	if err != nil {
		return err
	}

	return nil
}

func verifyConanInstall(c *conanInstaller, installDir string) error {
	// 1. ensure .pc file exists
	_, err := os.Stat(filepath.Join(installDir, c.pkg.Name+".pc"))
	if err != nil {
		return errors.New(".pc file does not exist: " + err.Error())
	}

	// 2. ensure pkg-config can find .pc file
	os.Setenv("PKG_CONFIG_PATH", installDir)
	defer os.Unsetenv("PKG_CONFIG_PATH")

	buildCmd := exec.Command("pkg-config", "--cflags", c.pkg.Name)
	out, err := buildCmd.CombinedOutput()
	if err != nil {
		return errors.New("pkg-config failed: " + err.Error() + " with output: " + string(out))
	}

	return nil
}
