package main

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type IsuConditionInflux struct {
	ID         int      
	JIAIsuUUID string  
	Timestamp  time.Time 
	IsSitting  bool     
	Message    string    
	Condition  string    
	IsDirty bool
	IsOverweight bool
	IsBroken bool
}

func InfluxClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	return c
}

func InsertConditions() {
	conditions := client.NewBatchPoints(client.)
}
func Influxdb() {
	c := InfluxClient()
	defer c.Close()

	q := client.NewQuery("SELECT count(value) FROM condition", "isu", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
