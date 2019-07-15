package als_gen

import (
	"bufio"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type AndroidKtExporter struct {
	OutDir   string
	Tmpl     *template.Template
	FileName string
}

func NewAndroidKtExporter(outDir string, tmpl *template.Template, fileName string) *AndroidKtExporter {
	return &AndroidKtExporter{OutDir: outDir, Tmpl: tmpl, FileName: fileName}
}


func (a *AndroidKtExporter) prepare() {
	_, err := os.Stat(a.OutDir)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(a.OutDir, os.FileMode(0755))
		return
	}
}

func (a *AndroidKtExporter) Export(als *Als) error {
	a.prepare()

	outPath := filepath.Join(a.OutDir, a.FileName)
	entity := &ExportEntity{
		Data:     als,
		FileName: a.FileName,
		FilePath: outPath,
		Date:     time.Now().Format("15:04:05 2006-01-02"),
	}

	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		return err
	}
	defer outFile.Close()
	outBuff := bufio.NewWriter(outFile)

	err = a.Tmpl.Execute(outBuff, entity)
	_ = outBuff.Flush()
	return err

}

