package utils

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"testing"
)

func TestDataSet(t *testing.T) {

	a := NewDataSet("test-a")

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	fmt.Println(a.Set())

	fmt.Println("remove X")
	a.Remove("X")
	fmt.Println(a.Set())

	fmt.Println("add X")
	a.Add("X")
	fmt.Println(a.Set())

	b := NewDataSet("test-b")
	for val := range a.Iter() {
		b.Add(val)
	}

	fmt.Println(b.Set())

	if !a.Equal(b.Set()) {
		t.Error("The sets are not equal after iterating (Iter) through the first set")
	}

	fmt.Println("union set a, b and WANG FEI")

	fmt.Println(a.Set().Union(b.Set()).Union(mapset.NewSet("Y", "X", "WANG", "FEI")))

	fmt.Println(a.Set().Contains("c", "a", "d", "b"))

}
