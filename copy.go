package copier

import (
	"reflect"
)

// MaxDepth 结构体最大复制深度
var MaxDepth uint8 = 5

// Copy 从src往dist复制, mapping可选，复制结构体时自定义属性名映射, distName => srcMame
func Copy(src, dist interface{}, mapping ...map[string]string) {
	var m map[string]string
	if len(mapping) > 0 {
		m = mapping[0]
	}
	apply(reflect.ValueOf(src), reflect.ValueOf(dist), 0, m)
}

func apply(src reflect.Value, dist reflect.Value, depth uint8, mapping map[string]string) {
	for src.Kind() == reflect.Ptr {
		if src.IsNil() {
			return
		}
		src = src.Elem()
	}
	for dist.Kind() == reflect.Ptr {
		if dist.IsNil() && dist.CanAddr() {
			dist.Set(reflect.New(dist.Type().Elem()))
		}
		dist = dist.Elem()
	}
	if !dist.CanAddr() {
		return
	}

	srcType := src.Type()
	distType := dist.Type()

	if srcType == distType {
		dist.Set(src)
	}

	converter := getConverter(srcType, distType)
	if converter != nil {
		distval := converter(src.Interface())
		dist.Set(reflect.ValueOf(distval))
		return
	}

	if copyInt(src, dist) {
		return
	}

	// if copyUint(src, dist) {
	// 	return
	// }

	if copyBool(src, dist) {
		return
	}

	if copyFloat(src, dist) {
		return
	}

	if copyStruct(src, dist, depth, mapping) {
		return
	}

	if copySlice(src, dist, depth) {
		return
	}
	// more support
}

func copyStruct(src reflect.Value, dist reflect.Value, depth uint8, mapping map[string]string) bool {
	if src.Kind() != reflect.Struct || dist.Kind() != reflect.Struct {
		return false
	}
	if depth >= MaxDepth {
		return true
	}
	r := newNameResolver(src.Type(), mapping)
	distType := dist.Type()
	for i := 0; i < distType.NumField(); i++ {
		fieldType := distType.Field(i)
		if fieldType.PkgPath != "" {
			continue
		}
		srcName := r.GetName(fieldType.Name)
		if srcName == "" {
			continue
		}
		srcField := src.FieldByName(srcName)
		if !srcField.IsValid() {
			continue
		}
		distField := dist.FieldByName(fieldType.Name)
		apply(srcField, distField, depth+1, mapping)
	}
	return true
}

func copyFloat(src reflect.Value, dist reflect.Value) bool {
	if isFloat(src.Kind()) && isFloat(dist.Kind()) {
		dist.SetFloat(src.Float())
		return true
	}
	return false
}

func isFloat(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

// func copyInt(src reflect.Value, dist reflect.Value) bool {
// 	if isInt(src.Kind()) && isInt(dist.Kind()) {
// 		dist.SetInt(src.Int())
// 	}
// 	return false
// }

func copyInt(src reflect.Value, dist reflect.Value) bool {
	if isInt(src.Kind()) {
		if isInt(dist.Kind()) {
			dist.SetInt(src.Int())
			return true
		} else if isUint(dist.Kind()) {
			x := src.Int()
			if x < 0 {
				x = 0
			}
			dist.SetUint(uint64(x))
			return true
		}
		return false
	} else if isUint(src.Kind()) {
		if isInt(dist.Kind()) {
			dist.SetInt(int64(src.Uint()))
			return true
		} else if isUint(dist.Kind()) {
			dist.SetUint(src.Uint())
			return true
		}
		return false
	}
	return false
}

func isInt(k reflect.Kind) bool {
	return k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 ||
		k == reflect.Int32 || k == reflect.Int64
}

// func copyUint(src reflect.Value, dist reflect.Value) bool {
// 	if isUint(src.Kind()) && isUint(dist.Kind()) {
// 		dist.SetUint(src.Uint())
// 	}
// 	return false
// }

func isUint(k reflect.Kind) bool {
	return k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 ||
		k == reflect.Uint32 || k == reflect.Uint64
}

func copySlice(src reflect.Value, dist reflect.Value, depth uint8) bool {
	if src.Kind() != reflect.Slice || dist.Kind() != reflect.Slice {
		return false
	}
	if src.Len() == 0 {
		return true
	}
	dist.Set(reflect.MakeSlice(dist.Type(), src.Len(), src.Len()))
	for i := 0; i < src.Len(); i++ {
		apply(src.Index(i), dist.Index(i), depth, nil)
	}
	return true
}

func copyBool(src reflect.Value, dist reflect.Value) bool {
	if src.Kind() == reflect.Bool {
		val := 0
		if src.Bool() {
			val = 1
		}
		if isInt(dist.Kind()) {
			dist.SetInt(int64(val))
			return true
		} else if isUint(dist.Kind()) {
			dist.SetUint(uint64(val))
			return true
		}
		return false
	} else if dist.Kind() == reflect.Bool {
		val := false
		if isInt(src.Kind()) {
			if src.Int() != 0 {
				val = true
			}
			dist.SetBool(val)
			return true
		} else if isUint(src.Kind()) {
			if src.Uint() != 0 {
				val = true
			}
			dist.SetBool(val)
			return true
		}
		return false
	}
	return false
}
