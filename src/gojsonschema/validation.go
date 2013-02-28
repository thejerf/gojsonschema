// @author       sigu-399
// @description  An implementation of JSON Schema, draft v4 - Go language
// @created      28-02-2013

package gojsonschema

import (
	"fmt"
	"reflect"
	"strconv"
)

type ValidationResult struct {
	valid         bool
	errorMessages []string
}

func (v *ValidationResult) IsValid() bool {
	return v.valid
}

func (v *ValidationResult) AddErrorMessage(message string) {
	v.errorMessages = append(v.errorMessages, message)
	v.valid = false
}

func (v *JsonSchemaDocument) Validate(document interface{}) ValidationResult {

	result := ValidationResult{valid: true}
	v.validateRecursive(v.rootSchema, document, &result)
	return result
}

func (v *JsonSchemaDocument) validateRecursive(currentSchema *JsonSchema, currentNode interface{}, result *ValidationResult) {

	fmt.Printf("Validation of schema %s\n", currentSchema.property)

	schProperty := currentSchema.property
	schTypes := currentSchema.types

	if currentNode == nil {
		if !schTypes.HasType(TYPE_NULL) {
			result.AddErrorMessage(fmt.Sprintf("%s must be of type %s", schProperty, schTypes.String()))
			return
		}
	} else {

		rValue := reflect.ValueOf(currentNode)
		rKind := rValue.Kind()

		var nextNode interface{}
		var ok bool

		fmt.Printf("Type %s\n", rKind.String())

		switch rKind {

		case reflect.Map:

			if !schTypes.HasType(TYPE_OBJECT) {
				result.AddErrorMessage(fmt.Sprintf("%s must be of type %s", schProperty, schTypes.String()))
				return
			}

			for _, pSchema := range currentSchema.propertiesChildren {
				castCurrentNode := currentNode.(map[string]interface{})
				nextNode, ok = castCurrentNode[pSchema.property]
				if !ok {
					result.AddErrorMessage(fmt.Sprintf("%s is required", pSchema.property))
					return
				}
				v.validateRecursive(pSchema, nextNode, result)
			}

		case reflect.String:

			if !schTypes.HasType(TYPE_STRING) {
				result.AddErrorMessage(fmt.Sprintf("%s must be of type %s", schProperty, schTypes.String()))
				return
			}

		case reflect.Float64:

			isInteger := true
			_, err := strconv.Atoi(fmt.Sprintf("%v", currentNode))
			if err != nil {
				isInteger = false
			}

			formatIsCorrect := schTypes.HasType(TYPE_NUMBER) || (isInteger && schTypes.HasType(TYPE_INTEGER))

			if !formatIsCorrect {
				result.AddErrorMessage(fmt.Sprintf("%s must be of type %s", schProperty, schTypes.String()))
				return
			}
		}
	}
}