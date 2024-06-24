package core

import (
	"github.com/985492783/sparrow-go/pkg/utils"
	"sync"
)

// TODO 回调，stream结束后调用manager移除所有注册项
type SwitcherItem struct {
	Type      string `json:"type"`
	Value     any    `json:"value"`
	FieldName string `json:"fieldName"`
	Desc      string `json:"desc"`
}

type fieldItem struct {
	AppName   string
	ClassName string
	Ip        string
	*SwitcherItem
}

type fieldMap *utils.MapMutex[string, *fieldItem]
type classNameMap *utils.MapMutex[string, fieldMap]
type appNameMap *utils.MapMutex[string, classNameMap]

// 1 app -> n class -> m field -> o ip

type NameSpace struct {
	dataMap *utils.MapMutex[string, appNameMap]
}
type ClientElement struct {
	fields     fieldMap
	classes    classNameMap
	apps       appNameMap
	nameSpaces *NameSpace
}
type ClientIp struct {
	ip       string
	mu       sync.Mutex
	elements []*ClientElement
}

var nameSpaceMap *utils.MapMutex[string, *NameSpace]
var clientMap *utils.MapMutex[string, *ClientIp]

func Register(clientId, namespace, appName, ip string, registry map[string]map[string]*SwitcherItem) {
	clientIp := clientMap.ComputeIfAbsent(clientId, func() *ClientIp {
		return &ClientIp{
			ip:       ip,
			elements: make([]*ClientElement, 16),
		}
	})
	clientIp.mu.Lock()
	defer clientIp.mu.Unlock()

	nameSpace := nameSpaceMap.ComputeIfAbsent(namespace, func() *NameSpace {
		return &NameSpace{
			dataMap: utils.NewMapMutex[string, appNameMap](),
		}
	})

	app := nameSpace.dataMap.ComputeIfAbsent(appName, func() appNameMap {
		return utils.NewMapMutex[string, classNameMap]()
	})

	appMap := (*utils.MapMutex[string, classNameMap])(app)
	for clazz, fieldM := range registry {
		cls := appMap.ComputeIfAbsent(clazz, func() classNameMap {
			return utils.NewMapMutex[string, fieldMap]()
		})
		classMap := (*utils.MapMutex[string, fieldMap])(cls)
		for fileName, field := range fieldM {
			fm := classMap.ComputeIfAbsent(fileName, func() fieldMap {
				return utils.NewMapMutex[string, *fieldItem]()
			})
			(*utils.MapMutex[string, *fieldItem])(fm).Put(ip, convertToFieldItem(appName, clazz, ip, field))

			//load in clientIp
			clientIp.elements = append(clientIp.elements, &ClientElement{
				fields:     fm,
				classes:    classMap,
				apps:       app,
				nameSpaces: nameSpace,
			})

		}
	}
}

func DeRegister(clientId string) {
	clientIp, ok := clientMap.Get(clientId)
	if !ok {
		return
	}
	clientIp.mu.Lock()
	defer clientIp.mu.Unlock()

}

func convertToFieldItem(appName, className, ip string, item *SwitcherItem) *fieldItem {
	return &fieldItem{
		AppName:      appName,
		ClassName:    className,
		Ip:           ip,
		SwitcherItem: item,
	}
}
