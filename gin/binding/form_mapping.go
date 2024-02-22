package binding

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iooikaak/frame/gin/internal/bytesconv"
	"github.com/iooikaak/frame/gin/internal/json"
)

var errUnknownType = errors.New("unknown type")

// scache struct reflect type cache.
var scache = &cache{
	data: make(map[reflect.Type]*sinfo),
}

type cache struct {
	data  map[reflect.Type]*sinfo
	mutex sync.RWMutex
}

func (c *cache) get(obj reflect.Type) (s *sinfo) {
	var ok bool
	c.mutex.RLock()
	if s, ok = c.data[obj]; !ok {
		c.mutex.RUnlock()
		s = c.set(obj)
		return
	}
	c.mutex.RUnlock()
	return
}

func (c *cache) set(obj reflect.Type) (s *sinfo) {
	s = new(sinfo)
	tp := obj.Elem()
	for i := 0; i < tp.NumField(); i++ {
		fd := new(field)
		fd.tp = tp.Field(i)
		tag := fd.tp.Tag.Get("form")
		fd.name, fd.option = parseTag(tag)
		if defV := fd.tp.Tag.Get("default"); defV != "" {
			dv := reflect.New(fd.tp.Type).Elem()
			setWithProperType1(fd.tp.Type.Kind(), []string{defV}, dv, fd.option)
			fd.hasDefault = true
			fd.defaultValue = dv
		}
		s.field = append(s.field, fd)
	}
	c.mutex.Lock()
	c.data[obj] = s
	c.mutex.Unlock()
	return
}

type sinfo struct {
	field []*field
}

type field struct {
	tp     reflect.StructField
	name   string
	option tagOptions

	hasDefault   bool          // if field had default value
	defaultValue reflect.Value // field default value
}

func mapUri(ptr interface{}, m map[string][]string) error {
	return mapFormByTag(ptr, m, "uri")
}

func mapForm(ptr interface{}, form map[string][]string) error {
	sinfo := scache.get(reflect.TypeOf(ptr))
	val := reflect.ValueOf(ptr).Elem()
	for i, fd := range sinfo.field {
		typeField := fd.tp
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := fd.name
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if structFieldKind == reflect.Struct {
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil {
					return err
				}
				continue
			}
		}
		inputValue, exists := form[inputFieldName]
		if !exists {
			// Set the field as default value when the input value is not exist
			if fd.hasDefault {
				structField.Set(fd.defaultValue)
			}
			continue
		}
		// Set the field as default value when the input value is empty
		if fd.hasDefault && inputValue[0] == "" {
			structField.Set(fd.defaultValue)
			continue
		}
		if _, isTime := structField.Interface().(time.Time); isTime {
			if err := setTimeField(inputValue[0], typeField, structField); err != nil {
				return err
			}
			continue
		}
		if err := setWithProperType1(typeField.Type.Kind(), inputValue, structField, fd.option); err != nil {
			return err
		}
	}
	return nil
}

var emptyField = reflect.StructField{}

func mapFormByTag(ptr interface{}, form map[string][]string, tag string) error {
	return mappingByPtr(ptr, formSource(form), tag)
}

// setter tries to set value on a walking by fields of a struct
type setter interface {
	TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSetted bool, err error)
}

type formSource map[string][]string

var _ setter = formSource(nil)

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form formSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByForm(value, field, form, tagValue, opt)
}

func mappingByPtr(ptr interface{}, setter setter, tag string) error {
	_, err := mapping(reflect.ValueOf(ptr), emptyField, setter, tag)
	return err
}

func mapping(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	if field.Tag.Get(tag) == "-" { // just ignoring this field
		return false, nil
	}

	var vKind = value.Kind()

	if vKind == reflect.Ptr {
		var isNew bool
		vPtr := value
		if value.IsNil() {
			isNew = true
			vPtr = reflect.New(value.Type().Elem())
		}
		isSetted, err := mapping(vPtr.Elem(), field, setter, tag)
		if err != nil {
			return false, err
		}
		if isNew && isSetted {
			value.Set(vPtr)
		}
		return isSetted, nil
	}

	if vKind != reflect.Struct || !field.Anonymous {
		ok, err := tryToSetValue(value, field, setter, tag)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if vKind == reflect.Struct {
		tValue := value.Type()

		var isSetted bool
		for i := 0; i < value.NumField(); i++ {
			sf := tValue.Field(i)
			if sf.PkgPath != "" && !sf.Anonymous { // unexported
				continue
			}
			ok, err := mapping(value.Field(i), tValue.Field(i), setter, tag)
			if err != nil {
				return false, err
			}
			isSetted = isSetted || ok
		}
		return isSetted, nil
	}
	return false, nil
}

type setOptions struct {
	isDefaultExists bool
	defaultValue    string
}

func tryToSetValue(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	var tagValue string
	var setOpt setOptions

	tagValue = field.Tag.Get(tag)
	tagValue, opts := head(tagValue, ",")

	if tagValue == "" { // default value is FieldName
		tagValue = field.Name
	}
	if tagValue == "" { // when field is "emptyField" variable
		return false, nil
	}

	var opt string
	for len(opts) > 0 {
		opt, opts = head(opts, ",")

		if k, v := head(opt, "="); k == "default" {
			setOpt.isDefaultExists = true
			setOpt.defaultValue = v
		}
	}

	return setter.TrySet(value, field, tagValue, setOpt)
}

func setByForm(value reflect.Value, field reflect.StructField, form map[string][]string, tagValue string, opt setOptions) (isSetted bool, err error) {
	vs, ok := form[tagValue]
	if !ok && !opt.isDefaultExists {
		return false, nil
	}

	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		return true, setSlice(vs, value, field)
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) != value.Len() {
			return false, fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		}
		return true, setArray(vs, value, field)
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}

		if len(vs) > 0 {
			val = vs[0]
		}
		return true, setWithProperType(val, value, field)
	}
}

func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value, field)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		}
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	default:
		return errUnknownType
	}
	return nil
}

func setWithProperType1(valueKind reflect.Kind, val []string, structField reflect.Value, option tagOptions) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val[0], 0, structField)
	case reflect.Int8:
		return setIntField(val[0], 8, structField)
	case reflect.Int16:
		return setIntField(val[0], 16, structField)
	case reflect.Int32:
		return setIntField(val[0], 32, structField)
	case reflect.Int64:
		return setIntField(val[0], 64, structField)
	case reflect.Uint:
		return setUintField(val[0], 0, structField)
	case reflect.Uint8:
		return setUintField(val[0], 8, structField)
	case reflect.Uint16:
		return setUintField(val[0], 16, structField)
	case reflect.Uint32:
		return setUintField(val[0], 32, structField)
	case reflect.Uint64:
		return setUintField(val[0], 64, structField)
	case reflect.Bool:
		return setBoolField(val[0], structField)
	case reflect.Float32:
		return setFloatField(val[0], 32, structField)
	case reflect.Float64:
		return setFloatField(val[0], 64, structField)
	case reflect.String:
		structField.SetString(val[0])
	case reflect.Slice:
		if option.Contains("split") {
			val = strings.Split(val[0], ",")
		}
		filtered := filterEmpty(val)
		switch structField.Type().Elem().Kind() {
		case reflect.Int64:
			valSli := make([]int64, 0, len(filtered))
			for i := 0; i < len(filtered); i++ {
				d, err := strconv.ParseInt(filtered[i], 10, 64)
				if err != nil {
					return err
				}
				valSli = append(valSli, d)
			}
			structField.Set(reflect.ValueOf(valSli))
		case reflect.String:
			valSli := make([]string, 0, len(filtered))
			for i := 0; i < len(filtered); i++ {
				valSli = append(valSli, filtered[i])
			}
			structField.Set(reflect.ValueOf(valSli))
		default:
			sliceOf := structField.Type().Elem().Kind()
			numElems := len(filtered)
			slice := reflect.MakeSlice(structField.Type(), len(filtered), len(filtered))
			for i := 0; i < numElems; i++ {
				if err := setWithProperType1(sliceOf, filtered[i:], slice.Index(i), ""); err != nil {
					return err
				}
			}
			structField.Set(slice)
		}
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil

	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, value reflect.Value, field reflect.StructField) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func filterEmpty(val []string) []string {
	filtered := make([]string, 0, len(val))
	for _, v := range val {
		if v != "" {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
