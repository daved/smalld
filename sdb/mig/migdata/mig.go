package migdata

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

var _mig_sql_0000_setup_adminareas_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x52\x4b\x6f\xe2\x30\x10\xbe\xe7\x57\xcc\x8d\xa0\x85\xd5\xc2\xa2\x4a\x55\x4f\x26\x71\x05\x52\x30\x34\x38\xa5\x3d\x21\x2b\x71\x5d\x4b\x7e\xd0\xd8\xa5\xe5\xdf\xd7\x4d\x79\x55\x44\xf8\x30\xd2\xcc\xf7\x98\x91\x67\xfa\x7d\xf8\xa3\xa5\xa8\x99\xe7\x50\x6c\xa2\xf3\x74\xe9\x43\xd4\xdc\xf8\x31\x17\xd2\x44\x51\x92\x63\x44\x31\x2c\xf1\x43\x81\x49\x82\x81\x55\x5a\x1a\x56\x73\xe6\xd6\x42\x56\x6b\xc7\xdf\x22\x08\x6f\x49\x51\x4e\x61\x35\xa5\x13\x18\x34\x85\x29\x09\xca\x19\x26\x14\xc6\xcf\xfb\x12\x99\xc3\x6c\x4a\x1e\x51\x56\xe0\x63\x8e\x9e\x4e\x79\x82\x92\x09\x86\xc1\xdd\xb1\x29\x45\xe3\xec\xbc\x23\xc4\x0d\x2f\xf4\x05\x69\x3c\x17\xbc\x0e\x1e\x14\x48\x91\x65\x90\xe2\x7b\x54\x64\x14\x0c\xff\xf4\x5b\xa6\xe2\xce\xe5\xa0\x9d\x6e\xaf\xd1\x5b\xa7\xd7\xc1\xa2\x7c\x65\x35\x2b\x7d\x30\xd9\xb2\x7a\x27\x8d\x88\x87\xff\xf6\x0c\xc5\x9c\x0f\xb0\x11\xbc\x85\x35\xb8\xdd\xb3\x4a\x5b\x71\x70\x9a\x29\x15\xa6\xf9\x29\xbd\x94\x41\xea\x5a\x44\xa3\x83\xf5\xc6\x3a\xcf\x54\x23\x6d\xb1\x3e\xb0\x0c\xd3\xed\xf8\x81\x20\xb8\xd5\x4d\xe0\xbe\xde\xc5\xb3\x77\xe5\xe5\xc2\xaa\x9d\xb0\xa6\x37\xfa\x3f\xbc\xe9\x46\xdd\xf0\x8d\x28\xa3\x38\xbf\xb6\x3a\x98\xaf\x08\x4e\xbf\x37\x74\xc2\xfe\x06\x2c\x68\x5b\x6f\x02\x9b\xea\x37\x92\xda\x0f\x73\xf5\x7c\xd2\x7c\xbe\xb8\xd8\xe3\x35\xfb\xaf\x00\x00\x00\xff\xff\x68\x75\xda\xf9\x9b\x02\x00\x00")

func mig_sql_0000_setup_adminareas_sql_bytes() ([]byte, error) {
	return bindata_read(
		_mig_sql_0000_setup_adminareas_sql,
		"mig/sql/0000_setup_adminareas.sql",
	)
}

func mig_sql_0000_setup_adminareas_sql() (*asset, error) {
	bytes, err := mig_sql_0000_setup_adminareas_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "mig/sql/0000_setup_adminareas.sql", size: 667, mode: os.FileMode(436), modTime: time.Unix(1455337123, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _mig_sql_0001_setup_locations_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x51\x41\x4f\xf2\x40\x10\xbd\xf7\x57\xcc\x8d\x92\x0f\xbe\x04\x35\x5e\x3c\x2d\xed\x18\x48\xca\x82\x65\x2b\x7a\x22\xb5\x9d\xd4\x4d\xda\x5d\x2c\x03\xa8\xbf\xde\xb5\x90\x82\x44\xd9\xc3\x24\xef\xcd\xbc\x99\xd9\x37\xfd\x3e\xfc\xab\x74\x51\xa7\x4c\x90\xac\xbc\x53\x38\x67\x17\x2b\x32\x3c\xa4\x42\x1b\xcf\x0b\x62\x14\x0a\x61\x8e\x0f\x09\xca\x00\xa1\xb4\x59\xca\xda\x9a\xf5\x52\xe7\xcb\x35\xbd\x79\xe0\xde\x5c\x89\x58\xc1\x62\xac\x46\x30\x68\x88\xb1\x74\xba\x09\x4a\x05\xc3\xe7\x03\x25\xa7\x30\x19\xcb\x47\x11\x25\xd8\x62\xf1\x74\xc4\x81\x08\x46\x08\x83\xbb\x76\xa4\x12\xc3\xe8\x64\x1e\xf8\x4d\x99\xce\x41\x1b\xa6\x82\x6a\xd7\x41\x81\x4c\xa2\x08\x42\xbc\x17\x49\xa4\xc0\xd0\x3b\x6f\xd3\xd2\xef\x9c\x2f\xd9\xe9\xf6\x1a\x71\x99\xbe\x50\x09\xec\xca\xf6\x38\xcd\x32\x30\x9b\x8a\x6a\x9d\xed\x89\x82\x6c\xd5\x04\xe2\xfa\xc3\x9f\x59\x37\xaa\x77\x73\x7d\x75\x7b\xd0\xd7\x94\x91\xde\x52\x0e\xac\x2b\x5a\x73\x5a\xad\x60\xa7\xf9\xd5\x6e\xb8\x61\xe0\xd3\x1a\x3a\x6e\x63\x77\x7e\xd7\xeb\xba\x0f\x89\x48\x61\xfc\xb7\x85\x30\x5d\x48\x0c\xbf\x9d\x6a\x33\xff\x75\xee\x74\xbf\xde\x05\x4d\xfe\x33\x13\xda\x9d\xb9\x78\xc2\x30\x9e\xce\xce\xdd\xbc\xd4\xfd\x2b\x00\x00\xff\xff\xa2\x30\x2c\xbc\x1e\x02\x00\x00")

func mig_sql_0001_setup_locations_sql_bytes() ([]byte, error) {
	return bindata_read(
		_mig_sql_0001_setup_locations_sql,
		"mig/sql/0001_setup_locations.sql",
	)
}

func mig_sql_0001_setup_locations_sql() (*asset, error) {
	bytes, err := mig_sql_0001_setup_locations_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "mig/sql/0001_setup_locations.sql", size: 542, mode: os.FileMode(436), modTime: time.Unix(1455337113, 0)}
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

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
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
	"mig/sql/0000_setup_adminareas.sql": mig_sql_0000_setup_adminareas_sql,
	"mig/sql/0001_setup_locations.sql": mig_sql_0001_setup_locations_sql,
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
	"mig": &_bintree_t{nil, map[string]*_bintree_t{
		"sql": &_bintree_t{nil, map[string]*_bintree_t{
			"0000_setup_adminareas.sql": &_bintree_t{mig_sql_0000_setup_adminareas_sql, map[string]*_bintree_t{
			}},
			"0001_setup_locations.sql": &_bintree_t{mig_sql_0001_setup_locations_sql, map[string]*_bintree_t{
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

