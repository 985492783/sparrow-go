package center

import (
	"github.com/985492783/sparrow-go/pkg/utils"
)

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

var nameSpaceMap *utils.MapMutex[string, *NameSpace]

func register(namespace, appName, ip string, registry map[string]map[string]*SwitcherItem) {
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
		}
	}
}

func convertToFieldItem(appName, className, ip string, item *SwitcherItem) *fieldItem {
	return &fieldItem{
		AppName:      appName,
		ClassName:    className,
		Ip:           ip,
		SwitcherItem: item,
	}
}
