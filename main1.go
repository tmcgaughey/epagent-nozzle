package main

import (
//    "crypto/tls"
//	"os"
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
	"strings"
)

type WMetric struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Value string `json:"value"`
}

type MetricFeed struct {
	Metrics []WMetric `json:"metrics"`
}

const EPAgent = "http://192.168.1.23:9191/apm/metricFeed"

func main1() {
	
	fmt.Println("URL:>",EPAgent)
	
	metriclist := new (MetricFeed)
	wmetrics := []WMetric{}
	
	mVal := WMetric{Type: "LongCounter", Name: "totalMetricsSent", Value: "0"}
	
	mVal.Name = "Firehose." + mVal.Name
		str := strings.Replace(mVal.Name, ".","|",-1)
		fmt.Println("str: " + str)
		idx := strings.LastIndexAny(str,"|")
		//fmt.Println("idx: " +idx)
		substring := str[idx:len(str)]
		fmt.Println("substring: " + substring)
		revsubstring := ":"+substring[1:len(substring)]
		fmt.Println("revsubstring: " + revsubstring)
		mVal.Name = strings.Replace(str, substring, revsubstring, -1)
		fmt.Printf("metric: %s\tvalue: %s", mVal.Name, mVal.Value)
	
	wmetrics = append(wmetrics, WMetric{
			Type: "StringEvent", 
			Name: "cf:Test", 
			Value: "test",
		})
	
	metriclist.Metrics = wmetrics
	
	//fmt.Println(wmetric)
	fmt.Println(metriclist)

    b, err := json.Marshal(metriclist)

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(string(b[:]))
	
	req, err := http.NewRequest("POST", EPAgent, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err !=nil {
		panic(err)
	}
	
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Haeders:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	
}