package main

import (
	"strings"
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

func TestInflux(t *testing.T) {
	WriteCondition()
	InsertConditions("111", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー")
	WriteCondition()
	PrintInfluxdb()
}

func TestCreatePoint(t *testing.T) {
	t.Log(CreatePoint("111", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー"))
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

	conditions := []IsuCondition{}

	if len(result.Results[0].Series) != 0 {
		for _, v := range result.Results[0].Series[0].Values {
			condition := IsuCondition{}
			timestamp, err := time.Parse("2006-01-02T15:04:05Z0700", v[0].(string))
			if err != nil {
				t.Fatal(err)
			}
			condition.Timestamp = timestamp
			condition.Condition = v[1].(string)
			condition.IsSitting = v[3].(bool)
			condition.JIAIsuUUID = v[4].(string)
			condition.Message = v[5].(string)
			conditions = append(conditions, condition)
		}
	}

	t.Logf("%+#v", conditions)
}
