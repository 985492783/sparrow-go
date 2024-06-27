package db

import (
	"maps"
	"sync"
)

type Properties struct {
	mu         sync.RWMutex
	data       map[string]any
	ns         string
	sparrow    string
	fileName   string
	datasource DataSource
}

func newProperties(ns, sparrow, fileName string) *Properties {
	return &Properties{
		data:     make(map[string]any),
		ns:       ns,
		sparrow:  sparrow,
		fileName: fileName,
	}
}

func (p *Properties) set(key string, value any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[key] = value
}

// SetSync 写入并同步修改datasource，如果datasource出现异常则回滚
func (p *Properties) SetSync(key string, value any) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	per, exist := p.data[key]

	p.data[key] = value
	err := p.datasource.updateData(p)
	//回滚
	if err != nil {
		if exist {
			p.data[key] = per
		} else {
			delete(p.data, key)
		}
		return err
	}
	return nil
}

func (p *Properties) SetAll(mapping map[string]any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	maps.Copy(p.data, mapping)
}

func (p *Properties) Get(key string) (any, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	d, ok := p.data[key]
	return d, ok
}

func (p *Properties) GetString(key string) (*string, bool) {
	if data, ok := p.Get(key); ok {
		return data.(*string), true
	}
	return nil, false
}
func (p *Properties) GetBool(key string) (*bool, bool) {
	if data, ok := p.Get(key); ok {
		return data.(*bool), true
	}
	return nil, false
}
func (p *Properties) GetInt(key string) (*int, bool) {
	if data, ok := p.Get(key); ok {
		return data.(*int), true
	}
	return nil, false
}
func (p *Properties) GetInt64(key string) (*int64, bool) {
	if data, ok := p.Get(key); ok {
		return data.(*int64), true
	}
	return nil, false

}
