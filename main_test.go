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
	testPackage := map[string]string{
		// "TARGET":       tempDir,
		"PACKAGE":      "testpackage",
		"VERSION":      "0.0.1",
		"ARCHITECTURE": "amd64",
		"MAINTAINER":   "someone",
		"DESCRIPTION":  "desc",
		"DEPENDS":      "dep1,dep2,dep3 (>= 7.0.0)",
		"CONFLICTS":    "confl1",
		"BREAKS":       "breaks1,breaks2,breaks3,breaks4",
		"REPLACES":     "replaces1,replaces2",
		"SECTION":      "cofefe",
		"PRIORITY":     "1023",
		"HOMEPAGE":     "http://internet",
		"POSTINST":     "test-data/postinst",
		"AUTO_PATH":    "test-data/deb-pkg",
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
		out, err := exec.Command("dpkg", "--info", targetFile).Output()
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println(string(out))
		_ = out // TODO: maybe test the output,
	}
	{
		out, err := exec.Command("dpkg", "--contents", targetFile).Output()
		if err != nil {
			t.Fatal(err)
		}
		// fmt.Println(string(out))
		_ = out // TODO: maybe test the output,
	}
}

func TestInstallRemovePackage(t *testing.T) {

	testPackage := map[string]string{
		// "TARGET":       tempDir,
		"PACKAGE":      "testpackage",
		"VERSION":      "0.0.1",
		"ARCHITECTURE": "amd64",
		"MAINTAINER":   "someone",
		"DESCRIPTION":  "desc",
		"HOMEPAGE":     "http://internet",
		"POSTINST":     "test-data/postinst",
		"AUTO_PATH":    "test-data/deb-pkg",
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
