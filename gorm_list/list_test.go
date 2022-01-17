package gorm_list

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

type TestReq struct {
	ListHandler
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

func (t *TestReq) InitQuery() *ListHandler {
	result := make(map[string]interface{})
	result["a"] = t.A  // 用户代码， 仅需要填写这个部分的代码即可
	result["b"] = t.B
	result["c"] = t.C

	t.ListHandler.QueryMap = result
	return &t.ListHandler
}

type TestDBStruct struct {
	A string `gorm:"column:a" json:"a" `
	B string `gorm:"column:b" json:"b"`
	C string `gorm:"column:c" json:"c"`
}

// TableName ...
func (*TestDBStruct) TableName() string {
	return "test"
}

var DB *gorm.DB

const (
	defaultActive  = 60
	defaultIdle    = 5
	defaultTimeout = 3600 * time.Second
	defaultAddr = ""
)

func init() {
	// 初始化MySQL连接
	db, err := gorm.Open(mysql.Open(defaultAddr), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	if err != nil {
		panic(fmt.Sprintf("init mysql connection failed, err: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("get sqlDB, err: %v", err))
	}

	sqlDB.SetMaxOpenConns(defaultActive)
	sqlDB.SetMaxIdleConns(defaultIdle)
	sqlDB.SetConnMaxLifetime(defaultTimeout)
}

func TestClient(t *testing.T) {
	req := TestReq{
		ListHandler: ListHandler{PageNum: 10, PageSize: 10},
		A: "a",
		B: "b",
	}

	result := TestDBStruct{}
	total, err := List(context.Background(), DB, &result, result.TableName(), req.InitQuery().MakeQueryOptions()...)
	if err != nil {
		t.Log(fmt.Errorf("list err %v", err))
	}

	t.Log(result)
	t.Log(total)
}

func TestClient1(t *testing.T) {
	req := TestReq{
		A: "a",
		B: "b",
	}

	result := TestDBStruct{}
	total, err := List(context.Background(), DB, &result, result.TableName(),
		WithPageNum(10),
		WithPageSize(10),
		WithQueryString("a = ?", req.A))
	if err != nil {
		t.Log(fmt.Errorf("list err %v", err))
	}

	t.Log(result)
	t.Log(total)
}