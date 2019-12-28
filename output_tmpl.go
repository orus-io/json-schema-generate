package generate

import "text/template"

var headerTmpl = template.Must(template.New("schema-generate").Parse(
	`// Code generated by schema-generate. DO NOT EDIT.

package {{.PackageName}}

import (
{{- range $path, $importName := .ImportPaths }}
	{{ $importName }}"{{ $path }}"
{{- end }}
)

`))

var mainTmpl = template.Must(template.New("schema-generate").Parse(
	`{{- with $top := . -}}

func ValueTypeToString(valueType jsoniter.ValueType) string {
	switch valueType {
	case jsoniter.StringValue:
		return "string"
	case jsoniter.NumberValue:
		return "number"
	case jsoniter.NilValue:
		return "nil"
	case jsoniter.BoolValue:
		return "bool"
	case jsoniter.ArrayValue:
		return "array"
	case jsoniter.ObjectValue:
		return "object"
	default:
		return "invalid"
	}
}

type commaTracker struct {
	stream *jsoniter.Stream
	started bool
}

func (t *commaTracker) More() {
	if t.started {
		t.stream.WriteMore()
	} else {
		t.started = true
	}
}

type isEmptyChecker interface {
	IsEmpty() bool
}

// IsEmpty reports whether v is zero struct
// Does not support cycle pointers for performance, so as json
func IsEmpty(v interface{}) bool {
	if i, ok := v.(isEmptyChecker); ok {
		return i.IsEmpty()
	}
	rv := {{ .Pkg "reflect" }}.ValueOf(v)
	return !rv.IsValid() || rv.IsZero()
}

var (
	jsonNullValue = []byte("null")
)

// OneOfStringNull is a 'string' or a 'null', and can be emptied
type OneOfStringNull struct {
	currentType jsoniter.ValueType
	stringValue string
}

// NewOneOfStringNull creates a empty OneOfStringNull
func NewOneOfStringNull() OneOfStringNull {
	return OneOfStringNull{jsoniter.InvalidValue, ""}
}

// NewOneOfStringNullString creates a OneOfStringNull of type string
func NewOneOfStringNullString(value string) OneOfStringNull {
	return OneOfStringNull{jsoniter.StringValue, value}
}

// NewOneOfStringNullNull creates a OneOfStringNull of type null
func NewOneOfStringNullNull() OneOfStringNull {
	return OneOfStringNull{jsoniter.NilValue, ""}
}

// IsEmpty returns true if the value is empty
func (value *OneOfStringNull) IsEmpty() bool {
	return value.currentType == jsoniter.InvalidValue
}

// IsNull returns true if the value is 'null'
func (value *OneOfStringNull) IsNull() bool {
	return value.currentType == jsoniter.NilValue
}

// IsString returns true if the value is a string
func (value *OneOfStringNull) IsString() bool {
	return value.currentType == jsoniter.StringValue
}

// StringValue returns the current value if IsString() is true, "" otherwise
func (value *OneOfStringNull) StringValue() string {
	if value.currentType == jsoniter.StringValue {
		return value.stringValue
	}
	return ""
}

// NullString returns the current value as a sql.NullString
func (value *OneOfStringNull) NullString() {{ .Pkg "database/sql" }}.NullString {
	return sql.NullString{
		Valid:  value.currentType == jsoniter.StringValue,
		String: value.stringValue,
	}
}

// MarshalJSONStream serializes to a jsoniter Stream
func (value OneOfStringNull) MarshalJSONStream(stream *{{ .Pkg "jsoniter" "github.com/json-iterator/go" }}.Stream) {
	if value.currentType == jsoniter.StringValue {
		stream.WriteString(value.stringValue)
	} else {
		stream.WriteNil()
	}
}

// MarshalJSON serialize to json
func (value OneOfStringNull) MarshalJSON() ([]byte, error) {
	switch value.currentType {
	case jsoniter.InvalidValue:
		return jsonNullValue, nil
	case jsoniter.NilValue:
		return jsonNullValue, nil
	case jsoniter.StringValue:
		return jsoniter.Marshal(value.stringValue)
	}
	return nil, {{ $top.Pkg "fmt" }}.Errorf(
		"OneOfStringNull unsupported type: %s",
		ValueTypeToString(value.currentType))
}

// UnmarshalJSON unserialize a OneOfStringNull from json
func (value *OneOfStringNull) UnmarshalJSON(data []byte) error {
	if {{ .Pkg "bytes" }}.Equal(data, jsonNullValue) {
		value.currentType = jsoniter.NilValue
	} else {
		if err := jsoniter.Unmarshal(data, &value.stringValue); err != nil {
			return err
		}
		value.currentType = jsoniter.StringValue
	}
	return nil
}

// OneOfNumberNull is a 'string' or a 'null', and can be emptied
type OneOfNumberNull struct {
	currentType jsoniter.ValueType
	numberValue float64
}

// NewOneOfNumberNull creates a empty OneOfNumberNull
func NewOneOfNumberNull() OneOfNumberNull {
	return OneOfNumberNull{jsoniter.InvalidValue, 0}
}

// NewOneOfNumberNullNumber creates a OneOfNumberNull of type number
func NewOneOfNumberNullNumber(value float64) OneOfNumberNull {
	return OneOfNumberNull{jsoniter.NumberValue, value}
}

// NewOneOfNumberNullNull creates a OneOfNumberNull of type null
func NewOneOfNumberNullNull() OneOfNumberNull {
	return OneOfNumberNull{jsoniter.NilValue, 0}
}

// IsEmpty returns true if the value is empty
func (value *OneOfNumberNull) IsEmpty() bool {
	return value.currentType == jsoniter.InvalidValue
}

// IsNull returns true if the value is 'null'
func (value *OneOfNumberNull) IsNull() bool {
	return value.currentType == jsoniter.NilValue
}

// IsNumber returns true if the value is a number
func (value *OneOfNumberNull) IsNumber() bool {
	return value.currentType == jsoniter.NumberValue
}

// NumberValue returns the current value if IsNumber() is true, 0 otherwise
func (value *OneOfNumberNull) NumberValue() float64 {
	if value.currentType == jsoniter.NumberValue {
		return value.numberValue
	}
	return 0
}

// MarshalJSON serialize to json
func (value OneOfNumberNull) MarshalJSON() ([]byte, error) {
	switch value.currentType {
	case jsoniter.InvalidValue:
		return jsonNullValue, nil
	case jsoniter.NilValue:
		return jsonNullValue, nil
	case jsoniter.NumberValue:
		return jsoniter.Marshal(value.numberValue)
	}
	return nil, fmt.Errorf(
		"OneOfNumberNull unsupported type: %s",
		ValueTypeToString(value.currentType))
}

// UnmarshalJSON unserialize a OneOfNumberNull from json
func (value *OneOfNumberNull) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNullValue) {
		value.currentType = jsoniter.NilValue
	} else {
		if err := jsoniter.Unmarshal(data, &value.numberValue); err != nil {
			return err
		}
		value.currentType = jsoniter.NumberValue
	}
	return nil
}

// OneOfBoolNull is a 'bool' or a 'null', and can be emptied
type OneOfBoolNull struct {
	currentType jsoniter.ValueType
	boolValue   bool
}

// NewOneOfBoolNull creates a empty OneOfBoolNull
func NewOneOfBoolNull() OneOfBoolNull {
	return OneOfBoolNull{jsoniter.InvalidValue, false}
}

// NewOneOfBoolNullBool creates a OneOfBoolNull of type number
func NewOneOfBoolNullBool(value bool) OneOfBoolNull {
	return OneOfBoolNull{jsoniter.BoolValue, value}
}

// NewOneOfBoolNullNull creates a OneOfBoolNull of type null
func NewOneOfBoolNullNull() OneOfBoolNull {
	return OneOfBoolNull{jsoniter.NilValue, false}
}

// IsEmpty returns true if the value is empty
func (value *OneOfBoolNull) IsEmpty() bool {
	return value.currentType == jsoniter.InvalidValue
}

// IsNull returns true if the value is 'null'
func (value *OneOfBoolNull) IsNull() bool {
	return value.currentType == jsoniter.NilValue
}

// IsBool returns true if the value is a bool
func (value *OneOfBoolNull) IsBool() bool {
	return value.currentType == jsoniter.BoolValue
}

// BoolValue returns the current value if IsBool() is true, false otherwise
func (value *OneOfBoolNull) BoolValue() bool {
	if value.currentType == jsoniter.BoolValue {
		return value.boolValue
	}
	return false
}

// MarshalJSON serialize to json
func (value OneOfBoolNull) MarshalJSON() ([]byte, error) {
	switch value.currentType {
	case jsoniter.InvalidValue:
		return jsonNullValue, nil
	case jsoniter.NilValue:
		return jsonNullValue, nil
	case jsoniter.BoolValue:
		return jsoniter.Marshal(value.boolValue)
	}
	return nil, fmt.Errorf(
		"OneOfBoolNull unsupported type: %s",
		ValueTypeToString(value.currentType))
}

// UnmarshalJSON unserialize a OneOfBoolNull from json
func (value *OneOfBoolNull) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNullValue) {
		value.currentType = jsoniter.NilValue
	} else {
		if err := jsoniter.Unmarshal(data, &value.boolValue); err != nil {
			return err
		}
		value.currentType = jsoniter.BoolValue
	}
	return nil
}

{{- range .Aliases }}

// {{ .Name }} ...
type {{ .Name }} {{ .Type }}
{{- end -}}

{{- range .Structs }}

// {{ $top.Comment .Name .Description }}
type {{ .Name }} struct {
	{{- range .Fields }}
	// {{ $top.Comment .Name .Description }}
	{{ .Name }} {{ .Type }} {{ $top.Backquote }}json:"{{ .JSONName }}{{ if not .Required }},omitempty{{ end }}"{{ $top.Backquote }}
	{{ end -}}
}
{{- end -}}

{{- range $struct := .Structs }}
{{ if .GenerateCode }}

// MarshalJSON serializes to JSON
func (s *{{ .Name }}) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	stream := jsoniter.ConfigDefault.BorrowStream(buf)
	s.MarshalJSONStream(stream)
	stream.Flush()
	err := stream.Error
	jsoniter.ConfigDefault.ReturnStream(stream)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s {{ .Name }}) MarshalJSONStream(stream *jsoniter.Stream) {
	{{- if .NoProp }}
	stream.WriteEmptyObject()
	{{- else }}
	stream.WriteObjectStart()
	ct := commaTracker{stream:stream}

	{{- range .Fields }}
	{{- if ne .JSONName "-" }}

	// Marshal the {{ .Name }} field
	{{- if and .Required .IsPointer }}

	// {{ .Name }} is required
	if s.{{ .Name }} == nil {
		stream.Error = {{ $top.Pkg "errors" }}.New("{{ .Name }} ({{ .JSONName }}) is a required")
		return
	}
	{{- end }}
	{{- if not .Required }}
	if !IsEmpty(s.{{ .Name }}) {
	{{- end }}
	ct.More()
	stream.WriteObjectField("{{ .JSONName }}")
	{{- if eq .Type "string" }}
	stream.WriteString(s.{{ .Name }})
	{{- else if eq .Type "OneOfStringNull" }}
	s.{{ .Name }}.MarshalJSONStream(stream)
	{{- else }}
	stream.WriteVal(s.{{ .Name }})
	if stream.Error != nil {
		return
	}
	{{- end }}
	{{- if not .Required }}
	}
	{{- end}}

	{{- end}}
	{{- end}}
	{{- if and .AdditionalType (ne .AdditionalType "false")}}
	for key, value := range s.AdditionalProperties {
		ct.More()
		stream.WriteObjectField(key)
		stream.WriteVal(value)
	}
	{{- end}}
	stream.WriteObjectEnd()
	{{- end}}
}

func (s *{{ .Name }}) UnmarshalJSON(data []byte) error {
	iter := jsoniter.ConfigDefault.BorrowIterator(data)
	s.UnmarshalJSONIterator(iter)
	err := iter.Error
	jsoniter.ConfigDefault.ReturnIterator(iter)
	return err
}

func (s *{{ .Name }}) UnmarshalJSONIterator(iter *jsoniter.Iterator) {
	{{- range .Fields }}
	{{- if .Required}}
	{{ .Name }}Received := false
	{{- end}}
	{{- end}}

	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		switch field {
		{{- range .Fields }}
		{{- if ne .JSONName "-" }}
		case "{{ .JSONName }}":
			{{- if and $top.AlwaysAcceptFalse (ne .Type "bool") (ne .Type "OneOfBoolNull")}}
			if iter.WhatIsNext() == jsoniter.BoolValue {
				if iter.ReadBool() {
					iter.ReportError("reading field {{ .JSONName }}", "{{ .JSONName }} is 'true', but the expected type is {{ .Type }}")
					return
				}
				// received 'false', which we accept and ignore for now
			}
			{{- end}}
			{{- if eq .Type "string" }}
			s.{{ .Name }} = iter.ReadString()
			{{- else if eq .Type "bool" }}
			s.{{ .Name }} = iter.ReadBool()
			{{ else }}
			iter.ReadVal(&s.{{ .Name }})
			{{- end}}
			if iter.Error != nil {
				return
			}
			{{- if .Required}}
			{{ .Name }}Received = true
			{{- end}}
		{{- end}}
		{{- end}}
		default:
			{{- if eq .AdditionalType "false" }}
			iter.ReportError("reading {{ .Name }}", "additional property not allowed: \"" + field + "\"")
			return
			{{- else if .AdditionalType }}
            if s.AdditionalProperties == nil {
                s.AdditionalProperties = make(map[string]{{ .AdditionalType }}, 0)
            }
            var additionalValue {{ .AdditionalType }}
			iter.ReadVal(&additionalValue)
			if iter.Error != nil {
				return
			}
            s.AdditionalProperties[field]= additionalValue
			{{- else }}
			// Ignore the additional property
			iter.Skip()
			{{- end }}
		}
	}

	{{- range .Fields }}
	{{- if .Required}}

	if !{{ .Name }}Received {
		iter.ReportError("validating {{ $struct.Name }}", "\"{{ .JSONName }}\" is required but was not present")
	}
	{{- end}}
	{{- end}}
}

{{- end -}}
{{- end -}}
{{- end -}}
`))
