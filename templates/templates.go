package templates

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _assets_templates_index_tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x53\xcd\x8e\xdb\x20\x10\xbe\xe7\x29\x46\xb4\x92\x1d\xa9\x82\x6e\xaf\xb5\xf3\x00\x3d\x6c\x0f\xab\xaa\x67\x82\x67\x1d\xb2\x18\x28\x43\x76\x37\x8a\xf2\xee\x1d\xec\x28\x71\xbc\x2d\x27\x8f\xbf\x9f\xf9\x83\x66\x97\x07\xb7\x59\x01\x9f\x66\x1b\xba\xe3\xf4\x39\x86\x64\x92\x8d\x19\x28\x99\xb6\x52\x4a\xef\xf5\xbb\xec\x43\xe8\x1d\xea\x68\x49\x9a\x30\x8c\xff\x94\xb3\x5b\x52\xfb\x3f\x07\x4c\x47\xf5\x20\x1f\xbe\xca\x6f\x97\x48\x0e\xd6\xcb\x3d\x55\x9b\x46\x4d\x56\x33\x6f\x67\xfd\x0b\x24\x74\xad\xa0\x7c\x74\x48\x3b\xc4\x2c\x60\x97\xf0\xb9\x15\x4a\x13\x61\x26\xc5\x19\x86\xe0\xa5\x21\x12\x90\x8f\x11\x5b\x91\xf1\x3d\xab\x31\x1e\xb0\xb3\xba\x15\xda\x39\x01\x6a\xe6\xdb\xd9\x57\xb0\x5d\x2b\xde\x92\x8e\x11\x93\xb8\x41\x33\xd0\x04\x77\x18\x3c\xdd\x81\xa7\x13\x24\xed\x7b\x04\xf9\xa8\x07\x24\x38\x9f\x67\xe0\x55\x6d\x1c\xd7\xd6\x8a\x68\xbd\xd8\x34\x76\xe8\xc7\xe1\x08\xd6\x4a\x16\x08\xc5\x9d\x32\x6b\x61\x8b\xbe\x63\x10\x60\x35\x2f\xe5\x9e\x77\x09\x97\xa3\xbf\x11\x5e\x75\x02\x03\x2d\x78\x7c\x83\xdf\xb8\x7d\x0a\xe6\x05\x73\x5d\x95\xc4\xd7\xf0\x57\x72\x9c\xa7\x5a\x7f\xbf\xf9\x18\x19\x7c\x88\xe8\x59\xf9\x7c\xf0\x26\xdb\xe0\xeb\xf5\xe9\xae\xb1\x42\xe1\x7e\x49\x73\xeb\x33\x56\x42\x8a\xc1\x13\x2e\xd8\xe5\xec\x29\x14\xc3\xcf\x32\xea\x44\xf8\xe3\xe9\xe7\xe3\x95\x2d\x3b\x9d\x35\x17\xb0\x94\x18\x06\x83\x43\xe9\x42\x5f\x17\xf9\xfa\x03\xa3\x34\x58\xe6\xc9\xbe\x75\x55\x26\xbb\xa9\xd6\x52\xe7\x9c\xea\x8a\x47\x5c\xc1\x97\x31\xad\xf4\xbc\x9b\x7f\xf8\x17\x75\x59\xcf\xa4\x2e\xb3\xbc\xaa\xc7\x8d\x15\xfd\xb8\xb4\xff\x68\x19\x62\x2d\xeb\x64\xb9\x35\xbe\xab\xb9\x80\xf5\xea\x03\x95\xcd\x3f\x5d\xee\x0e\xfb\xc7\x84\x23\x97\xc5\x0b\xdb\xf3\x2d\x3c\xcf\x76\x3c\x5f\x6a\xa3\xa6\xc7\xd6\xa8\xe9\xfd\xad\xfe\x06\x00\x00\xff\xff\xa3\x49\x7a\x02\x88\x03\x00\x00")

func assets_templates_index_tmpl_bytes() ([]byte, error) {
	return bindata_read(
		_assets_templates_index_tmpl,
		"assets/templates/index.tmpl",
	)
}

func assets_templates_index_tmpl() (*asset, error) {
	bytes, err := assets_templates_index_tmpl_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "assets/templates/index.tmpl", size: 904, mode: os.FileMode(436), modTime: time.Unix(1419478396, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	"assets/templates/index.tmpl": assets_templates_index_tmpl,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"assets": &_bintree_t{nil, map[string]*_bintree_t{
		"templates": &_bintree_t{nil, map[string]*_bintree_t{
			"index.tmpl": &_bintree_t{assets_templates_index_tmpl, map[string]*_bintree_t{
			}},
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
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

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

