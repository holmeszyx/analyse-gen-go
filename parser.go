package als_gen

import (
	"fmt"
	"github.com/pelletier/go-toml"
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
	alsConf, err := toml.LoadFile(alsFile)
	if err != nil {
		return nil, err
	}

	als := Als{
		Events: make([]Event, 0, 16),
	}

	// parse config
	confTree := alsConf.Get(KEY_CONFIG).(*toml.Tree)
	if confTree != nil {
		conf := als.Config
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
	for _, key := range alsConf.Keys() {
		if key != KEY_CONFIG && key != KEY_Params {
			eventKeys = append(eventKeys, key)
		}
	}

	for _, key := range eventKeys {
		eve := alsConf.Get(key).(*toml.Tree)
		c := eve.Get(KEY_EVENT_COMMENT)
		if c == nil || (c.(string) == "") {
			// ignore
			fmt.Printf("The event(%s) not have comment(c). ignore it.. \n", key)
			continue
		}

		event := Event{
			Params: make([]Field, 0, 4),
		}
		event.Name = key
		for _, p := range eve.Keys() {
			if p == KEY_EVENT_COMMENT {
				event.Comment = eve.Get(p).(string)
			} else {
				// params
				if paramName, ok := aliasMap(p); ok {
					event.Params = append(event.Params, Field{
						Name:    paramName,
						Comment: eve.Get(p).(string),
					})
				}
			}
		}
		als.Events = append(als.Events, event)
	}

	return &als, nil
}
