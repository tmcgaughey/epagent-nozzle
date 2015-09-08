package main

import (
//    "crypto/tls"
//	"os"
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
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

func main() {
	
	fmt.Println("URL:>",EPAgent)
	
	metriclist := new (MetricFeed)
	wmetrics := []WMetric{}
	
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