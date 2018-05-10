package uscan

//Possible Result values
const (
	ResultStatusNewerPackageAvailable = "newer package available"
)

// Result of the uscan command. Uscan produces this XML output when ran with
// the --dehs flag.
type Result struct {
	Package         string `xml:"package"`
	DebianUVersion  string `xml:"debian-uversion"`
	UpstreamVersion string `xml:"upstream-version"`
	UpstreamURL     string `xml:"upstream-url"`
	Status          string `xml:"status"`
	Target          string `xml:"target"`
	TargetPath      string `xml:"target-path"`
	Message         string `xml:"Message"`
}
