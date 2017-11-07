package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/drone-plug/drone-plugins-go/plug/plugtest"
)

func TestBuildPackage(t *testing.T) {
	p := NewPlugin()
	pt := plugtest.New(t, p)

	tempDir, err := ioutil.TempDir("", "dpkbdeb")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	// pt.SetDebug()
	pt.SetPluginVars(map[string]string{
		"target":      tempDir,
		"package":     "testpackage",
		"version":     "0.0.1",
		"arch":        "amd64",
		"maintainer":  "someone",
		"description": "desc",
		"depends":     "dep1,dep2,dep3 (>= 7.0.0)",
		"conflicts":   "confl1",
		"replaces":    "replaces1,replaces2",
		"section":     "optional",
		"homepage":    "http://internet",
		"postinst":    "test-data/postinst",
		"auto":        "test-data/deb-pkg",
		"files": plugtest.JSON(map[string]string{
			"/usr/share/testpackage_main.go":   "main.go",
			"/usr/share/testpackage_drone.yml": ".drone.yml",
		}),
	})

	pt.AssertSuccess()
	pt.AssertOutput("ascd")

	// targetFile := fmt.Sprintf("%s/testpackage-0.0.1-amd64.deb", tempDir)

	// {
	// 	out, err := exec.Command("dpkg", "--info", targetFile).CombinedOutput()
	// 	fmt.Println(string(out))
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	_ = out // TODO: maybe test the output,
	// }
	// {
	// 	out, err := exec.Command("dpkg", "--contents", targetFile).CombinedOutput()
	// 	fmt.Println(string(out))
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	_ = out // TODO: maybe test the output,
	// }
}
