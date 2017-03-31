package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"pkg.re/essentialkaos/ek.v7/arg"
	"pkg.re/essentialkaos/ek.v7/env"
	"pkg.re/essentialkaos/ek.v7/fmtc"
	"pkg.re/essentialkaos/ek.v7/fmtutil"
	"pkg.re/essentialkaos/ek.v7/fsutil"
	"pkg.re/essentialkaos/ek.v7/path"
	"pkg.re/essentialkaos/ek.v7/usage"
	"pkg.re/essentialkaos/ek.v7/usage/update"

	. "github.com/essentialkaos/shdoc/parser"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SHDoc"
	VER  = "0.3.2"
	DESC = "Tool for viewing and exporting docs for shell scripts"
)

const (
	ARG_OUTPUT   = "o:output"
	ARG_TEMPLATE = "t:template"
	ARG_NAME     = "n:name"
	ARG_NO_COLOR = "nc:no-color"
	ARG_HELP     = "h:help"
	ARG_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var argMap = arg.Map{
	ARG_OUTPUT:   {},
	ARG_TEMPLATE: {Value: "html"},
	ARG_NAME:     {},
	ARG_NO_COLOR: {Type: arg.BOOL},
	ARG_HELP:     {Type: arg.BOOL, Alias: "u:usage"},
	ARG_VER:      {Type: arg.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	args, errs := arg.Parse(argMap)

	if len(errs) != 0 {
		fmtc.Println("Arguments parsing errors:")

		for _, err := range errs {
			fmtc.Printf("  %s\n", err.Error())
		}

		os.Exit(1)
	}

	if arg.GetB(ARG_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if arg.GetB(ARG_VER) {
		showAbout()
		return
	}

	if arg.GetB(ARG_HELP) || len(args) == 0 {
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

// process start source processing
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

	if arg.GetS(ARG_NAME) != "" {
		doc.Title = arg.GetS(ARG_NAME)
	}

	if arg.GetS(ARG_OUTPUT) == "" {
		if pattern == "" {
			simpleRender(doc)
		} else {
			findInfo(doc, pattern)
		}
	} else {
		renderTemplate(doc)
	}
}

// findInfo search geven pattern in entities names
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

// simpleRender print all document info to console
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

// renderTemplate read template and render to file
func renderTemplate(doc *Document) {
	projectDir := env.Get().GetS("GOPATH")
	templateFile := path.Join(
		projectDir, "src/github.com/essentialkaos/shdoc/templates",
		arg.GetS(ARG_TEMPLATE)+".tpl",
	)

	if !fsutil.CheckPerms("FRS", templateFile) {
		printErrorAndExit("Can't read template %s - file does not exist or empty", templateFile)
	}

	outputFile := arg.GetS(ARG_OUTPUT)

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

// renderConstant print constant info to console
func renderConstant(c *Variable) {
	fmtc.Printf("{s-}%4d:{!} {m*}%s{!} {s}={!} %s "+getVarTypeDesc(c.Type)+"\n", c.Line, c.Name, c.Value)
	fmtc.Printf("      %s\n", c.UnitedDesc())
}

// renderMethod print variable info to console
func renderVariable(v *Variable) {
	fmtc.Printf("{s-}%4d:{!} {c*}%s{!} {s}={!} %s "+getVarTypeDesc(v.Type)+"\n", v.Line, v.Name, v.Value)
	fmtc.Printf("      %s\n", v.UnitedDesc())
}

// renderMethod print method info to console
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

// getVarTypeDesc return type description
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

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo("", "file")

	info.AddOption(ARG_OUTPUT, "Path to output file", "file")
	info.AddOption(ARG_TEMPLATE, "Name of template", "name")
	info.AddOption(ARG_NAME, "Overwrite default name", "name")
	info.AddOption(ARG_NO_COLOR, "Disable colors in output")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

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

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "Essential Kaos",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/shdoc", update.GitHubChecker},
	}

	about.Render()
}
