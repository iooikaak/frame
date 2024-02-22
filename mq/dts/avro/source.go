// Code generated by github.com/actgardner/gogen-avro/v7. DO NOT EDIT.
/*
 * SOURCE:
 *     Record.avsc
 */
package avro

import (
	"io"

	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
)

type Source struct {
	SourceType SourceType `json:"sourceType"`
	// source datasource version information
	Version string `json:"version"`
}

const SourceAvroCRC64Fingerprint = "Z\xe3G\xb0O\u007f\xb4z"

func NewSource() *Source {
	return &Source{}
}

func DeserializeSource(r io.Reader) (*Source, error) {
	t := NewSource()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func DeserializeSourceFromSchema(r io.Reader, schema string) (*Source, error) {
	t := NewSource()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func writeSource(r *Source, w io.Writer) error {
	var err error
	err = writeSourceType(r.SourceType, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Version, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Source) Serialize(w io.Writer) error {
	return writeSource(r, w)
}

func (r *Source) Schema() string {
	return "{\"fields\":[{\"name\":\"sourceType\",\"type\":{\"name\":\"SourceType\",\"namespace\":\"com.alibaba.dts.formats.avro\",\"symbols\":[\"MySQL\",\"Oracle\",\"SQLServer\",\"PostgreSQL\",\"MongoDB\",\"Redis\",\"DB2\",\"PPAS\",\"DRDS\",\"HBASE\",\"HDFS\",\"FILE\",\"OTHER\"],\"type\":\"enum\"}},{\"doc\":\"source datasource version information\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"com.alibaba.dts.formats.avro.Source\",\"type\":\"record\"}"
}

func (r *Source) SchemaName() string {
	return "com.alibaba.dts.formats.avro.Source"
}

func (_ *Source) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Source) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Source) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Source) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Source) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Source) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Source) SetString(v string)   { panic("Unsupported operation") }
func (_ *Source) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Source) Get(i int) types.Field {
	switch i {
	case 0:
		return &SourceTypeWrapper{Target: &r.SourceType}
	case 1:
		return &types.String{Target: &r.Version}
	}
	panic("Unknown field index")
}

func (r *Source) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Source) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ *Source) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Source) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Source) Finalize()                        {}

func (_ *Source) AvroCRC64Fingerprint() []byte {
	return []byte(SourceAvroCRC64Fingerprint)
}
