package tgbotapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNonEmpty(t *testing.T) {
	params := make(Params)
	params.AddNonEmpty("value", "value")
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], "value")
	params.AddNonEmpty("test", "")
	assert.Len(t, params, 1)
	assert.Equal(t, params["test"], "")
}

func TestAddNonZero(t *testing.T) {
	params := make(Params)
	params.AddNonZero("value", 1)
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], "1")
	params.AddNonZero("test", 0)
	assert.Len(t, params, 1)
	assert.Equal(t, params["test"], "")
}

func TestAddNonZero64(t *testing.T) {
	params := make(Params)
	params.AddNonZero64("value", 1)
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], "1")
	params.AddNonZero64("test", 0)
	assert.Len(t, params, 1)
	assert.Equal(t, params["test"], "")
}

func TestAddBool(t *testing.T) {
	params := make(Params)
	params.AddBool("value", true)
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], "true")
	params.AddBool("test", false)
	assert.Len(t, params, 1)
	assert.Equal(t, params["test"], "")
}

func TestAddNonZeroFloat(t *testing.T) {
	params := make(Params)
	params.AddNonZeroFloat("value", 1)
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], "1.000000")
	params.AddNonZeroFloat("test", 0)
	assert.Len(t, params, 1)
	assert.Equal(t, params["test"], "")
}

func TestAddInterface(t *testing.T) {
	params := make(Params)
	data := struct {
		Name string `json:"name"`
	}{
		Name: "test",
	}
	params.AddInterface("value", data)
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], `{"name":"test"}`)
	params.AddInterface("test", nil)
	assert.Len(t, params, 1)
	assert.Equal(t, params["test"], "")

	params = make(Params)
	var test *string = nil
	params.AddInterface("test", test)
	assert.Len(t, params, 0)

	badData := struct {
		Value chan interface{}
	}{
		Value: make(chan interface{}),
	}
	err := params.AddInterface("value", badData)
	assert.Error(t, err)
}

func TestAddFirstValid(t *testing.T) {
	params := make(Params)
	params.AddFirstValid("value", 0, "", "test")
	assert.Len(t, params, 1)
	assert.Equal(t, params["value"], "test")
	params.AddFirstValid("value2", 3, "test")
	assert.Len(t, params, 2)
	assert.Equal(t, params["value2"], "3")

	params.AddFirstValid("value")
	assert.Equal(t, params["value"], "test")

	params = make(Params)
	err := params.AddFirstValid("value", struct {
		Value string
	}{
		Value: "test",
	})
	assert.NoError(t, err)
	err = params.AddFirstValid("value", make(chan interface{}))
	assert.Error(t, err)
	assert.Equal(t, `{"Value":"test"}`, params["value"])
}
