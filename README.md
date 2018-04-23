# autodeb

autodeb tries to automatically update Debian packages to newer upstream versions or to backport them.

autodeb is the concretization Lucas Nussbaum's GSOC 2018 proposed project titled "Automatic Packages for Everything (backports, new upstream versions, etc.)". The [project proposal](https://wiki.debian.org/SummerOfCode2018/Projects) can be found in the Debian Wiki.

## Available executables

- ``list-packages-with-newer-upstream-versions``: lists source packages that have newer upstream versions available

- ``update-random-package``: find a package that needs updating to a newer upstream version and try updating it automatically. On success, the output is moved to the current directory.

- ``autodeb-server``: This is the server component of the system. It provides a web interface, a REST API and dput-compatible interface.

- ``autodeb-runner``: TODO. This executable is not yet available.

## Getting started

### 1. Setup Go

Note that you might want to get a recent version of the go compiler from a backports repository.

```shell
$ apt-get install golang-go git make
$ export GOPATH=~/go
$ go get -u golang.org/x/lint/golint
```

### 2. Clone the project

```shell
$ mkdir -p $GOPATH/src/salsa.debian.org/autodeb-team/
$ git clone https://salsa.debian.org/autodeb-team/autodeb.git $GOPATH/src/salsa.debian.org/autodeb-team/autodeb
$ cd $GOPATH/src/salsa.debian.org/autodeb-team/autodeb
```

### 3. Build the project

```shell
$ make get-deps
$ make
```

### 4. Run any of the scripts

Note that runtime dependencies of the scripts include:
 - devscripts
 - sbuild

```shell
$ ./list-packages-with-newer-upstream-versions
$ ./update-random-package
$ ./autodeb-server
```
