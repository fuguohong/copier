English | [简体中文](README-CN.md)

# copier
copy tool for copying structs, slices, and basic variables.


## Background

When writing Go code, you may come across a large number of structs or slices that have the same or very similar structures and need to be converted between each other. This is especially true when using RPC, where data needs to be passed between different vertical layers, and the data structures in different layers are almost identical but cannot be used directly due to different protos generating different structs. In such cases, you end up writing a lot of conversion code like `a.name = b.name`, which can be painful when dealing with a large number of fields and nested structures. To address this, the copier tool was developed for struct copying.

During struct copying, the field types defined on both sides may differ to represent the same information. For example, one side may define a field as `uint64`, while the other side defines it as `time`. Therefore, copier supports custom type conversion rules to allow for cross-type conversions.

## Features

- Copy fields with the same name in structs
- Recursive copying of structs, with a default maximum depth of 5 levels
- Copy slices
- Conversion between numerical types. Note: Precision loss may occur when copying from higher precision to lower precision.
- Cross-type conversions. Built-in conversions: uint64 <-> time; int <-> bool; time -> string; conversions between int and float with different precisions; 
- Custom conversion rules
- Custom field name mapping

## Installation

```
go get github.com/fuguohong/copier
```

## Usage

**Copying Structs**

```go
type a struct {
  Attr    int
  Sub     *sub
  BaseURL string
  AName   string
}

type sub struct {
  Name string
}

type b struct {
  Attr    int64
  Sub     sub
  BaseURL string
  BName   string
}

func main() {
  src := &a{Attr: 1, Sub: &sub{Name: "fgh"}, BaseURL: "a"}
  dist := &b{}
  copier.Copy(src, dist)
  // dist: b{Attr: 1, Sub: sub{Name: "fgh"}, BaseUrl: "a"}

  // Attaching model to the main proto structure by fetching sub information
  proto := &a{Attr: 2}
  sub := service.FindSomething() // get &sub{Name: "test"}
  copier.Copy(sub, &proto.Sub)
  // proto.Sub: {Name: "test"}

  // Custom field name mapping; distAttrName => srcAttrName
  src := &a{AName: "a"}
  dist := &b{}
  copier.CopyWithMapping(src, dist, map[string]string{"BName": "AName"})
  // dist: b{BName: "a"}
}
```

**Copying Slices**

```go
src := []*a{{Attr: 1}, {Attr: 2}}
var dist []*b
copier.Copy(src, &dist)
// dist: []*b{{Attr: 1}, {Attr: 2}}
```

**Cross-Type Conversion**

```go
type time1 struct {
	CreateAt time.Time
}

type time2 struct {
	CreateAt uint64
}

func main() {
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

**Registering Conversion Rules**

```go
// Same type rules will override
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

## Notes

The destination address must be a valid address. If you are unsure whether the destination is a valid address, add an `&` to the `dist` parameter.

```go
src := &soft{Name: "TestAddr"}
var dist *soft
Copy(src, dist) // Failed, dist is not a valid address, dist is nil
Copy(src, &dist) // Correct, dist.Name: "TestAddr"
```

## TODO
- map < - > struct （no demand）

## test cover
```
ok      github.com/fuguohong/copier     0.329s  coverage: 100.0% of statements
```