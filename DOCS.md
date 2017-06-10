Use this plugin to quickly create debian dpkg (.deb) packages.

## Config

The following parameters are used to configure the plugin:

### Required fields

* **package** - The name of your package
* **version** - Must adhere to debian version syntax
* **architecture** - CPU arch for your binaries, or "all"

Supported values are: all, amd64, arm64, armel, armhf, i386, mips, mipsel,
powerpc, ppc64el, s390xa. Additionally, the value 386 will automatically be
translated to i386 (for Go/GOARCH).

* **maintainer** - Your Name <email@example.com>
* **description** - Brief explanation of your package

### Optional fields

* **depends** - Other packages you depend on. E.g: "python" or "curl (>= 7.0.0)"
* **conflicts** - Packages your package are not compatible with
* **breaks** - Packages your package breaks
* **replaces** - Packages your package replaces
* **section** - section (default "default")
* **priority** - priority (default "extra")
* **homepage** - URL to your project homepage or source repository, if you have one

For more details on how to specify various config options, refer to the
debian package specification:

- https://www.debian.org/doc/debian-policy/ch-controlfields.html
- https://www.debian.org/doc/manuals/debian-faq/ch-pkg_basics.en.html

### File and path fields

* **auto_path** - auto path (default "deb-pkg")

drone-deb will automatically include all files under **auto_path**. For
example, the following files will be automatically included and installed to
their corresponding paths:

    deb-pkg/etc/mysqld/my.conf  -> /etc/mysqld/my.conf
    deb-pkg/usr/bin/mysqld      -> /usr/bin/mysqld

You can override this behavior by setting autoPath to - (dash character) and /
or by using the Files map to create a custom source -> dest mapping.

* **files** - Additional files to add from outside of the **auto_path**.

### Control script fields

By default drone-deb will use any of these files if they are present at the
root **auto_path** root level:

- preinst
- postinst
- prerm
- postrm

You can override this behavior by setting the following plugin options:

* **preinst** - preinst script
* **postinst** - postinst script
* **prerm** - prerm script
* **postrm** - postrm script

### Build options fields

* **preserve_symlinks** - By default contents of symlink targets are copied. This
    option writes symlinks to the archive instead
* **upgrade_configs** - Indicates whether apt should replace files under /etc when
    installing a new package version. By default these files are not upgraded
* **target** - target directory to create .deb file in

## Examples

Common examples of package builds:


```yaml
pipeline:

  simple:
    image: plugins/deb
    description: This is a simple package example only using auto_path for everything.
    package: simple-example
    version: 0.0.1
    maintainer: Thomas Frössman<thomasf@jossystem.se>
    homepage: https://example.com
    architecture: all
    auto_path: contrib/debian/deb-pkg/

  no_auto:
    description: Same as the simple example without auto_path.
    image: plugins/deb
    package: no-auto-example
    version: 0.0.2
    maintainer: Thomas Frössman<thomasf@jossystem.se>
    architecture: all
    postrm: contrib/debian/deb-pkg/postrm
    files:
      contrib/debian/deb-pkg/etc/testdata.ini: etc/testdata.ini

  files-example:
    description: An example package which uses files to specify adding additional files and control scripts from outside of the auto_path.
    image: plugins/deb
    package: files-example
    version: 1.0.0
    maintainer: Thomas Frössman<thomasf@jossystem.se>
    architecture: all
    auto_path: contrib/debian/deb-pkg/
    postinst: contrib/debian/postinst
    files:
      dirA/somefile: usr/share/my-package-somefile
      .drone.yml: usr/share/my-package-drone.yml

  dependencies-example:
    description: An example package which uses files to specify adding additional files and control scripts from outside of the auto_path.
    image: plugins/deb
    package: deps-example
    version: 1.0.0
    maintainer: Thomas Frössman<thomasf@jossystem.se>
    architecture: all
    auto_path: contrib/debian/deb-pkg/
    depends:
      - python (>= 3.5)
      - go (>= 1.8)
    breaks:
      - simple-example
    replaces:
      - files-example
```
