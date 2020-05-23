package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// tag name
const (
	constructorTag = "constructor"
	getterTag      = "getter"
	setterTag      = "setter"
)

// field elements required to create a constructor
type field struct {
	FieldName         string
	FieldType         string
	ConstructorIgnore bool
	GetterIgnore      bool
	SetterIgnore      bool
}

// searchFiled extract the elements necessary to create a constructor from ast.StructType.
func searchFiled(structType *ast.StructType) ([]field, error) {
	var fields []field

	for _, list := range structType.Fields.List {
		f := field{}

		if list.Tag == nil {
			f.FieldName = list.Names[0].Name
			f.FieldType = types.ExprString(list.Type)
			f.ConstructorIgnore = false
			f.GetterIgnore = false
			f.SetterIgnore = false
		} else {
			value := strings.Trim(list.Tag.Value, "`")

			// if the string "goor" is not included
			if !strings.Contains(value, cliName) {
				continue
			}

			// by getting the position where the tag key name "goor" starts and adding +5 to it, the
			// then you can get "constructor:-" of "goor: "constructor:-".
			tag := value[strings.Index(value, cliName)+5:]
			tag = strings.Trim(tag, "\"")

			var tags []string
			if strings.Contains(tag, ";") {
				tags = strings.Split(tag, ";")
			} else {
				tags = []string{tag}
			}

			for _, t := range tags {
				switch {
				case strings.Contains(t, constructorTag):
					if t[len(constructorTag)+1:] == "-" {
						f.ConstructorIgnore = true
					}
				case strings.Contains(t, getterTag):
					if t[len(getterTag)+1:] == "-" {
						f.GetterIgnore = true
					}
				case strings.Contains(t, setterTag):
					if t[len(setterTag)+1:] == "-" {
						f.SetterIgnore = true
					}
				default:
					// don't do anything.
				}
			}

			f.FieldName = list.Names[0].Name
			f.FieldType = types.ExprString(list.Type)
		}
		fields = append(fields, f)
	}

	if len(fields) == 0 {
		return fields, fmt.Errorf("fields is not found")
	}

	return fields, nil
}
