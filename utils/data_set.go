package utils

import (
	"errors"
	"github.com/deckarep/golang-set"
)

/*
进行数据集的操作处理
*/
type DataSet struct {
	name string
	set  mapset.Set
}

func NewDataSet(name string) *DataSet {
	return &DataSet{
		name: name,
		set:  mapset.NewSet(),
	}
}

func (d *DataSet) Set() mapset.Set {
	return d.set
}

func (d *DataSet) Add(element interface{}) error {
	ret := d.set.Add(element)
	if !ret {
		return errors.New("DataSet Add error.")
	}
	return nil
}

func (d *DataSet) Remove(element interface{}) error {
	d.set.Remove(element)
	return nil
}


// 集合元素数量
func (d *DataSet) Num() int {
	return d.set.Cardinality()
}

func (d *DataSet) Clear() {
	d.set.Clear()
}

func (d *DataSet) Iter() <-chan interface{} {
	return d.set.Iter()
}

func (d *DataSet) Equal(other mapset.Set) bool {
	return d.set.Equal(other)
}
