package paladin_test

import (
	"context"
	"fmt"
	"os"

	"github.com/iooikaak/frame/config/paladin"
	"github.com/iooikaak/frame/config/paladin/apollo"

	"gopkg.in/yaml.v3"

	"github.com/BurntSushi/toml"
)

type exampleConf struct {
	Bool    bool
	Int     int64
	Float   float64
	String  string
	Strings []string
	TestKey string `toml:"test_key"`
}

type exampleConfYaml struct {
	Bool    bool
	Int     int64
	Float   float64
	String  string
	Strings []string
	TestKey string `yaml:"test_key"`
}

func (e *exampleConf) Set(text string) error {
	var ec exampleConf
	if err := toml.Unmarshal([]byte(text), &ec); err != nil {
		return err
	}
	*e = ec
	return nil
}

func (e *exampleConfYaml) Set(text string) error {
	var ec exampleConfYaml
	if err := yaml.Unmarshal([]byte(text), &ec); err != nil {
		return err
	}
	*e = ec
	return nil
}

// ExampleClient is an example client usage.
// exmaple.toml:
/*
	bool = true
	int = 100
	float = 100.1
	string = "text"
	strings = ["a", "b", "c"]
*/
//nolint
func ExampleFileClient() {
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	var ec exampleConf
	// var setter
	if err := paladin.Watch("example.toml", &ec); err != nil {
		panic(err)
	}
	if err := paladin.Get("example.toml").UnmarshalTOML(&ec); err != nil {
		panic(err)
	}
	// use exampleConf
	// watch event key
	go func() {
		for event := range paladin.WatchEvent(context.TODO(), "key") {
			fmt.Println(event)
		}
	}()
}

// ExampleApolloClient is an example client for apollo driver usage.
//nolint
func ExampleApolloClient() {
	/*
		pass flags or set envs that apollo needs, for example:

		```
		export APOLLO_APP_ID=SampleApp
		export APOLLO_CLUSTER=default
		export APOLLO_CACHE_DIR=/tmp
		export APOLLO_META_ADDR=localhost:8080
		export APOLLO_NAMESPACES=example.yml
		```
	*/

	if err := paladin.Init(apollo.PaladinDriverApollo); err != nil {
		panic(err)
	}
	var ec exampleConf
	// var setter
	if err := paladin.Watch("example.yml", &ec); err != nil {
		panic(err)
	}
	if err := paladin.Get("example.yml").UnmarshalYAML(&ec); err != nil {
		panic(err)
	}
	// use exampleConf
	// watch event key
	go func() {
		for event := range paladin.WatchEvent(context.TODO(), "key") {
			fmt.Println(event)
		}
	}()
}

// ExampleMap is an example map usage.
// exmaple.toml:
/*
	bool = true
	int = 100
	float = 100.1
	string = "text"
	strings = ["a", "b", "c"]

	[object]
	string = "text"
	bool = true
	int = 100
	float = 100.1
	strings = ["a", "b", "c"]
*/
//nolint
func ExampleApolloMap() {
	var (
		m    paladin.TOML
		strs []string
	)
	// paladin toml
	if err := paladin.Watch("example.toml", &m); err != nil {
		panic(err)
	}
	// value string
	s, err := m.Get("string").String()
	if err != nil {
		s = "default"
	}
	fmt.Println(s)
	// value bool
	b, err := m.Get("bool").Bool()
	if err != nil {
		b = false
	}
	fmt.Println(b)
	// value int
	i, err := m.Get("int").Int64()
	if err != nil {
		i = 100
	}
	fmt.Println(i)
	// value float
	f, err := m.Get("float").Float64()
	if err != nil {
		f = 100.1
	}
	fmt.Println(f)
	// value slice
	if err = m.Get("strings").Slice(&strs); err == nil {
		fmt.Println(strs)
	}
}

//nolint
func ExampleMapToml() {
	var (
		m      paladin.TOML
		strs   []string
		strMap map[string]interface{}
	)
	{
		_ = os.Setenv("APOLLO_APP_ID", "GO")
		_ = os.Setenv("APOLLO_CLUSTER", "default")
		_ = os.Setenv("APOLLO_CACHE_DIR", "/tmp")
		_ = os.Setenv("APOLLO_META_ADDR", "10.1.2.134:8080")
		_ = os.Setenv("APOLLO_NAMESPACES", "application,data")
	}
	if err := paladin.Init(apollo.PaladinDriverApollo); err != nil {
		panic(err)
	}

	// paladin toml
	if err := paladin.Watch("application", &m); err != nil {
		panic(err)
	}
	// value string
	s, err := m.Get("string").String()
	if err != nil {
		s = "default"
	}
	fmt.Println(s)
	// value bool
	b, err := m.Get("bool").Bool()
	if err != nil {
		b = false
	}
	fmt.Println(b)
	// value int
	i, err := m.Get("int").Int64()
	if err != nil {
		i = 100
	}
	fmt.Println(i)
	// value float
	f, err := m.Get("float").Float64()
	if err != nil {
		f = 100.1
	}
	fmt.Println(f)
	// value slice
	if err = m.Get("strings").Slice(&strs); err == nil {
		fmt.Println(strs)
	}
	if strMap, err = m.Get("strMap").Map(); err == nil {
		fmt.Println(strMap)
	}
}

//nolint
func ExampleMapYaml() {
	var (
		m      paladin.YAML
		strs   []string
		strMap map[string]interface{}
	)
	{
		_ = os.Setenv("APOLLO_APP_ID", "GO")
		_ = os.Setenv("APOLLO_CLUSTER", "default")
		_ = os.Setenv("APOLLO_CACHE_DIR", "/tmp")
		_ = os.Setenv("APOLLO_META_ADDR", "10.1.2.134:8080")
		_ = os.Setenv("APOLLO_NAMESPACES", "application,data")
	}
	if err := paladin.Init(apollo.PaladinDriverApollo); err != nil {
		panic(err)
	}

	// paladin toml
	if err := paladin.Watch("data", &m); err != nil {
		panic(err)
	}
	// value string
	s, err := m.Get("string").String()
	if err != nil {
		s = "default"
	}
	fmt.Println(s)
	// value bool
	b, err := m.Get("bool").Bool()
	if err != nil {
		b = false
	}
	fmt.Println(b)
	// value int
	i, err := m.Get("int").Int64()
	if err != nil {
		i = 100
	}
	fmt.Println(i)
	// value float
	f, err := m.Get("float").Float64()
	if err != nil {
		f = 100.1
	}
	fmt.Println(f)
	// value slice
	if err = m.Get("strings").Slice(&strs); err == nil {
		fmt.Println(strs)
	}
	if strMap, err = m.Get("strMap").Map(); err == nil {
		fmt.Println(strMap)
	}
}
