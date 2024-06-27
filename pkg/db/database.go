package db

import (
	"errors"
	"github.com/985492783/sparrow-go/pkg/config"
	"net/url"
	"sync"
)

// Database 缓存
// 首先从数据源中加载持久化数据后，返回给上层
// 上层调用setSync更改数据同时写入持久化数据源
// 监听持久化数据源，如有更新，与data进行比对更新持久化数据并更新自身-最终一致性
type Database struct {
	datasource DataSource
	data       map[string]*Properties
	mu         sync.RWMutex
}

type DataSource interface {
	getData(ns, sparrow, fileName string) *Properties
	updateData(properties *Properties) error
}

type dBConfig struct {
	proto string
	host  string
	path  string
	query map[string]string
	// basic auth暂时不支持
	user *url.Userinfo
}

func NewDatabase(config *config.SparrowConfig) (*Database, error) {
	endpoint := config.DatabaseConfig.Endpoint
	parse, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	db := &dBConfig{
		proto: parse.Scheme,
		host:  parse.Host,
		path:  parse.Path,
	}
	var datasource DataSource
	var er error
	switch db.proto {
	case "file":
		datasource, er = newFileDB(db)
	default:
		er = errors.New("unsupported database protocol")
	}
	if er != nil {
		return nil, er
	}
	return &Database{
		datasource: datasource,
		data:       make(map[string]*Properties),
	}, nil
}

func (database *Database) GetData(ns, sparrow, fileName string) *Properties {
	key := ns + "@@" + sparrow + "@@" + fileName
	database.mu.RLock()
	property, ok := database.data[key]
	database.mu.RUnlock()
	if ok {
		return property
	}

	database.mu.Lock()
	defer database.mu.Unlock()
	data := database.datasource.getData(ns, sparrow, fileName)
	data.datasource = database.datasource

	database.data[key] = data
	return data
}
