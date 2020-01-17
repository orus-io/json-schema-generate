package generate

import (
	"strings"
	"text/template"
)

var streamMarshallerTypes = []string{
	"EmptyBool",
	"EmptyNumber",
	"EmptyString",
	"OneOfBoolNull",
	"OneOfNumberNull",
	"OneOfStringNull",
}

var iteratorUnmashallerTypes = []string{
	"EmptyBool",
	"EmptyNumber",
	"EmptyString",
	"OneOfBoolNull",
	"OneOfNumberNull",
	"OneOfStringNull",
}

var funcs = template.FuncMap{
	// comment outputs a string the '// ' in front of each line
	"comment": func(s ...string) string {
		lines := strings.Split(strings.Join(s, " "), "\n")
		return strings.Join(lines, "\n// ")
	},
	// fieldName creates a field name from a data type
	"fieldName": func(s string) string {
		if strings.HasPrefix(s, "*") {
			s = s[1:]
		}
		return strings.ToLower(s[0:1]) + s[1:]
	},
	// capitalize returns the string with the first letter uppered
	"capitalize": func(s string) string {
		return strings.ToUpper(s[0:1]) + s[1:]
	},
	// DeferedType removes a leading '*' from a data type
	"deferedType": func(s string) string {
		if strings.HasPrefix(s, "*") {
			s = s[1:]
		}
		return s
	},
	// isplainjsontype returns true if the given json type is a pod
	"isplainjsontype": func(t string) bool {
		switch t {
		case "bool":
			return true
		case "string":
			return true
		case "number":
			return true
		case "integer":
			return true
		case "null":
			return true
		}
		return false
	},
	// oneOfContainsJsonType returns true if the given oneOf contains the passed jsontype
	"oneOfContainsJsonType": func(o OneOf, t string) bool {
		for _, ot := range o.Types {
			if ot.JSONType == t {
				return true
			}
		}
		return false
	},
	// ispointer returns true if the given type starts with "*"
	"ispointer": func(t string) bool {
		return t[0] == '*'
	},
	// isStreamMarshaller returns true if the given type is known to have a MarshalJSONStream function
	"isStreamMarshaller": func(t string) bool {
		for _, s := range streamMarshallerTypes {
			if s == t {
				return true
			}
		}
		return false
	},
	// isIteratorUnmarshaller returns true if the given type is known to have a UnmarshalJSONIterator function
	"isIteratorUnmarshaller": func(t string) bool {
		for _, s := range iteratorUnmashallerTypes {
			if s == t {
				return true
			}
		}
		return false
	},
}

var headerTmpl = template.Must(template.New("schema-generate").Funcs(funcs).Parse(
	`// Code generated by schema-generate. DO NOT EDIT.

package {{.PackageName}}

import (
{{- range $path, $importName := .ImportPaths }}
	{{ $importName }}"{{ $path }}"
{{- end }}
)

`))

var mainTmpl = template.Must(template.New("schema-generate").Funcs(funcs).Parse(
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

// NewEmptyString creates a non-empty EmptyString
func NewEmptyString(s string) EmptyString {
	return EmptyString{s, true}
}

// EmptyString is string or nothing
type EmptyString struct {
	String string
	Valid bool // Valid is true if String is not empty
}

func (s EmptyString) IsEmpty() bool {
	return !s.Valid
}

func (s EmptyString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return jsoniter.Marshal(s.String)
	}
	return []byte("\"\""), nil
}

func (s *EmptyString) Set(value string) {
	s.String = value
	s.Valid = true
}

func (s *EmptyString) Unset() {
	s.String = ""
	s.Valid = false
}

func (s EmptyString) MarshalJSONStream(stream *jsoniter.Stream) {
	if s.Valid {
		stream.WriteString(s.String)
	} else {
		stream.WriteString("")
	}
}

func (s *EmptyString) UnmarshalJSONIterator(iter *jsoniter.Iterator) {
	s.String = iter.ReadString()
	s.Valid = iter.Error == nil
}

func (s *EmptyString) UnmarshalJSON(data []byte) error {
	if err := jsoniter.Unmarshal(data, &s.String); err != nil {
		return err
	}
	s.Valid = true
	return nil
}

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

{{- range $oneOf := .OneOfs }}

type {{ .Name }}Enum = int

const (
	{{ $oneOf.Name }}EnumNotSet {{ $oneOf.Name }}Enum = iota
	{{- range $i, $type := .Types }}
	{{ $oneOf.Name }}Enum{{ .ShortType }}
	{{- end }}
)

type {{ .Name }} struct {
	Type {{ $oneOf.Name }}Enum

	value interface{}
}

func (o {{ $oneOf.Name }}) IsNotSet() bool {
	return o.Type == {{ $oneOf.Name }}EnumNotSet
}

{{- range .Types }}

func (o {{ $oneOf.Name }}) Is{{ .ShortType }}() bool {
	return o.Type == {{ $oneOf.Name }}Enum{{ .ShortType }}
}

{{- if ne "nil" .Type}}
func (o {{ $oneOf.Name }}) {{ .ShortType }}() {{ .Type }} {
	return o.value.({{ .Type }})
}
{{- end }}

func (o *{{ $oneOf.Name }}) Set{{ .ShortType }}(
	{{- if ne "nil" .Type}}v {{ .Type }}{{ end -}}
) {
	{{- if ne "nil" .Type }}
	o.value = v
	{{- end }}
	o.Type = {{ $oneOf.Name }}Enum{{ .ShortType }}
}
{{- end }}

func (o {{ $oneOf.Name }}) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	stream := jsoniter.ConfigDefault.BorrowStream(buf)
	o.MarshalJSONStream(stream)
	err := stream.Error
	jsoniter.ConfigDefault.ReturnStream(stream)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (o {{ $oneOf.Name }}) MarshalJSONStream(stream *jsoniter.Stream) {
	switch o.Type {
	{{- range .Types }}
	case {{ $oneOf.Name }}Enum{{ .ShortType }}:
		{{- if eq "bool" .Type }}
		stream.WriteBool(o.value.(bool))
		{{- else if eq "string" .Type }}
		stream.WriteString(o.value.(string))
		{{- else if eq "int" .Type }}
		stream.WriteInt(o.value.(int))
		{{- else if eq "float64" .Type }}
		stream.WriteFloat64(o.value.(float64))
		{{- else if eq "nil" .Type }}
		stream.WriteNil()
		{{- else }}
		stream.WriteVal(o.value)
		{{- end }}
	{{- end }}
	}
}

func (o *{{ $oneOf.Name }}) UnmarshalJSON(data []byte) error {
	iter := jsoniter.ConfigDefault.BorrowIterator(data)
	o.UnmarshalJSONIterator(iter)
	err := iter.Error
	jsoniter.ConfigDefault.ReturnIterator(iter)
	return err
}

func (o *{{ $oneOf.Name }}) UnmarshalJSONIterator(iter *jsoniter.Iterator) {
	switch t := iter.WhatIsNext(); t {

	{{- range .Types }}
	{{- if eq "nil" .Type }}
	case jsoniter.NilValue:
		iter.ReadNil()
		o.SetNil()
	{{- else if eq "bool" .Type }}
	case jsoniter.BoolValue:
		o.SetBool(iter.ReadBool())
	{{- else if eq "string" .Type }}
	case jsoniter.StringValue:
		o.SetString(iter.ReadString())
	{{- end }}
	{{- end }}

	{{- if and (oneOfContainsJsonType . "integer") (oneOfContainsJsonType . "number") }}
	case jsoniter.NumberValue:
		b := iter.SkipAndReturnBytes()
		if iter.Error == io.EOF {
			iter.Error = nil
		}

		subiter := jsoniter.ConfigDefault.BorrowIterator(b)
		defer jsoniter.ConfigDefault.ReturnIterator(subiter)

		i := subiter.ReadInt()
		if subiter.Error == nil || subiter.Error == {{ $top.Pkg "io" }}.EOF {
			o.SetInt(i)
			return
		}
		subiter.Error = nil

		subiter.ResetBytes(b)
		f := subiter.ReadFloat64()
		if subiter.Error == nil || subiter.Error == io.EOF {
			o.SetFloat64(f)
			return
		}

		iter.Error = subiter.Error
		
	{{- else if oneOfContainsJsonType . "integer" }}
	case jsoniter.NumberValue:
		o.SetInt(iter.ReadFloat64())
	{{- else if oneOfContainsJsonType . "number" }}
	case jsoniter.NumberValue:
		o.SetFloat64(iter.ReadFloat64())
	{{- end }}
	{{- if oneOfContainsJsonType . "object" }}
	case jsoniter.ObjectValue:
		any := iter.ReadAny()

		{{- range .Types }}
		{{- if eq "object" .JSONType }}

		{ // attempt to read a {{ .Type }}
			var value {{ deferedType .Type }}
			any.ToVal(&value)
			if any.LastError() == nil {
				o.Set{{ .ShortType }}({{ if ispointer .Type }}&{{ end }}value)
				return
			}
		}
		{{- end }}
		{{- end }}

		iter.Error = any.LastError()
	{{- end }}
	}
}

{{- end }}

{{- range .Aliases }}

// {{ .Name }} ...
type {{ .Name }} {{ .Type }}
{{- end -}}

{{- range .Structs }}

// {{ comment .Name .Description }}
type {{ .Name }} struct {
	{{- range .Fields }}
	// {{ comment .Name .Description }}
	{{ .Name }} {{ .Type }} {{ $top.Backquote }}json:"{{ .JSONName }}{{ if not .Required }},omitempty{{ end }}"{{ $top.Backquote }}
{{ end }}
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
	{{- else if isStreamMarshaller .Type }}
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
			{{- else if isIteratorUnmarshaller .Type }}
			s.{{ .Name }}.UnmarshalJSONIterator(iter)
			{{- else }}
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
