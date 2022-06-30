package gorm

import (
	"database/sql"
	"gorm/dialect"
	"gorm/log"
	"gorm/session"
)

type Engine struct {
	db *sql.DB

	dialect dialect.Dialect
}

type TxFunc func(s *session.Session) (interface{}, error)

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	// 检测打开数据库是否成功
	if err != nil {
		log.Error(err)
		return
	}
	// 检测连接是否存活
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	d, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s 未发现", driver)
	}
	e = &Engine{db: db, dialect: d}
	log.Info("Connect database success")
	return
}

// 关闭所有连接
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			s.Rollback()
			panic(p)
		} else if err != nil {
			s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}
