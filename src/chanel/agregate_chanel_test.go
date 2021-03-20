package chanel

import (
	"strings"
	"testing"
)

func TestSimpleAgregate(t *testing.T){
	ag := NewAgregateChanel(2)
	ag.Add("test")
	if val := ag.Get(); !strings.EqualFold(val,"test") {
		t.Error("Must find value test")
	}
}

func TestDuplicateAgregate(t *testing.T){
	ag := NewAgregateChanel(2)
	ag.Add("test")
	ag.Add("test")
	ag.Add("different")
	ag.Add("super")
	if val := ag.Get(); !strings.EqualFold(val,"test") {
		t.Error("Must find value test")
	}
	if val := ag.Get(); !strings.EqualFold(val,"different") {
		t.Error("Must find value different (remove duplicate)",val)
	}
}

func TestDuplicateAgregateWithTimeout(t *testing.T){
	ag := NewAgregateChanel(2)
	ag.Add("test")
	ag.Add("test")
	if val := ag.Get(); !strings.EqualFold(val,"test") {
		t.Error("Must find value test")
	}
}

