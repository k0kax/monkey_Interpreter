package object

//各种对象

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	STRING_OBJ       = "STRING"       //字符串
	BUILTIN_OBJ      = "BUILTIN"      //内置函数
	ARRAY_OBJ        = "ARRAY"        //数组
	HASH_OBJ         = "HASH"         //哈希表
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

// 字符串
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// 内置函数
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// 数组
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ","))
	out.WriteString("]")

	return out.String()
}

// 哈希表的外层键值
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// 布尔
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// 整数
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// 字符串
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()} //可能哈希碰撞
}

type Hashable interface { //外层键值函数接口
	HashKey() HashKey
}

// 哈希表的真实键值对
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ","))
	out.WriteString("}")

	return out.String()
}
