package flagparser

import (
	"fmt"
	"io"
	"strings"
)

var data flagData

func Help() bool {
	return data.Help()
}

func FprintUsage(writer io.Writer) (int, error) {
	return data.FprintUsage(writer)
}

func PrintUsage() (int, error) {
	return data.PrintUsage()
}

func FprintVersion(writer io.Writer) (int, error) {
	return data.FprintVersion(writer)
}

func PrintVersion() (int, error) {
	return data.PrintVersion()
}

func FprintLicense(writer io.Writer) (int, error) {
	return data.FprintLicense(writer)
}

func PrintLicense() (int, error) {
	return data.PrintLicense()
}

func FprintReport(writer io.Writer) (int, error) {
	return data.FprintReport(writer)
}

func PrintReport() (int, error) {
	return data.PrintReport()
}

func FprintLF(writer io.Writer) (int, error) {
	return data.FprintLF(writer)
}

func PrintLF() (int, error) {
	return data.PrintLF()
}

func Version() bool {
	return data.Version()
}

func License() bool {
	return data.License()
}

func Report() bool {
	return data.Report()
}

func NotRunMode() bool {
	return Help() || Version() || License() || Report()
}

func NotRunModeOption() string {
	if !NotRunMode() {
		return ""
	}

	var result strings.Builder

	if data.Help() {
		result.WriteString(fmt.Sprintf("%s%s, ", OptionPrefix, data.HelpName))
	}

	if data.Version() {
		result.WriteString(fmt.Sprintf("%s%s, ", OptionPrefix, data.VersionName))
	}

	if data.License() {
		result.WriteString(fmt.Sprintf("%s%s, ", OptionPrefix, data.LicenseName))
	}

	if data.Report() {
		result.WriteString(fmt.Sprintf("%s%s, ", OptionPrefix, data.ReportName))
	}

	return strings.TrimSuffix(result.String(), ", ")
}

func ConfigFile() string {
	return data.ConfigFile()
}

func NotRunAutoReload() bool {
	return data.NotAutoReloadData
}

func RunAutoReload() bool {
	return !NotRunAutoReload()
}

func SetOutput(writer io.Writer) {
	data.SetOutput(writer)
}
