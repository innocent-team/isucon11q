package main

import (
	"fmt"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/labstack/gommon/log"
)

const INFLUX_WRITE_SPAN = 500 * time.Millisecond

const (
	fTime = "time"
	fJIAIsuUUID = "jiaIsuUUID"
	fCondition = "condition"
	fMessage = "message"
	fIsSitting = "isSitting"
	fConditionLevel = "conditionLevel"
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
		fJIAIsuUUID: jiaIsuUUID,
	}
	conditionLevel, err := calculateConditionLevel(condition)
	if err != nil {
		return nil, fmt.Errorf("Error condition level: %w", err)
	}
	fields := map[string]interface{}{
		fIsSitting:      isSitting,
		fCondition:      condition,
		fMessage:        message,
		fConditionLevel: conditionLevel,
	}
	point, err := client.NewPoint("condition", tags, fields, timestamp)
	if err != nil {
		return nil, fmt.Errorf("Error New Point: %w", err)
	}
	return point, nil
}

var conditionPoints client.BatchPoints

func InsertConditions(jiaIsuUUID string, timestamp time.Time, isSitting bool, scondition string, message string) error {
	log.Print("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
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
		log.Printf("%#+v", conditionPoints)
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
	if len(result.Series) == 0 {
		return res
	}
	m := columnMap(result.Series[0].Columns)
	for _, v := range result.Series[0].Values {
		timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[m[fTime]].(string))
		if err != nil {
			log.Printf("error: timestamp  %v", err)
			continue
		}
		condition := InfluxCondition {
			Timestamp: timestamp,
			Condition: v[m[fCondition]].(string),
			ConditionLevel: v[m[fConditionLevel]].(string),
			IsSitting: v[m[fIsSitting]].(bool),
			JIAIsuUUID: v[m[fJIAIsuUUID]].(string),
			Message: v[m[fMessage]].(string),
		}

		res = append(res, condition)
	}
	return res
}

func columnMap(columns []string) map[string]int {
	m := map[string]int{}
	for i, v := range columns {
		m[v] = i
	}
	return m
}

func getLastCondtionsByIsuList(isuList []Isu) (map[string]InfluxCondition, error) {
	var builder strings.Builder
	builder.WriteString(`SELECT last(*) FROM "condition" WHERE "jiaIsuUUID" =~ /`)
    for i, id := range isuList {
		if i != 0 {
			builder.WriteString(`|`)
		}
		builder.WriteString(id.JIAIsuUUID)
    }
	builder.WriteString(`/ GROUP BY "jiaIsuUUID" ORDER BY "time" DESC `)
	q := client.NewQueryWithParameters(builder.String(), "isu", "", client.Params{})
	c := InfluxClient()
	result, err := c.Query(q)
	if err != nil {
		return nil, fmt.Errorf("Errro query: %w", err)
	}

	influxConditionsMap := map[string]InfluxCondition{}
	for _, row := range result.Results[0].Series {
		m := columnMap(row.Columns)
		jiaIsuUUID := row.Tags[fJIAIsuUUID]
		for _, v := range row.Values {
			timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[m["time"]].(string))
			if err != nil {
				log.Printf("error: timestamp  %v", err)
				continue
			}
			condition := InfluxCondition {
				Timestamp: timestamp,
				Condition: v[m["last_condition"]].(string),
				ConditionLevel: v[m["last_conditionLevel"]].(string),
				IsSitting: v[m["last_isSitting"]].(bool),
				JIAIsuUUID: jiaIsuUUID,
				Message: v[m["last_message"]].(string),
			}
			influxConditionsMap[jiaIsuUUID] = condition
		}
	}
	return influxConditionsMap, nil
}