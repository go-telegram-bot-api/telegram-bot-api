package tgbotapi

import (
	"testing"
)

func assertLen(t *testing.T, params Params, l int) {
	actual := len(params)
	if actual != l {
		t.Fatalf("Incorrect number of params, expected %d but found %d\n", l, actual)
	}
}

func assertEq(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("Values did not match, a: %v, b: %v\n", a, b)
	}
}

func TestAddNonEmpty(t *testing.T) {
	params := make(Params)
	params.AddNonEmpty("value", "value")
	assertLen(t, params, 1)
	assertEq(t, params["value"], "value")
	params.AddNonEmpty("test", "")
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddNonZero(t *testing.T) {
	params := make(Params)
	params.AddNonZero("value", 1)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "1")
	params.AddNonZero("test", 0)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddNonZero64(t *testing.T) {
	params := make(Params)
	params.AddNonZero64("value", 1)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "1")
	params.AddNonZero64("test", 0)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddBool(t *testing.T) {
	params := make(Params)
	params.AddBool("value", true)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "true")
	params.AddBool("test", false)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddNonZeroFloat(t *testing.T) {
	params := make(Params)
	params.AddNonZeroFloat("value", 1)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "1.000000")
	params.AddNonZeroFloat("test", 0)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddInterface(t *testing.T) {
	params := make(Params)
	data := struct {
		Name string `json:"name"`
	}{
		Name: "test",
	}
	params.AddInterface("value", data)
	assertLen(t, params, 1)
	assertEq(t, params["value"], `{"name":"test"}`)
	params.AddInterface("test", nil)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddFirstValid(t *testing.T) {
	params := make(Params)
	params.AddFirstValid("value", 0, "", "test")
	assertLen(t, params, 1)
	assertEq(t, params["value"], "test")
	params.AddFirstValid("value2", 3, "test")
	assertLen(t, params, 2)
	assertEq(t, params["value2"], "3")
}
