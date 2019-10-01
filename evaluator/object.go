package evaluator

import (
	"fmt"
	"ast"
	"bytes"
	"strings"
	"hash/fnv"
)

const (
	HASH_OBJ = "HASH"
	ARRAY_OBJ = "ARRAY"
	BUILTIN_OBJ = "BUILTIN"
	STRING_OBJ = "STRING"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE_OBJ"

	ERROR_OBJ = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
)

type Hashable interface {
	HashKey()HashKey
}

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

//整型数
type Integer struct {
	Value int64
}

func (i *Integer)Inspect()string{
	return fmt.Sprintf("%d",i.Value)
}
func (i *Integer)Type()ObjectType{
	return INTEGER_OBJ
}
func (i *Integer)HashKey()HashKey{
	return HashKey{Type:i.Type(),Value:uint64(i.Value)}
}
//Boolean
type Boolean struct {
	Value bool
}
func (b *Boolean)Inspect()string{
	return fmt.Sprintf("%t",b.Value)
}
func (b *Boolean)Type()ObjectType{
	return BOOLEAN_OBJ
}
func (b *Boolean)HashKey()HashKey{
	var v uint64
	if b.Value{
		v = 1
	}else {
		v = 0
	}

	return HashKey{Type:b.Type(),Value:v}
}

var (
	TRUE = &Boolean{Value:true}
	FALSE = &Boolean{Value:false}
)

//NULL
type Null struct {

}
func (n *Null)Type()ObjectType{
	return NULL_OBJ
}
func (n *Null)Inspect()string{
	return "null"
}

var (
	NULL = &Null{}
)

//return object
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue)Type()ObjectType{
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue)Inspect()string{
	return rv.Value.Inspect()
}

//error obj
type Error struct {
	Message string
}
func (e *Error)Type()ObjectType{
	return ERROR_OBJ
}
func (e *Error)Inspect()string{
	return "ERROR:"+e.Message
}

//环境
type Environment struct {
	outer *Environment
	store map[string]Object
}
func NewEnvironment()*Environment{
	s := make(map[string]Object)
	return &Environment{store:s,outer:nil}
}

func (e *Environment)Get(name string)(Object,bool){
	obj,ok := e.store[name]
	if !ok && e.outer != nil{
		obj,ok = e.outer.Get(name)
	}

	return obj,ok
}

func (e *Environment)Set(name string,obj Object)Object{
	e.store[name] = obj
	return obj
}

func NewEnclosedEnvironment(outer *Environment)*Environment{
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Function struct {
	Parameter []*ast.Indetifier
	Body *ast.BlockStatement
	Env *Environment
}
func (f *Function)Type()ObjectType{
	return FUNCTION_OBJ
}
func (f *Function)Inspect()string{
	var out bytes.Buffer

	params := []string{}
	for _,p := range f.Parameter{
		params = append(params,p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params,","))
	out.WriteString("){\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}


type StringObject struct {
	Value string
}
func (s StringObject)Type()ObjectType{
	return STRING_OBJ
}
func (s *StringObject)Inspect()string{
	return s.Value
}

func (s *StringObject)HashKey()HashKey{
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type:s.Type(),Value:h.Sum64()}
}

//
type BuiltInFunction func(args ...Object)Object
type Builtin struct {
	Fn BuiltInFunction
}

func (b *Builtin)Type()ObjectType{
	return BUILTIN_OBJ
}
func (b *Builtin)Inspect()string{
	return "builtin function"
}

//array object
type Array struct {
	Element []Object
}
func (a *Array)Type()ObjectType{
	return ARRAY_OBJ
}
func (a *Array)Inspect()string{
	var out bytes.Buffer

	ele := []string{}
	for _, e := range a.Element{
		ele = append(ele,e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(ele,","))
	out.WriteString("]")

	return out.String()
}

//hash
type HashKey struct {
	Type ObjectType
	Value uint64
}

type HashPair struct {
	Key Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash)Type()ObjectType{
	return HASH_OBJ
}
func (h *Hash)Inspect()string{
	var out bytes.Buffer

	s := []string{}
	for _,v := range h.Pairs{
		s = append(s,v.Key.Inspect() + ":" + v.Value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(s,","))
	out.WriteString("}")

	return out.String()
}

