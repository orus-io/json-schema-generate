package test

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	emptytypes "github.com/orus-io/json-schema-generate/test/emptytypes_gen"
	"github.com/stretchr/testify/assert"
)

func TestEmptyTypes(t *testing.T) {
	for _, tt := range []struct {
		p emptytypes.Product
		j string
	}{
		{emptytypes.Product{Id: 1, Name: "n1"}, `{"id": 1, "name": "n1"}`},
		{emptytypes.Product{Id: 2, Name: "n2", Description: emptytypes.NewEmptyString("hi")}, `{"id": 2, "name": "n2", "description": "hi"}`},
	} {
		if s, err := jsoniter.MarshalToString(&tt.p); assert.NoError(t, err) {
			assert.JSONEq(t, tt.j, s)
		}
		var p emptytypes.Product
		if err := jsoniter.UnmarshalFromString(tt.j, &p); assert.NoError(t, err) {
			assert.Equal(t, tt.p, p)
		}
	}
}
