package main

import (
	"context"
	"flag"
	"os"

	"github.com/drone-plug/drone-deb/deb"
	"github.com/drone-plug/drone-plugins-go/plug"
	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

const build = "0" // build number set at compile-time

func main() {

	c := &struct {
		// required
		// Name         string

		Auto        string // Defaults to "contrib/debian"
		ControlFile string
		Preinst     string
		Postinst    string
		Prerm       string
		Postrm      string
		Files       map[string]string
		Conffiles   []string

		// output

		Target string
	}{
		Files:     make(map[string]string),
		Conffiles: []string{"/etc/**"},
	}

	var bp deb.BinaryControl
	// Binary Debian Control File - Required fields
	flag.StringVar(&bp.Package, "package", "", "package package")
	flag.Var((*VersionFlag)(&bp.Version), "version", "package version")
	flag.Var((*ArchFlag)(&bp.Arch), "arch", "package architectures")

	flag.StringVar(&bp.Maintainer, "maintainer", "", "package maintainer")
	flag.StringVar(&bp.Description, "description", "", "package description")

	// Optional Fields

	flag.Var((*DependencyFlag)(&bp.Breaks), "breaks", "package Breaks")
	flag.Var((*DependencyFlag)(&bp.Conflicts), "conflicts", "package Conflicts")
	flag.Var((*DependencyFlag)(&bp.Depends), "depends", "package Depends")
	flag.Var((*DependencyFlag)(&bp.PreDepends), "pre-depends", "package Pre-Depends")
	flag.Var((*DependencyFlag)(&bp.Recommends), "recommends", "package Recommends")
	flag.Var((*DependencyFlag)(&bp.Replaces), "replaces", "package Replaces")
	flag.Var((*DependencyFlag)(&bp.Suggests), "suggests", "package Suggests")

	flag.StringVar(&bp.Section, "section", "default", "package section")
	flag.StringVar(&bp.Priority, "priority", "extra", "package priority")
	flag.StringVar(&bp.Homepage, "homepage", "", "package homepage")

	// Files
	plug.StringSliceVar(&c.Conffiles, "conf_files", "config files")
	flag.StringVar(&c.Auto, "auto", "contrib/debian", "auth path")
	// flag.Var(&c.Files, "files", "package files")
	flag.StringVar(&c.Preinst, "preinst", "", "package preinst script")
	flag.StringVar(&c.Postinst, "postinst", "", "package postinst script")
	flag.StringVar(&c.Prerm, "prerm", "", "package prerm script")
	flag.StringVar(&c.Postrm, "postrm", "", "package postrm script")

	flag.StringVar(&c.Target, "target", "", "target directory")

	plug.Run(func(ctx context.Context) error {

		if c.ControlFile != "" {
			f, err := os.Open(c.ControlFile)
			if err != nil {
				plug.Fatal(err)
			}
			defer f.Close()
			var bpc deb.BinaryControlTemplate
			if err := control.Unmarshal(&bpc, f); err != nil {
				plug.Fatal(err)
			}
			bpc.AddMissing(&bp)
		}

		if err := bp.Validate(); err != nil {
			for _, e := range err.(deb.MissingFieldsError).Refs {
				plug.Usagef(e, "must be set")
			}
			return plug.ErrUsageError
		}

		if err := control.Marshal(os.Stdout, &bp); err != nil {
			plug.Fatal(err)
		}

		return nil
	})
}

// func OLD() {
// 	d := debpkg.New()
// 	defer d.Close()

// 	d.SetName(c.Name)
// 	d.SetVersion(c.Version)
// 	d.SetArchitecture(c.Architecture)
// 	d.SetMaintainer(c.Maintainer)
// 	d.SetDescription(c.Description)

// 	// d.SetDepends(strings.Join(c.Depends, ", "))
// 	// d.SetConflicts(strings.Join(c.Conflicts, ", "))
// 	// d.SetReplaces(strings.Join(c.Replaces, ", "))

// 	d.SetSection(c.Section)
// 	d.SetHomepage(c.Homepage)

// 	controlFs := vfs.NewNameSpace()
// 	dataFs := vfs.NewNameSpace()

// 	maintScripts := []string{"preinst", "postinst", "prerm", "postrm"}
// 	if c.Auto != "" {
// 		controlFileMap := make(map[string]string)
// 		fsOrDie(dataFs.BindSafe(
// 			"/", vfs.SafeExclude(vfs.SafeOS(c.Auto), maintScripts...),
// 			"/", vfs.BindAfter))

// 		for _, cf := range maintScripts {
// 			cfp := filepath.Join(c.Auto, cf)
// 			if _, err := os.Stat(cfp); err == nil {
// 				controlFileMap[cf] = cfp
// 			}
// 		}
// 		fsOrDie(controlFs.BindSafe("/", vfs.SafeMap(controlFileMap), "/", vfs.BindBefore))
// 	}

// 	if c.Files != nil && len(c.Files) > 0 {
// 		fsOrDie(dataFs.BindSafe("/", vfs.SafeFileMap(c.Files), "/", vfs.BindAfter))
// 	}

// 	controlFileMap := make(map[string]string)
// maint:
// 	for dst, src := range map[string]string{
// 		"preinst":  c.Preinst,
// 		"postinst": c.Postinst,
// 		"prerm":    c.Prerm,
// 		"postrm":   c.Postrm,
// 	} {
// 		src = strings.TrimSpace(src)
// 		if src == "" {
// 			continue maint
// 		}
// 		controlFileMap[dst] = src
// 	}

// 	{
// 		var conffilesBuf bytes.Buffer
// 		err := vfs.Walk("/", dataFs, func(p string, info os.FileInfo, err error) error {
// 			// log.Println("p", p, info)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			if info.IsDir() {
// 				err := d.AddEmptyDirectory(p)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				return nil
// 			}
// 			if opr, ok := info.(vfs.OSPather); ok {
// 				op := opr.OSPath()
// 				err := d.AddFile(op, p)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			} else {
// 				log.Fatalln("expected all files to be of OSPather type!", p)
// 			}

// 			for _, pattern := range c.Conffiles {
// 				ok, err := filepath.Match(pattern, p)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				if ok {
// 					fmt.Fprintln(&conffilesBuf, p)
// 				}
// 			}
// 			return nil
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if conffilesBuf.Len() > 0 {
// 			conffilesBuf.WriteString("\n")
// 			controlFileMap["conffiles"] = conffilesBuf.String()
// 		}
// 		if len(controlFileMap) > 0 {
// 			fsOrDie(controlFs.BindSafe("/", vfs.SafeMap(controlFileMap), "/", vfs.BindBefore))
// 		}
// 	}
// 	log.Println(d.GetFilename())
// 	if err := d.Write(filepath.Join(c.Target, d.GetFilename())); err != nil {
// 		log.Fatal(err)
// 	}
// 	d.Close()
// }

// Files .
type Files struct {
	Auto     string // Defaults to "contrib/debian"
	Control  string
	Preinst  string
	Postinst string
	Prerm    string
	Postrm   string
	// Files     plugins.StringMapFlag
	// Conffiles plugins.StringSliceFlag
}

type ArchFlag dependency.Arch

func (s *ArchFlag) String() string {
	return dependency.Arch(*s).String()
}

func (s *ArchFlag) Set(value string) error {
	if value == "386" {
		value = "i386"
	}
	a, err := dependency.ParseArch(value)
	if err != nil {
		return err
	}
	*s = ArchFlag(*a)
	return nil
}

// StringSliceFlag is a flag type which
type DependencyFlag dependency.Dependency

func (s *DependencyFlag) String() string {
	str, err := dependency.Dependency(*s).MarshalControl()
	if err != nil {
		plug.Println(err)
		return ""
	}
	return str

}

func (s *DependencyFlag) Set(value string) error {
	d, err := dependency.Parse(value)
	if err != nil {
		return err
	}
	*s = DependencyFlag(*d)
	return nil
}

func (d *DependencyFlag) Value() dependency.Dependency {
	return dependency.Dependency(*d)
}

type VersionFlag version.Version

func (s *VersionFlag) String() string {
	return version.Version(*s).String()
}

func (s *VersionFlag) Set(value string) error {
	d, err := version.Parse(value)
	if err != nil {
		return err
	}
	*s = VersionFlag(d)
	return nil
}

func (d *VersionFlag) Value() version.Version {
	return version.Version(*d)
}
