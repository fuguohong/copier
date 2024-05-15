package copier

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type NameResolverTest struct {
	suite.Suite
}

func TestNameResolver(t *testing.T) {
	suite.Run(t, new(NameResolverTest))
}

type nameStruct struct {
	Id        int
	ID        int
	Id2       int
	FirstName int
	Lastname  int
	HomeUrl   int
	BaseURL   int
}

func (t *NameResolverTest) TestGetName() {
	n := nameStruct{}
	r := newNameResolver(reflect.ValueOf(n).Type(), nil)

	t.Equal(r.GetName("Id"), "Id")
	t.Equal(r.GetName("ID"), "ID")

	t.Equal(r.GetName("Id2"), "Id2")
	t.Equal(r.GetName("ID2"), "Id2")

	t.Equal(r.GetName("FirstName"), "FirstName")
	t.Equal(r.GetName("Firstname"), "FirstName")

	t.Equal(r.GetName("Lastname"), "Lastname")
	t.Equal(r.GetName("LastName"), "Lastname")

	t.Equal(r.GetName("HomeUrl"), "HomeUrl")
	t.Equal(r.GetName("HomeURL"), "HomeUrl")

	t.Equal(r.GetName("BaseURL"), "BaseURL")
	t.Equal(r.GetName("BaseUrl"), "BaseURL")

	t.Equal(r.GetName("FullName"), "")
}

func (t *NameResolverTest) TestGetNameWithMapping() {
	n := nameStruct{}
	r := newNameResolver(reflect.ValueOf(n).Type(), map[string]string{
		"FullName": "FirstName",
	})

	t.Equal(r.GetName("Id"), "Id")
	t.Equal(r.GetName("ID"), "ID")

	t.Equal(r.GetName("Id2"), "Id2")
	t.Equal(r.GetName("ID2"), "Id2")

	t.Equal(r.GetName("FirstName"), "FirstName")
	t.Equal(r.GetName("Firstname"), "FirstName")
	t.Equal(r.GetName("FullName"), "FirstName")
}
