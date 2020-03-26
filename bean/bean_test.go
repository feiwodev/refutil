package bean

import (
	"encoding/json"
	"fmt"
	"testing"
)

// ------------------------------------------------------
// Created by fei wo at 2020/3/26
// ------------------------------------------------------
// Copyright©2020-2030
// ------------------------------------------------------
// blog: http://www.feiwo.xyz
// ------------------------------------------------------
// email: zhuyongluck@qq.com
// ------------------------------------------------------
//  bean testing
// ------------------------------------------------------

type Person struct {
	Name string
	Age int
	Address struct{
		City string
		Area string
	}
}

type PersonDto struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Address struct{
		City string `json:"city"`
		Area string `json:"area"`
	} `json:"address"`
}


func TestCopyProperties(t *testing.T) {
	dto := PersonDto{
		Name: "feiwo",
		Age:  12,
		Address: struct {
			City string `json:"city"`
			Area string `json:"area"`
		}{City:"changsha",Area:"yuelu"},
	}

	p := new(Person)
	err := CopyProperties(&dto, p)
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func BenchmarkCopyProperties(b *testing.B) {
	dto := PersonDto{
		Name: "feiwo",
		Age:  12,
		Address: struct {
			City string `json:"city"`
			Area string `json:"area"`
		}{City:"changsha",Area:"yuelu"},
	}

	p := new(Person)
	b.StartTimer()
	for i := 0 ; i < b.N ; i++  {
		err := CopyProperties(&dto, p)
		if err != nil {
			fmt.Print(err)
		}
	}
	b.StopTimer()
	fmt.Println(p)
}

func TestCopyPropertiesIgnoreDefaultVal(t *testing.T) {

	dto := PersonDto{
		Name: "feiwo",
		Age:  12,
		Address: struct {
			City string `json:"city"`
			Area string `json:"area"`
		}{Area:"yuelu"},
	}

	p := new(Person)
	p.Address.City = "shenz"
	err := CopyPropertiesIgnoreDefaultVal(&dto, p,false)
	if err != nil {
		t.Error(err)
	}
	bytes, _ := json.Marshal(p)
	t.Log(string(bytes))
}

func TestCopyPropertiesIgnoreFilter(t *testing.T) {
	dto := PersonDto{
		Name: "feiwo",
		Age:  12,
		Address: struct {
			City string `json:"city"`
			Area string `json:"area"`
		}{City:"changsha",Area:"yuelu"},
	}

	p := new(Person)
	p.Name = "xiaoqi"
	err := CopyPropertiesIgnoreField(&dto, p, "Name")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}
