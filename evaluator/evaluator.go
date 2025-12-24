package evaluator

//求值器

import (
	"fmt"
	"monkey_Interpreter/ast"
	"monkey_Interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	//语句
	case *ast.Program:
		return evalProgram(node.Statements, env)

	//表达式语句
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	//整形字面量
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	//布尔型字面量
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	//前缀表达式
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	//中缀表达式
	case *ast.InfixExpression:

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	//块
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	//if语句
	case *ast.IfExpression:
		return evalIfExpression(node, env)

	//return语句
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	//let语句
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	//标识符
	case *ast.Identifier:
		return evalIdentifier(node, env)

	//函数
	case *ast.FunctionLiteral:
		params := node.Parameters                                         //取参数
		body := node.Body                                                 //取函数体
		return &object.Function{Parameters: params, Env: env, Body: body} //返回一个包含参数、函数体、绑定的环境 的对象

	//函数调用
	case *ast.CallExpression:
		function := Eval(node.Function, env) //获取调用的函数
		if isError(function) {               //检查函数合法性
			return function
		}

		args := evalExpression(node.Arguments, env) //求值所有实参，得到具体值
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		//调用函数，并传入求值后得到的实参
		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

// 原始布尔转为布尔对象
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

// 语句求值
func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// 前缀求值 中转函数
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// 非!的求值
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// 测试-前缀表达式 仅用于整数
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// 中缀表达式 中转函数
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ: //左右值都是整数
		return evalIntegerInfixExpression(operator, left, right) //整数的处理

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case operator == "&&":
		// 左值为假则直接返回左值
		if !isTruthy(left) {
			return left
		}
		// 左值为真则返回右值
		return right
	case operator == "||":
		// 左值为真则直接返回左值
		if isTruthy(left) {
			return left
		}
		// 左值为假则返回右值
		return right

	case left.Type() != right.Type(): //对象不同
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	//字符串拼接
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// 整数的中缀操作符处理
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value   //取左值
	rightVal := right.(*object.Integer).Value //取右值

	switch operator {
	//四则运算操作
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	//布尔操作
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// if选择语句的求值
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env) //处理条件
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) { //条件对 首选项
		return Eval(ie.Consequence, env) //执行结果
	} //首选项条件不对，且备选结果不为空

	//执行中间选项
	for _, al := range ie.Alternatives {
		alter_obj := evalElifExpression(al, env)
		if alter_obj != NULL { //任何中间选项，最先正确执行的，返回它的执行结果
			return alter_obj
		}
	}

	//执行最后选项
	if ie.LastAlternative != nil {
		return Eval(ie.LastAlternative, env)
	}

	//首选项错误，有或没有中间选项，没有else的最后结果
	return NULL
}

// elif的处理 和前面差不多
func evalElifExpression(ef *ast.ElIfExpression, env *object.Environment) object.Object {
	alternnative_condition := Eval(ef.Condition, env)
	if isError(alternnative_condition) {
		return alternnative_condition
	}

	if isTruthy(alternnative_condition) {
		return Eval(ef.Consequence, env)
	} else {
		return NULL
	}
}

// 辅助函数 判断条件是否成立
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// block求值的辅助函数
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// 生成报错信息 生成*object.Error
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// 辅助函数，用与验证obj是不是错误
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false

}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	//取环境绑定的值
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	//如果没有绑定的值，则在内置函数环境中查找
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found:" + node.Value)
}

// 辅助函数：求值所有实参
func evalExpression(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// 执行函数调用
func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function: //自定义的函数
		extendedEnv := extendFunction(fn, args) //创建局部环境
		evaluated := Eval(fn.Body, extendedEnv) //执行函数，也就是配合环境执行函数体的内容
		return unwarpReturnValue(evaluated)     //解包，返回函数执行后的结果

	case *object.Builtin: //内置函数
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

// 辅助函数：创建函数局部环境，绑定形参与实参
func extendFunction(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	//遍历所有形参，并绑定到实参
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// 辅助函数：解包返回值
func unwarpReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

// 字符串拼接
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown opearator:%s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}
