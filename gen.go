package main

import (
	"fmt"
	"google.golang.org/api/compute/v1"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type Field struct {
	Name       string
	PrettyType string
	Type       string
	Complex    bool
	Array      bool
}

func main() {
	r := reflect.TypeOf(compute.InstanceGroupsSetNamedPortsRequest{})
	typeName := r.Name()
	fields := make([]*Field, 0, r.NumField())
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		newField := &Field{
			Name:       field.Name,
			PrettyType: field.Type.String(),
			Type:       extractType(field.Type),
		}

		switch field.Type.Kind() {
		case reflect.Slice:
			elem := field.Type.Elem()
			newField.Array = true
			if elem.Kind() == reflect.Ptr {
				newField.Complex = true
				newField.PrettyType = getPrettyType(elem)
			}

		case reflect.Ptr:
			newField.Complex = true
			newField.PrettyType = getPrettyType(field.Type)
		}

		fields = append(fields, newField)
	}

	t := template.Must(template.New("temp").Parse(temp))
	data := struct {
		Type             string
		SharedFields     []*Field
		ProductionFields []*Field
	}{
		Type:             typeName,
		SharedFields:     fields,
		ProductionFields: fields,
	}
	t.Execute(os.Stdout, data)

}
func getPrettyType(k reflect.Type) string {
	chunks := strings.Split(k.String(), ".")
	return chunks[len(chunks)-1]
}

func extractType(t reflect.Type) string {
	// Complex types will always be pointers
	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Slice {
		return t.String()
	}

	if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Ptr {
		return t.String()
	}

	var isPtr bool
	var isSlice bool
	if t.Kind() == reflect.Slice {
		isSlice = true
		isPtr = t.Elem().Kind() == reflect.Ptr
	} else if t.Kind() == reflect.Ptr {
		isPtr = true
	}

	chunks := strings.Split(t.String(), ".")
	finalType := chunks[len(chunks)-1]
	if isPtr {
		finalType = fmt.Sprintf("*%s", finalType)
	}

	if isSlice {
		finalType = fmt.Sprintf("[]%s", finalType)
	}

	return finalType
}

const temp = `
package shared

import (
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

type {{ .Type }} struct {
  {{- range .SharedFields}}
  	{{ .Name }} {{ .Type }}
  {{- end }}
}

func (s *{{ .Type }}) ToProduction() *compute.{{ .Type }} {
  if s == nil {
    return nil
  }

  return &compute.{{ .Type }}{
  {{- range .ProductionFields -}}
  {{ if .Complex -}}
    {{ if .Array }}
        {{ .Name }}: {{ .PrettyType }}ArrayToProduction(s.{{ .Name }}),
    {{- else }}
        {{ .Name }}: s.{{ .Name }}.ToProduction(),
    {{- end }}
  {{- else }}
  	{{ .Name }}: s.{{ .Name }},
  {{- end }}
  {{- end }}
  }
}

func {{ .Type }}FromProduction(s *compute.{{ .Type }}) *{{ .Type }} {
  if s == nil {
    return nil
  }

  return &{{ .Type }}{
  {{- range .ProductionFields -}}
  {{ if .Complex -}}
    {{ if .Array }}
        {{ .Name }}: {{ .PrettyType }}ArrayFromProduction(s.{{ .Name }}),
    {{- else }}
        {{ .Name }}: {{ .PrettyType }}FromProduction(s.{{ .Name }}),
    {{- end }}
  {{- else }}
  	{{ .Name }}: s.{{ .Name }},
  {{- end }}
  {{- end }}
  }
}

func {{ .Type }}ArrayToProduction(items []*{{ .Type }}) []*compute.{{ .Type }} {
	arr := make([]*compute.{{ .Type }}, 0, len(items))
	for _, v := range items {
		arr = append(arr, v.ToProduction())
	}
	return arr
}

func {{ .Type }}ArrayFromProduction(items []*compute.{{ .Type }}) []*{{ .Type }} {
	arr := make([]*{{ .Type }}, 0, len(items))
	for _, v := range items {
		arr = append(arr, {{ .Type }}FromProduction(v))
	}
	return arr
}

`
