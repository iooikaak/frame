package paladin

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	Map
}

func (y *YAML) Set(text string) error {
	if err := y.UnmarshalText(text); err != nil {
		return err
	}
	return nil
}

func (y *YAML) UnmarshalText(text string) error {
	raws := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(text), &raws); err != nil {
		return err
	}
	values := map[string]*Value{}
	for k, v := range raws {
		k = KeyNamed(k)
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Map:
			buf := bytes.NewBuffer(nil)
			err := yaml.NewEncoder(buf).Encode(v)
			// b, err := toml.Marshal(v)
			if err != nil {
				return err
			}
			// NOTE: value is map[string]interface{}
			values[k] = NewValue(v, buf.String())
		case reflect.Slice:
			raw := map[string]interface{}{
				k: v,
			}
			buf := bytes.NewBuffer(nil)
			err := yaml.NewEncoder(buf).Encode(raw)
			// b, err := toml.Marshal(raw)
			if err != nil {
				return err
			}
			// NOTE: value is []interface{}
			values[k] = NewValue(v, buf.String())
		case reflect.Bool:
			b := v.(bool)
			values[k] = NewValue(b, strconv.FormatBool(b))
		case reflect.Int64:
			i := v.(int64)
			values[k] = NewValue(i, strconv.FormatInt(i, 10))
		case reflect.Int:
			i := v.(int)
			values[k] = NewValue(i, strconv.Itoa(i))
		case reflect.Float64:
			f := v.(float64)
			values[k] = NewValue(f, strconv.FormatFloat(f, 'f', -1, 64))
		case reflect.String:
			s := v.(string)
			values[k] = NewValue(s, s)
		default:
			return errors.Errorf("UnmarshalYAML: unknown kind(%v)", rv.Kind())
		}
	}
	y.Store(values)
	return nil
}
