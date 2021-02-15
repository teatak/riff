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

var _schemaNodeGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\xb1\x6a\xc4\x30\x0c\x86\xf7\x3c\x85\xf2\x1a\x1e\x9b\xa9\x4b\x08\x78\x2c\x1d\x4c\xa2\x24\x86\x54\x32\x92\x1a\x28\xe5\xde\xfd\xc0\xce\x70\xf1\x71\xdc\x71\x9b\xf5\x59\xfc\xe2\xff\x1a\xfb\x4b\x08\x3d\x4f\x08\xff\x0d\x00\xc0\x14\x2c\x74\x48\x86\xe2\xc0\x9b\x44\x5a\xda\xcc\x63\xaa\x66\xf5\xb8\xcd\x0e\x3e\x98\x37\x0c\x54\x20\x85\x1f\x3c\xaf\xad\x66\x69\x60\x31\x07\x9f\x64\x05\x49\x1a\x2b\xa2\x28\x7b\x1c\x51\x1d\x7c\xf5\xa8\xe6\xcb\xd8\x7e\x97\x4f\x0a\xc9\xaf\x6c\xe7\x5c\xb5\x60\xf9\x54\x30\x2c\x64\x47\xd1\xc8\x74\xc4\x5e\x9a\xa3\x19\xaa\xdd\xb4\x1b\x99\xe6\xb8\xdc\x47\xfd\x6a\xc7\x64\x48\xd5\x95\x87\x32\x5e\x2b\x5f\x3b\x4b\x4f\x45\xbc\xd5\x35\x47\x0b\x2f\x82\xaa\x0e\x86\xe3\x95\x1d\x5c\x03\x00\x00\xff\xff\x8c\xc2\xf4\x9d\xdf\x01\x00\x00")

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

	info := bindataFileInfo{name: "schema/node.graphqls", size: 479, mode: os.FileMode(420), modTime: time.Unix(1613382030, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaSchemaGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x92\xc1\x6e\x83\x30\x0c\x86\xef\x79\x0a\xf7\x56\xa4\x3d\x41\x6e\xb4\xe5\x50\xa9\xdd\x3a\xc2\x76\x99\x38\x44\x60\x20\xd2\x48\xba\x24\x54\xaa\x26\xde\x7d\x4a\x48\x18\x74\xb7\x9d\xb0\x3f\xfb\x37\xf6\x0f\xa6\xea\xb0\xe7\xf0\x4d\x00\x00\xbe\x06\xd4\x77\x0a\xaf\xee\xe1\x41\x3f\x58\x6e\x85\x92\x14\xce\x21\x22\x23\x21\xf6\x7e\xc5\xa9\x29\xe8\xa4\xaa\x71\x2b\x79\x8f\x14\x98\xd5\x42\xb6\x9b\x84\xc2\xb3\xaa\x71\xae\x1a\x0a\x1f\x0e\x94\x9e\x68\xd1\x34\x14\x72\xd1\x34\x3e\x35\xa8\x6f\xa8\x17\x0a\x07\x44\xf5\x30\xf2\xc9\x58\x6e\x7d\xca\x2d\x26\x14\xd8\xd4\xb4\x14\xb8\xb7\x04\x5c\xce\x8b\xc6\xcd\xc3\xae\xf1\xa4\xd0\xb7\x5d\x48\xcf\xeb\xd2\x51\x5e\x07\xbb\x29\x93\xbf\x95\x70\x05\xb6\xc2\x58\x8c\x83\x3c\x8b\xf7\xfe\x1a\xe6\x8e\x9a\x26\xcd\x0d\x95\x92\x8d\x68\xe7\xc3\x3c\x4f\x28\xec\x94\xfa\x44\x2e\x7d\x3a\xc8\x7f\x4f\x5f\x99\xf6\x38\x7b\x24\x04\xe5\xd0\x4f\x26\x06\x47\xd2\xd3\xf1\x3d\xf3\x11\x7b\x63\x97\x6c\x5f\xf8\xf8\x90\xa5\x87\x50\x3e\xcd\xb2\x7d\x5f\x07\x11\x2b\xd2\xbc\x08\xd1\xcb\xc5\x07\x79\x36\xc1\x68\xbc\xfb\xbe\xa1\xbb\x15\x76\xa7\xb9\xac\xba\xf5\x62\xad\xb0\xac\xe3\x6b\x76\x43\x6d\xfc\x1f\x17\xe1\x48\x7e\x02\x00\x00\xff\xff\xec\xd3\xcb\x82\xa4\x02\x00\x00")

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

	info := bindataFileInfo{name: "schema/schema.graphqls", size: 676, mode: os.FileMode(420), modTime: time.Unix(1613315642, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaServiceGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\xc1\x6a\xc4\x30\x0c\x44\xef\xfe\x0a\xe5\x37\x7c\xec\x9e\xf6\xd0\x65\x61\x8f\xa5\x07\x63\xab\xc1\xb0\x91\x8c\x2c\x17\x4a\xe9\xbf\x97\x75\xe3\xa4\x4e\x4a\xa0\xb7\x68\xe2\x19\x3f\x8f\x22\xa5\xa2\xf0\x5c\xd4\x69\x64\xba\xa1\xbc\x47\x8f\xe7\x2a\x7e\x1a\x00\x00\x3f\x05\x0b\xa7\x29\x0c\x75\x22\x37\xa1\x85\x9b\x4a\xa4\xf1\x47\x89\xa9\x9f\x13\x8b\x5a\x38\x93\x0e\xe6\xcb\x98\x3e\xfe\xc2\xa1\xcb\x3e\xf6\xea\x47\x42\xb8\x60\xd6\x99\xaa\x01\x31\xbd\xc5\xb1\x37\x66\x75\x5a\xf2\x89\x49\x91\xb4\xff\xf5\x2f\xe2\x96\x55\x1d\x4e\x71\x3e\x20\x3c\x0a\xe6\x6c\xe1\x3a\x7f\xad\x7c\x4d\x69\x70\x45\xa4\x22\x2c\x71\xca\xea\xee\xbf\xe6\x48\xd7\x25\xee\x89\xf9\x8e\x8e\xd6\xb4\xfe\xa5\x7b\x74\xe2\x80\xd9\xc2\xcb\xa3\x95\x47\x99\xc3\xeb\x62\xdd\xac\xf0\xcf\xed\xa1\x08\xcb\x71\x19\xfb\x3b\xb7\xf5\x14\xef\x37\xf0\xdf\x01\x00\x00\xff\xff\xc3\x38\x9b\x89\x44\x02\x00\x00")

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

	info := bindataFileInfo{name: "schema/service.graphqls", size: 580, mode: os.FileMode(420), modTime: time.Unix(1613379092, 0)}
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
