package session

import (
	"fmt"
	"gorm/log"
	"gorm/schema"
	"reflect"
	"strings"
)

// 解析操作是比较耗时的，因此将解析的结果保存在成员变量refTable中，
// 即使 Model() 被调用多次，如果传入的结构体名称不发生变化，则不会更新 refTable 的值。
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// 对外提供访问refTable的方法，并对refTable是否为空进行判断
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model 未被设置")
	}
	return s.refTable
}

// 根据结构体创建表
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var colums []string
	for _, v := range table.Fields {
		colums = append(colums, fmt.Sprintf("%s %s %s", v.Name, v.Type, v.Tag))
	}
	desc := strings.Join(colums, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// 删除与表名相同的表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP Table IF EXISTS %s;", s.RefTable().Name)).Exec()
	return err
}

// 判断数据库中该表是否存在
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
