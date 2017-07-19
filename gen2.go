package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type UpdateTypeEnum uint8

const (
	None UpdateTypeEnum = iota
	Partial
	Full
)

var UpdateTypeEnumOrdered = []UpdateTypeEnum{
	None,
	Partial,
	Full,
}

const (
	v1     ComputeApiGenVersion = "v1"
	v0beta                      = "v0beta"
)

type ComputeApiGenVersion string

var OrderedComputeApiVersions = []ComputeApiGenVersion{
	v0beta,
	v1,
}

type UpdateType struct {
	LowestVersion ComputeApiGenVersion
	Type          UpdateTypeEnum
}

type UpdateInfo struct {
	Versions []VersionInfo
}

type VersionInfo struct {
	ClientName  string
	ServiceName string
}

type GenResource struct {
	Type           string
	LowestVersion  ComputeApiGenVersion
	HighestVersion ComputeApiGenVersion
	UpdateTypes    []UpdateType
	ExtraParams    []string
}

type GenResource2 struct {
	Type                 string
	PluralType           string
	Versions             []VersionInfo
	UpdateInfo           UpdateInfo
	ExtraParamsSignature string
	ExtraParamsCall      string
}

func main() {
	t := template.Must(template.New("temp").Parse(multiversionServiceTemplate))
	data := &GenResource{
		Type:          "Address",
		LowestVersion: v1,
		ExtraParams:   []string{"region"},
	}

	realData := GenResource2{
		Type: data.Type,
	}

	if strings.HasSuffix(realData.Type, "s") {
		realData.PluralType = data.Type + "es"
	} else {
		realData.PluralType = data.Type + "s"
	}

	for _, version := range OrderedComputeApiVersions {
		clientName := ""
		if version == v0beta {
			clientName = "Beta"
		}

		realData.Versions = append(realData.Versions, VersionInfo{ClientName: clientName, ServiceName: fmt.Sprintf("%v", version)})
		if version == data.LowestVersion {
			break
		}
	}

	if len(data.UpdateTypes) > 0 {
		var updateVersion ComputeApiGenVersion
		match := false
		for _, v := range data.UpdateTypes {
			if v.Type == Full {
				updateVersion = v.LowestVersion
				match = true
			}
		}
		if match == true {
			for _, version := range OrderedComputeApiVersions {
				clientName := ""
				if version == v0beta {
					clientName = "Beta"
				}

				realData.UpdateInfo.Versions = append(realData.UpdateInfo.Versions, VersionInfo{ClientName: clientName, ServiceName: fmt.Sprintf("%v", version)})
				if version == updateVersion {
					break
				}
			}
		}
	}

	if len(data.ExtraParams) > 0 {
		eCall := ""

		for _, param := range data.ExtraParams {
			eCall += fmt.Sprintf(" %s,", param)
		}

		realData.ExtraParamsCall = eCall
		realData.ExtraParamsSignature = eCall[:len(eCall)-1] + " string,"
	}

	t.Execute(os.Stdout, realData)
}

const multiversionServiceTemplate = `
package google

import (
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (s *ComputeMultiversionService) Insert{{ $.Type }}(project string,{{ $.ExtraParamsSignature }} resource *computeBeta.{{ $.Type }}, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	{{ range .Versions }}
		case {{ .ServiceName  }}:
		{{ .ServiceName }}Resource := &compute{{ .ClientName }}.{{ $.Type }}{}
		err := Convert(resource, {{ .ServiceName }}Resource)
		if err != nil {
			return nil, err
		}

		{{ .ServiceName }}Op, err := s.{{ .ServiceName }}.{{ $.PluralType }}.Insert(project,{{ $.ExtraParamsCall }} {{ .ServiceName }}Resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert({{ .ServiceName }}Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) Get{{ $.Type }}(project string,{{ $.ExtraParamsSignature }} resource string, version ComputeApiVersion) (*computeBeta.{{ $.Type }}, error) {
	res := &computeBeta.{{ .Type }}{}
	switch version {
	{{ range .Versions }}
		case {{ .ServiceName }}:
		r, err := s.{{ .ServiceName }}.{{ $.PluralType }}.Get(project,{{ $.ExtraParamsCall }} resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(r, res)
		if err != nil {
			return nil, err
		}

		return res, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) Delete{{ $.Type }}(project string,{{ $.ExtraParamsSignature }} resource string, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	{{ range .Versions }}
		case {{ .ServiceName }}:
		{{ .ServiceName }}Op, err := s.{{ .ServiceName }}.{{ $.PluralType }}.Delete(project,{{ $.ExtraParamsCall }} resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert({{ .ServiceName }}Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}

{{ if ne (len .UpdateInfo.Versions) 0 }}
func (s *ComputeMultiversionService) Update{{ $.Type }}(project string,{{ $.ExtraParamsSignature }} resourceName string, resource *computeBeta.{{ $.Type }}, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	{{ range .UpdateInfo.Versions }}
		case {{ .ServiceName }}:
		{{ .ServiceName }}Resource := &compute{{ .ClientName }}.{{ $.Type }}{}
		err := Convert(resource, {{ .ServiceName }}Resource)
		if err != nil {
			return nil, err
		}
		{{ .ServiceName }}Op, err := s.{{ .ServiceName }}.{{ $.PluralType }}.Update(project,{{ $.ExtraParamsCall }} resourceName, {{ .ServiceName }}Resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert({{ .ServiceName }}Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}
{{ end }}
`
