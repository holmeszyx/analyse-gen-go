package main

import (
	gen "als-gen"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type cmder struct {
	configFile     string
	templateFile   string
	userPackageDir bool
	outBase        string
	outFilename    string
}

var cmd = cmder{

}

func main() {

	flag.StringVar(&(cmd.configFile), "c", "als.toml", "the analyse define file (*.toml)")
	flag.StringVar(&(cmd.templateFile), "t", "", "The template to generate source code.")
	flag.StringVar(&(cmd.outBase), "out-base-dir", ".", "The generator out put dir.")
	flag.StringVar(&(cmd.outFilename), "o", "St.kt", "The generator out put file name.")
	flag.BoolVar(&(cmd.userPackageDir), "use-pkg", true, "use package to dir which relative to outBaseDir")
	flag.Parse()

	if !isFileExists(cmd.configFile) {
		fmt.Println("analyse define file", cmd.configFile, "not found.")
		os.Exit(1)
	}

	if !isFileExists(cmd.templateFile) {
		fmt.Println("template file", cmd.templateFile, "not found.")
		os.Exit(1)
	}


	parser := &gen.Parser{}
	als, err := parser.Parse(cmd.configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	tempFile, err := os.Open(cmd.templateFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	text, _ := ioutil.ReadAll(tempFile)
	// fmt.Println(string(text))
	_ = tempFile.Close()

	temp := template.New("temps").Funcs(gen.GetTmplFuncMap())
	temp, err = temp.Parse(string(text))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	outDir := cmd.outBase
	if cmd.userPackageDir {
		pkg := als.Config.Package
		dirs := strings.Split(pkg, ".")
		dir := filepath.Join(dirs...)
		outDir = filepath.Join(cmd.outBase, dir)
	}

	export := gen.NewAndroidKtExporter(outDir, temp, cmd.outFilename)
	err = export.Export(als)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

}

func isFileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
