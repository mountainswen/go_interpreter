package evaluator

var builtins = map[string]*Builtin{
	"len":&Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1{
				return newError("wrong number of arguments.got =%d," +
					"want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *StringObject:
				return &Integer{Value:int64(len(arg.Value))}
			case *Array:
				return &Integer{Value:int64(len(arg.Element))}
			default:
				return newError("the type not support " +
					"len func")
			}

			return NULL
		},
	},
	"print":{
		Fn: func(args ...Object) Object {

			if len(args) != 1{
				return newError("wrong number of arguments.got =%d," +
					"want=1", len(args))
			}

			return &StringObject{Value:args[0].Inspect()}
		},
	},
}