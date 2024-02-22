package test

import "net/url"

type TestParam struct {
	val url.Values
}

func NewTestParam() *TestParam {
	return &TestParam{
		val: make(map[string][]string),
	}
}

func (p *TestParam) Add(key, value string) {
	p.val.Add(key, value)
}

func (p *TestParam) Encode() string {
	return p.val.Encode()
}

func (p *TestParam) Values() url.Values {
	return p.val
}
