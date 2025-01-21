package flagparser

import (
	"flag"
	"fmt"
	resource "github.com/SongZihuan/auto-aliyun-cdn-ssl"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"io"
	"reflect"
	"strings"
)

const OptionIdent = "  "
const OptionPrefix = "--"
const UsagePrefixWidth = 10

type flagData struct {
	flagReady  bool
	flagSet    bool
	flagParser bool

	HelpData               bool
	HelpName               string
	HelpUsage              string
	VersionData            bool
	VersionName            string
	VersionUsage           string
	LicenseData            bool
	LicenseName            string
	LicenseUsage           string
	ReportData             bool
	ReportName             string
	ReportUsage            string
	ConfigFileData         string
	ConfigFileName         string
	ConfigFileUsage        string
	NotAutoReloadData      bool
	NotAutoReloadName      string
	NotAutoReloadShortName string
	NotAutoReloadUsage     string

	Usage string
}

func initData() {
	data = flagData{
		flagReady:  false,
		flagSet:    false,
		flagParser: false,

		HelpData:               false,
		HelpName:               "help",
		HelpUsage:              fmt.Sprintf("Show usage of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		VersionData:            false,
		VersionName:            "version",
		VersionUsage:           fmt.Sprintf("Show version of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		LicenseData:            false,
		LicenseName:            "license",
		LicenseUsage:           fmt.Sprintf("Show license of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		ReportData:             false,
		ReportName:             "report",
		ReportUsage:            fmt.Sprintf("Show how to report questions/errors of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		ConfigFileData:         "config.yaml",
		ConfigFileName:         "config",
		ConfigFileUsage:        fmt.Sprintf("%s", "The location of the running configuration file of the backend service. The option is a string, the default value is config.yaml in the running directory."),
		NotAutoReloadData:      false,
		NotAutoReloadName:      "not-auto-reload",
		NotAutoReloadShortName: "",
		NotAutoReloadUsage:     fmt.Sprintf("%s", "Disable automatic detection of configuration file changes and reloading of system programs. This feature is enabled by default. This feature consumes a certain amount of performance. If your performance is not enough, you can choose to disable it."),
		Usage:                  "",
	}

	data.ready()
}

func (d *flagData) writeUsage() {
	if len(d.Usage) != 0 {
		return
	}

	if d.isFlagSet() || d.isFlagParser() {
		panic("flag is parser")
	}

	var result strings.Builder
	result.WriteString(utils.FormatTextToWidth(fmt.Sprintf("Usage of %s:", utils.GetArgs0Name()), utils.NormalConsoleWidth))
	result.WriteString("\n")

	val := reflect.ValueOf(*d)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		if !strings.HasSuffix(field.Name, "Data") {
			continue
		}

		option := field.Name[:len(field.Name)-4]
		optionName := ""
		optionShortName := ""
		optionUsage := ""

		if utils.HasFieldByReflect(typ, option+"Name") {
			var ok bool
			optionName, ok = val.FieldByName(option + "Name").Interface().(string)
			if !ok {
				panic("can not get option name")
			}
		}

		if utils.HasFieldByReflect(typ, option+"ShortName") {
			var ok bool
			optionShortName, ok = val.FieldByName(option + "ShortName").Interface().(string)
			if !ok {
				panic("can not get option short name")
			}
		} else if len(optionName) > 1 {
			optionShortName = optionName[:1]
		}

		if utils.HasFieldByReflect(typ, option+"Usage") {
			var ok bool
			optionUsage, ok = val.FieldByName(option + "Usage").Interface().(string)
			if !ok {
				panic("can not get option usage")
			}
		}

		var title string
		var title1 string
		var title2 string
		if field.Type.Kind() == reflect.Bool {
			var optionData bool
			if utils.HasFieldByReflect(typ, option+"Data") {
				var ok bool
				optionData, ok = val.FieldByName(option + "Data").Interface().(bool)
				if !ok {
					panic("can not get option data")
				}
			}

			if optionData == true {
				panic("bool option can not be true")
			}

			if optionName != "" {
				title1 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(optionName, utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			}

			if optionShortName != "" {
				title2 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(optionShortName, utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			}
		} else if field.Type.Kind() == reflect.String {
			var optionData string
			if utils.HasFieldByReflect(typ, option+"Data") {
				var ok bool
				optionData, ok = val.FieldByName(option + "Data").Interface().(string)
				if !ok {
					panic("can not get option data")
				}
			}

			if optionName != "" && optionData != "" {
				title1 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s string, default: '%s'", optionName, optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			} else if optionName != "" && optionData == "" {
				title1 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s string", optionName), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			}

			if optionShortName != "" && optionData != "" {
				title2 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s string, default: '%s'", optionShortName, optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			} else if optionShortName != "" && optionData == "" {
				title2 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s string", optionShortName), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			}
		} else if field.Type.Kind() == reflect.Uint {
			var optionData uint
			if utils.HasFieldByReflect(typ, option+"Data") {
				var ok bool
				optionData, ok = val.FieldByName(option + "Data").Interface().(uint)
				if !ok {
					panic("can not get option data")
				}
			}

			if optionName != "" && optionData != 0 {
				title1 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s number, default: %d", optionName, optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			} else if optionName != "" && optionData == 0 {
				title1 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s number", optionName), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			}

			if optionShortName != "" && optionData != 0 {
				title2 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s number, default: %d", optionShortName, optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			} else if optionShortName != "" && optionData == 0 {
				title2 = fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s number", optionShortName), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			}
		} else {
			panic("error flag type")
		}

		if title1 == "" && title2 == "" {
			continue
		} else if title1 != "" && title2 == "" {
			title = title1
		} else if title1 == "" {
			title = title2
		} else {
			title = fmt.Sprintf("%s\n%s", title1, title2)
		}

		result.WriteString(title)
		result.WriteString("\n")

		usegae := utils.FormatTextToWidthAndPrefix(optionUsage, UsagePrefixWidth, utils.NormalConsoleWidth)
		result.WriteString(usegae)
		result.WriteString("\n\n")
	}

	d.Usage = strings.TrimRight(result.String(), "\n")
}

func (d *flagData) setFlag() {
	if d.isFlagSet() {
		return
	}

	flag.BoolVar(&d.HelpData, data.HelpName, data.HelpData, data.HelpUsage)
	flag.BoolVar(&d.HelpData, data.HelpName[0:1], data.HelpData, data.HelpUsage)

	flag.BoolVar(&d.VersionData, data.VersionName, data.VersionData, data.VersionUsage)
	flag.BoolVar(&d.VersionData, data.VersionName[0:1], data.VersionData, data.VersionUsage)

	flag.BoolVar(&d.LicenseData, data.LicenseName, data.LicenseData, data.LicenseUsage)
	flag.BoolVar(&d.LicenseData, data.LicenseName[0:1], data.LicenseData, data.LicenseUsage)

	flag.BoolVar(&d.ReportData, data.ReportName, data.ReportData, data.ReportUsage)
	flag.BoolVar(&d.ReportData, data.ReportName[0:1], data.ReportData, data.ReportUsage)

	flag.StringVar(&d.ConfigFileData, data.ConfigFileName, data.ConfigFileData, data.ConfigFileUsage)
	flag.StringVar(&d.ConfigFileData, data.ConfigFileName[0:1], data.ConfigFileData, data.ConfigFileUsage)

	flag.Usage = func() {
		_, _ = d.PrintUsage()
	}
	d.flagSet = true
}

func (d *flagData) parser() {
	if d.flagParser {
		return
	}

	if !d.isFlagSet() {
		panic("flag not set")
	}

	flag.Parse()
	d.flagParser = true
}

func (d *flagData) ready() {
	if d.isReady() {
		return
	}

	d.writeUsage()
	d.setFlag()
	d.parser()
	d.flagReady = true
}

func (d *flagData) isReady() bool {
	return d.isFlagSet() && d.isFlagParser() && d.flagReady
}

func (d *flagData) isFlagSet() bool {
	return d.flagSet
}

func (d *flagData) isFlagParser() bool {
	return d.flagParser
}

func (d *flagData) Help() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.HelpData
}

func (d *flagData) FprintUsage(writer io.Writer) (int, error) {
	return fmt.Fprintf(writer, "%s\n", d.Usage)
}

func (d *flagData) PrintUsage() (int, error) {
	return d.FprintUsage(flag.CommandLine.Output())
}

func (d *flagData) Version() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.VersionData
}

func (d *flagData) FprintVersion(writer io.Writer) (int, error) {
	version := utils.FormatTextToWidth(fmt.Sprintf("Version of %s: %s", utils.GetArgs0Name(), resource.Version), utils.NormalConsoleWidth)
	return fmt.Fprintf(writer, "%s\n", version)
}

func (d *flagData) PrintVersion() (int, error) {
	return d.FprintVersion(flag.CommandLine.Output())
}

func (d *flagData) FprintLicense(writer io.Writer) (int, error) {
	title := utils.FormatTextToWidth(fmt.Sprintf("License of %s:", utils.GetArgs0Name()), utils.NormalConsoleWidth)
	license := utils.FormatTextToWidth(resource.License, utils.NormalConsoleWidth)
	return fmt.Fprintf(writer, "%s\n%s\n", title, license)
}

func (d *flagData) PrintLicense() (int, error) {
	return d.FprintLicense(flag.CommandLine.Output())
}

func (d *flagData) FprintReport(writer io.Writer) (int, error) {
	// 不需要title
	report := utils.FormatTextToWidth(resource.Report, utils.NormalConsoleWidth)
	return fmt.Fprintf(writer, "%s\n", report)
}

func (d *flagData) PrintReport() (int, error) {
	return d.FprintReport(flag.CommandLine.Output())
}

func (d *flagData) FprintLF(writer io.Writer) (int, error) {
	return fmt.Fprintf(writer, "\n")
}

func (d *flagData) PrintLF() (int, error) {
	return d.FprintLF(flag.CommandLine.Output())
}

func (d *flagData) License() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.LicenseData
}

func (d *flagData) Report() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.ReportData
}

func (d *flagData) ConfigFile() string {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.ConfigFileData
}

func (d *flagData) SetOutput(writer io.Writer) {
	flag.CommandLine.SetOutput(writer)
}
