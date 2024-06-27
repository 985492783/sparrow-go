package db

import (
	"encoding/json"
	"github.com/985492783/sparrow-go/pkg/config"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	con, err := config.LoadConfig("../../config.toml")
	if err != nil {
		t.Error(err)
	}
	_, err = NewDatabase(con)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFileDB_Read(t *testing.T) {
	config := &dBConfig{
		proto: "file",
		path:  "../../testdata",
	}
	fileDB, err := newFileDB(config)
	if err != nil {
		t.Error(err)
	}
	properties := fileDB.getData("public", "switcher", "test@@com.sparrow.test.Node")
	if properties.fileName != "test@@com.sparrow.test.Node" {
		t.Error("wrong file name")
	}
	properties.set("string_fieldA", "abcdefg")
	properties.set("string_fieldB", "!@#$%^&*()_+")
	properties.set("int_fieldC", 10000000)
	properties.set("float_fieldD", 1111.12546)
	properties.set("bool_fieldE", true)
	mp := map[string]string{
		"string_fieldA": "abcdefg",
		"string_fieldB": "!@#$%^&*()_+",
	}
	marshal, err := json.Marshal(mp)
	if err != nil {
		t.Error(err)
	}
	properties.set("json_fieldF", string(marshal))
	err = fileDB.updateData(properties)
	if err != nil {
		t.Error(err)
	}

	p2 := fileDB.getData("public", "switcher", "test@@com.sparrow.test.Node")
	if len(p2.data) != len(properties.data) {
		t.Error()
	}
}
