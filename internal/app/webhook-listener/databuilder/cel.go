package databuilder

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/ext"
)

// evaluateCELExpression creates a CEL environment and evaluates the given expression
// It returns an error if the expression cannot be compiled, the program cannot be created, the value of the expression
// is nil, the value is not a string, or the value is not a boolean
// It returns the value of the expression if it is successful
// The celExpression is the CEL expression to be evaluated
// The jsonData is the JSON data to be used in the evaluation. It is a map of string keys and interface values from the webhook payload
func evaluateCELExpression(celExpresion string, jsonData map[string]interface{}) (string, error) {
	// Create the CEL environment
	env, err := cel.NewEnv(
		// Extensible functions and types
		ext.Strings(), ext.Encoders(), ext.Math(), ext.Sets(), ext.Lists(),
		// Declaration of variable 'data'
		cel.Declarations(
			decls.NewVar("data", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)
	if err != nil {
		return "", err
	}

	// Compile the CEL expression
	ast, issues := env.Compile(celExpresion)
	if issues != nil && issues.Err() != nil {
		return "", issues.Err()
	}

	// Create the CEL program
	program, err := env.Program(ast)
	if err != nil {
		return "", err
	}

	// Extract data from the jsonData
	out, _, err := program.Eval(map[string]interface{}{
		"data": jsonData,
	})
	if err != nil {
		return "", err
	}

	value := out.Value()
	// Check if the value is nil
	if value == nil {
		return "", errors.New(fmt.Sprintf("expression '%s' did not return a value", celExpresion))
	}
	// Check if the value is a boolean and convert it to a string if it is
	if _, ok := value.(bool); ok {
		value = strconv.FormatBool(value.(bool))
	}
	// Check if the value is a string
	if _, ok := value.(string); !ok {
		return "", errors.New(fmt.Sprintf("expression '%s' did not return a string value", celExpresion))
	}

	return value.(string), nil
}
