package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	err := exec.Command("go", "build", "-o", "drone-deb-test").Run()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	_ = os.Remove("drone-deb-test")
	os.Exit(code)
}

func TestBuildPackage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "dpkbdeb")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	testPackage := map[string]string{
		"TARGET":      tempDir,
		"NAME":        "testpackage",
		"VERSION":     "0.0.1",
		"ARCH":        "amd64",
		"MAINTAINER":  "someone",
		"DESCRIPTION": "desc",
		"DEPENDS":     "dep1,dep2,dep3 (>= 7.0.0)",
		"CONFLICTS":   "confl1",
		"REPLACES":    "replaces1,replaces2",
		"SECTION":     "optional",
		"HOMEPAGE":    "http://internet",
		"POSTINST":    "test-data/postinst",
		"AUTO":        "test-data/deb-pkg",
		"FILES":       `{"/usr/share/testpackage_main.go":"main.go", "/usr/share/testpackage_drone.yml":".drone.yml"}`,
	}

	targetFile := fmt.Sprintf("%s/testpackage-0.0.1-amd64.deb", tempDir)

	var env []string
	env = append(env, fmt.Sprintf("PLUGIN_TARGET=%s", tempDir))
	for k, v := range testPackage {
		env = append(env, fmt.Sprintf("PLUGIN_%s=%v", k, v))
	}
	cmd := exec.Command("./drone-deb-test")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	{
		out, err := exec.Command("dpkg", "--info", targetFile).CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			t.Fatal(err)
		}
		_ = out // TODO: maybe test the output,
	}
	{
		out, err := exec.Command("dpkg", "--contents", targetFile).CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			t.Fatal(err)
		}
		_ = out // TODO: maybe test the output,
	}
}

func TestInstallRemovePackage(t *testing.T) {
	t.Skip()

	testPackage := map[string]string{
		// "TARGET":       tempDir,
		"PACKAGE":      "testpackage",
		"VERSION":      "0.0.1",
		"ARCHITECTURE": "amd64",
		"MAINTAINER":   "someone",
		"DESCRIPTION":  "desc",
		"HOMEPAGE":     "http://internet",
		"POSTINST":     "test-data/postinst",
		"AUTO":         "test-data/deb-pkg",
		"FILES":        `{"main.go":"/usr/share/testpackage_main.go", ".drone.yml":"/usr/share/testpackage_drone.yml"}`,
	}

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	targetFile := fmt.Sprintf("%s/testpackage-0.0.1-amd64.deb", tempDir)

	var env []string
	env = append(env, fmt.Sprintf("PLUGIN_TARGET=%s", tempDir))
	for k, v := range testPackage {
		env = append(env, fmt.Sprintf("PLUGIN_%s=%v", k, v))
	}
	cmd := exec.Command("./drone-deb-test")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	{
		cmd := exec.Command("dpkg", "-i", targetFile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
	}
	{
		cmd := exec.Command("apt-get", "remove", "-y", "testpackage")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			t.Fatal(err)
		}
	}

}
