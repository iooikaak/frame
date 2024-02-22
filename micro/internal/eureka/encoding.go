package eureka

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/micro/go-micro/v2/registry"
)

func encode(buf []byte) string {
	var b bytes.Buffer
	defer b.Reset()

	w := zlib.NewWriter(&b)
	if _, err := w.Write(buf); err != nil {
		return ""
	}
	_ = w.Close()

	return hex.EncodeToString(b.Bytes())
}

func decode(d string) []byte {
	hr, err := hex.DecodeString(d)
	if err != nil {
		return nil
	}

	br := bytes.NewReader(hr)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return nil
	}

	rBuf, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil
	}
	_ = zr.Close()
	return rBuf
}

func encodeEndpoints(en []*registry.Endpoint) string {
	endpoint, _ := json.Marshal(en)
	return encode(endpoint)
}

func decodeEndpoints(d string) []*registry.Endpoint {
	en := make([]*registry.Endpoint, 0)
	_ = json.Unmarshal(decode(d), &en)

	return en
}

func encodeNodes(en []*registry.Node) string {
	endpoint, _ := json.Marshal(en)
	return encode(endpoint)
}

func decodeNodes(d string) []*registry.Node {
	en := make([]*registry.Node, 0)
	_ = json.Unmarshal(decode(d), &en)

	return en
}

func encodeValue(value *Value) string {
	instances := value.Val
	b, _ := json.Marshal(instances)
	return encode(b)
}
