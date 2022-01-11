package ceph_test

import (
	"github.com/chrisamti/ceph-rest-client/ceph"
	"testing"
)

func TestPathJoin(t *testing.T) {
	var a, c string
	var b *string

	a = "alpha"
	c = "gamma"

	expectedResult := "alpha/gamma"

	joinResult := ceph.PathJoin(a, b, c)

	t.Logf("joined strings: %s", joinResult)
	if expectedResult != joinResult {
		t.Errorf("expected '%s' - got '%s'", expectedResult, joinResult)
	}

	// test with ptr string set
	expectedResult = "alpha/beta/gamma"
	b = new(string)
	*b = "beta"

	joinResult = ceph.PathJoin(a, b, c)

	t.Logf("joined strings: %s", joinResult)
	if expectedResult != joinResult {
		t.Errorf("expected '%s' - got '%s'", expectedResult, joinResult)
	}
}

func TestStaticCounter(t *testing.T) {
	a := ceph.StaticCounter()

	t.Log(a())
	t.Log(a())

	staticValue := a()
	if staticValue != 3 {
		t.Errorf("expected static counter to be 3 - got %d", staticValue)
	}
}
