package evaluator

import (
	"ast"
	"fmt"
	)

func Eval(node ast.Node,env *Environment) Object {
	switch node := node.(type) {
	case *ast.HashLiteral:
		return evalHashLiteral(node,env)

	case *ast.IndexExpression:
		left := Eval(node.Left,env)
		index := Eval(node.Index,env)
		return evalIndexExpression(left,index)

	case *ast.ArrayLiteral:

		array := &Array{}
		o := evalExpression(node.Element,env)

		array.Element = o

		return array

	case *ast.CallExpression:
		function := Eval(node.Function,env) //get function object
		args := evalExpression(node.Arguments,env)

		return applyFunction(function,args)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &Function{Parameter:params,
		Body:body,
		Env:env}

	case *ast.Indetifier:
		return evalIdentifier(node,env)

	case *ast.LetStatement:
		val := Eval(node.Value,env)
		env.Set(node.Name.TokenLiteral(),val)

	case *ast.Program:
		return evalStatements(node.Statements,env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression,env)

	case *ast.IntergerLiteral:
		return &Integer{Value:node.Value}

	case *ast.StringLiteral:
		return &StringObject{Value:node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right,env)
		op := node.Operator
		return evalprefixExpression(op,right,env)

	case *ast.InfixExpression:
		left := Eval(node.Left,env)
		right := Eval(node.Right,env)
		op := node.Operator

		return evalInfixExpression(op,left,right)

	case *ast.IfExpression:
		return evalIfExpression(node,env)

	case *ast.BlockStatement:
		return evalStatements(node.Statements,env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue,env)
		return &ReturnValue{Value:val}
	}

	return nil
}

func evalHashLiteral(node *ast.HashLiteral,env *Environment)Object{
	pairs := make(map[HashKey]HashPair)

	for knode,vnode := range node.Pairs{
		key := Eval(knode,env)

		hashKey,ok := key.(Hashable)
		if !ok{
			return newError("invalid hash key")
		}

		value := Eval(vnode,env)

		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key:key,Value:value}
	}

	return &Hash{Pairs:pairs}
}

func evalIndexExpression(left,index Object)Object{
	switch  {
	case left.Type() == ARRAY_OBJ &&
		index.Type() == INTEGER_OBJ:
			return evalArrayIndexExpression(left,index)

	case left.Type() == HASH_OBJ:
		return evalHashIndexExpression(left,index)
	default:
		return newError("index operator not supported")
	}
}

func evalHashIndexExpression(left Object,index Object)Object{
	hashObj := left.(*Hash)

	key,ok := index.(Hashable)
	if !ok{
		return newError("invalid hash key")
	}

	pair,ok := hashObj.Pairs[key.HashKey()]
	if !ok{
		return NULL
	}

	return pair.Value
}
func evalArrayIndexExpression(array,index Object)Object{
	arrayObj := array.(*Array)
	idx := index.(*Integer).Value

	max := int64(len(arrayObj.Element)) - 1

	if idx < 0 || idx >max{
		return newError("index out of range")
	}

	return arrayObj.Element[idx]
}


func evalStatements(stmts []ast.Statement,env *Environment)Object{
	var result Object

	for _,statement := range stmts{
		result = Eval(statement,env)
		if returnVal,ok := result.(*ReturnValue);ok{
			return returnVal.Value
		}

		if result,ok := result.(*Error);ok{
			return result
		}

	}

	return result
}

func nativeBoolToBooleanObj(input bool)Object{
	if input{
		return TRUE
	}

	return FALSE
}

func evalprefixExpression(op string,right Object,env *Environment)Object{
	switch op {
	case "!":
		return evalBangOperatorPrefix(right)
	case "-":
		return evalminuxOperatorPrefix(right)
	default:
		return newError("unkown operator:%s%s",op,right.Type())
	}
}

func evalBangOperatorPrefix(right Object)Object{
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	}

	return NULL
}

func evalminuxOperatorPrefix(right Object)Object{
	if right.Type() != INTEGER_OBJ{
		return newError("unkown operator:%s",right.Type())
	}

	value := right.(*Integer).Value
	return &Integer{Value:-value}
}

func evalInfixExpression(operator string,
	left,right Object)Object{
	switch  {
	case left.Type() == INTEGER_OBJ &&
		right.Type() == INTEGER_OBJ:
			return evalIntegerInfixExpression(operator,
				left,right)
	case left.Type() == BOOLEAN_OBJ &&
		right.Type() == BOOLEAN_OBJ:
		return evalBoolInfixExpression(operator,left,right)
	case left.Type() == STRING_OBJ &&
		right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator,left,right)
	case left.Type() != right.Type():
		return newError("type missmatch:%s%s%s:",left.Type(),operator,right.Type())
	default:
		return NULL
	}
}

func evalStringInfixExpression(operator string,
	left,right Object)Object{
		leftVal := left.(*StringObject).Value
		rightVal := right.(*StringObject).Value

	switch operator {
	case "+":
		return &StringObject{Value:leftVal+rightVal}
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string,
	left,right Object)Object{
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value:leftVal + rightVal}
	case "-":
		return &Integer{Value:leftVal - rightVal}
	case "*":
		return &Integer{Value:leftVal * rightVal}
	case "/":
		return &Integer{Value:leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObj(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObj(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObj(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalBoolInfixExpression(op string,
	left,right Object)Object{
		leftVal := left.(*Boolean).Value
		rightVal := right.(*Boolean).Value

	switch op {
	case "!=":
		return nativeBoolToBooleanObj(leftVal != rightVal)
	case "==":
		return nativeBoolToBooleanObj(leftVal == rightVal)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression,env *Environment)Object{
	condition := Eval(ie.Condition,env)
	if isTurthy(condition){
		return Eval(ie.Consequence,env)
	}else if ie.Alternative != nil{
		return Eval(ie.Alternative,env)
	}else {
		return NULL
	}
}

func isTurthy(obj Object)bool{
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		if _,ok := obj.(*Integer);ok{
			if obj.(*Integer).Value == 0{
				return false
			}
		}

		return true
	}
}

func newError(format string,a ...interface{})Object{
	return &Error{Message:fmt.Sprintf(format,a...)}
}

func evalIdentifier(
	node *ast.Indetifier,
	env *Environment)Object{
		val,ok := env.Get(node.TokenLiteral())

		if ok{
			return val
		}

		val,ok = builtins[node.TokenLiteral()]
		if ok{
			return val
		}

		if !ok{
			object := newError("identifier not found:%s", node.TokenLiteral())
			return object
		}

		return val
}

func applyFunction(fn Object,args []Object)Object{
	function, ok := fn.(*Function)
	if ok{
		extendedEnv := extendFunctionEnv(function,args)
		evaluated := Eval(function.Body,extendedEnv)
		return unwarapReturnValue(evaluated)
	}

	//看一下是不是builtin function
	function1,ok := fn.(*Builtin)
	if ok{
		return function1.Fn(args...)
	}

	return newError("not a function:%s",fn.Type())

}

func extendFunctionEnv(fn *Function,
	args []Object)*Environment{

	env := NewEnclosedEnvironment(fn.Env)

	for i,param := range fn.Parameter{ //将
		env.Set(param.TokenLiteral(),args[i])
	}

	return env
}

func unwarapReturnValue(obj Object)Object{
	if returnValue,ok := obj.(*ReturnValue);ok{
		return returnValue.Value
	}

	return obj
}

func evalExpression(exps []ast.Expression,env *Environment)[]Object{
	var result []Object

	for _, e := range exps{
		evaluated := Eval(e,env)

		result = append(result,evaluated)
	}

	return result
}