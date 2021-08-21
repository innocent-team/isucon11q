package main

import (
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
	"time"
)

func TestInflux(t *testing.T) {
	WriteCondition()
	InsertConditions("111", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー", "なまいき")
	WriteCondition()
    PrintInfluxdb()
}

func TestCreatePoint(t *testing.T) {
	t.Log(CreatePoint("111", time.Now(), true, "is_dirty=false,is_overweight=false,is_broken=false", "へろー", "なまいき"))
}
