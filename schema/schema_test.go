package schema

import (
	"gorm/dialect"
	"gorm/log"
	"testing"
)

type User struct {
	Name string `gorm:"PRIMARY KEY"`
	Age  int
}

func TestParse(t *testing.T) {
	sqlite, _ := dialect.GetDialect("sqlite3")
	schema := Parse(&User{}, sqlite)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("解析User结构体错误")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("主键解析错误")
	}
}

func TestRecordValues(t *testing.T) {
	sqlite, _ := dialect.GetDialect("sqlite3")
	schema := Parse(&User{}, sqlite)
	vars := schema.RecordValues(&User{Name: "hmj", Age: 22})
	log.Info(vars...)
	if len(vars) == 0 {
		t.Fatal("解析失败")
	}
}
