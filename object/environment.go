package object

//
type Environment struct {
	store map[string]Object
	outer *Environment //上一层环境
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// 获取Object对象
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// 设置Object对象
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// 外层环境
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
