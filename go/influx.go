package main

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

const INFLUX_WRITE_SPAN = 500 * time.Millisecond

const (
	rowIndexTimestamp      = 0
	rowIndexCondition      = 1
	rowIndexConditionLevel = 2
	rowIndexIsSitting      = 3
	rowIndexJIAIsuUUID     = 4
	rowIndexMessage        = 5
)

var influxAddr string

func InfluxClient() client.Client {
	if influxAddr == "" {
		influxAddr = getEnv("INFLUX_ADDR", "http://localhost:8086")
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxAddr,
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	return c
}

func CreatePoint(jiaIsuUUID string, timestamp time.Time, isSitting bool, condition string, message string) (*client.Point, error) {
	tags := map[string]string{
		"jiaIsuUUID": jiaIsuUUID,
	}
	conditionLevel, err := calculateConditionLevel(condition)
	if err != nil {
		return nil, fmt.Errorf("Error condition level: %w", err)
	}
	fields := map[string]interface{}{
		"isSitting":      isSitting,
		"condition":      condition,
		"message":        message,
		"conditionLevel": conditionLevel,
	}
	point, err := client.NewPoint("condition", tags, fields, timestamp)
	if err != nil {
		return nil, fmt.Errorf("Error New Point: %w", err)
	}
	return point, nil
}

var conditionPoints client.BatchPoints

func InsertConditions(jiaIsuUUID string, timestamp time.Time, isSitting bool, scondition string, message string) error {
	point, err := CreatePoint(jiaIsuUUID, timestamp, isSitting, scondition, message)
	if err != nil {
		return fmt.Errorf("Error CreatePoint: %w", err)
	}
	conditionPoints.AddPoint(point)
	WriteCondition()
	return nil
}

// conditionPointsを初期化 + Write
func WriteCondition() {
	if conditionPoints != nil && len(conditionPoints.Points()) > 0 {
		c := InfluxClient()
		defer c.Close()
		err := c.Write(conditionPoints)
		if err != nil {
			fmt.Println("Error Influx Write: ", err.Error())
		}
	}
	var err error
	conditionPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database: "isu",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}
}

func StartInfluxCondition() {
	WriteCondition()
	go func() {
		for {
			WriteCondition()
			time.Sleep(INFLUX_WRITE_SPAN)
		}
	}()
}

func PrintInfluxdb() {
	c := InfluxClient()
	defer c.Close()

	q := client.NewQuery("SELECT * FROM condition", "isu", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		fmt.Printf("%#+v\n", response.Results)
	}
}

type InfluxCondition struct {
	Timestamp time.Time
	Condition string
	IsSitting bool
	JIAIsuUUID string
	Message string
	ConditionLevel string
}

func ResultInfluxConditons(result client.Result) []InfluxCondition {
	res := []InfluxCondition{}
	for _, v := range result.Series[0].Values {
		timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[rowIndexTimestamp].(string))
		if err != nil {
			log.Print("error: timestamp  %v", err)
			continue
		}
		condition := InfluxCondition {
			Timestamp: timestamp,
			Condition: v[rowIndexCondition].(string),
			ConditionLevel: v[rowIndexConditionLevel].(string),
			IsSitting: v[rowIndexIsSitting].(bool),
			JIAIsuUUID: v[rowIndexJIAIsuUUID].(string),
			Message: v[rowIndexMessage].(string),
		}

		res = append(res, condition)
	}
	return res
}