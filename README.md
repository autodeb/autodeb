# autopkgupdate

autopkgupdate tries to automatically update Debian packages to newer upstream versions or to backport them.

autopkgupdate is the concretization Lucas Nussbaum's GSOC 2018 proposed project titled "Automatic Packages for Everything (backports, new upstream versions, etc.)". The [project proposal](https://wiki.debian.org/SummerOfCode2018/Projects) can be found in the Debian Wiki.

## Available scripts

- ``list-packages-with-newer-upstream-versions``: lists source packages that have newer upstream versions available

- ``update-random-package``: find a package that needs updating to a newer upstream version and try updating it automatically. On success, the output is moved to the current directory.

## Getting started

### 1. Setup Go

```shell
$ sudo apt install golang-go
```

### 2. Clone the project

```shell
$ export GOPATH=~/go
$ mkdir -p $GOPATH/src/salsa.debian.org/aviau/
$ git clone git@salsa.debian.org:aviau/autopkgupdate.git $GOPATH/src/salsa.debian.org/aviau/autopkgupdate
$ cd $GOPATH/src/salsa.debian.org/aviau/autopkgupdate
```

### 3. Build the project

```shell
$ make all
```

### 4. Run any of the scripts

Note that runtime dependencies of the scripts include:
 - devscripts
 - sbuild

```shell
$ ./list-packages-with-newer-upstream-versions
$ ./update-random-package
```
