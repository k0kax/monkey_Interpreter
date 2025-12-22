package object

//各种对象

import (
	"bytes"
	"fmt"
	"monkey_Interpreter/ast"
	"strings"
)

// 对象类型
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"      //整数
	BOOLEAN_OBJ      = "BOOLEAN"      //布尔
	NULL_OBJ         = "NULL"         //空值
	RETURN_VALUE_OBJ = "RETURN_VALUE" //返回值
	ERROR_OBJ        = "ERROR"        //错误
	FUNCTION_OBJ     = "FUNCTION"     //函数
)

// 对象接口
type Object interface {
	Type() ObjectType //对象类型
	Inspect() string  //辅助调试用
}

// 整数对象
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// 布尔型对象
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// 空值对象
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// 用来存放返回值对象
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// 存放错误
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// 函数对象
type Function struct {
	Parameters []*ast.Identifier   //形参
	Body       *ast.BlockStatement //函数体
	Env        *Environment        //函数定义时的环境，闭包用
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(")")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
