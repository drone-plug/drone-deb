package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cbednarski/mkdeb/deb"
	"github.com/facebookgo/flagenv"
)

var build = "0" // build number set at compile-time

func main() {
	ps := deb.DefaultPackageSpec()
	var target string

	flag.StringVar(&target, "target", ".", "target path")

	// Binary Debian Control File - Required fields
	flag.StringVar(&ps.Package, "package", "", "package package")
	flag.StringVar(&ps.Version, "version", "", "package version")
	flag.StringVar(&ps.Architecture, "architecture", "", "package architecture")
	flag.StringVar(&ps.Maintainer, "maintainer", "", "package maintainer")
	flag.StringVar(&ps.Description, "description", "", "package description")

	// Optional Fields
	flag.Var((*StringSliceFlag)(&ps.Depends), "depends", "package depends")
	flag.Var((*StringSliceFlag)(&ps.Conflicts), "conflicts", "package conflicts")
	flag.Var((*StringSliceFlag)(&ps.Breaks), "breaks", "package breaks")
	flag.Var((*StringSliceFlag)(&ps.Replaces), "replaces", "package replaces")
	flag.StringVar(&ps.Section, "section", "default", "package section")
	flag.StringVar(&ps.Priority, "priority", "extra", "package priority")
	flag.StringVar(&ps.Homepage, "homepage", "", "package homepage")

	// Control Scripts
	flag.StringVar(&ps.Preinst, "preinst", "", "package preinst script")
	flag.StringVar(&ps.Postinst, "postinst", "", "package postinst script")
	flag.StringVar(&ps.Prerm, "prerm", "", "package prerm script")
	flag.StringVar(&ps.Postrm, "postrm", "", "package postrm script")

	// Build time options
	flag.StringVar(&ps.AutoPath, "auto_path", "deb-pkg", "auth path")
	flag.Var((*JsonStringMapFlag)(&ps.Files), "files", "package files")
	flag.StringVar(&ps.TempPath, "temp_path", "", "temp path")
	flag.BoolVar(&ps.UpgradeConfigs, "upgrade-configs", false, "upgrade configs")
	flag.BoolVar(&ps.PreserveSymlinks, "preserve-symlinks", false, "preserve symlinks")

	flagenv.Prefix = "plugin_"
	flagenv.Parse()
	flag.Parse()

	if ps.Architecture == "386" {
		ps.Architecture = "i386"
	}

	err := ps.Build(target)
	if err != nil {
		fmt.Println(err)
		data, err := json.MarshalIndent(&ps, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		os.Exit(1)
	}
}

type StringSliceFlag []string

func (s *StringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *StringSliceFlag) Set(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}

type JsonStringMapFlag map[string]string

func (s *JsonStringMapFlag) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (s *JsonStringMapFlag) Set(value string) error {
	var m map[string]string
	err := json.Unmarshal([]byte(value), &m)
	if err != nil {
		return err
	}
	ss := *s
	for k, v := range m {
		ss[k] = v
	}
	return nil
}
