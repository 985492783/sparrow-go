package core

import (
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

type fieldMap map[string]*fieldItem
type classNameMap map[string]fieldMap
type appNameMap map[string]classNameMap

// 1 app -> n class -> m field -> o ip

type NameSpace struct {
	mu      sync.RWMutex
	dataMap map[string]appNameMap
}

type ClientElement struct {
	field      string
	class      string
	appName    string
	nameSpaces *NameSpace
}
type ClientIp struct {
	ip       string
	elements []*ClientElement
	mu       sync.Mutex
}

var nameSpaceMap sync.Map
var clientMap sync.Map

func Register(clientId, namespace, appName, ip string, registry map[string]map[string]*SwitcherItem) {

	client, _ := clientMap.LoadOrStore(clientId, &ClientIp{
		ip: ip,
	})

	client.(*ClientIp).mu.Lock()
	defer client.(*ClientIp).mu.Unlock()

	// 防止被DeRegister删除
	clientMap.LoadOrStore(clientId, client)
	ns, _ := nameSpaceMap.LoadOrStore(namespace, &NameSpace{
		dataMap: make(map[string]appNameMap),
	})

	nameSpace := ns.(*NameSpace)
	nameSpace.mu.Lock()
	defer nameSpace.mu.Unlock()

	app := putIfAbsent(nameSpace.dataMap, appName, func() appNameMap {
		return make(appNameMap)
	})

	appMap := (map[string]classNameMap)(app)
	for clazz, fieldM := range registry {
		cls := putIfAbsent(appMap, clazz, func() classNameMap {
			return make(classNameMap)
		})
		classMap := (map[string]fieldMap)(cls)
		for fileName, field := range fieldM {
			fm := putIfAbsent(classMap, fileName, func() fieldMap {
				return make(fieldMap)
			})
			(map[string]*fieldItem)(fm)[ip] = convertToFieldItem(appName, clazz, ip, field)
		}
	}
}

func DeRegister(clientId string) {
	if client, ok := clientMap.Load(clientId); ok {
		clientIp := client.(*ClientIp)
		clientIp.mu.Lock()
		defer clientIp.mu.Unlock()

		for _, element := range clientIp.elements {
			deRegister(clientIp.ip, element)
		}
		// 移除clientMap
		clientMap.Delete(clientId)
	}
}

func deRegister(ip string, element *ClientElement) {
	element.nameSpaces.mu.Lock()
	defer element.nameSpaces.mu.Unlock()

	app, ok := element.nameSpaces.dataMap[element.appName]
	if !ok {
		return
	}
	appMap := (map[string]classNameMap)(app)

	cls, ok := appMap[element.class]
	if !ok {
		return
	}

	classMap := (map[string]fieldMap)(cls)
	fm, ok := classMap[element.field]
	if !ok {
		return
	}

	delete(fm, ip)
	if len(fm) == 0 {
		delete(classMap, element.field)
		if len(classMap) == 0 {
			delete(appMap, element.class)
			if len(appMap) == 0 {
				delete(element.nameSpaces.dataMap, element.appName)
			}
		}
	}
}

func putIfAbsent[T any](mp map[string]T, key string, fn func() T) T {
	app, ok := mp[key]
	if !ok {
		t := fn()
		mp[key] = t
		return t
	}
	return app
}

func convertToFieldItem(appName, className, ip string, item *SwitcherItem) *fieldItem {
	return &fieldItem{
		AppName:      appName,
		ClassName:    className,
		Ip:           ip,
		SwitcherItem: item,
	}
}

func GetNs() []string {
	ns := make([]string, 0)
	nameSpaceMap.Range(func(key, value any) bool {
		ns = append(ns, key.(string))
		return true
	})
	return ns
}

func GetJSON(namespace string) any {
	ns, ok := nameSpaceMap.Load(namespace)
	if !ok {
		return ""
	}
	return ns.(*NameSpace).dataMap
}
