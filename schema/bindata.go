// Code generated for package schema by go-bindata DO NOT EDIT. (@generated)
// sources:
// schema/node.graphqls
// schema/schema.graphqls
// schema/service.graphqls
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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaNodeGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\x41\x4b\x04\x31\x0c\x85\xef\xf3\x2b\xb2\x7f\xa3\x47\xf7\xe4\x65\x59\xa8\x37\xf1\x50\x66\xb3\x33\x85\x31\x29\xc9\x73\x40\xc4\xff\x2e\x4e\x8b\x38\x15\x51\xf6\xd6\x7c\x0d\xef\xe5\xbd\x01\xaf\x85\xe9\xa4\x17\xa6\xb7\x81\x88\xe8\x92\x90\x8e\x2c\x60\x0b\x14\x61\x59\xa6\xc3\xc6\x73\xe9\x66\x8f\xbc\x5c\x03\xdd\xa9\x2e\x9c\xa4\x42\x49\xcf\xbc\x5f\x9b\x81\x72\x56\x43\xa0\x7b\x41\x45\x56\xc6\x8e\x38\xdb\x9a\x47\xf6\x40\x8f\x27\x76\xc4\x3a\x1e\x9e\xea\xa7\xa4\x12\x67\xc5\x5e\xd7\x91\xb0\x59\x25\x70\x25\x2b\x9b\x67\x95\x26\xfb\x3e\xb4\x64\xec\xf8\x96\x6e\x54\xb9\xe6\xe9\xa7\xd4\x8b\x1f\x55\xc0\xd2\xb9\xfc\x5a\xc6\xff\xc2\xf7\x9d\x95\x3f\x8b\xb8\x29\x6b\xdb\x31\x3c\xe4\xcf\x0b\xbe\x58\x31\x9d\x8c\xdd\x03\x9d\xdb\x6b\xeb\xe5\x23\x00\x00\xff\xff\xd4\xfb\x0d\x9c\xf3\x01\x00\x00")

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

	info := bindataFileInfo{name: "schema/node.graphqls", size: 499, mode: os.FileMode(420), modTime: time.Unix(1615649436, 0)}
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

	info := bindataFileInfo{name: "schema/schema.graphqls", size: 676, mode: os.FileMode(420), modTime: time.Unix(1615649436, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaServiceGraphqls = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\xc1\x6a\xc3\x30\x10\x44\xef\xfe\x8a\xf5\x6f\xe8\xd8\x9c\x72\x68\x08\xa4\xb7\xd2\x83\xb0\xb7\x46\x10\xef\x8a\xd5\xa8\x50\x4a\xff\xbd\xd8\x8d\x94\xca\x29\x81\xdc\xbc\x63\xcf\xf8\xed\x6c\x90\x98\x41\xcf\x19\x1e\x41\xe5\xc4\xf6\x11\x06\xde\xaf\xe2\x57\x47\x44\x34\xcc\xa3\xa3\xdd\x3c\xf6\xeb\x24\x7e\x66\x47\x27\x58\x90\xe9\x57\x09\xb1\x9d\xa3\x1a\x1c\xed\x05\x7d\xf7\xdd\x75\x6d\xfc\x41\xc7\x26\xfb\xbe\x17\x9f\x91\xe9\xc0\x09\x17\xaa\x02\xa4\xf2\x1e\xa6\xd6\x98\xe0\x91\xd3\x4e\x05\x2c\x68\x5f\x3d\x44\x5c\xb2\x56\x87\x07\x57\xc5\xf0\x12\x96\x9c\xfa\x55\x34\x9d\x8c\x53\x72\x74\xbc\x3c\x5d\x99\x8b\x52\x80\xb3\xd9\x8a\x55\xcd\x50\xf8\xf3\x9f\x39\xc8\xb1\xc6\x3d\xa9\x9e\xd9\xcb\x35\xad\xdd\xfe\x76\x1d\xd1\x91\x93\xa3\xd7\xa5\xa9\xa5\xe0\xfe\xad\x5a\x37\x67\xfd\xf7\xa2\x6c\xa6\x76\xbf\xa0\xdb\x7f\x6e\x2b\xcb\xc3\xb0\x81\xff\x09\x00\x00\xff\xff\x0a\x5f\x7f\x40\x58\x02\x00\x00")

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

	info := bindataFileInfo{name: "schema/service.graphqls", size: 600, mode: os.FileMode(420), modTime: time.Unix(1615649436, 0)}
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
