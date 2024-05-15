package copier

import (
	"reflect"
	"strings"
)

type nameResolver struct {
	structValue reflect.Type
	nameMap     map[string]string
	customMap   map[string]string
}

// GetName 获取属性名
func (n *nameResolver) GetName(name string) string {
	if n.customMap != nil {
		targetName, ok := n.customMap[name]
		if ok {
			return targetName
		}
	}
	targetName, _ := n.nameMap[name]
	if targetName != "" {
		return targetName
	}
	targetName, _ = n.nameMap[strings.ToLower(name)]
	return targetName
}

func newNameResolver(structValue reflect.Type, customMap map[string]string) *nameResolver {
	n := &nameResolver{
		structValue: structValue,
		nameMap:     make(map[string]string),
		customMap:   customMap,
	}

	for i := 0; i < structValue.NumField(); i++ {
		name := structValue.Field(i).Name
		n.nameMap[name] = name
		n.nameMap[strings.ToLower(name)] = name
	}

	return n
}
