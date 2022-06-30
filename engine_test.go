package gorm

import (
	"gorm/log"
	"gorm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `gorm:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	e, err1 := NewEngine("sqlite3", "gee.db")
	s := e.NewSession().Model(&User{})
	ok := s.HasTable()
	if err1 != nil || !ok {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	row, err := s.Insert(user1, user2, user3)
	if row == 0 || err != nil {
		t.Fatal("插入数据失败")
	}
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age DESC").One(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count()
	t.Log(affected, count)
	if affected != 1 || count != 2 {
		t.Fatal("failed to delete or count")
	}
}

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	log.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	t.Helper()
	e, _ := NewEngine("sqlite3", "gee.db")
	defer e.Close()
	s := e.NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.One(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}
