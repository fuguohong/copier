package copier

import (
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"testing"
	"time"
)

type time1 struct {
	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt uint64
}

type time2 struct {
	CreateAt uint64
	UpdateAt *timestamppb.Timestamp
	DeleteAt *timestamppb.Timestamp
}

type CopyTest struct {
	suite.Suite
}

func TestCoPy(t *testing.T) {
	suite.Run(t, new(CopyTest))
}

func (t *CopyTest) TestCopyTime() {
	tfmt := "2006-01-02 15:04:05"
	tm, _ := time.Parse(tfmt, "2022-01-01 00:00:00")
	tmstr := "2022-01-01 08:00:00"
	tmint := tm.Unix()

	t.Run("time 2 uint64", func() {
		src := &time1{CreateAt: time.Unix(123456, 0)}
		dist := &time2{}
		Copy(src, dist)
		t.Equal(dist.CreateAt, uint64(123456))
		t.NotNil(dist.UpdateAt)

		t.True(dist.UpdateAt.AsTime().IsZero())
		t.True(dist.DeleteAt.AsTime().IsZero())
	})

	t.Run("uint64 2 time", func() {
		src := &time2{CreateAt: 1234567}
		dist := &time1{}
		Copy(src, dist)
		t.Equal(dist.CreateAt.Unix(), int64(1234567))

		t.True(dist.UpdateAt.IsZero())
		t.True(dist.DeleteAt == 0)
	})

	t.Run("uint64 2 timestamp", func() {
		src := &time1{DeleteAt: 123456}
		dist := &time2{}
		Copy(src, dist)
		t.Equal(dist.DeleteAt.Seconds, int64(123456))

		t.True(dist.CreateAt == 0)
		t.True(dist.UpdateAt.AsTime().IsZero())
	})

	t.Run("time 2 timestamp", func() {
		src := &time1{UpdateAt: time.Unix(123456, 0)}
		dist := &time2{}
		Copy(src, dist)
		t.Equal(dist.UpdateAt.Seconds, int64(123456))

		t.True(dist.CreateAt == 0)
		t.True(dist.DeleteAt.AsTime().IsZero())
	})

	t.Run("timestamp 2 uint64", func() {
		tmp := time.Unix(123456, 0)
		src := &time2{DeleteAt: timestamppb.New(tmp)}
		dist := &time1{}
		Copy(src, dist)
		t.Equal(dist.DeleteAt, uint64(123456))

		t.True(dist.CreateAt.IsZero())
		t.True(dist.UpdateAt.IsZero())
	})

	t.Run("timestamp 2 time", func() {
		src := timestamppb.New(tm)
		var dist time.Time
		Copy(src, &dist)
		t.Equal(dist.Unix(), tmint)
		t.Equal(dist.Format(tfmt), tmstr)
	})

	t.Run("timestamp 2 string", func() {
		src := timestamppb.New(tm)
		var dist string
		Copy(src, &dist)
		t.Equal(dist, tmstr)
	})

	t.Run("time 2 string", func() {
		var dist string
		Copy(tm, &dist)
		t.Equal(dist, tmstr)

		src, _ := time.Parse("2006-01-02 15:04:05-0700", "2022-01-01 00:00:00+0800")
		dist = ""
		Copy(src, &dist)
		t.Equal(dist, "2022-01-01 00:00:00")
	})
}

type attach struct {
	Url  string
	Type bool
}

type project struct {
	Name      string
	Price     float32
	CreateAt  uint64
	Att       attach
	Tags      []int64
	Subp      *project
	Version   uint8
	Time      time.Time
	Attaches  []*attach
	ProjectId int64
}

type soft struct {
	Name     string
	Price    float64
	Subp     *soft
	Att      attach
	CreateAT time.Time
	Tags     []int
	Version  uint32
	Time     time.Time
	Attaches []*attach
	SoftId   int64
}

func (t *CopyTest) TestCopyStruct() {
	src := &soft{
		Name:  "fgh_test",
		Price: 123.56,
		Subp: &soft{
			Name:  "inner pointer",
			Price: 888,
		},
		Att: attach{
			Url:  "inner struct",
			Type: true,
		},
		CreateAT: time.Unix(123456, 0),
		Tags:     []int{1, 2, 3},
		Version:  uint32(6),
		Time:     time.Unix(123456, 0),
		Attaches: []*attach{{Url: "1"}, {Url: "2"}},
		SoftId:   666,
	}

	dist := &project{}
	Copy(src, dist)
	t.Equal(dist.Name, "fgh_test")
	t.Equal(dist.Price, float32(123.56))
	t.NotNil(dist.Subp)
	t.Equal(dist.Subp.Name, "inner pointer")
	t.Equal(dist.Subp.Price, float32(888))
	t.Equal(dist.Att.Url, "inner struct")
	t.Equal(dist.Att.Type, true)
	t.Equal(dist.CreateAt, uint64(123456))
	t.Equal(len(dist.Tags), 3)
	t.Equal(dist.Tags[1], int64(2))
	t.Equal(dist.Time.Unix(), int64(123456))
	t.Equal(len(dist.Attaches), 2)
	t.Equal(dist.ProjectId, int64(0))

	dist = &project{}
	Copy(src, dist, map[string]string{
		"ProjectId": "SoftId",
	})
	t.Equal(dist.Name, "fgh_test")
	t.Equal(dist.CreateAt, uint64(123456))
	t.Equal(dist.ProjectId, int64(666))
}

func (t *CopyTest) TestCircleStruct() {
	s1 := &soft{Name: "1"}
	s2 := &soft{Name: "2"}
	s1.Subp = s2
	s2.Subp = s1

	var dist *soft
	Copy(s1, &dist)

	t.NotNil(dist)
	t.Equal(dist.Name, "1")
	t.Equal(dist.Subp.Name, "2")
}

func (t *CopyTest) TestCopySlice() {
	src := []*soft{{Name: "s1"}, {Name: "s2"}}
	var dist []*project
	Copy(src, &dist)
	t.Equal(len(dist), 2)
	t.Equal(dist[0].Name, "s1")
	t.Equal(dist[1].Name, "s2")
}

func (t *CopyTest) TestRegisterNewRule() {
	RegisterConverter(reflect.TypeOf(time.Time{}), reflect.TypeOf(""),
		func(v interface{}) interface{} {
			r := v.(time.Time)
			return r.Format(time.RFC3339)
		})
	src := time.Unix(1641288319, 0)
	var dist string
	Copy(src, &dist)
	t.Equal(dist, "2022-01-04T17:25:19+08:00")
}

func (t *CopyTest) TestMultiPointer() {
	src := soft{Name: "TestMultiPointer"}
	var dist *project
	Copy(src, &dist)
	t.Equal(dist.Name, "TestMultiPointer")

	dist = &project{}
	Copy(src, &dist)
	t.Equal(dist.Name, "TestMultiPointer")
}

func (t *CopyTest) TestAddr() {
	src := &soft{Name: "TestAddr"}
	var dist *soft
	Copy(src, dist)
	t.Nil(dist)

	Copy(src, &dist)
	t.NotNil(dist)
	t.Equal(dist.Name, "TestAddr")

	dist = &soft{}
	Copy(src, dist)
	t.NotNil(dist)
	t.Equal(dist.Name, "TestAddr")
}

func (t *CopyTest) TestCopyInt() {
	t.Run("int to uint and int", func() {
		src := 6
		var dist uint16
		Copy(src, &dist)
		t.Equal(dist, uint16(6))

		src = -32
		Copy(src, &dist)
		t.Equal(dist, uint16(0))

		var dist2 int16
		Copy(src, &dist2)
		t.Equal(dist2, int16(-32))
	})

	t.Run("uint to int and uint16", func() {
		src := uint32(16)
		var dist int
		Copy(src, &dist)
		t.Equal(dist, int(16))

		var dist2 uint16
		Copy(src, &dist2)
		t.Equal(dist2, uint16(16))
	})
}

func (t *CopyTest) TestCopyBool() {
	t.Run("int to bool", func() {
		src := 6
		var dist bool
		Copy(src, &dist)
		t.True(dist)

		src = 0
		Copy(src, &dist)
		t.False(dist)

		src2 := uint64(32)
		Copy(src2, &dist)
		t.True(dist)

		src2 = 0
		Copy(src2, &dist)
		t.False(dist)
	})

	t.Run("bool to int", func() {
		src := true
		var dist int
		var dist2 uint32
		Copy(src, &dist)
		Copy(src, &dist2)
		t.Equal(dist, 1)
		t.Equal(dist2, uint32(1))

		src = false
		Copy(src, &dist)
		Copy(src, &dist2)
		t.Equal(dist, 0)
		t.Equal(dist2, uint32(0))
	})
}
