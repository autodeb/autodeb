# autopkgupdate

autopkgupdate tries to automatically update Debian packages to newer upstream versions or to backport them.

autopkgupdate is the concretization Lucas Nussbaum's GSOC 2018 proposed project titled "Automatic Packages for Everything (backports, new upstream versions, etc.)". The [project proposal](https://wiki.debian.org/SummerOfCode2018/Projects) can be found in the Debian Wiki.

## Available scripts

- ``list-packages-with-newer-upstream-versions``: lists source packages that have newer upstream versions available

## Getting started

### Setup Go

```shell
$ sudo apt install golang-go
$ export GOPATH=~/go
```

### Get autopkgupdate

```shell
$ mkdir -p $GOPATH/src/salsa.debian.org/aviau/
$ git clone git@salsa.debian.org:aviau/autopkgupdate.git $GOPATH/src/salsa.debian.org/aviau/autopkgupdate
```

### Build the scripts

```shell
$ cd $GOPATH/src/salsa.debian.org/aviau/autopkgupdate
$ make list-packages-to-update
```
