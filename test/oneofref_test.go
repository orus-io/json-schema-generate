package test

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	oneofref "github.com/orus-io/json-schema-generate/test/oneofref_gen"
	"github.com/stretchr/testify/assert"
)

func TestOneOfRef(t *testing.T) {
	var d oneofref.DataType
	assert.True(t, d.IsNotSet())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString(`{"name": "foo"}`, &d)) &&
		assert.True(t, d.IsThing()) &&
		assert.Equal(t, "foo", d.Thing().Name)
}
