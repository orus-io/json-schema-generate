package test

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	oneof "github.com/orus-io/json-schema-generate/test/oneof_gen"
	"github.com/stretchr/testify/assert"
)

func TestOneOf(t *testing.T) {
	var d oneof.ComplexdataType
	assert.True(t, d.IsNotSet())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString("null", &d)) &&
		assert.True(t, d.IsNil())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString("true", &d)) &&
		assert.True(t, d.IsBool()) &&
		assert.True(t, d.Bool())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString("false", &d)) &&
		assert.True(t, d.IsBool()) &&
		assert.False(t, d.Bool())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString("\"16\"", &d)) &&
		assert.True(t, d.IsString()) &&
		assert.Equal(t, "16", d.String())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString("16", &d)) &&
		assert.True(t, d.IsInt()) &&
		assert.Equal(t, 16, d.Int())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString("16.2", &d)) &&
		assert.True(t, d.IsFloat64()) &&
		assert.Equal(t, 16.2, d.Float64())
	_ = assert.NoError(t, jsoniter.UnmarshalFromString(`{"age": 12}`, &d)) &&
		assert.True(t, d.IsNotSoAnonymous()) &&
		assert.Equal(t, 12, d.NotSoAnonymous().Age)
}
