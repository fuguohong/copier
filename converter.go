package copier

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"time"
)

// Converter type convert func
type Converter func(interface{}) interface{}

type typeConverter struct {
	src  reflect.Type
	dist reflect.Type
	fn   Converter
}

var converts = make([]typeConverter, 0)

// RegisterConverter 注册新的类型转换函数
func RegisterConverter(src reflect.Type, dist reflect.Type, fn Converter) {
	for i, c := range converts {
		if c.src == src && c.dist == dist {
			converts[i] = typeConverter{
				src:  src,
				dist: dist,
				fn:   fn,
			}
			return
		}
	}
	converts = append(converts, typeConverter{
		src:  src,
		dist: dist,
		fn:   fn,
	})
}

func init() {
	var (
		TypeUint64    = reflect.TypeOf(uint64(0))
		TypeTime      = reflect.TypeOf(time.Time{})
		TypeTimestamp = reflect.TypeOf(*timestamppb.Now())
		TypeString    = reflect.TypeOf("")

		ZeroTime   = time.Time{}
		TimeFormat = "2006-01-02 15:04:05"
	)
	RegisterConverter(TypeTime, TypeUint64, func(v interface{}) interface{} {
		r := v.(time.Time)
		if r.Unix() <= 0 {
			return uint64(0)
		}
		return uint64(r.Unix())
	})

	RegisterConverter(TypeUint64, TypeTime, func(v interface{}) interface{} {
		r := v.(uint64)
		if r == 0 {
			return ZeroTime
		}
		return time.Unix(int64(r), 0)
	})

	RegisterConverter(TypeTimestamp, TypeUint64, func(v interface{}) interface{} {
		r := v.(timestamppb.Timestamp)
		if r.Seconds <= 0 {
			return uint64(0)
		}
		return uint64(r.Seconds)
	})

	RegisterConverter(TypeUint64, TypeTimestamp, func(v interface{}) interface{} {
		r := v.(uint64)
		if r == 0 {
			return *timestamppb.New(ZeroTime)
		}
		return timestamppb.Timestamp{
			Seconds: int64(r),
			Nanos:   0,
		}
	})

	RegisterConverter(TypeTime, TypeTimestamp, func(v interface{}) interface{} {
		r := v.(time.Time)
		return *timestamppb.New(r)
	})

	RegisterConverter(TypeTimestamp, TypeTime, func(v interface{}) interface{} {
		r := v.(timestamppb.Timestamp)
		return r.AsTime().Local()
	})

	RegisterConverter(TypeTime, TypeString, func(v interface{}) interface{} {
		r := v.(time.Time)
		return r.Local().Format(TimeFormat)
	})

	RegisterConverter(TypeTimestamp, TypeString, func(v interface{}) interface{} {
		r := v.(timestamppb.Timestamp)
		return r.AsTime().Local().Format(TimeFormat)
	})

}

func getConverter(src reflect.Type, dist reflect.Type) Converter {
	for _, c := range converts {
		if c.src == src && c.dist == dist {
			return c.fn
		}
	}
	return nil
}
