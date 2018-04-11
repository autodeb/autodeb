# autopkgupdate

autopkgupdate tries to automatically update Debian packages to newer upstream versions or to backport them.

autopkgupdate is the concretization Lucas Nussbaum's GSOC 2018 proposed project titled "Automatic Packages for Everything (backports, new upstream versions, etc.)". The [project proposal](https://wiki.debian.org/SummerOfCode2018/Projects) can be found in the Debian Wiki.

## Available scripts

- ``list-packages-with-newer-upstream-versions``: lists source packages that have newer upstream versions available

- ``update-random-package``: find a package that needs updating to a newer upstream version and try updating it automatically. On success, the output is moved to the current directory.

## Getting started

### 1. Clone the project

```shell
$ export GOPATH=~/go
$ mkdir -p $GOPATH/src/salsa.debian.org/aviau/
$ git clone git@salsa.debian.org:aviau/autopkgupdate.git $GOPATH/src/salsa.debian.org/aviau/autopkgupdate
```

### 2. Build the project

#### Building inside a docker container

This is useful if you don't want to spend time setting up a Go environement.

```shell
$ make docker-all
```

#### Building without a docker container

##### Setup Go

```shell
$ sudo apt install golang-go
```

##### Build the scripts

```shell
$ cd $GOPATH/src/salsa.debian.org/aviau/autopkgupdate
$ make all
```

### 3. Run any of the scripts

Note that runtime dependencies of the scripts include:
 - devscripts
 - sbuild

```shell
$ ./list-packages-with-newer-upstream-versions
$ ./update-random-package
```
