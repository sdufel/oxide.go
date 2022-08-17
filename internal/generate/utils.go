package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

func openGeneratedFile(filename string) (*os.File, error) {
	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %v", err)
	}

	p := filepath.Join(cwd, filename)

	// Create the generated files.
	// Open the file for writing.
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating %q: %v", p, err)
	}

	// Add the header to the package.
	fmt.Fprintf(f, "// Code generated by `%s`. DO NOT EDIT.\n\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(f, "package oxide")
	fmt.Fprintln(f, "")

	return f, nil
}

func isLocalEnum(v *openapi3.SchemaRef) bool {
	return v.Ref == "" && v.Value.Type == "string" && len(v.Value.Enum) > 0
}

func isLocalObject(v *openapi3.SchemaRef) bool {
	return v.Ref == "" && v.Value.Type == "object" && len(v.Value.Properties) > 0
}

// formatStringType converts a string schema to a valid Go type.
func formatStringType(t *openapi3.Schema) string {
	if t.Format == "date-time" {
		return "*time.Time"
	} else if t.Format == "date" {
		return "*time.Time"
	} else if t.Format == "time" {
		return "*time.Time"
	} else if t.Format == "email" {
		return "string"
	} else if t.Format == "hostname" {
		return "string"
	} else if t.Format == "ipv4" {
		return "string"
	} else if t.Format == "ipv6" {
		return "string"
	} else if t.Format == "uri" {
		return "string"
	} else if t.Format == "uuid" {
		return "string"
	} else if t.Format == "uuid3" {
		return "string"
	}

	return "string"
}

// toLowerFirstLetter returns the given string with the first letter converted to lower case.
func toLowerFirstLetter(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// makeSingular returns the given string but singular.
func makeSingular(s string) string {
	if strings.HasSuffix(s, "Status") {
		return s
	}
	return strings.TrimSuffix(s, "s")
}

// makePlural returns the given string but plural.
func makePlural(s string) string {
	singular := makeSingular(s)
	if strings.HasSuffix(singular, "s") {
		return singular + "es"
	}

	return singular + "s"
}

func trimStringFromSpace(s string) string {
	if idx := strings.Index(s, " "); idx != -1 {
		return s[:idx]
	}
	return s
}

func containsMatchFirstWord(s []string, str string) bool {
	for _, v := range s {
		if trimStringFromSpace(v) == trimStringFromSpace(str) {
			return true
		}
	}

	return false
}

func isPageParam(s string) bool {
	return s == "nextPage" || s == "pageToken" || s == "limit"
}

// printProperty converts an object's property name to a valid Go identifier.
func printProperty(p string) string {
	c := strcase.ToCamel(p)
	return c
}

func printPropertyLower(p string) string {
	s := strcase.ToLowerCamel(printProperty(p))
	return s
}

// printType converts a schema type to a valid Go type.
func printType(property string, r *openapi3.SchemaRef) string {
	s := r.Value
	t := s.Type

	// If we have a reference, just use that.
	if r.Ref != "" {
		ref := getReferenceSchema(r)
		// Just use the type of the reference.
		if ref == "Name" {
			return "string"
		}

		return ref
	}

	// See if we have an allOf.
	if s.AllOf != nil {
		if len(s.AllOf) > 1 {
			fmt.Printf("[WARN] TODO: allOf for %q has more than 1 item\n", property)
			return "TODO"
		}

		return printType(property, s.AllOf[0])
	}

	if t == "string" {
		reference := getReferenceSchema(r)
		if reference != "" {
			return reference
		}

		return formatStringType(s)
	} else if t == "integer" {
		return "int"
	} else if t == "number" {
		return "float64"
	} else if t == "boolean" {
		return "bool"
	} else if t == "array" {
		reference := getReferenceSchema(s.Items)
		if reference != "" {
			return fmt.Sprintf("[]%s", reference)
		}

		// TODO: handle if it is not a reference.
		return "[]string"
	} else if t == "object" {
		// Most likely this is a local object, we will handle it.
		return strcase.ToCamel(property)
	}

	fmt.Printf("[WARN] TODO: skipping type %q for %q, marking as interface{}\n", t, property)
	return "interface{}"
}

func compareFiles(expected, actual string) error {
	f1, err := ioutil.ReadFile(expected)
	if err != nil {
		return err
	}

	f2, err := ioutil.ReadFile(actual)
	if err != nil {
		return err
	}

	if !bytes.Equal(f1, f2) {
		return fmt.Errorf("%v is not equal to %v", expected, actual)
	}
	return nil
}