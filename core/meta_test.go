package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	BoolKey   = MetaDataKey("bool")
	StringKey = MetaDataKey("string")
	IntKey    = MetaDataKey("int")
	FloatKey  = MetaDataKey("float64")
)

func TestCommon(t *testing.T) {
	meta := NewMetaData()
	assert.True(t, meta.SetBool(BoolKey, true))
	_, ok := meta.GetBool("unknown")
	assert.False(t, ok)
	assert.True(t, meta.Del(BoolKey))
	assert.False(t, meta.Del(BoolKey))
	_, nok := meta.GetBool(BoolKey)
	assert.False(t, nok)
}

func TestBoolValues(t *testing.T) {
	meta := NewMetaData()
	assert.True(t, meta.SetBool(BoolKey, true))
	assert.False(t, meta.SetBool(BoolKey, true))
	val, ok := meta.GetBool(BoolKey)
	assert.True(t, val)
	assert.True(t, ok)
	meta.SetInt(BoolKey, 1)
	_, nok := meta.GetBool(BoolKey)
	assert.False(t, nok)
	_, nok = meta.GetBool("unknown")
	assert.False(t, nok)
}

func TestStringValues(t *testing.T) {
	meta := NewMetaData()
	assert.True(t, meta.SetString(StringKey, "aaaaa"))
	assert.False(t, meta.SetString(StringKey, "zzzzz"))
	val, ok := meta.GetString(StringKey)
	assert.Equal(t, val, "zzzzz")
	assert.True(t, ok)
	meta.SetBool(StringKey, true)
	_, nok := meta.GetString(StringKey)
	assert.False(t, nok)
	_, nok = meta.GetString("unknown")
	assert.False(t, nok)
}

func TestIntValues(t *testing.T) {
	meta := NewMetaData()
	assert.True(t, meta.SetInt(IntKey, 1000000))
	assert.False(t, meta.SetInt(IntKey, 2000000))
	val, ok := meta.GetInt(IntKey)
	assert.Equal(t, val, 2000000)
	assert.True(t, ok)
	meta.SetBool(IntKey, true)
	_, nok := meta.GetInt(IntKey)
	assert.False(t, nok)
	_, nok = meta.GetInt("unknown")
	assert.False(t, nok)
}

func TestFloatValues(t *testing.T) {
	meta := NewMetaData()
	assert.True(t, meta.SetFloat(FloatKey, 1.000))
	assert.False(t, meta.SetFloat(FloatKey, -1.000))
	val, ok := meta.GetFloat(FloatKey)
	assert.Equal(t, val, -1.000)
	assert.True(t, ok)
	meta.SetBool(FloatKey, true)
	_, nok := meta.GetFloat(FloatKey)
	assert.False(t, nok)
	_, nok = meta.GetFloat("unknown")
	assert.False(t, nok)
}
