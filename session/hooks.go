package session

const (
	BEFORE_QUERY  = "BeforeQuery"
	AFTER_QUERY   = "AfterQuery"
	BEFORE_UPDATE = "BeforeUpdate"
	AFTER_UPDATE  = "AfterUpdate"
	BEFORE_DELETE = "BeforeDelete"
	AFTER_DELETE  = "AfterDelete"
	BEFORE_INSERT = "BeforeInsert"
	AFTER_INSERT  = "AfterInsert"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}

type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}

type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}

type IAfterDelete interface {
	AfterDelete(s *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

type IAfterInster interface {
	AfterInstert(s *Session) error
}

func (s *Session) CallMethod(method string, value interface{}) {
	if value == nil {
		value = s.RefTable().Model
	}
	switch method {
	case BEFORE_QUERY:
		if i, ok := value.(IBeforeQuery); ok {
			i.BeforeQuery(s)
		}
	case AFTER_QUERY:
		if i, ok := value.(IAfterQuery); ok {
			i.AfterQuery(s)
		}
	case BEFORE_UPDATE:
		if i, ok := value.(IBeforeUpdate); ok {
			i.BeforeUpdate(s)
		}
	case AFTER_UPDATE:
		if i, ok := value.(IAfterUpdate); ok {
			i.AfterUpdate(s)
		}
	case BEFORE_INSERT:
		if i, ok := value.(IBeforeInsert); ok {
			i.BeforeInsert(s)
		}
	case AFTER_INSERT:
		if i, ok := value.(IAfterInster); ok {
			i.AfterInstert(s)
		}
	case BEFORE_DELETE:
		if i, ok := value.(IBeforeDelete); ok {
			i.BeforeDelete(s)
		}
	case AFTER_DELETE:
		if i, ok := value.(IAfterDelete); ok {
			i.AfterDelete(s)
		}
	default:
		panic("unsupport hook type")
	}
}
