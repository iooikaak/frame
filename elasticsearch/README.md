
# 快速开始
* es6

```go
        //获取es
        	es, err = New(&ElasticConfig{
        		Addrs:              []string{"http://10.120.1.226:9200"},
        		Username:           "",
        		Password:           "",
        		HealthcheckEnabled: true,
        		SnifferEnabled:     false,
        	})
        	if err != nil {
        		panic("NewClient err ", err)
        	}
        
        	b, err := es.IndexExists(index).Do(ctx)
        	if !b {
        		// Create index
        		createIndex, err := es.CreateIndex(index).Body(testMapping).IncludeTypeName(true).Do(ctx)
        		if err != nil {
        			panic(err)
        		}
        		if createIndex == nil {
        			xlog.Errorf("expected result to be != nil; got: %v", createIndex)
        		}
        	}
        
        	//更多操作方式请看官方文档
        	//官方文档https://olivere.github.io/elastic/

```


* es7

```go
       //获取es
       	es, err = NewV7(&ElasticConfig{
       		Addrs:              []string{"http://127.0.0.1:9200"},
       		Username:           "",
       		Password:           "",
       		HealthcheckEnabled: true,
       		SnifferEnabled:     false,
       	})
       	if err != nil {
       		xlog.Errorf("NewClient err ", err)
       	}
       
       	b, err := es.IndexExists(index).Do(ctx)
       	if !b {
       		// Create index
       		createIndex, err := es.CreateIndex(index).Body(testMapping).Do(ctx)
       		if err != nil {
       			xlog.Error(err)
       		}
       		if createIndex == nil {
       			xlog.Errorf("expected result to be != nil; got: %v", createIndex)
       		}
       	}
       
       	//更多操作方式请看官方文档
       	//官方文档https://olivere.github.io/elastic/
       

```