package center

import "sync"

type MapMutex[K comparable, V any] struct {
	mmap map[K]V
	lock sync.RWMutex
}

func (m *MapMutex[K, V]) Load() {

}

type fieldItem struct {
}

type fieldMap MapMutex[string, fieldItem]
type classNameMap MapMutex[string, fieldMap]
type appKeyMap MapMutex[string, classNameMap]

// 1 app -> n class -> m field -> o ip

type NameSpace struct {
	dataMap MapMutex[string, appKeyMap]
}

var nameSpaceMap sync.Map
