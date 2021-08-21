package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/labstack/gommon/log"
)

const INFLUX_WRITE_SPAN = 500 * time.Millisecond

const (
	fTime           = "time"
	fJIAIsuUUID     = "jiaIsuUUID"
	fCondition      = "condition"
	fMessage        = "message"
	fIsSitting      = "isSitting"
	fConditionLevel = "conditionLevel"
	fCharacter = "character"
	fIsuID = "isuID"
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

func CreatePoint(isuID int,jiaIsuUUID string, timestamp time.Time, isSitting bool, condition string, message string, character string) (*client.Point, error) {
	tags := map[string]string{
		fJIAIsuUUID: jiaIsuUUID,
		fCharacter: character,
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
		fIsuID: isuID,
	}
	point, err := client.NewPoint("condition", tags, fields, timestamp)
	if err != nil {
		return nil, fmt.Errorf("Error New Point: %w", err)
	}
	return point, nil
}

var conditionPoints client.BatchPoints

func InsertConditions(isuID int, jiaIsuUUID string, timestamp time.Time, isSitting bool, scondition string, message string, character string) error {
	point, err := CreatePoint(isuID, jiaIsuUUID, timestamp, isSitting, scondition, message, character)
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
	Timestamp      time.Time
	Condition      string
	IsSitting      bool
	JIAIsuUUID     string
	Message        string
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
		condition := InfluxCondition{
			Timestamp:      timestamp,
			Condition:      v[m[fCondition]].(string),
			ConditionLevel: v[m[fConditionLevel]].(string),
			IsSitting:      v[m[fIsSitting]].(bool),
			JIAIsuUUID:     v[m[fJIAIsuUUID]].(string),
			Message:        v[m[fMessage]].(string),
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
	for _, id := range isuList {
		builder.WriteString(`SELECT * FROM "condition" WHERE "jiaIsuUUID" = '`)
		builder.WriteString(id.JIAIsuUUID)
		builder.WriteString(`' ORDER BY "time" DESC LIMIT 1; `)
	}
	q := client.NewQueryWithParameters(builder.String(), "isu", "", client.Params{})
	c := InfluxClient()
	result, err := c.Query(q)
	if err != nil {
		return nil, fmt.Errorf("Errro query: %w", err)
	}

	influxConditionsMap := map[string]InfluxCondition{}
	for _, res := range result.Results {
		if len(res.Series) == 0 {
			continue
		}
		row := res.Series[0]
		m := columnMap(row.Columns)
		for _, v := range row.Values {
			timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[m["time"]].(string))
			if err != nil {
				log.Printf("error: timestamp  %v", err)
				continue
			}
			id := v[m[fJIAIsuUUID]].(string)
			log.Printf("!!!!!!!!!!!!!!!%#v", res)
			condition := InfluxCondition{
				Timestamp:      timestamp,
				Condition:      v[m[fCondition]].(string),
				ConditionLevel: v[m[fConditionLevel]].(string),
				IsSitting:      v[m[fIsSitting]].(bool),
				JIAIsuUUID:     id,
				Message:        v[m[fMessage]].(string),
			}
			influxConditionsMap[id] = condition
		}
	}
	return influxConditionsMap, nil
}

func getTrendByCharacterType(character string) (TrendResponse, error) {
	res := TrendResponse{
				Character: character,
				Info:      []*TrendCondition{},
				Warning:   []*TrendCondition{},
				Critical:  []*TrendCondition{},
			}
	c := InfluxClient()
	defer c.Close()

	q := client.NewQueryWithParameters(`SELECT last(*) FROM condition WHERE character = $character GROUP BY jiaIsuUUID ORDER BY time DESC` , "isu", "", client.Params{
		"character": character,
	})
	resp, err := c.Query(q)
	if err != nil {
		return res, err
	}
	if resp.Err != "" {
		return res, err
	}
	for _, row := range resp.Results[0].Series {
		m := columnMap(row.Columns)
		for _, v :=  range row.Values {
			timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[m["time"]].(string))
			if err != nil {
				log.Printf("error: timestamp  %v", err)
				continue
			}
			id, err :=  v[m["last_isuID"]].(json.Number).Int64()
			if err != nil {
				log.Printf("error: number  %v", err)
				continue
			}
			level :=  v[m["last_conditionLevel"]].(string)
			cond := &TrendCondition{
				ID: int(id),
				Timestamp: timestamp.Unix(),
			}
			fmt.Printf("Type: %s\n", level)
			switch level {
				case "info":
					res.Info = append(res.Info, cond)
				case "warning":
					res.Warning = append(res.Warning, cond)
				case "critical":
					res.Critical = append(res.Critical, cond)
			}
		}
	}
	return res, nil
}