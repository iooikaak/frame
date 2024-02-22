package pipeline

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/iooikaak/frame/net/metadata"
	xtime "github.com/iooikaak/frame/time"
	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	conf := &Config{
		MaxSize:  3,
		Interval: xtime.Duration(time.Millisecond * 20),
		Buffer:   3,
		Worker:   10,
	}
	type recv struct {
		mirror string
		values map[string][]interface{}
	}
	var runs []recv
	do := func(c context.Context, values map[string][]interface{}) {
		runs = append(runs, recv{
			mirror: metadata.String(c, metadata.Mirror),
			values: values,
		})
	}
	split := func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}
	p := NewPipeline(conf)
	p.Do = do
	p.Split = split
	p.Start()
	err := p.Add(context.Background(), "1", 1)
	assert.Nil(t, err)
	err = p.Add(context.Background(), "1", 2)
	assert.Nil(t, err)
	err = p.Add(context.Background(), "11", 3)
	assert.Nil(t, err)
	err = p.Add(context.Background(), "2", 3)
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 60)
	mirrorCtx := metadata.NewContext(context.Background(), metadata.MD{metadata.Mirror: "1"})
	err = p.Add(mirrorCtx, "2", 3)
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 60)
	p.SyncAdd(mirrorCtx, "5", 5)
	time.Sleep(time.Millisecond * 60)
	err = p.Close()
	assert.Nil(t, err)
	expt := []recv{
		{
			mirror: "",
			values: map[string][]interface{}{
				"1":  {1, 2},
				"11": {3},
			},
		},
		{
			mirror: "",
			values: map[string][]interface{}{
				"2": {3},
			},
		},
		{
			mirror: "1",
			values: map[string][]interface{}{
				"2": {3},
			},
		},
		{
			mirror: "1",
			values: map[string][]interface{}{
				"5": {5},
			},
		},
	}
	if !reflect.DeepEqual(runs, expt) {
		t.Errorf("expect get %+v,\n got: %+v", expt, runs)
	}
}

func TestPipelineSmooth(t *testing.T) {
	conf := &Config{
		MaxSize:  100,
		Interval: xtime.Duration(time.Second),
		Buffer:   100,
		Worker:   10,
	}
	type result struct {
		ts time.Time
	}
	var results []result
	do := func(c context.Context, values map[string][]interface{}) {
		results = append(results, result{
			ts: time.Now(),
		})
	}
	split := func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}
	p := NewPipeline(conf)
	p.Do = do
	p.Split = split
	p.Start()
	for i := 0; i < 10; i++ {
		err := p.Add(context.Background(), strconv.Itoa(i), 1)
		assert.Nil(t, err)
	}
	time.Sleep(time.Millisecond * 1500)
	if len(results) != conf.Worker {
		t.Errorf("expect results equal worker")
		t.FailNow()
	}
	for i, r := range results {
		if i > 0 {
			if r.ts.Sub(results[i-1].ts) < time.Millisecond*20 {
				t.Errorf("expect runs be smooth")
				t.FailNow()
			}
		}
	}
}
