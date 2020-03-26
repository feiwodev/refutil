package bean

import (
	"errors"
	"reflect"
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
//  like spring bean util. copy object properties
// ------------------------------------------------------

var (
	SourceOrTargetNotPtrError = errors.New("源和目标结构体必须传入指针")
)

// 字段信息
type fieldInfo struct {
	Type reflect.Kind
	Name string
	Val reflect.Value
}

// CopyProperties: 复制结构字段值，忽略零值
//
// source: 源结构体指针
//
// target: 目标结构体指针
func CopyProperties(source interface{}, target interface{}) error {
	return CopyPropertiesIgnoreDefaultVal(source, target, true)
}

// CopyPropertiesIgnoreDefaultVal: 复制结构字段值，判断是否忽略零值
//
// source: 源结构体指针
//
// target: 目标结构体指针
//
// ignoreDefaultVal: 是否忽略零值，true:忽略， false: 不忽略
// 忽略零值：则不会拷贝源结构体中字段为零值的字段
// 不忽略零值：则不会不管源结构体中的字段值是否为零值，全部拷贝
func CopyPropertiesIgnoreDefaultVal(source interface{}, target interface{}, ignoreDefaultVal bool) error {
	return CopyPropertiesIgnoreFilter(source, target, ignoreDefaultVal)
}

// CopyPropertiesIgnoreField: 忽略零值，忽略字段
//
// source: 源结构体指针
//
// target: 目标结构体指针
//
// ignoreField: 忽略字段
func CopyPropertiesIgnoreField(source interface{}, target interface{}, ignoreField ...string) error {
	return CopyPropertiesIgnoreFilter(source, target, true, ignoreField...)
}

// CopyPropertiesIgnoreFilter: 复制结构字段值，判断是否忽略零值，并忽略指定字段
//
// source: 源结构体指针
//
// target: 目标结构体指针
//
// ignoreDefaultVal: 是否忽略零值
//
// ignoreField: 忽略字段
func CopyPropertiesIgnoreFilter(source interface{}, target interface{}, ignoreDefaultVal bool, ignoreField ...string) error {
	// 判断传入的是否都是指针，防止传入结构体，导致大量的复制
	if !hasInPtr(source) || !hasInPtr(target) {
		return SourceOrTargetNotPtrError
	}

	sourceFieldInfos := make([]*fieldInfo,0,10)
	targetFieldInfos := make([]*fieldInfo,0,10)
	getFields(getRefValue(source),&sourceFieldInfos)
	getFields(getRefValue(target), &targetFieldInfos)

	copyProp(sourceFieldInfos,targetFieldInfos,ignoreDefaultVal, ignoreField...)
	return nil
}


// hasInPtr: 判断传入类型是否是指针类型
//
// in : 任意类型
func hasInPtr(in interface{}) bool {
	typeOf := reflect.TypeOf(in)
	return typeOf.Kind() == reflect.Ptr
}

// getRefValue: 反射拿到值对象
//
// in : 任意类型
func getRefValue(in interface{}) reflect.Value {
	return reflect.ValueOf(in).Elem()
}

// hasValStructAndPtr：判断是否是结构体或接口
//
// val: 反射出的字段值结构
func hasValStructAndPtr(val reflect.Value) bool {
	if val.Kind() == reflect.Struct || val.Kind() == reflect.Ptr {
		return true
	}
	return false
}

// getFields: 递归获取Struct中的字段
//
// val: 反射字段Value
//
// fieldValues: Struct反射出的字段集合
func getFields(val reflect.Value, fieldValues *[]*fieldInfo)  {
	numField := val.NumField()
	for i := 0; i < numField ; i++  {
		field := val.Field(i)
		if !hasValStructAndPtr(field) {
			filedInfo := &fieldInfo{
				Type: field.Type().Kind(),
				Name: val.Type().Field(i).Name,
				Val:  field,
			}
			*fieldValues = append(*fieldValues, filedInfo)
		}
		if hasValStructAndPtr(field) {
			getFields(field, fieldValues)
		}
	}
}


//	copyProp：两个Struct交换值
//
//	sourceFieldInfo: 源Struct字段信息
//
//	targetFieldInfo: 目标Struct字段信息
//
//	ignoreDefaultVal: 是否忽略零值
//
//	ignoreField: 忽略字段可变参数
func copyProp(sourceFieldInfo []*fieldInfo, targetFieldInfo []*fieldInfo, ignoreDefaultVal bool, ignoreField ...string)  {
	for _, source := range sourceFieldInfo {
		for _, target := range targetFieldInfo {
			if source.Name == target.Name && source.Type == target.Type {
				// 有忽略字段
				if len(ignoreField) > 0 {
					if !hasIgnoreField(target.Name, ignoreField) {
						setVal(ignoreDefaultVal,source,target)
					}
				}else {
					// 无忽略字段
					setVal(ignoreDefaultVal,source,target)
				}
			}
		}
	}
}

// setVal: 交换字段值
//
// ignoreDefaultVal: 是否忽略零值
//
// source: 源Struct字段信息
//
// target: 目标Struct字段信息
func setVal(ignoreDefaultVal bool, source *fieldInfo, target *fieldInfo)  {
	// 是否忽略零值
	if ignoreDefaultVal {
		if !source.Val.IsZero() {
			target.Val.Set(source.Val)
		}
	}else {
		// 不忽略则全量赋值
		target.Val.Set(source.Val)
	}
}

// hasIgnoreField: 判断字段是否是忽略中的字段
//
// fieldName: 字段名字
//
// ignoreField: 忽略字段
func hasIgnoreField(fieldName string, ignoreField []string) bool {
	for _, name := range ignoreField {
		if fieldName == name{
			return true
		}
	}
	return false
}