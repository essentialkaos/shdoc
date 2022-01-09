package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/path"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/update"

	. "github.com/essentialkaos/shdoc/parser"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SHDoc"
	VER  = "0.8.1"
	DESC = "Tool for viewing and exporting docs for shell scripts"
)

const (
	OPT_OUTPUT   = "o:output"
	OPT_TEMPLATE = "t:template"
	OPT_NAME     = "n:name"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_OUTPUT:   {},
	OPT_TEMPLATE: {Value: "html"},
	OPT_NAME:     {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Init is main function
func Init() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		fmtc.Println("Arguments parsing errors:")

		for _, err := range errs {
			fmtc.Printf("  %s\n", err.Error())
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	switch len(args) {
	case 1:
		process(args[0], "")
	case 2:
		process(args[0], args[1])
	default:
		showUsage()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// process starts source processing
func process(file string, pattern string) {
	if !fsutil.IsExist(file) {
		printErrorAndExit("File %s does not exist", file)
	}

	if !fsutil.IsReadable(file) {
		printErrorAndExit("File %s is not readable", file)
	}

	if !fsutil.IsNonEmpty(file) {
		printErrorAndExit("File %s is empty", file)
	}

	doc, errs := Parse(file)

	if len(errs) != 0 {
		printError("Shell script docs parsing errors:")

		for _, err := range errs {
			printError("  %s", err.Error())
		}

		os.Exit(1)
	}

	if !doc.IsValid() {
		printWarn("File %s doesn't contains documentation", file)
		os.Exit(2)
	}

	if options.GetS(OPT_NAME) != "" {
		doc.Title = options.GetS(OPT_NAME)
	}

	if options.GetS(OPT_OUTPUT) == "" {
		if pattern == "" {
			simpleRender(doc)
		} else {
			findInfo(doc, pattern)
		}
	} else {
		renderTemplate(doc)
	}
}

// findInfo searches geven pattern in entities names
func findInfo(doc *Document, pattern string) {
	fmtc.NewLine()

	if doc.Constants != nil {
		for _, c := range doc.Constants {
			if strings.Contains(c.Name, pattern) {
				renderConstant(c)
				fmtc.NewLine()
			}
		}
	}

	if doc.Variables != nil {
		for _, v := range doc.Variables {
			if strings.Contains(v.Name, pattern) {
				renderVariable(v)
				fmtc.NewLine()
			}
		}
	}

	if doc.Methods != nil {
		for _, m := range doc.Methods {
			if strings.Contains(m.Name, pattern) {
				renderMethod(m, true)
				fmtc.NewLine()
			}
		}
	}
}

// simpleRender prints all document info to console
func simpleRender(doc *Document) {
	if doc.HasAbout() {
		fmtutil.Separator(false, "ABOUT")

		for _, l := range doc.About {
			fmtc.Printf("  %s\n", l)
		}
	}

	if doc.HasConstants() {
		fmtutil.Separator(false, "CONSTANTS")

		totalConstants := len(doc.Constants)

		for i, c := range doc.Constants {
			renderConstant(c)

			if i < totalConstants-1 {
				fmtc.NewLine()
			}
		}
	}

	if doc.HasVariables() {
		fmtutil.Separator(false, "GLOBAL VARIABLES")

		totalVariables := len(doc.Variables)

		for i, v := range doc.Variables {
			renderVariable(v)

			if i < totalVariables-1 {
				fmtc.NewLine()
			}
		}
	}

	if doc.HasMethods() {
		fmtutil.Separator(false, "METHODS")

		totalMethods := len(doc.Methods)

		for i, m := range doc.Methods {
			renderMethod(m, false)

			if i < totalMethods-1 {
				fmtc.NewLine()
				fmtc.NewLine()
			}
		}
	}

	fmtutil.Separator(false)
}

// renderTemplate reads template and render to file
func renderTemplate(doc *Document) {
	projectDir := os.Getenv("GOPATH")
	templateFile := path.Join(
		projectDir, "src/github.com/essentialkaos/shdoc/templates",
		options.GetS(OPT_TEMPLATE)+".tpl",
	)

	if !fsutil.CheckPerms("FRS", templateFile) {
		printErrorAndExit("Can't read template %s - file does not exist or empty", templateFile)
	}

	outputFile := options.GetS(OPT_OUTPUT)

	if fsutil.IsExist(outputFile) {
		os.Remove(outputFile)
	}

	fd, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	defer fd.Close()

	tpl, err := ioutil.ReadFile(templateFile)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	t := template.New("Template")
	t, err = t.Parse(string(tpl[:]))

	err = t.Execute(fd, doc)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	fmtutil.Separator(false, doc.Title)

	fmtc.Printf("  {*}Constants:{!} %d\n", len(doc.Constants))
	fmtc.Printf("  {*}Variables:{!} %d\n", len(doc.Variables))
	fmtc.Printf("  {*}Methods:{!}   %d\n", len(doc.Methods))
	fmtc.NewLine()
	fmtc.Printf(
		"  {*}Output:{!} %s {s-}(%s){!}\n", outputFile,
		fmtutil.PrettySize(fsutil.GetSize(outputFile)),
	)

	fmtutil.Separator(false)
}

// renderConstant prints constant info to console
func renderConstant(c *Variable) {
	fmtc.Printf("{s-}%4d:{!} {m*}%s{!} {s}={!} %s "+getVarTypeDesc(c.Type)+"\n", c.Line, c.Name, c.Value)
	fmtc.Printf("      %s\n", c.UnitedDesc())
}

// renderMethod prints variable info to console
func renderVariable(v *Variable) {
	fmtc.Printf("{s-}%4d:{!} {c*}%s{!} {s}={!} %s "+getVarTypeDesc(v.Type)+"\n", v.Line, v.Name, v.Value)
	fmtc.Printf("      %s\n", v.UnitedDesc())
}

// renderMethod prints method info to console
func renderMethod(m *Method, showExamples bool) {
	fmtc.Printf("{s-}%4d:{!} {b*}%s{!} {s}-{!} %s\n", m.Line, m.Name, m.UnitedDesc())

	if len(m.Arguments) != 0 {
		fmtc.NewLine()

		for _, a := range m.Arguments {
			switch {
			case a.IsOptional:
				fmtc.Printf("  {s-}%2s.{!} %s "+getVarTypeDesc(a.Type)+" {s-}[Optional]{!}\n", a.Index, a.Desc)
			case a.IsWildcard:
				fmtc.Printf("  {s-}%2s.{!} %s\n", a.Index, a.Desc)
			default:
				fmtc.Printf("  {s-}%2s.{!} %s "+getVarTypeDesc(a.Type)+"\n", a.Index, a.Desc)
			}
		}
	}

	if m.ResultCode {
		fmtc.NewLine()
		fmtc.Printf("  {*}Code:{!} 0 - ok, 1 - not ok\n")
	}

	if m.ResultEcho != nil {
		fmtc.NewLine()
		fmtc.Printf("  {*}Echo:{!} %s "+getVarTypeDesc(m.ResultEcho.Type)+"\n", strings.Join(m.ResultEcho.Desc, " "))
	}

	if m.Example != nil && showExamples {
		fmtc.NewLine()
		fmtc.Println("  {*}Example:{!}")
		fmtc.NewLine()

		for _, l := range m.Example {
			fmtc.Printf("    %s\n", l)
		}
	}
}

// getVarTypeDesc returns type description
func getVarTypeDesc(t VariableType) string {
	switch t {
	case VAR_TYPE_STRING:
		return "{b}(String){!}"
	case VAR_TYPE_NUMBER:
		return "{y}(Number){!}"
	case VAR_TYPE_BOOLEAN:
		return "{g}(Boolean){!}"
	default:
		return ""
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit prints error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "file")

	info.AddOption(OPT_OUTPUT, "Path to output file", "file")
	info.AddOption(OPT_TEMPLATE, "Name of template", "name")
	info.AddOption(OPT_NAME, "Overwrite default name", "name")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"script.sh",
		"Parse shell script and show docs in console",
	)

	info.AddExample(
		"script.sh -t markdown -o my_script.md",
		"Parse shell script and save docs using given export template",
	)

	info.AddExample(
		"script.sh someEntity",
		"Parse shell script and show docs for some constant, variable or method",
	)

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "shdoc"))
	case "fish":
		fmt.Printf(fish.Generate(info, "shdoc"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "shdoc"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout shows info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/shdoc", update.GitHubChecker},
	}

	about.Render()
}
