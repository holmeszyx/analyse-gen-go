package als_gen

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

const (
	KEY_CONFIG        = "Config"
	KEY_Params        = "Params"
	KEY_EVENT_COMMENT = "c"
)

type Field struct {
	Name    string
	Comment string
}

type Event struct {
	Field
	Params []Field
}

type Config struct {
	Package string
	Dir     string
}

type Params map[string]string

type Als struct {
	Config Config
	Alias  Params
	Events []Event
}

type Parser struct {
}

func (p *Parser) Parse(alsFile string) (*Als, error) {
	alsContent, err := func(file string) (string, error) {
		f, err := os.Open(file)
		if err != nil {
			return "", err
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}(alsFile)

	if err != nil {
		return nil, err
	}

	alsConf, err := toml.Load(alsContent)
	if err != nil {
		return nil, err
	}

	// keep the tables and keys write order
	order := &orderKeeper{}
	order.parse(alsContent)

	als := Als{
		Events: make([]Event, 0, 16),
	}

	// parse config
	confTree := alsConf.Get(KEY_CONFIG).(*toml.Tree)
	if confTree != nil {
		conf := als.Config
		//todo nil handler
		conf.Package = confTree.Get("package").(string)
		conf.Dir = confTree.Get("dir").(string)
		als.Config = conf
	}
	// parse params alias
	paramsTree := alsConf.Get(KEY_Params).(*toml.Tree)
	aliasReverseMap := make(Params)
	if paramsTree != nil {
		e := paramsTree.Unmarshal(&(als.Alias))
		if e != nil {
			fmt.Println(e)
		}
		for k, v := range als.Alias {
			aliasReverseMap[v] = k
		}
	}

	// convert alias to origin param name
	aliasMap := func(alias string) (string, bool) {
		if als.Alias == nil {
			return alias, false
		}
		if a, ok := aliasReverseMap[alias]; ok {
			if a != "" {
				return a, true
			} else {
				return alias, true
			}
		} else if _, ok := als.Alias[alias]; ok {
			return alias, true
		}
		return "", false
	}

	// parse events
	var eventKeys []string = make([]string, 0, 16)
	var paramsKeys []string = make([]string, 0, 8)
	var paramsOrder map[string]int = make(map[string]int)

	// use order.tableNames to instead of range alsConf.Keys()
	for index := 0; index < len(order.tableNames); index++ {
		table := order.tableNames[index]

		key := table.Name
		if key != KEY_CONFIG && key != KEY_Params {
			eventKeys = append(eventKeys, key)
		} else if key == KEY_Params {
			// load params order list
			start := table.Position
			var end int32 = math.MaxInt32
			if index != (len(order.tableNames) - 1) {
				end = order.tableNames[index+1].Position
			}

			for pi := 0; pi < len(order.keyNames); pi++ {
				theKey := order.keyNames[pi]
				if theKey.Position > start && theKey.Position < end {
					paramsOrder[theKey.Name] = len(paramsKeys)
					paramsKeys = append(paramsKeys, theKey.Name)
				}
			}

		}
	}

	for _, key := range eventKeys {
		// parse properties of each event

		evp := alsConf.Get(key).(*toml.Tree)
		c := evp.Get(KEY_EVENT_COMMENT)
		if c == nil || (c.(string) == "") {
			// ignore
			fmt.Printf("The event(%s) not have comment(c). ignore it.. \n", key)
			continue
		}

		event := Event{
			Params: make([]Field, 0, 4),
		}
		event.Name = key

		for _, p := range evp.Keys() {
			if p == KEY_EVENT_COMMENT {
				event.Comment = evp.Get(p).(string)
			} else {
				// params
				if paramName, ok := aliasMap(p); ok {
					event.Params = append(event.Params, Field{
						Name:    paramName,
						Comment: evp.Get(p).(string),
					})
				}
			}
		}
		if len(event.Params) > 0 {
			// sort with the written order
			sort.Slice(event.Params, func(i, j int) bool {
				pi := event.Params[i].Name
				pj := event.Params[j].Name

				oi := math.MaxInt32
				oj := math.MaxInt32

				if o, ok := paramsOrder[pi]; ok {
					oi = o
				}
				if o, ok := paramsOrder[pj]; ok {
					oj = o
				}

				return oi < oj
			})
		}

		als.Events = append(als.Events, event)
	}

	return &als, nil
}
