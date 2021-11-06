package types

import (
	"math"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var UInt = graphql.NewScalar(graphql.ScalarConfig{
	Name: "UInt",
	Description: "The `UInt` scalar type represents non-fractional unsigned whole numeric " +
		"values. UInt can represent values between 0 and 2^31 - 1. ",
	Serialize:  serializeUInt,
	ParseValue: parseUInt,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if uIntValue, err := strconv.ParseUint(valueAST.Value, 10, 32); err == nil {
				return uIntValue
			}
		}
		return nil
	},
})

func parseUInt(value interface{}) interface{} {
	switch value := value.(type) {
	case float64:
		if value < 0 || value > math.MaxInt32 {
			return nil
		}
		return uint64(value)
	case *float64:
		if value == nil {
			return nil
		}
		return parseUInt(*value)
	default:
		return nil
	}
}

func serializeUInt(value interface{}) interface{} {
	switch value := value.(type) {
	case uint64:
		return value
	case *uint64:
		if value == nil {
			return nil
		}
		return serializeUInt(*value)
	default:
		return nil
	}
}
