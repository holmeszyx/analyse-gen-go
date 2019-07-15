package als_gen

import (
	"fmt"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := &Parser{}
	als, err := parser.Parse("testdata/als.toml")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Configs: ", als.Config)
	fmt.Println("Params: ", als.Alias)
	fmt.Println("Events:", len(als.Events))
	for _, v := range als.Events {
		fmt.Println(v)
	}
}
