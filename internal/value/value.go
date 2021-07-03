package value

type Kind int

const (
	V_INTEGER Kind = iota
	V_BOOLEAN
	V_FUNCTION
	V_EMPTY
	V_CONS
	V_SYMBOL
	V_STRING
	V_NIL
	V_REFERENCE
	V_ARRAY
	V_DICT
)

type Value interface {
	Kind() Kind
	Display() string
	Str() string
	GetInt() int
	GetString() string
	GetHead() Value
	GetTail() Value
	GetValue() Value
	SetValue(Value)
	GetArray() []Value
	GetMap() map[string]Value
	Apply([]Value, interface{}) (Value, error)
	// internal methods
	IsTrue() bool
	IsEqual(Value) bool
	//isEq() bool    -- don't think we need pointer equality for now - = is enough?
}

func IsEqual(v1 Value, v2 Value) bool {
	return v1.IsEqual(v2)
}

func IsTrue(v Value) bool {
	return v.IsTrue()
}

func IsAtom(v Value) bool {
	k := v.Kind()
	return k == V_INTEGER || k == V_BOOLEAN ||
		k == V_SYMBOL || k == V_STRING
}

func IsSymbol(v Value) bool {
	return v.Kind() == V_SYMBOL
}

func IsCons(v Value) bool {
	return v.Kind() == V_CONS
}

func IsEmpty(v Value) bool {
	return v.Kind() == V_EMPTY
}

func IsNumber(v Value) bool {
	return v.Kind() == V_INTEGER
}

func IsBool(v Value) bool {
	return v.Kind() == V_BOOLEAN
}

func IsRef(v Value) bool {
	return v.Kind() == V_REFERENCE
}

func IsString(v Value) bool {
	return v.Kind() == V_STRING
}

func IsFunction(v Value) bool {
	k := v.Kind()
	return k == V_FUNCTION
}

func IsNil(v Value) bool {
	return v.Kind() == V_NIL
}

func IsArray(v Value) bool {
	return v.Kind() == V_ARRAY
}

func IsDict(v Value) bool {
	return v.Kind() == V_DICT
}

func Classify(v Value) string {
	switch v.Kind() {
	case V_INTEGER:
		return "int"
	case V_BOOLEAN:
		return "bool"
	case V_EMPTY:
		return "list"
	case V_CONS:
		return "list"
	case V_SYMBOL:
		return "symbol"
	case V_FUNCTION:
		return "fun"
	case V_STRING:
		return "string"
	case V_NIL:
		return "nil"
	case V_REFERENCE:
		return "reference"
	case V_ARRAY:
		return "array"
	case V_DICT:
		return "dict"
	default:
		return "?"
	}
}
