package pstore_test

import (
	"os"
	"strings"
	"testing"

	"github.com/saenai255/pstore"
)

type TestStruct struct {
	Value string
}

const TEST_ASSETS_PATH = "test_assets"

func assetCacheFileExists(t *testing.T, name, key string) {
	file := TEST_ASSETS_PATH + "/" + name + "_" + key + ".pcache"
	_, err := os.Stat(file)

	if err != nil {
		t.Errorf("cache file %v does not exist: %v", file, err)
	}
}

func assetCacheFileNotExists(t *testing.T, name, key string) {
	file := TEST_ASSETS_PATH + "/" + name + "_" + key + ".pcache"
	_, err := os.Stat(file)

	if err == nil {
		t.Errorf("cache file exists: %v", file)
	}
}

func Test_SetPrimitive(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	assetCacheFileExists(t, "CacheName", "KeyName")
}

func Test_SetStruct(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", TestStruct{Value: "value"}); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	assetCacheFileExists(t, "CacheName", "KeyName")
}

func TestSet_InMem(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")
	p.SaveToDiskOnSet = false

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	assetCacheFileNotExists(t, "CacheName", "KeyName")
}

func Test_GetPrimitive(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value string
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if value != "Value" {
		t.Errorf("got wrong value: %v", value)
	}
}

func Test_GetStruct(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", TestStruct{Value: "value"}); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value TestStruct
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if value.Value != "value" {
		t.Errorf("got wrong value: %v", value)
	}

}

func Test_GetStructPtr(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", &TestStruct{Value: "value"}); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value TestStruct
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if value.Value != "value" {
		t.Errorf("got wrong value: %v", value)
	}
}

func Test_GetStructPtrPtr(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	it := &TestStruct{Value: "value"}
	if err := p.Set("KeyName", &it); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value TestStruct
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if value.Value != "value" {
		t.Errorf("got wrong value: %v", value)
	}
}

func Test_GetPrimitivePtr(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value string
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if value != "Value" {
		t.Errorf("got wrong value: %v", value)
	}
}

func Test_GetArray(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", []int{1, 2, 3}); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value []int
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if len(value) != 3 {
		t.Errorf("got wrong value: %v", value)
	}

	if value[0] != 1 || value[1] != 2 || value[2] != 3 {
		t.Errorf("got wrong value: %v", value)
	}
}

func Test_GetMap(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", map[string]int{"key": 1}); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	var value map[string]int
	if err := p.Get("KeyName", &value); err != nil {
		t.Errorf("failed to get key: %v", err)
	}

	if len(value) != 1 {
		t.Errorf("got wrong value: %v", value)
	}

	if value["key"] != 1 {
		t.Errorf("got wrong value: %v", value)
	}
}

func Test_Has(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	if found, err := p.Has("KeyName"); err != nil || !found {
		t.Errorf("key not found")
	}
}

func Test_Keys(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("Key1", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	if err := p.Set("Key2", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	keys, err := p.Keys()
	if err != nil {
		t.Errorf("failed to get keys: %v", err)
	}

	if len(keys) != 2 {
		t.Errorf("got wrong keys: %v", keys)
	}

	// Order is not guaranteed
	if keys[0] != "Key1" && keys[1] != "Key1" || keys[0] != "Key2" && keys[1] != "Key2" {
		t.Errorf("got wrong keys: %v", keys)
	}
}

func Test_Delete(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	if err := p.Delete("KeyName"); err != nil {
		t.Errorf("failed to delete key: %v", err)
	}

	assetCacheFileNotExists(t, "CacheName", "KeyName")
}

func Test_SaveToDisk(t *testing.T) {
	t.Cleanup(ClearTestAssets)

	p := pstore.New(TEST_ASSETS_PATH, "CacheName")
	p.SaveToDiskOnSet = false

	if err := p.Set("KeyName", "Value"); err != nil {
		t.Errorf("failed to set key: %v", err)
	}

	assetCacheFileNotExists(t, "CacheName", "KeyName")

	if err := p.SaveToDisk(); err != nil {
		t.Errorf("failed to save to disk: %v", err)
	}

	assetCacheFileExists(t, "CacheName", "KeyName")
}

func ClearTestAssets() {
	files, err := os.ReadDir(TEST_ASSETS_PATH)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".pcache") {
			os.Remove(TEST_ASSETS_PATH + "/" + file.Name())
		}
	}
}
