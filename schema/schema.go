package schema

import (
	"go/ast"
	"gorm/dialect"
	"reflect"
)

type Field struct {
	// 属性名字
	Name string

	// 属性类型
	Type string

	// 属性的约束条件（tag）
	Tag string
}

type Schema struct {
	// 被映射对象的model
	Model interface{}

	// 表名
	Name string

	// 所有的字段
	Fields []*Field

	// 所有的字段名称
	FieldNames []string

	// 记录字段名和 Field 的映射关系，方便之后直接使用，无需遍历 Fields
	fieldMap map[string]*Field
}

// 对外暴露Field
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// 解析结构体字段
func Parse(target interface{}, dialect dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(target)).Type()
	schema := &Schema{
		Model:    target,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 判断是否是嵌入字段或是否是大写字母开头
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("gorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// func Parse(target interface{}, dialect dialect.Dialect) *Schema {
// 	tv := reflect.ValueOf(target).Elem()
// 	tp := tv.Type()
// 	schema := &Schema{
// 		Model:    target,
// 		Name:     tp.Name(),
// 		fieldMap: make(map[string]*Field),
// 	}
// 	for i := 0; i < tp.NumField(); i++ {
// 		p := tp.Field(i)
// 		if !p.Anonymous && ast.IsExported(p.Name) {
// 			field := &Field{
// 				Name: p.Name,
// 				Type: dialect.DataTypeOf(reflect.New(p.Type).Elem()),
// 			}
// 			if v, ok := p.Tag.Lookup("gorm"); ok {
// 				field.Tag = v
// 			}
// 			schema.Fields = append(schema.Fields, field)
// 			schema.FieldNames = append(schema.FieldNames, p.Name)
// 			schema.fieldMap[p.Name] = field
// 		}
// 	}
// 	return schema
// }

// 获取传入结构体的值
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fields []interface{}
	for _, filed := range s.Fields {
		fields = append(fields, destValue.FieldByName(filed.Name).Interface())
	}
	return fields
}
