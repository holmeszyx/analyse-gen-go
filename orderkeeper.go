package als_gen

import "unicode"

// keep key or table order

//[enter_]
//c = "123"

// line start with '['. end with ']', escape

const tokenTabStart byte = '['
const tokenTabEnd byte = ']'
const tokenAssign byte = '='
const tokenComment byte = '#'
const tokenLineEnd byte = '\n'

type tokenParser interface {
	Parent() tokenParser
	SetParent(parser tokenParser)

	Start(token byte) bool
	End(token byte) bool
	Eat(c byte) bool
	Content() []byte
}

type element struct {
	Name     string
	Position int32
}

func newElement(name string, position int32) *element {
	return &element{Name: name, Position: position}
}

// record the tables and keys written order
type orderKeeper struct {
	tableNames []element
	keyNames   []element

	lineStart bool
	lineNum   int32
}

// parse the order
func (o *orderKeeper) parse(content string) {
	o.tableNames = make([]element, 0, 16)
	o.keyNames = make([]element, 0, 16)

	raw := []byte(content)

	table := newTableParser(o)
	keys := newKeysParser(o)

	parsers := []tokenParser{
		table,
		keys,
	}

	var tp tokenParser = nil
	o.lineStart = true
	o.lineNum = 1

	for index := 0; index < len(raw); index++ {
		b := raw[index]

		if tp == nil && unicode.IsSpace(rune(b)) {
			// continue
		} else {
			if tp != nil {
				if !tp.End(b) {
					tp.Eat(b)
				} else {
					// o.tableNames = append(o.tableNames, string(tp.Content()))
					tp = nil
				}
			} else {
				for _, parser := range parsers {
					if parser.Start(b) {
						tp = parser
						break
					}
				}
			}
			if o.lineStart {
				o.lineStart = false
			}
		}

		if b == tokenLineEnd {
			o.lineStart = true
			o.lineNum++
		}

	}
}

// to parser table name
type tableParser struct {
	keeper  *orderKeeper
	parent  tokenParser
	content []byte
}

func newTableParser(keeper *orderKeeper) *tableParser {
	return &tableParser{keeper: keeper}
}

func (t *tableParser) Parent() tokenParser {
	return t.parent
}

func (t *tableParser) SetParent(parser tokenParser) {
	t.parent = parser
}

func (t *tableParser) Content() []byte {
	return t.content
}

func (t *tableParser) Start(token byte) bool {
	ok := t.keeper.lineStart && token == tokenTabStart
	if ok {
		t.content = make([]byte, 0, 16)
	}
	return ok
}

func (t *tableParser) End(token byte) bool {
	ok := token == tokenTabEnd || token == tokenLineEnd
	if ok {
		t.keeper.tableNames = append(t.keeper.tableNames, element{
			Name:     string(t.Content()),
			Position: t.keeper.lineNum,
		})
	}
	return ok
}

func (t *tableParser) Eat(c byte) bool {
	if t.content != nil {
		t.content = append(t.content, c)
		return true
	}
	return false
}

// parser the key name
type keysParser struct {
	keeper   *orderKeeper
	parent   tokenParser
	content  []byte
	Assigned bool
}

func newKeysParser(keeper *orderKeeper) *keysParser {
	return &keysParser{keeper: keeper}
}

func (k *keysParser) Parent() tokenParser {
	return k.parent
}

func (k *keysParser) SetParent(parser tokenParser) {
	k.parent = parser
}

func (k *keysParser) Start(token byte) bool {
	ok := k.keeper.lineStart && token != tokenComment && token != tokenTabStart && token != tokenTabEnd
	if ok {
		k.content = make([]byte, 0, 16)
		k.Assigned = false

		// eat first letter
		k.Eat(token)
	}
	return ok
}

func (k *keysParser) End(token byte) bool {
	ok := token == tokenLineEnd
	if ok {
		k.keeper.keyNames = append(k.keeper.keyNames, element{
			Name:     string(k.Content()),
			Position: k.keeper.lineNum,
		})
	}
	return ok
}

func (k *keysParser) Eat(c byte) bool {
	if c == tokenAssign {
		// end content
		k.Assigned = true
		return true
	}
	if k.Assigned {
		return false
	}

	if k.content != nil && !unicode.IsSpace(rune(c)) {
		k.content = append(k.content, c)
		return true
	}
	return false
}

func (k *keysParser) Content() []byte {
	return k.content
}
