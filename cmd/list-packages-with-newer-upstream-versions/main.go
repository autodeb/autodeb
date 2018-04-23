package main

import (
	"fmt"

	"salsa.debian.org/autodeb-team/autodeb/internal/udd"
)

func main() {
	packages, _ := udd.PackagesWithNewerUpstreamVersions()

	for _, pkg := range packages {
		fmt.Println(pkg.Package)
	}
}
