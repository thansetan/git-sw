package gitconfig

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Section represents section name (and maybe subsection) of a config.
type Section struct {
	Name, Subsection string
}

// NewSection converts the given string to a Section.
// Only alphanumeric characters and '-' are allowed for the section name.
// To add a subsection, include a '.' after the section name.
// Subsection names can contain any character except newline and null bytes.
// Example: "url.git@github.com" results in Name = "git", Subsection = "git@github.com".
func NewSection(s string) (Section, error) {
	var sec Section
	ix := strings.Index(s, ".")
	if ix > -1 {
		sec.Name = s[:ix]
		sec.Subsection = s[ix+1:]
	} else {
		sec.Name = s
	}
	if !sec.isValidName() {
		return Section{}, ErrInvalidSection
	}
	if !sec.isValidSubsection() {
		return Section{}, ErrInvalidSubsection
	}

	return sec, nil
}

func (s Section) String() string {
	if len(s.Subsection) == 0 {
		return fmt.Sprintf("[%s]", s.Name)
	}
	return fmt.Sprintf("[%s \"%s\"]", s.Name, s.Subsection)
}

// DottedString joins section name and subsection with a dot (.).
func (s Section) DottedString() string {
	if len(s.Subsection) == 0 {
		return s.Name
	}
	return s.Name + "." + s.Subsection
}

func (s Section) isValidName() bool {
	if len(s.Name) == 0 {
		return false
	}
	return !strings.ContainsFunc(s.Name, func(r rune) bool {
		return !(isAlnum(r) || r == '-' || r == '.')
	})
}

func (s Section) isValidSubsection() bool {
	return !strings.ContainsFunc(s.Subsection, func(r rune) bool {
		return r == '\n' || r == 0
	})
}

// VariableName represents config variable name.
type VariableName string

func (vn VariableName) isValid() bool {
	return len(vn) > 0 && isAlpha(vn[0]) && !strings.ContainsFunc(string(vn), func(r rune) bool {
		return !(isAlnum(r) || r == '-')
	})
}

// Key represents key of a variable. It consists of
// section, subsection (if any), and variable name.
type Key struct {
	Section Section
	Name    VariableName
}

func (k Key) String() string {
	return fmt.Sprintf("%s.%s", k.Section.DottedString(), k.Name)
}

// Value represents value of a config variable.
type Value struct{ v interface{} }

func (val Value) Value() interface{} {
	return val.v
}

// ValidateValue validates whether s is a valid value for .gitconfig or not.
// .gitconfig values can contain any characters and may span multiple lines.
// '\' indicates that the config value continues on the next line. There MUST NOT be
// any characters after '\'. If '\' is needed as part of the value, it MUST be
// escaped to '\\'.
func ValidateValue(s string) error {
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if i == len(s)-1 {
				break
			}
			switch s[i+1] {
			case '\\', '"':
				i++
				continue
			case '\n', '\t', '\b':
				continue
			default:
				return ErrInvalidVariableValue
			}
		}
	}

	return nil
}

func (val Value) String() string {
	var quoted bool

	s, ok := val.v.(string)
	if !ok {
		return fmt.Sprintf("%v", val.v)
	}

	b := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' {
			quoted = !quoted
			continue
		}
		if ch == '\\' {
			if i < len(s)-1 {
				switch s[i+1] {
				case '"':
					ch = '"'
					i++
				case '\\':
					i++
				case 'n':
					ch = '\n'
					i++
				case 'b':
					ch = '\b'
					i++
				case 't':
					ch = '\t'
					i++
				}
			}
		}
		b = append(b, ch)
	}

	return string(b)
}

// TODO: maintain insertion order ?? how ??
//
//	add a slice to store the keys ??
type GitConfig struct {
	data *orderedMap[Section, *orderedMap[VariableName, []Value]]
}

func (GitConfig) isValidValues(vals ...interface{}) ([]Value, error) {
	values := make([]Value, 0, len(vals))
	for i := range vals {
		t := reflect.ValueOf(vals[i])
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
		case reflect.String:
			err := ValidateValue(vals[i].(string))
			if err != nil {
				return nil, err
			}
		default:
			return nil, ErrInvalidValueType
		}
		values = append(values, Value{vals[i]})
		// type assertion fails to catch local type
		// switch val.(type) {
		// case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string, bool:
		// default:
		// 	return false
		// }
	}
	return values, nil
}

func (g GitConfig) sectionExists(section Section) bool {
	_, ok := g.data.get(section)
	return ok
}

func (g GitConfig) keyExists(section Section, key VariableName) bool {
	_, ok := g.data.mustGet(section).get(key)
	return ok
}

func (GitConfig) splitKey(k string) (Section, VariableName, error) {
	ix := strings.LastIndexByte(k, '.')
	if ix == -1 {
		return Section{}, "", ErrInvalidKey
	}

	varName := VariableName(k[ix+1:])
	if !varName.isValid() {
		return Section{}, "", ErrInvalidVariableName
	}

	section, err := NewSection(k[:ix])
	if err != nil {
		return Section{}, "", err
	}

	return section, varName, nil
}

func (g GitConfig) get(section Section, key VariableName) ([]Value, error) {
	if !g.sectionExists(section) {
		return nil, ErrKeyNotFound
	}

	if !g.keyExists(section, key) {
		return nil, ErrKeyNotFound
	}

	return g.data.mustGet(section).mustGet(key), nil
}

func (g *GitConfig) add(section Section, name VariableName, vals ...Value) {
	if !g.sectionExists(section) {
		g.data.put(section, newOrderedMap[VariableName, []Value]())
	}

	if !g.keyExists(section, name) {
		g.data.mustGet(section).put(name, make([]Value, 0))
	}

	g.data.mustGet(section).mustGetNode(name).val.val = append(g.data.mustGet(section).mustGet(name), vals...)
}

func (g *GitConfig) set(section Section, name VariableName, vals ...Value) {
	if !g.sectionExists(section) {
		g.data.put(section, newOrderedMap[VariableName, []Value]())
	}

	if !g.keyExists(section, name) {
		g.data.mustGet(section).put(name, make([]Value, 0))
	}

	g.data.mustGet(section).put(name, vals)
}

func (g *GitConfig) unset(section Section, name VariableName) error {
	if !g.sectionExists(section) {
		return ErrKeyNotFound
	}

	if !g.keyExists(section, name) {
		return ErrKeyNotFound
	}

	if g.data.mustGet(section).len() == 1 {
		g.data.remove(section)
	} else {
		g.data.mustGet(section).remove(name)
	}

	return nil
}

// New creates a new GitConfig.
func New() *GitConfig {
	return &GitConfig{
		data: newOrderedMap[Section, *orderedMap[VariableName, []Value]](),
	}
}

// Get retrieves value of a given key, if the key contains multiple values,
// the last value is returned.
func (g GitConfig) Get(key string) (Value, error) {
	section, varName, err := g.splitKey(key)
	if err != nil {
		return Value{}, err
	}

	data, err := g.get(section, varName)
	if err != nil {
		return Value{}, err
	}

	if len(data) == 0 {
		return Value{}, nil
	}

	return data[len(data)-1], nil
}

// GetAll retrieves all values of a given key.
func (g GitConfig) GetAll(key string) ([]Value, error) {
	section, varName, err := g.splitKey(key)
	if err != nil {
		return nil, err
	}

	data, err := g.get(section, varName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Set assigns vals to a given key. If the key already exists, the current value is
// replaced. To add new values to an existing key, use Add().
func (g *GitConfig) Set(key string, vals ...interface{}) error {
	if len(vals) == 0 {
		return ErrEmptyValue
	}

	values, err := g.isValidValues(vals...)
	if err != nil {
		return err
	}

	section, varName, err := g.splitKey(key)
	if err != nil {
		return err
	}

	g.set(section, varName, values...)

	return nil
}

// Add appends (or creates if it doesn't exists yet) the vals to a given key.
func (g *GitConfig) Add(key string, vals ...interface{}) error {
	if len(vals) == 0 {
		return ErrEmptyValue
	}

	values, err := g.isValidValues(vals...)
	if err != nil {
		return err
	}

	section, varName, err := g.splitKey(key)
	if err != nil {
		return err
	}

	g.add(section, varName, values...)

	return nil
}

// Unset removes given key from config file. If the section contains only one variable,
// the section is removed, if it contains more than one variables,
// only variable with the given key is removed.
func (g *GitConfig) Unset(key string) error {
	section, varName, err := g.splitKey(key)
	if err != nil {
		return err
	}

	err = g.unset(section, varName)
	if err != nil {
		return err
	}

	return nil
}

// Save writes the current configuration to path.
// If the file already exists, it will be overwritten.
func (g GitConfig) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sections := g.data.keys()
	for i := range sections {
		_, err = f.WriteString(fmt.Sprintf("%s\n", sections[i]))
		if err != nil {
			return err
		}
		variables := g.data.mustGet(sections[i]).keys()
		for j := range variables {
			values := g.data.mustGet(sections[i]).mustGet(variables[j])
			for k := range values {
				_, err = f.WriteString(fmt.Sprintf("\t%s = %s\n", variables[j], values[k].Value()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Keys returns slice of all keys in the order they're
// inserted.
func (g GitConfig) Keys() []Key {
	keys := make([]Key, 0)

	for _, section := range g.data.keys() {
		variables := g.data.mustGet(section)
		if cap(keys)-len(keys) < variables.len() {
			keys = append(keys[:cap(keys)], make([]Key, variables.len())...)[:len(keys)]
		}

		for _, name := range variables.keys() {
			keys = append(keys, Key{section, name})
		}
	}

	return keys
}
