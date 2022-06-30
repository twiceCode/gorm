package session

import (
	"database/sql"
	"gorm/clause"
	"gorm/dialect"
	"gorm/log"
	"gorm/schema"
	"strings"
)

type Session struct {
	// 保存数据库的连接状态
	db *sql.DB

	// 构建SQL语句
	sql strings.Builder

	// 占位符对应的值
	sqlVars []interface{}

	// 指定不同数据库的映射
	dialect dialect.Dialect

	// 操作的表
	refTable *schema.Schema

	// sql语句生成器
	clause clause.Clause

	// 支持事务
	tx *sql.Tx
}

// 作为 *sql.DB和 *sql.Tx的公共父接口
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// 兼容前面的对*sql.DB的调用
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// 初始化Session结构体
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// 清理sql语句
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// 通过该方法可以构建SQL语句
func (s *Session) Raw(sql string, val ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, val...)
	return s
}

// 执行一个sql语句
func (s *Session) Exec() (result sql.Result, err error) {
	// 执行完默认清除语句
	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// 执行一个查询的sql语句
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 获取多条查询结果
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
