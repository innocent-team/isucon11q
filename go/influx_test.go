package main

import (
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/labstack/gommon/log"
)

func TestInflux(t *testing.T) {
	WriteCondition()
	InsertConditions("222", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー", "やっかみ")
	InsertConditions("111", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー", "なまいき")
	WriteCondition()
	PrintInfluxdb()
}

func TestCreatePoint(t *testing.T) {
	t.Log(CreatePoint("111", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー", "なまいき"))
}

func TestIsuConditions(t *testing.T) {
	TestInflux(t)
	q := client.NewQueryWithParameters(`SELECT * FROM "condition" WHERE "jiaIsuUUID" = $jiaIsuUUID AND "time" < $endTime ORDER BY "time" DESC`, "isu", "", client.Params{
		"jiaIsuUUID": "111",
		"endTime":    time.Now(),
	})
	c := InfluxClient()
	result, err := c.Query(q)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+#v", result)

	conditions := ResultInfluxConditons(result.Results[0])
	t.Logf("%+#v", conditions)
}

func TestInfluxByIDs(t *testing.T) {
	TestInflux(t)
	isuList := []Isu{{JIAIsuUUID: "111"}, {JIAIsuUUID: "222"}}
	
	influxConditionsMap, err := getLastCondtionsByIsuList(isuList)
	if err != nil {
		t.Fatal(err)
	}
	responseList := []GetIsuListResponse{}
	var formattedCondition *GetIsuConditionResponse
	t.Logf("%+#v", influxConditionsMap)
	for _, isu := range isuList {
		if condition, ok := influxConditionsMap[isu.JIAIsuUUID]; ok {
			formattedCondition = &GetIsuConditionResponse{
				JIAIsuUUID:     condition.JIAIsuUUID,
				IsuName:        isu.Name,
				Timestamp:      condition.Timestamp.Unix(),
				IsSitting:      condition.IsSitting,
				Condition:      condition.Condition,
				ConditionLevel: condition.ConditionLevel,
				Message:        condition.Message,
			}
		}

		res := GetIsuListResponse{
			ID:                 isu.ID,
			JIAIsuUUID:         isu.JIAIsuUUID,
			Name:               isu.Name,
			Character:          isu.Character,
			LatestIsuCondition: formattedCondition}
		responseList = append(responseList, res)
	}
	
	t.Logf("Response %v", responseList)
}

func TestIsuGraphResponse(t *testing.T) {
	TestInflux(t)
	q := client.NewQueryWithParameters(`
	SELECT * FROM "condition"
	WHERE "jiaIsuUUID" = $jiaIsuUUID
	ORDER BY "time" ASC`, "isu", "", client.Params{
		"jiaIsuUUID": "111",
		"endTime":    time.Now(),
	})
	c := InfluxClient()
	influxResp, err := c.Query(q)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+#v", influxResp)

	conditions := []IsuCondition{}

	if len(influxResp.Results[0].Series) != 0 {
		m := columnMap(influxResp.Results[0].Series[0].Columns)
		for _, v := range influxResp.Results[0].Series[0].Values {
			condition := IsuCondition{}
			timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[m[fTime]].(string))
			if err != nil {
				log.Print(err)
				continue
			}
			condition.Timestamp = timestamp
			condition.Condition = v[m[fCondition]].(string)
			condition.IsSitting = v[m[fIsSitting]].(bool)
			condition.JIAIsuUUID = v[m[fJIAIsuUUID]].(string)
			condition.Message = v[m[fMessage]].(string)
			conditions = append(conditions, condition)
		}
	}
	t.Logf("%+#v", conditions)
}