package als_gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestOrder_Parser(t *testing.T) {
	order := &orderKeeper{}
	file, err := os.Open("testdata/als.toml")
	if err != nil {
		t.Error(err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	order.parse(string(text))

	fmt.Println()
	fmt.Println("Tables:")
	fmt.Println()
	for _, table := range order.tableNames {
		fmt.Println(table)
	}

	fmt.Println()
	fmt.Println("Keys:")
	fmt.Println()

	for _, keys := range order.keyNames {
		fmt.Println(keys)
	}

}
