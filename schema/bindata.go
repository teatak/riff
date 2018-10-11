// Code generated by go-bindata.
// sources:
// schema/node.graphqls
// schema/schema.graphqls
// schema/service.graphqls
// DO NOT EDIT!

package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaNodeGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\xb1\x6a\xc4\x30\x10\x44\x7b\x7f\xc5\xfa\x37\x54\xc6\x55\x1a\x13\x50\x19\x52\x08\x7b\x6d\x0b\x9c\x5d\xa1\x9d\x18\xc2\x71\xff\x7e\x60\xb9\x38\xeb\x38\x38\xae\xdc\x27\x31\xc3\x9b\x06\xff\x89\xa9\xd7\x91\xe9\xd2\x10\x11\x8d\x01\xa1\x63\x01\x67\x47\x1e\x39\xca\xdc\xee\x3c\xa6\xea\x36\xcf\xeb\xe4\xe8\x43\x75\xe5\x20\x05\x4a\xf8\xe5\xf3\xb7\x05\x48\x5f\x9a\xe1\xe8\x53\x50\x50\x4e\x43\x45\x8c\xf3\x16\x07\x36\x47\xdf\x3d\x1b\x7c\x39\xdb\x9f\xf2\x28\x21\xf9\x45\x71\xce\x35\x04\xec\x55\x01\x5c\xc8\xc6\xd9\xa2\xca\x11\x7b\x6d\x0e\x33\x36\xdc\xd9\x0d\x2a\x53\x9c\x1f\xa3\xfe\xac\x53\x01\x4b\xd5\xf2\x74\x8c\xd7\xe4\xeb\xcd\x52\xa5\xfd\xae\xd9\x2d\x00\x00\xff\xff\x6b\x5f\x25\x3f\xb5\x01\x00\x00")

func schemaNodeGraphqlsBytes() ([]byte, error) {
	return bindataRead(
		_schemaNodeGraphqls,
		"schema/node.graphqls",
	)
}

func schemaNodeGraphqls() (*asset, error) {
	bytes, err := schemaNodeGraphqlsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/node.graphqls", size: 437, mode: os.FileMode(420), modTime: time.Unix(1539249430, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaSchemaGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x92\xb1\x6e\x83\x30\x10\x86\x77\x3f\xc5\x65\x0b\x52\x9f\xc0\x1b\x49\x18\x22\x25\x6d\x8a\x69\x97\x8a\xc1\x82\x03\x2c\x15\x93\xda\x26\x52\x54\xf1\xee\x15\xe6\x8c\x80\x4e\xfc\xfe\xee\xff\x8f\x3b\x83\x2d\x1a\x6c\x25\xfc\x32\x00\x80\x9f\x1e\xcd\x93\xc3\xfb\xf8\xf0\xa0\xed\x9d\x74\xaa\xd3\x1c\xae\xa4\xd8\xc0\x98\x7b\xde\x71\x32\x51\x4e\x77\x25\xee\xb5\x6c\x91\x83\x70\x46\xe9\x7a\x17\x71\x78\xed\x4a\x9c\xab\x96\xc3\xd7\x08\x72\x4f\x8c\xaa\x2a\x0e\xa9\xaa\x2a\x7f\xb4\x68\x1e\x68\x16\x89\x11\xa8\x62\xd3\xf2\xc5\x3a\xe9\xfc\x51\x3a\x8c\x38\x88\xc9\xb4\x0c\x8c\x6f\x21\x9c\xcf\x83\x86\xc9\x69\xd6\xb0\x12\xf9\xf6\x8b\xe8\x75\x5d\x3a\xeb\x7b\xef\x76\x79\xf4\xbf\x42\x5b\x60\xad\xac\xc3\x4d\x23\x0e\xe9\x8a\x4f\x5d\x22\x0e\x87\xae\xfb\x46\xa9\x7d\xb4\xd7\x9b\xf0\xf6\xf2\x82\x79\x60\x0c\x75\xdf\x4e\x4b\xd3\x06\xf1\xe5\xfc\x99\x78\x25\x3e\xc4\x2d\x39\x66\x5e\x9f\x92\xf8\x44\xe5\xcb\x1c\x3b\xb6\x25\x85\x44\x16\xa7\x19\xa9\xb7\x9b\x17\x69\x32\xc1\x70\x51\xe3\xf7\x20\x77\xad\xdc\xc1\x48\x5d\x34\xf3\x4c\x01\x8b\x46\xae\xd9\x03\x8d\xf5\x7f\x48\x80\x03\xfb\x0b\x00\x00\xff\xff\x55\xb2\x9e\xcb\x54\x02\x00\x00")

func schemaSchemaGraphqlsBytes() ([]byte, error) {
	return bindataRead(
		_schemaSchemaGraphqls,
		"schema/schema.graphqls",
	)
}

func schemaSchemaGraphqls() (*asset, error) {
	bytes, err := schemaSchemaGraphqlsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/schema.graphqls", size: 596, mode: os.FileMode(420), modTime: time.Unix(1539226389, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaServiceGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\xc1\x6a\xc5\x20\x10\x45\xf7\x7e\xc5\xf8\x1b\x2e\x9b\x55\x16\x0d\xa5\x59\x96\x2e\x44\xa7\x22\x34\x33\xa2\x63\xa1\x94\xfc\x7b\x49\xd2\x06\x4c\xc2\x83\xf7\x96\x57\xb9\xe7\x1e\x98\x48\xa9\x0a\x3c\x57\xb1\x12\x99\x46\xcc\x5f\xd1\x61\xbf\x3e\xfe\x28\x00\x00\x37\x79\x03\xdd\xe4\xf5\x9a\xc8\x4e\x68\x60\x94\x1c\x29\x6c\x2f\x31\xb5\x39\x71\x16\x03\x3d\x89\x56\xb3\x52\x1b\xfe\x15\x43\x2c\x82\x17\xf4\xbb\x78\x4b\x2c\x62\xa5\x96\x17\x1b\xf6\xda\x32\x23\xdf\x09\x61\xc0\x22\x7f\x13\xff\xee\x4c\x1f\x31\xb4\xbc\x0d\xd0\x31\x09\x92\xb4\x5f\x0f\xc9\xac\x0d\x2b\xa8\x77\x8f\xd6\xe1\x0c\x25\xf6\x58\x0c\xbc\x2d\xbe\x03\x7b\xd4\xef\x7b\xf5\x70\x87\xcb\x13\x60\xce\x9c\x6f\x6b\x9e\x37\x8f\xe2\xd5\x39\x2c\xc5\xc0\x13\xf3\x27\x5a\xd2\x6a\xfe\x0d\x00\x00\xff\xff\x9c\x43\xcd\x6c\x09\x02\x00\x00")

func schemaServiceGraphqlsBytes() ([]byte, error) {
	return bindataRead(
		_schemaServiceGraphqls,
		"schema/service.graphqls",
	)
}

func schemaServiceGraphqls() (*asset, error) {
	bytes, err := schemaServiceGraphqlsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/service.graphqls", size: 521, mode: os.FileMode(420), modTime: time.Unix(1539249420, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schema/node.graphqls":    schemaNodeGraphqls,
	"schema/schema.graphqls":  schemaSchemaGraphqls,
	"schema/service.graphqls": schemaServiceGraphqls,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"schema": &bintree{nil, map[string]*bintree{
		"node.graphqls":    &bintree{schemaNodeGraphqls, map[string]*bintree{}},
		"schema.graphqls":  &bintree{schemaSchemaGraphqls, map[string]*bintree{}},
		"service.graphqls": &bintree{schemaServiceGraphqls, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
