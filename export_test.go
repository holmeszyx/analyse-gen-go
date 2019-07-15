package als_gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func TestAndroidKt_Export(t *testing.T) {
	parser := &Parser{}
	als, err := parser.Parse("testdata/als.toml")
	if err != nil {
		t.Error(err)
	}
	pkg := als.Config.Package
	dirs := strings.Split(pkg, ".")
	dir := filepath.Join(dirs...)

	tempFile, err := os.Open(filepath.Join("testdata", "umeng-android-kt.tmpl"))
	if err != nil {
		t.Error(err)
	}
	text, _ := ioutil.ReadAll(tempFile)
	// fmt.Println(string(text))
	_ = tempFile.Close()

	temp := template.New("umeng").Funcs(GetTmplFuncMap())
	temp, err = temp.Parse(string(text))
	if err != nil {
		t.Error(err)
	}


	export := NewAndroidKtExporter(filepath.Join("testdata", dir), temp, "St.kt")
	err = export.Export(als)
	if err != nil {
		t.Error(err)
	}
}
