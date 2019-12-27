package generate

import (
	"bytes"
	"io"
	"sort"
	"strings"
)

func getOrderedFieldNames(m map[string]Field) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func getOrderedStructNames(m map[string]Struct) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

// OutputData contains all the data necessary for the template
type OutputData struct {
	ImportPaths map[string]string

	PackageName string
	Structs     []Struct
	Aliases     []Field
	Backquote   string
}

// Pkg ...
func (d *OutputData) Pkg(name string, path ...string) string {
	var realPath string
	var importName string
	if len(path) == 0 {
		if i := strings.LastIndex(name, "/"); i != -1 {
			realPath = name
			name = name[i+1:]
		} else {
			realPath = name
		}
	} else {
		realPath = path[0]
		importName = name + " "
	}
	d.ImportPaths[realPath] = importName
	return name
}

// Comment outputs a string the '// ' in front of each line
func (OutputData) Comment(s ...string) string {
	lines := strings.Split(strings.Join(s, " "), "\n")
	return strings.Join(lines, "\n// ")
}

// NoProp returns true if the struct has no property
func (s Struct) NoProp() bool {
	return len(s.Fields) == 0 && (s.AdditionalType == "" || s.AdditionalType == "false")
}

// IsPointer returns true if the type is a pointer
func (f Field) IsPointer() bool {
	return strings.HasPrefix(f.Type, "*")
}

// Output generates code and writes to w.
func Output(w io.Writer, g *Generator, pkg string) {
	structs := g.Structs
	aliases := g.Aliases

	data := OutputData{
		ImportPaths: make(map[string]string),

		PackageName: cleanPackageName(pkg),
		Backquote:   "`",
	}

	for _, k := range getOrderedStructNames(structs) {
		data.Structs = append(data.Structs, structs[k])
	}
	for _, k := range getOrderedFieldNames(aliases) {
		data.Aliases = append(data.Aliases, aliases[k])
	}

	codeBuf := new(bytes.Buffer)

	if err := mainTmpl.Execute(codeBuf, &data); err != nil {
		panic(err)
	}

	if err := headerTmpl.Execute(w, &data); err != nil {
		panic(err)
	}

	w.Write(codeBuf.Bytes())
}

func cleanPackageName(pkg string) string {
	pkg = strings.Replace(pkg, ".", "", -1)
	pkg = strings.Replace(pkg, "_", "", -1)
	pkg = strings.Replace(pkg, "-", "", -1)
	return pkg
}
