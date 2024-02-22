package eureka

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	// create eureka client
	client := NewClient(&Config{
		DefaultZone:                    []string{"http://peer1:1111", "http://peer2:2222", "http://peer3:3333"},
		App:                            "go-example",
		Port:                           10000,
		RenewalIntervalInSecs:          1,
		DurationInSecs:                 3,
		RollDiscoveriesIntervalSeconds: 3,
		DataCenterName:                 "discovery",
		Metadata: map[string]string{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        "0",
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})
	// start client, register、heartbeat、GetApplications
	client.Start()

	// http server
	http.HandleFunc("/v1/services", func(writer http.ResponseWriter, request *http.Request) {
		// full applications from eureka server
		urls, _ := client.GetService("go-example")
		b, _ := json.Marshal(urls)
		_, _ = writer.Write(b)
	})

	http.HandleFunc("/v2/services", func(writer http.ResponseWriter, request *http.Request) {
		// full applications from eureka server
		urls := client.getCenterUrl()
		sli := make([]string, 0)
		if urls != nil {
			for i := 0; i < urls.Len(); i++ {
				sli = append(sli, urls.Value.(string))
				urls = urls.Next()
			}
		}
		b, _ := json.Marshal(sli)
		_, _ = writer.Write(b)
	})

	// start http server
	if err := http.ListenAndServe(":10000", nil); err != nil {
		fmt.Println(err)
	}
}

func TestUnRegister(t *testing.T) {
	err := UnRegister("http://127.0.0.1:8761", "mf.demoadmintst.micro", "192.168.94.231:mf.demoadmintst.micro:10087")
	t.Logf("%v", err)
}
