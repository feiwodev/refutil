## 结构体字段值赋值工具

**为什么会有这个工具？**

在平时开发API的时候，表结构映射的字段与我们想要输出的字段数不一致，有些字段我们是需要隐藏，不输出给前端的。一般的，前端传递过来的字段，也会与我们表结构字段不一致，通常我们会有两个结构来表示。

```go
// 注册示例
// 数据库结构
type User struct {
	Name string
	Mobile string
	State int
  ...
}
// 前端传递的数据结构
type Register struct {
  Name string `json:"name"`
  Mobile  string `json:"mobile"`
  ...
}

r := Register{}
u := User{Name:r.Name,Mobile:r.Mobile,State:1,...}
```

一般情况，我们会手动的创建一个`User`结构体，将`Register`赋值给`User`结构体，这样无疑是可行的，但当我们编写成百上千个接口时，结构体就会变得很多，结构体小的时候，我们还好，当结构体字段很多时，手动去赋值就会变得很没有效率。故而开发这个小工具。

**API函数**

```go
// CopyProperties: 复制结构字段值，忽略零值
//
// source: 源结构体指针
//
// target: 目标结构体指针
func CopyProperties(source interface{}, target interface{}) error 

// CopyPropertiesIgnoreDefaultVal: 复制结构字段值，判断是否忽略零值
//
// source: 源结构体指针
//
// target: 目标结构体指针
//
// ignoreDefaultVal: 是否忽略零值，true:忽略， false: 不忽略
// 忽略零值：则不会拷贝源结构体中字段为零值的字段
// 不忽略零值：则不会不管源结构体中的字段值是否为零值，全部拷贝
func CopyPropertiesIgnoreDefaultVal(source interface{}, target interface{}, ignoreDefaultVal bool) error 

// CopyPropertiesIgnoreField: 忽略零值，忽略字段
//
// source: 源结构体指针
//
// target: 目标结构体指针
//
// ignoreField: 忽略字段
func CopyPropertiesIgnoreField(source interface{}, target interface{}, ignoreField ...string) error 

// CopyPropertiesIgnoreFilter: 复制结构字段值，判断是否忽略零值，并忽略指定字段
//
// source: 源结构体指针
//
// target: 目标结构体指针
//
// ignoreDefaultVal: 是否忽略零值
//
// ignoreField: 忽略字段
func CopyPropertiesIgnoreFilter(source interface{}, target interface{}, ignoreDefaultVal bool, ignoreField ...string) error
```



**测试**

```go
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
```



