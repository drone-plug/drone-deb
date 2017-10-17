package deb

import (
	"fmt"
	"strconv"
	"strings"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

// control .
type BinaryControl struct {
	control.Paragraph

	Package string `required:"true"`
	// Source string
	Version  version.Version `required:"true"`
	Section  string
	Priority string
	Arch     dependency.Arch `control:"Architecture" required:"true"`
	// Essential bool
	// InstalledSize int             `control:"Installed-Size"`
	Homepage    string
	Description string `required:"true"`
	Depends     dependency.Dependency
	Recommends  dependency.Dependency
	Maintainer  string `required:"true"`
	Suggests    dependency.Dependency
	Breaks      dependency.Dependency
	Replaces    dependency.Dependency
	Conflicts   dependency.Dependency
	PreDepends  dependency.Dependency `control:"Pre-Depends"`
}

func (b *BinaryControl) SetInstalledSize(bytes int) {
	// Convert size from bytes to kilobytes. If there is a remainder, round up.
	if bytes%1024 > 0 {
		bytes = bytes/1024 + 1
	} else {
		bytes = bytes / 1024
	}

	b.Set("Installed-Size", strconv.Itoa(bytes))

}

// MissingFieldsError .
type MissingFieldsError struct {
	Names []string
	Refs  []interface{}
}

func (m MissingFieldsError) Error() string {
	return fmt.Sprintf("These required fields are missing: %s", strings.Join(m.Names, ", "))
}

func (p *BinaryControl) Validate() error {
	missing := []string{}
	var refs []interface{}
	if p.Package == "" {
		missing = append(missing, "package")
		refs = append(refs, &p.Package)

	}
	var v version.Version
	if p.Version == v {
		missing = append(missing, "version")
	}
	var a dependency.Arch
	if p.Arch == a {
		missing = append(missing, "arch")
	}
	if p.Maintainer == "" {
		missing = append(missing, "maintainer")
	}
	if p.Description == "" {
		missing = append(missing, "description")
	}
	if len(missing) > 0 {
		err := MissingFieldsError{
			Names: missing,
			Refs: refs,
		}
		return err
	}
	return nil
}

// control .
type BinaryControlTemplate struct {
	control.Paragraph

	Package     string
	Version     version.Version
	Arch        dependency.Arch `control:"Architecture"`
	Maintainer  string
	Section     string
	Priority    string
	Homepage    string
	Description string
	Depends     dependency.Dependency
	Recommends  dependency.Dependency
	Suggests    dependency.Dependency
	Breaks      dependency.Dependency
	Replaces    dependency.Dependency
	Conflicts   dependency.Dependency
	PreDepends  dependency.Dependency `control:"Pre-Depends"`
}

func (t BinaryControlTemplate) AddMissing(b *BinaryControl) {
	var a dependency.Arch
	for _, f := range []struct {
		tpl dependency.Arch
		dst *dependency.Arch
	}{
		{t.Arch, &b.Arch},
	} {
		if f.tpl != a && *f.dst == a {
			*f.dst = f.tpl
		}
	}
	for _, f := range []struct {
		tpl dependency.Dependency
		dst *dependency.Dependency
	}{
		{t.Depends, &b.Depends},
		{t.Recommends, &b.Recommends},
		{t.Suggests, &b.Suggests},
		{t.Breaks, &b.Breaks},
		{t.Replaces, &b.Replaces},
		{t.Conflicts, &b.Conflicts},
		{t.PreDepends, &b.PreDepends},
	} {
		if len(f.tpl.Relations) > 0 && len(f.dst.Relations) == 0 {
			*f.dst = f.tpl
		}
	}
	var v version.Version
	for _, f := range []struct {
		tpl version.Version
		dst *version.Version
	}{
		{t.Version, &b.Version},
	} {
		if f.tpl != v && *f.dst == v {
			*f.dst = f.tpl
		}
	}
	for _, f := range []struct {
		tpl string
		dst *string
	}{
		{t.Package, &b.Package},
		{t.Maintainer, &b.Maintainer},
		{t.Section, &b.Section},
		{t.Priority, &b.Priority},
		{t.Homepage, &b.Homepage},
		{t.Description, &b.Description},
	} {
		if f.tpl != "" && *f.dst == "" {
			*f.dst = f.tpl
		}
	}
}
