package copier

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type CopyTest struct {
	suite.Suite
}

func TestCoPy(t *testing.T) {
	suite.Run(t, new(CopyTest))
}

func (t *CopyTest) TestCopyTime() {

	t.Run("time 2 uint64", func() {
		src := time.Unix(123456, 0)
		var dist uint64
		Copy(src, &dist)
		t.Equal(uint64(123456), dist)

		src2 := time.Time{}
		var dist2 uint64
		Copy(src2, &dist2)
		t.Equal(uint64(0), dist2)
	})

	t.Run("uint64 2 time", func() {
		src := uint64(1234567)
		dist := time.Time{}
		Copy(src, &dist)
		t.Equal(dist.Unix(), int64(1234567))
	})

	t.Run("time 2 string", func() {
		src, _ := time.Parse(time.RFC3339, "2024-05-01T09:08:00+08:00")
		var dist string
		Copy(src, &dist)
		t.Equal(dist, "2024-05-01T09:08:00+08:00")
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
	CreateAt time.Time
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
		CreateAt: time.Unix(123456, 0),
		Tags:     []int{1, 2, 3},
		Version:  uint32(6),
		Time:     time.Unix(123456, 0),
		Attaches: []*attach{{Url: "1"}, {Url: "2"}},
		SoftId:   666,
	}

	dist := &project{}
	Copy(src, dist)
	t.Equal("fgh_test", dist.Name)
	t.Equal(float32(123.56), dist.Price)
	t.NotNil(dist.Subp)
	t.Equal("inner pointer", dist.Subp.Name)
	t.Equal(float32(888), dist.Subp.Price)
	t.Equal("inner struct", dist.Att.Url)
	t.Equal(true, dist.Att.Type)
	t.Equal(uint64(123456), dist.CreateAt)
	t.Equal(3, len(dist.Tags))
	t.Equal(int64(2), dist.Tags[1])
	t.Equal(int64(123456), dist.Time.Unix())
	t.Equal(2, len(dist.Attaches))
	t.Equal(int64(0), dist.ProjectId)
}

func (t *CopyTest) TestCircleStruct() {
	s1 := &soft{Name: "1"}
	s2 := &soft{Name: "2"}
	s1.Subp = s2
	s2.Subp = s1

	var dist *soft
	Copy(s1, &dist)

	t.NotNil(dist)
	t.Equal("1", dist.Name)
	t.Equal("2", dist.Subp.Name)
}

func (t *CopyTest) TestMapping() {
	src := project{
		Name:      "namea",
		ProjectId: 123,
	}

	var dist *soft
	CopyWithMapping(src, &dist, map[string]string{
		"SoftId": "ProjectId",
		"Name":   "",
	})
	t.Equal(int64(123), dist.SoftId)
	t.Equal("", dist.Name)

	var dist2 *soft
	CopyWithMapping(src, &dist2, nil)
	t.Equal(int64(0), dist2.SoftId)
	t.Equal("namea", dist2.Name)
}

func (t *CopyTest) TestCopySlice() {
	src := []*soft{{Name: "s1"}, {Name: "s2"}}
	var dist []*project
	Copy(src, &dist)
	t.Equal(2, len(dist))
	t.Equal("s1", dist[0].Name)
	t.Equal("s2", dist[1].Name)
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
	t.Equal("2022-01-04T17:25:19+08:00", dist)
}

func (t *CopyTest) TestMultiPointer() {
	src := soft{Name: "TestMultiPointer"}
	var dist *project
	Copy(src, &dist)
	t.Equal("TestMultiPointer", dist.Name)

	dist = &project{}
	Copy(src, &dist)
	t.Equal("TestMultiPointer", dist.Name)
}

func (t *CopyTest) TestAddr() {
	src := &soft{Name: "TestAddr"}
	var dist *soft
	Copy(src, dist)
	t.Nil(dist)

	Copy(src, &dist)
	t.NotNil(dist)
	t.Equal("TestAddr", dist.Name)

	dist = &soft{}
	Copy(src, dist)
	t.NotNil(dist)
	t.Equal("TestAddr", dist.Name)

	var src2 *soft
	var dist2 *soft
	Copy(src2, dist2)
	t.Nil(dist2)
}

func (t *CopyTest) TestCopyInt() {
	t.Run("int to uint and int", func() {
		src := 6
		var dist uint16
		Copy(src, &dist)
		t.Equal(uint16(6), dist)

		src = -32
		Copy(src, &dist)
		t.Equal(uint16(0), dist)

		var dist2 int16
		Copy(src, &dist2)
		t.Equal(int16(-32), dist2)
	})

	t.Run("uint to int and uint16", func() {
		src := uint32(16)
		var dist int
		Copy(src, &dist)
		t.Equal(int(16), dist)

		var dist2 uint16
		Copy(src, &dist2)
		t.Equal(uint16(16), dist2)
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

		src3 := ""
		var dist3 bool
		Copy(src3, &dist3)
		t.False(dist3)
	})

	t.Run("bool to int", func() {
		src := true
		var dist int
		var dist2 uint32
		Copy(src, &dist)
		Copy(src, &dist2)
		t.Equal(1, dist)
		t.Equal(uint32(1), dist2)

		src = false
		Copy(src, &dist)
		Copy(src, &dist2)
		t.Equal(0, dist)
		t.Equal(uint32(0), dist2)
	})
}
