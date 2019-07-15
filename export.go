package als_gen

import (
	"strings"
	"text/template"
	"unicode"
)

type Exporter interface {
	// export
	Export(als *Als) error
}

type ExportEntity struct {
	Data     *Als
	FileName string
	FilePath string
	Date     string
}

func GetTmplFuncMap() template.FuncMap {

	funcMap := template.FuncMap{
		"firstCap": FirstCap,
		"toFuncName": ToFuncName,
	}

	return funcMap

}

func FirstCap(text string) string {
	if len(text) > 0 {
		return strings.ToUpper(text[0:1]) + text[1:]
	}
	return text
}

func ToFuncName(name string) string {
	raw := []byte(name)
	finalData := make([]byte, 0, len(raw))
	shouldCap := false
	for _, c := range raw {
		if !unicode.IsSpace(rune(c)) && c != '-' && c != '_' {
			if shouldCap {
				c = byte(unicode.ToUpper(rune(c)))
				shouldCap = false
			}
			finalData = append(finalData, c)
		} else {
			shouldCap = true
		}
	}
	return string(finalData)
}
