package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thomasf/drone-plugins-go/plugins"
	"github.com/thomasf/vfs"
	"github.com/xor-gate/debpkg"
)

var build = "0" // build number set at compile-time

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := &struct {
		// required
		Name         string
		Version      string
		Architecture string
		Maintainer   string
		Description  string
		Depends      plugins.StringSliceFlag
		Conflicts    plugins.StringSliceFlag
		Replaces     plugins.StringSliceFlag
		Section      string // Defaults to "default"
		Priority     string // Defaults to "extra"
		Homepage     string
		Auto         string // Defaults to "contrib/debian"
		Preinst      string
		Postinst     string
		Prerm        string
		Postrm       string
		Files        plugins.StringMapFlag
		Conffiles    plugins.StringSliceFlag

		// output

		Target string
	}{
		Files:     make(map[string]string),
		Conffiles: []string{"/etc/**"},
	}

	// Binary Debian Control File - Required fields
	flag.StringVar(&c.Name, "name", "", "package package")
	flag.StringVar(&c.Version, "version", "", "package version")
	flag.StringVar(&c.Architecture, "arch", "", "package architecture")
	flag.StringVar(&c.Maintainer, "maintainer", "", "package maintainer")
	flag.StringVar(&c.Description, "description", "", "package description")

	// Optional Fields
	flag.Var(&c.Depends, "depends", "package depends")
	flag.Var(&c.Conflicts, "conflicts", "package conflicts")
	flag.Var(&c.Replaces, "replaces", "package replaces")
	flag.StringVar(&c.Section, "section", "default", "package section")
	flag.StringVar(&c.Priority, "priority", "extra", "package priority")
	flag.StringVar(&c.Homepage, "homepage", "", "package homepage")

	// Files
	flag.StringVar(&c.Auto, "auto", "contrib/debian", "auth path")
	flag.Var(&c.Files, "files", "package files")
	flag.StringVar(&c.Preinst, "preinst", "", "package preinst script")
	flag.StringVar(&c.Postinst, "postinst", "", "package postinst script")
	flag.StringVar(&c.Prerm, "prerm", "", "package prerm script")
	flag.StringVar(&c.Postrm, "postrm", "", "package postrm script")

	flag.StringVar(&c.Target, "target", "", "target directory")

	plugins.Parse()

	if c.Architecture == "386" {
		c.Architecture = "i386"
	}

	d := debpkg.New()
	defer d.Close()

	d.SetName(c.Name)
	d.SetVersion(c.Version)
	d.SetArchitecture(c.Architecture)
	d.SetMaintainer(c.Maintainer)
	d.SetDescription(c.Description)

	d.SetDepends(strings.Join(c.Depends, ", "))
	d.SetConflicts(strings.Join(c.Conflicts, ", "))
	d.SetReplaces(strings.Join(c.Replaces, ", "))
	d.SetSection(c.Section)
	d.SetHomepage(c.Homepage)

	controlFs := vfs.NewNameSpace()
	dataFs := vfs.NewNameSpace()

	if c.Auto != "" {
		controlFileMap := make(map[string]string)
		dataFs.Bind(
			"/", vfs.Exclude(vfs.OS(c.Auto), "/preinst", "/postinst", "/prerm", "/postrm"),
			"/", vfs.BindAfter)

		for _, cf := range []string{"preinst", "postinst", "prerm", "postrm"} {
			cfp := filepath.Join(c.Auto, cf)
			if _, err := os.Stat(cfp); err == nil {
				controlFileMap[cf] = cfp
			}
		}
	}

	if c.Files != nil && len(c.Files) > 0 {
		fmap := make(map[string]string)
	fls:
		for dst, src := range c.Files {
			fi, err := os.Stat(src)
			if err != nil {
				log.Fatalf("error loading Files: %v", err)
			}
			if fi.IsDir() {
				dataFs.Bind(dst, vfs.OS(src), "/", vfs.BindAfter)
				continue fls
			}

			fmap[strings.TrimLeft(dst, "/")] = src
		}
		dataFs.Bind("/", vfs.FileMap(fmap), "/", vfs.BindAfter)
	}

	controlFileMap := make(map[string]string)
maint:
	for dst, src := range map[string]string{
		"preinst":  c.Preinst,
		"postinst": c.Postinst,
		"prerm":    c.Prerm,
		"postrm":   c.Postrm,
	} {
		src = strings.TrimSpace(src)
		if src == "" {
			continue maint
		}
		fi, err := os.Stat(src)
		if err != nil {
			log.Fatal(err)
		}
		if fi.IsDir() {
			log.Fatalf("maint script %s must be a file, not a directory:", src)
		}
		controlFileMap[dst] = src
	}

	{
		var conffilesBuf bytes.Buffer
		err := vfs.Walk("/", dataFs, func(p string, info os.FileInfo, err error) error {
			// log.Println("p", p, info)
			if err != nil {
				log.Fatal(err)
			}
			if info.IsDir() {
				err := d.AddEmptyDirectory(p)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			}
			if opr, ok := info.(vfs.OSPather); ok {
				op := opr.OSPath()
				err := d.AddFile(op, p)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatalln("expected all files to be of OSPather type!", p)
			}

			for _, pattern := range c.Conffiles {
				ok, err := filepath.Match(pattern, p)
				if err != nil {
					log.Fatal(err)
				}
				if ok {
					fmt.Fprintln(&conffilesBuf, p)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		if conffilesBuf.Len() > 0 {
			conffilesBuf.WriteString("\n")
			controlFileMap["conffiles"] = conffilesBuf.String()
		}
		if len(controlFileMap) > 0 {
			controlFs.Bind("/", vfs.Map(controlFileMap), "/", vfs.BindBefore)
		}
	}
	log.Println(d.GetFilename())
	if err := d.Write(filepath.Join(c.Target, d.GetFilename())); err != nil {
		log.Fatal(err)
	}
	d.Close()
}
