package main

import (
	"fmt"

	"salsa.debian.org/aviau/autopkgupdate/internal/udd"
)

func main() {
	packages, _ := udd.PackagesWithNewerUpstreamVersions()

	for _, pkg := range packages {
		fmt.Println(pkg.Package)
	}
}
