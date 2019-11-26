package als_gen

import (
	"bytes"
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
		"firstCap":   FirstCap,
		"toFuncName": ToFuncName,
		"safeHan":    SafeHan,
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

// SafeHan Only "number", "chinese", "_-.+".
// others replace with '_'
func SafeHan(text string) string {
	if text == "" {
		return text
	}

	runeText := []rune(text)
	finalData := bytes.NewBuffer(make([]byte, 0, len(text)))
	for _, r := range runeText {
		if unicode.IsNumber(r) ||
			unicode.IsLetter(r) ||
			unicode.Is(unicode.Han, r) ||
			r == '-' || r == '_' || r == '.' || r == '+' {
			finalData.WriteRune(r)
		} else {
			finalData.WriteRune('_')
		}
	}
	return finalData.String()
}
