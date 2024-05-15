# copier
便捷复制工具，复制结构体，切片，基础变量等

[TOC]

## 背景

编写go代码时可能会遇到大量结构相同或非常相似的struct/slice需要相互转换，尤其是在使用rpc的时候，在各个垂直分层中传递数据，可能不同层级的数据结构几乎是一致的，但是因为是不同的proto,
生成的也是不同的结构体，不能直接使用。 这时候就要写大量a.name = b.name的转换代码，当字段数量很多并且结构存在嵌套时非常痛苦。 顾编写了这个工具用于做结构复制。

在做结构复制时，为了表示相同的信息，两边定义的字段类型却可能不同；例如时间，一边定义了uint64，一边定义为time， 所以需要支持自定义的类型转换规则，允许跨类型转换



## 特性

- 复制相同名称的结构体字段
- 递归复制结构体，默认最大5层
- 复制切片
- 数值类型互换。注意：高精度往低精度复制时精度丢失问题
- 跨类型互转。内置uint64 <-> time;不同精度的int、float转换;int<->bool;time -> string
- 自定义转化规则
- 自定义字段名映射


## 安装
```
go get github.com/fuguohong/copier
```



## 使用说明

**复制结构体**

```go
type a struct{
  Attr int
  Sub *sub
  BaseURL string
  AName string
}

type sub struct{
  string Name
}

type b struct{
  Attr int64
  Sub sub
  BaseURL string
  BName string
}

func main(){
  src := &a{Attr: 1, Sub: &sub{Name: "fgh"}, BaseURL: "a"}
  dist := &b{}
  copier.Copy(src, dist)
  // dist b{Attr: 1, Sub: sub{Name: "fgh"}, BaseUrl: "a"}
  
  // 查找子信息获得了model，往proto主结构附加
  proto := &a{Attr: 2}
  sub := das.find() // get &sub{Name: "test"}
  copier.Copy(sub, &proto.Sub)
  // proto.Sub : {Name: "test"}
  
  
  // 自定义字段名映射; distAttrName => srcAttrName
  src := &a{AName: "a"}
  dist := &b{}
  copier.Copy(src, dist, map[string]string{"BName": "AName"})
  // dist b{BName: "a"}
}
```

**复制切片**

```go
src := []*a{{Attr: 1}, {Attr: 2}}
var dist []*b
copier.Copy(src, &dist)
// dist : []*b{{Attr: 1}, {Attr: 2}}
```

**跨类型转换**

```go
type time1 struct {
	CreateAt time.Time
}

type time2 struct {
	CreateAt uint64
}

func main(){
    src := &time1{CreateAt: time.Unix(123456, 0)}
    dist := &time2{}
    Copy(src, dist)
    // dist.CreateAt: uint64(123456)
    
    src, _ = time.Parse("2006-01-02 15:04:05", "2022-01-05 11:20:00")
    var dist string
    Copy(src, &dist)
   // dist: "2022-01-05T11:20:00+08:00"
}

```

**注册转换规则**

```go
// 同类型规则会覆盖	
RegisterConverter(reflect.TypeOf(time.Time{}), reflect.TypeOf(""),
    func(v interface{}) interface{} {
        r := v.(time.Time)
        return r.Format(time.RFC3339)
    })

src := time.Unix(1641288319, 0)
var dist string
Copy(src, &dist)
// dist: 2022-01-04T17:25:19+08:00
```


## 注意事项

目标地址必须是可取地址的，如果不确定目标是否可取地址，dist参数统一加上&即可

```go
src := &soft{Name: "TestAddr"}
var dist *soft
Copy(src, dist) // 失败，dist不可取地址， dist为nil
Copy(src, &dist) // 正确，dist.Name: "TestAddr"
```


## TODO
- map < - > struct （暂无需求）

## test cover
```
ok      github.com/fuguohong/copier     0.329s  coverage: 100.0% of statements
```

