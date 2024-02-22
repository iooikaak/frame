package elasticsearch

import (
	"context"
	"encoding/json"
	"testing"

	//"github.com/iooikaak/frame/elastic"
	v7 "github.com/olivere/elastic/v7"
)

type tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

var testMapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text"
				}
		}
	}
}
`

//func TestSearch(t *testing.T) {
//
//	var (
//		queryBuilder *elastic.BoolQuery
//		err          error
//		es           *elastic.Client
//		//tr           *tracer.Tracer
//		//span         tracer.TraceOperator
//		index = "elastic-test"
//		doc   = "doc"
//		ctx   = context.Background()
//	)
//
//	//获取es
//	es, err = New(&ElasticConfig{
//		//Addrs:              []string{"http://10.120.1.226:9200"},
//		Addrs:              []string{"http://192.168.1.16:9200"},
//		Username:           "",
//		Password:           "",
//		HealthcheckEnabled: true,
//		SnifferEnabled:     false,
//	})
//	if err != nil {
//		t.Fatal("NewClient err ", err)
//	}
//
//	b, err := es.IndexExists(index).Do(ctx)
//	if err != nil {
//		t.Fatal("NewClient err ", err)
//	}
//	if !b {
//		// Create index
//		createIndex, err := es.CreateIndex(index).Body(testMapping).IncludeTypeName(true).Do(ctx)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if createIndex == nil {
//			t.Errorf("expected result to be != nil; got: %v", createIndex)
//		}
//	}
//
//	//更多操作方式请看官方文档，README.md，里有各种文档地址和说明
//	//官方文档https://olivere.github.io/elastic/
//
//	t.Run("insert/update/delete", func(t *testing.T) {
//		tests := []elastic.BulkableRequest{
//			elastic.NewBulkIndexRequest().Index(index).Type(doc).Id("1").
//				Doc(tweet{User: "mf", Message: "Welcome to Golang and Elasticsearch."}),
//
//			elastic.NewBulkIndexRequest().Index(index).Type(doc).Id("2").
//				Doc(tweet{User: "mf", Message: "Dancing all night long. Yeah."}),
//		}
//
//		res, err := es.Bulk().Add(tests...).Do(ctx)
//		if err != nil {
//			t.Error(err)
//		}
//		t.Log(res.Succeeded())
//
//		//更新
//		upres, err := es.Update().Index(index).Type(doc).Id("1").
//			Doc(map[string]interface{}{"message": "update_message"}).
//			DetectNoop(true).
//			FetchSource(true).
//			Do(ctx)
//
//		if err != nil {
//			t.Fatal("error:", err)
//		}
//
//		if upres.GetResult == nil {
//			t.Fatal("expected GetResult != nil")
//		}
//
//		data, err := json.Marshal(upres.GetResult.Source)
//		if err != nil {
//			t.Fatalf("expected to marshal body as JSON, got: %v", err)
//		}
//
//		t.Log(string(data))
//
//		//删除
//		delres, err := es.Delete().Index(index).Type(doc).Id("1").Do(ctx)
//		if err != nil {
//			t.Fatal("error:", err)
//		}
//
//		if want, have := "deleted", delres.Result; want != have {
//			t.Errorf("expected Result = %q; got %q", want, have)
//		}
//	})
//
//	t.Run("search", func(t *testing.T) {
//		tests := []elastic.BulkableRequest{
//			elastic.NewBulkIndexRequest().Index(index).Type(doc).Id("3").
//				Doc(tweet{User: "mf", Message: "Welcome to Golang and Elasticsearch.3333"}),
//
//			elastic.NewBulkIndexRequest().Index(index).Type(doc).Id("4").
//				Doc(tweet{User: "mf", Message: "Dancing all night long. Yeah.444444"}),
//		}
//		res, err := es.Bulk().Add(tests...).Do(ctx)
//		if err != nil {
//			t.Fatalf("Bulk error: %v", err)
//		}
//
//		t.Log(res.Succeeded())
//
//		//查询语句太多，挑啦个复合查询，比较有代表性的查询
//		queryBuilder = elastic.NewBoolQuery().Must(elastic.NewTermQuery("user", "mf"))
//
//		searchResult, err := es.Search(index).
//			Type(doc).
//			Query(queryBuilder).
//			Size(100).
//			Do(ctx)
//
//		if err != nil {
//			t.Error("ES query err:", err)
//		} else {
//			t.Log("TotalHits:", searchResult.TotalHits())
//			if searchResult.TotalHits() > 0 {
//				ress := make(map[string]interface{})
//				for _, v := range searchResult.Hits.Hits {
//					err := json.Unmarshal(*v.Source, &ress)
//					if err != nil {
//						t.Log(err)
//					}
//					t.Log(ress)
//				}
//			}
//		}
//	})
//
//}

func TestSearchV7(t *testing.T) {

	var (
		queryBuilder *v7.BoolQuery
		err          error
		es           *v7.Client
		//tr           *tracer.Tracer
		//span         tracer.TraceOperator
		index = "elastic-test"
		doc   = "doc"
		ctx   = context.Background()
	)

	//获取es
	es, err = NewV7(&ElasticConfig{
		Addrs:              []string{"http://127.0.0.1:9200"},
		Username:           "",
		Password:           "",
		HealthcheckEnabled: true,
		SnifferEnabled:     false,
	})
	if err != nil {
		t.Fatal("NewClient err ", err)
	}

	b, err := es.IndexExists(index).Do(ctx)
	if err != nil {
		t.Fatal("NewClient err ", err)
	}

	if !b {
		// Create index
		createIndex, err := es.CreateIndex(index).Body(testMapping).Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if createIndex == nil {
			t.Errorf("expected result to be != nil; got: %v", createIndex)
		}
	}

	//更多操作方式请看官方文档，README.md，里有各种文档地址和说明
	//官方文档https://olivere.github.io/elastic/

	t.Run("insert/update/delete", func(t *testing.T) {
		tests := []v7.BulkableRequest{
			v7.NewBulkIndexRequest().Index(index).Id("1").
				Doc(tweet{User: "mf", Message: "Welcome to Golang and Elasticsearch."}),

			v7.NewBulkIndexRequest().Index(index).Id("2").
				Doc(tweet{User: "mf", Message: "Dancing all night long. Yeah."}),
		}

		res, err := es.Bulk().Add(tests...).Do(ctx)
		if err != nil {
			t.Error(err)
		}

		t.Log(res.Succeeded())

		//更新
		upres, err := es.Update().Index(index).Id("1").
			Doc(map[string]interface{}{"message": "update_message"}).
			DetectNoop(true).
			FetchSource(true).
			Do(ctx)

		if err != nil {
			t.Fatal("error:", err)
		}

		if upres.GetResult == nil {
			t.Fatal("expected GetResult != nil")
		}

		data, err := json.Marshal(upres.GetResult.Source)
		if err != nil {
			t.Fatalf("expected to marshal body as JSON, got: %v", err)
		}

		t.Log(string(data))

		//删除
		delres, err := es.Delete().Index(index).Id("1").Do(ctx)
		if err != nil {
			t.Fatal("error:", err)
		}

		if want, have := "deleted", delres.Result; want != have {
			t.Errorf("expected Result = %q; got %q", want, have)
		}
	})

	t.Run("search", func(t *testing.T) {
		tests := []v7.BulkableRequest{
			v7.NewBulkIndexRequest().Index(index).Type(doc).Id("3").
				Doc(tweet{User: "mf", Message: "Welcome to Golang and Elasticsearch.3333"}),

			v7.NewBulkIndexRequest().Index(index).Type(doc).Id("4").
				Doc(tweet{User: "mf", Message: "Dancing all night long. Yeah.444444"}),
		}
		res, err := es.Bulk().Add(tests...).Do(ctx)
		if err != nil {
			t.Fatalf("Bulk error: %v", err)
		}

		t.Log(res.Succeeded())

		//查询语句太多，挑啦个复合查询，比较有代表性的查询
		queryBuilder = v7.NewBoolQuery().Must(v7.NewTermQuery("user", "mf"))

		searchResult, err := es.Search(index).
			Query(queryBuilder).
			Size(100).
			Do(ctx)

		if err != nil {
			t.Error("ES query err:", err)
		} else {
			t.Log("TotalHits:", searchResult.TotalHits())
			if searchResult.TotalHits() > 0 {
				ress := make(map[string]interface{})
				for _, v := range searchResult.Hits.Hits {
					err := json.Unmarshal(v.Source, &ress)
					if err != nil {
						t.Log(err)
					}
					t.Log(ress)
				}
			}
		}
	})

}
