package dialect

import "reflect"

type Dialect interface {
	// 将 Go语言的类型转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// 返回某个表是否存在的SQL语句，参数是表名(table)。
	TableExistSQL(tableName string) (string, []interface{})
}

var dialectMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (d Dialect, ok bool) {
	d, ok = dialectMap[name]
	return
}

// 随笔：return做了那几件事
// 1. 若为匿名返回参数，将创建该类型的临时变量，保存返回值
// 2. 执行defer函数，若是匿名函数则返回值不会再次改变，若是有名则会改变
// 3. 返回值（临时变量）或者 有名变量
