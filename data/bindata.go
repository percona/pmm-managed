// Code generated by go-bindata. DO NOT EDIT.
// sources:
// iatemplates/mongodb_connections_memory_usage.yml (789B)
// iatemplates/mongodb_high_memory_usage.yml (723B)
// iatemplates/mongodb_restarted.yml (683B)
// iatemplates/mysql_down.yml (541B)
// iatemplates/mysql_restarted.yml (731B)
// iatemplates/mysql_too_many_connections.yml (905B)
// iatemplates/node_high_cpu_load.yml (807B)
// iatemplates/node_low_free_memory.yml (805B)
// iatemplates/node_swap_filled_up.yml (819B)
// iatemplates/postgresql_down.yml (558B)
// iatemplates/postgresql_restarted.yml (753B)
// iatemplates/postgresql_too_many_connections.yml (890B)

package data

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

var _iatemplatesMongodb_connections_memory_usageYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xc1\x8e\x9b\x30\x10\xbd\xf3\x15\x4f\x51\x23\x25\x55\x49\xb3\xab\xee\xc5\x52\x2b\xb5\xea\x35\xb7\xde\xa2\x08\x19\x18\xc0\x12\xb6\xd1\xd8\xd0\x45\x29\xff\x5e\xd9\x0b\x88\xb4\x5a\x2e\xd8\x8f\x99\xc7\x9b\xf7\x26\x4d\xd3\xc4\x93\xee\x5a\xe9\xc9\x89\x04\x48\x61\xa4\x26\x81\x4e\xeb\x4c\x5b\x53\xdb\x32\xcf\x0a\x6b\x0c\x15\x5e\x59\xe3\x32\x4d\xda\xf2\x98\xf5\x4e\xd6\x94\x00\xc0\x40\xec\x94\x35\x02\x4f\xf1\xea\x7a\xad\x25\x8f\x02\x97\x58\x88\xde\x51\x89\x7c\xc4\x25\x70\xfd\xfc\x81\x0d\x57\xac\xa7\xd7\x8e\x05\xfe\xa4\xf1\x12\xdb\x43\xf5\xc1\xd8\x92\xb2\xa0\xe4\x88\xc3\x22\xc3\xb9\xad\x92\x7b\x38\x67\x7e\xec\xe8\xeb\xae\xe8\x99\xc9\xf8\xdd\x74\xc4\x47\x3c\x9d\x9f\xbf\xcc\xaf\x99\xf4\x33\xac\x79\xa4\x8c\xe7\x79\x94\x0b\xe9\x5f\xd6\xcb\x36\xcb\x47\x4f\xee\x38\xf7\x04\x82\xf3\x7c\xfe\x86\xeb\x15\x27\xdf\x30\xb9\xc6\xb6\x25\x6e\xb7\xf8\xa1\x93\x2c\x75\x34\x2d\x3c\x8b\x71\x6b\xd9\x8c\x6f\x2c\xf9\x8e\x8e\xb8\x20\xe3\x65\x4d\xa8\xd8\xea\xe0\x46\xa5\xea\x9e\xa9\x84\x96\xaf\x4a\xf7\x7a\xed\xea\x8d\xf2\x02\xbb\xfd\x6e\x45\xc2\xac\x02\x55\x6b\xa5\x5f\x31\x96\xa6\x26\x81\xeb\xf9\x53\xd0\x7b\x5b\xf1\x41\xb6\x3d\x09\x3c\xbf\x44\xa4\xb2\x2c\xf0\xf2\xc6\xed\x68\x20\x56\x7e\x14\xf8\x2d\xd9\x28\x53\x47\x54\x1a\x63\xbd\x8c\xc6\x8a\xe4\x1f\xd9\x4b\x74\x8d\xaa\x1b\xe8\x25\xd6\x30\x43\x3e\x6e\xf3\xc4\xe1\x7e\xc7\x87\x56\xe6\xd4\xba\x93\x23\x1e\x54\xf1\x66\x38\xa6\x69\xb1\xb5\x24\x57\xb0\xea\x7c\xdc\x98\x35\x75\x20\x74\x46\xcd\x98\xa6\x3d\x6c\xb5\xfc\xe7\xa0\x2d\x13\x7c\x23\xcd\x7f\x19\xec\x8f\x50\x2e\xee\xd7\xca\x92\x8f\x78\x5f\xc2\x83\x54\x6b\xb6\x95\xeb\x6a\x04\xa5\xa7\xe4\x6f\x00\x00\x00\xff\xff\xf7\x4a\xff\x2e\x15\x03\x00\x00")

func iatemplatesMongodb_connections_memory_usageYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesMongodb_connections_memory_usageYml,
		"iatemplates/mongodb_connections_memory_usage.yml",
	)
}

func iatemplatesMongodb_connections_memory_usageYml() (*asset, error) {
	bytes, err := iatemplatesMongodb_connections_memory_usageYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/mongodb_connections_memory_usage.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xa2, 0x60, 0xfb, 0xda, 0x65, 0xcf, 0x88, 0x87, 0x88, 0xf3, 0x74, 0x1b, 0x4c, 0x9d, 0x51, 0x1c, 0xc, 0x66, 0x7d, 0x78, 0x3f, 0xd0, 0x88, 0x23, 0x53, 0x36, 0xc6, 0xe4, 0x29, 0xf8, 0xf7, 0x27}}
	return a, nil
}

var _iatemplatesMongodb_high_memory_usageYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xc1\xaa\xdb\x30\x10\xbc\xe7\x2b\x86\x47\x03\x49\xa9\xd3\xbc\xd2\x42\xd1\xa1\xd0\xd2\x6b\x6e\xbd\x85\x60\xe4\x78\x6d\x0b\xbc\x92\x59\xc9\xe9\x33\xaf\xf9\xf7\x22\xd9\x56\x53\x4a\x7d\xb1\xb4\xda\x19\xcd\xce\xa8\x28\x8a\x4d\x20\x1e\x7a\x1d\xc8\xab\x0d\x50\xc0\x6a\x26\x85\x81\xb9\x64\x67\x5b\x57\x57\x65\x67\xda\xae\x64\x62\x27\x53\x39\x7a\xdd\xd2\x06\x00\x6e\x24\xde\x38\xab\xf0\x9c\xb6\x7e\x64\xd6\x32\x29\x9c\x52\x23\x46\x4f\x35\xaa\x09\xa7\x48\xf2\xfd\x5b\xea\xa1\x97\x41\x14\x7e\x15\x69\x93\x20\xb1\x63\x67\x5d\x4d\x65\xbc\x76\x8f\xdd\x7a\xa7\xf7\xf1\xc6\x52\xc8\x9b\x9a\x6c\xc0\x5b\x3c\x1f\x3f\x7c\x5c\x7e\xfb\x85\xe0\x3d\x9c\xfd\x1b\x9e\xd6\x8b\xd4\x13\xf1\x0f\x17\x74\x5f\x56\x53\x20\xbf\x62\x22\xc3\x71\x59\x7f\xc1\xf9\x8c\x43\xe8\x84\x7c\xe7\xfa\x1a\x97\x4b\x3a\x18\xb4\x68\x4e\x6e\xc4\x6f\x75\x24\xb7\x2d\xf5\x87\x91\xbf\x62\x20\xb9\x92\x0d\xba\x25\x34\xe2\x18\x57\x67\x1b\xd3\x8e\x42\x35\x58\xbf\x18\x1e\x39\xa3\x46\x6b\x82\xc2\xd3\xf6\x29\x57\xc2\x34\x90\x42\xd3\x3b\x1d\x72\x4d\xb4\x6d\x49\xe1\x7c\x7c\x17\xf5\x5e\x72\xfd\xa6\xfb\x91\x14\x3e\xcf\x23\x34\x4e\x14\x3e\xcd\xdc\x9e\x6e\x24\x26\x4c\x0a\x3f\xb5\x58\x63\xdb\x54\xd5\xd6\xba\xa0\x83\x71\x36\x0f\xf4\x27\xa9\x39\x1a\xc4\x7c\xc1\x6b\x6c\x71\x86\xdd\xeb\x2b\xde\xf4\xba\xa2\xde\x1f\x3c\xc9\xcd\x5c\x67\x87\x71\xbf\xaf\x3e\xd6\xe4\xaf\x62\x86\x90\x9e\x40\x8e\x14\x88\xc8\x24\x12\xf7\xfb\x16\xae\x59\x89\x77\xec\x84\x10\x3a\x6d\xff\x31\x7d\xbb\x87\xf1\xe9\xc1\x64\x96\x6a\xc2\xff\x25\xc4\xd4\x1f\x4e\x73\xfe\x51\xdd\x61\xf3\x3b\x00\x00\xff\xff\x90\xbd\xde\xac\xd3\x02\x00\x00")

func iatemplatesMongodb_high_memory_usageYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesMongodb_high_memory_usageYml,
		"iatemplates/mongodb_high_memory_usage.yml",
	)
}

func iatemplatesMongodb_high_memory_usageYml() (*asset, error) {
	bytes, err := iatemplatesMongodb_high_memory_usageYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/mongodb_high_memory_usage.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xfa, 0x20, 0xc1, 0xb8, 0xfd, 0xab, 0xd5, 0xcf, 0x8d, 0xcf, 0xee, 0x8e, 0x95, 0x6e, 0x39, 0x95, 0xef, 0x42, 0x5a, 0x59, 0xb8, 0xc, 0xdf, 0x9d, 0x59, 0x17, 0x84, 0xc9, 0xae, 0x56, 0xd8, 0x7}}
	return a, nil
}

var _iatemplatesMongodb_restartedYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x91\xcd\x8e\xd3\x40\x10\x84\xef\x7e\x8a\x3a\x70\x00\x09\x3b\x8e\x56\x08\x34\xe2\xc2\x0a\x6e\x1b\x81\xb4\xdc\xa2\xc8\xea\xd8\x9d\x78\x90\xe7\x47\xd3\x9d\x84\x68\xc9\xbb\x23\x4f\x62\xef\x2e\x27\xe6\xd6\xdd\xa5\x4f\x35\x55\x65\x59\x16\xca\x2e\x0e\xa4\x2c\xa6\x00\x4a\x78\x72\x6c\x10\x9d\x6b\x5c\xf0\xfb\xd0\x6d\x9b\xc4\xa2\x94\x94\xbb\x02\x00\x8e\x9c\xc4\x06\x6f\xb0\xcc\xa3\x1c\x9c\xa3\x74\x36\x58\x8d\xea\xaf\xf7\x78\xad\xe6\xdf\x31\x19\xfc\x29\xf3\x00\x4c\x48\xeb\x45\xc9\xb7\xdc\x1c\xa2\x5a\xc7\x8d\x70\x1b\x7c\x27\x37\xd5\x67\xac\xd7\xa8\xb4\x4f\x2c\x7d\x18\x3a\x6c\x36\xf9\x10\x29\x91\xcb\x2e\xc7\x37\x39\x9d\x65\xb7\xfd\x0b\x4b\x8f\x57\x2a\x76\x29\x38\x68\xcf\x18\x48\x74\x32\x38\xcb\x0f\xde\xaa\x81\xcc\xb3\x9e\x23\x1b\xec\x86\x40\xcf\x9a\x44\x7e\xcf\x06\xeb\xfa\x3d\x96\x9f\xea\x7a\x33\x1f\x8e\x34\x1c\xd8\xe0\xae\xae\xf3\x6a\x17\x92\xc1\xb2\xbe\xc2\x84\x8f\x9c\xac\x9e\x0d\x4e\x94\xbc\xf5\xfb\xbc\x25\xef\x83\x92\xda\xe0\xe7\xaf\x74\x2c\x6d\xb2\x51\x73\xac\x73\x56\xc0\xcf\xde\x0a\xa6\x7a\x60\x05\x3e\x28\x76\xd6\x5b\xe9\xb9\x33\xe8\x55\xa3\x98\xc5\xe2\x97\x4d\x54\x45\x4e\x6d\xf0\x54\xb5\xc1\x2d\xb6\x29\x9c\x84\x17\x3f\x56\xab\xf2\xe3\xf2\xee\xc3\xcc\x9b\x2a\x9a\xd2\xc7\xd3\x13\xde\x0c\xb4\xe5\x41\x2a\xe1\x74\xb4\x2d\x37\x63\xa8\xb8\x5c\x70\x22\x79\xae\x32\x0b\xf3\x57\xc7\xd3\xad\x2c\xd0\x3e\x54\x33\xfb\xe1\xcb\xfd\xb7\x87\x47\xf3\x02\x89\xcb\xa5\xf8\xa7\x90\xd5\xf7\xab\x81\xd7\xf0\xb7\xff\xe1\xe7\x5d\xf1\x37\x00\x00\xff\xff\xf5\x30\x77\xf4\xab\x02\x00\x00")

func iatemplatesMongodb_restartedYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesMongodb_restartedYml,
		"iatemplates/mongodb_restarted.yml",
	)
}

func iatemplatesMongodb_restartedYml() (*asset, error) {
	bytes, err := iatemplatesMongodb_restartedYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/mongodb_restarted.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xdd, 0x5a, 0xfb, 0xc9, 0x88, 0xf6, 0x85, 0x33, 0x43, 0xa2, 0x8e, 0x39, 0x5, 0xd2, 0x52, 0x74, 0x17, 0x65, 0x56, 0xf7, 0x2f, 0x35, 0x5a, 0x19, 0x18, 0x32, 0xb7, 0x3e, 0x8d, 0xe6, 0xeb, 0xf5}}
	return a, nil
}

var _iatemplatesMysql_downYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x51\x4b\x6f\xda\x40\x10\xbe\xfb\x57\x7c\x42\x3d\x80\x54\xdb\x45\x15\xaa\xb4\x12\x87\x56\xed\x0d\xda\x46\x44\xca\xd1\x5a\xd6\x43\xbc\x89\xf7\x91\xdd\x01\x62\x11\xff\xf7\x88\x05\x13\xc4\x29\xb7\xdd\x79\x7c\xaf\xc9\xf3\x3c\x63\x32\xbe\x95\x4c\x51\x64\x40\x0e\x2b\x0d\x09\x78\x63\x2a\xd3\xc5\x97\xb6\xaa\xdd\xde\x66\x00\xb0\xa3\x10\xb5\xb3\x02\xd3\xf4\x8d\x5b\x63\x64\xe8\x04\x96\xdd\xea\x6e\x81\xdf\xff\x1e\xfe\xa6\x3a\xbd\xfa\x20\xf0\x96\xa7\x4f\x1a\xc3\xba\xc3\x38\x52\xd8\x69\x45\xd5\x11\xfe\x2b\xac\xab\x87\xe7\xd0\xe0\xce\xd3\x04\xe3\x13\xe9\xd6\x4f\x30\x9f\xe3\x5b\x02\xf1\x32\x48\x93\xd4\x01\x1b\x17\x04\x66\xf1\xa4\x80\x76\x14\x34\x77\x02\x2a\x68\xd6\x4a\xb6\xa9\xdc\xca\x35\xb5\xe7\x71\x80\xb5\x7a\x26\x16\xf8\xbf\x5c\xe6\x3f\xa6\xdf\x67\x43\xb9\xd1\xf6\x31\x56\xec\xaa\xda\x09\x8c\x3a\x8a\xa3\xd4\x91\xd6\x3a\x96\xac\x9d\xbd\x20\xd4\x14\x55\xd0\x9e\x93\xf7\x8b\x2f\xe0\xbe\xd1\x11\x43\x78\xd0\x11\xd6\x31\x36\xda\xea\xd8\x50\x2d\xd0\x30\xfb\x28\xca\xf2\x49\x07\x59\x78\x0a\xca\x59\x59\x28\x67\xca\x75\x70\xfb\x48\xe5\x8d\x20\x9c\x73\xd4\x36\xb2\xb4\x8a\x70\x38\xe0\xcb\xc9\x4a\x71\x9d\x1d\xfa\xfe\xc8\x75\x8c\xbb\xb8\xac\x2e\x7e\xfe\xfa\xb3\x58\x89\xab\x1d\xf4\xfd\xc7\x01\x6e\xef\x84\xf1\x27\x58\x26\xd9\x7b\x00\x00\x00\xff\xff\x15\xee\xc0\xf4\x1d\x02\x00\x00")

func iatemplatesMysql_downYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesMysql_downYml,
		"iatemplates/mysql_down.yml",
	)
}

func iatemplatesMysql_downYml() (*asset, error) {
	bytes, err := iatemplatesMysql_downYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/mysql_down.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x62, 0xa3, 0x8c, 0xd2, 0xb2, 0x4, 0x58, 0x9b, 0x84, 0xf, 0xfd, 0x1e, 0x35, 0x5d, 0x7b, 0xbb, 0xe, 0xff, 0x1a, 0x15, 0xc4, 0xae, 0x7d, 0xa1, 0xeb, 0xa1, 0x98, 0x39, 0x4d, 0x2a, 0x13, 0x63}}
	return a, nil
}

var _iatemplatesMysql_restartedYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\xcb\x8e\xd3\x40\x10\xbc\xe7\x2b\x4a\x2b\x0e\x20\xe1\xc4\xd1\x0a\x81\x46\x5c\x40\xe2\x96\x48\x40\xb8\x45\x91\xd5\xb1\x3b\xf1\xc0\x3c\xcc\x74\x27\xc1\x5a\xf2\xef\x28\x13\xc7\xec\xee\x69\x7d\xeb\xaa\x52\xb9\xba\x6b\x8a\xa2\x98\x28\xfb\xce\x91\xb2\x98\x09\x50\x20\x90\x67\x83\xce\xfb\xca\xf7\xf2\xdb\x55\x89\x45\x29\x29\x37\x13\x00\x38\x72\x12\x1b\x83\xc1\x3c\x8f\x72\xf0\x9e\x52\x6f\xb0\xec\x57\xdf\x16\xf8\xfe\x44\xcb\x7f\xba\x64\xf0\xb7\xc8\x03\x70\xb5\xdb\xbb\xb8\x25\x57\x89\x92\x1e\xa4\x3a\x74\x6a\x3d\x0f\x82\x8f\x58\xaf\x31\xd5\x36\xb1\xb4\xd1\x35\xd8\x6c\x32\xd1\x51\x22\x9f\xc3\x5d\xbe\x5b\xc0\x51\x36\xe0\x8f\xb2\xac\xb8\x8e\xa1\x11\xec\x52\xf4\xd0\x96\xe1\x48\x14\xc3\x1e\xa3\xfc\x10\xac\x1a\xc8\x38\x6b\xdf\xb1\xc1\xce\x45\xfa\xaf\x49\x14\xf6\x6c\xb0\x2e\xdf\x62\xfe\xa1\x2c\x37\x23\x71\x24\x77\x60\x83\xfb\xb2\xcc\xd0\x2e\x26\x83\x79\x79\x35\x13\x3e\x72\xb2\xda\x1b\x9c\x28\x05\x1b\xf6\x19\x75\xb4\x65\x37\x6e\xa1\xb6\xfe\xc5\x6a\xf0\x75\xb9\x2c\xde\xcf\xef\xdf\xdd\xe0\xd6\x86\xbd\x54\x1a\xab\x26\x1a\xdc\xf5\x2c\x77\x99\xa1\x10\xa2\x92\xda\x18\x46\x87\x86\xa5\x4e\xb6\xd3\x5c\xc6\x78\x63\xe0\x47\x6b\x05\xb7\x4a\x61\x05\x21\x2a\x76\x36\x58\x69\xb9\x31\x68\x55\x3b\x31\xb3\xd9\x4f\x9b\x68\xda\x71\xaa\x63\xa0\x69\x1d\xfd\x6c\x9b\xe2\x49\x78\xf6\x2c\x10\x86\x62\x6d\x10\xa5\x50\x33\x1e\x1e\xf0\xea\xba\xca\x54\x38\x1d\x6d\xcd\xd5\xa5\x0f\x9c\xcf\x38\x91\x60\x7c\x2c\x59\x98\xaf\x74\xa1\x64\x68\x84\xf6\x71\x3a\x3a\x2f\x3e\x7d\xfe\xb2\x58\x99\x47\x96\x38\x9f\x27\xcf\xba\xbc\xfe\xfe\xa9\xf5\xeb\x17\xa4\x79\x33\xf9\x17\x00\x00\xff\xff\xe5\xfd\x3c\xbf\xdb\x02\x00\x00")

func iatemplatesMysql_restartedYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesMysql_restartedYml,
		"iatemplates/mysql_restarted.yml",
	)
}

func iatemplatesMysql_restartedYml() (*asset, error) {
	bytes, err := iatemplatesMysql_restartedYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/mysql_restarted.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x3c, 0xb1, 0xa3, 0xe8, 0x64, 0xe1, 0x54, 0x1b, 0x55, 0x94, 0xb7, 0x6a, 0x40, 0xad, 0x9c, 0x3f, 0x32, 0xca, 0x24, 0x64, 0xb9, 0xed, 0xcd, 0x57, 0x9c, 0x6b, 0x22, 0x91, 0x43, 0xb1, 0x1a, 0xc0}}
	return a, nil
}

var _iatemplatesMysql_too_many_connectionsYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x53\x51\x6b\xdb\x4c\x10\x7c\xf7\xaf\x18\xc2\x17\x48\x3e\x6a\x3b\xa1\x84\x96\x83\x16\x52\xc8\x9b\x0d\x2d\x69\xfb\x62\x8c\x58\x4b\x6b\xe9\x52\xdd\xae\x7a\xb7\x72\x22\xd2\xfc\xf7\x22\x39\x56\xdd\xb4\x54\x4f\xba\xb9\xd5\xee\xcc\xec\x68\x3a\x9d\x4e\x8c\x43\x53\x93\x71\x72\x13\x60\x0a\xa1\xc0\x0e\x4d\x08\x59\xe8\xd2\xf7\x3a\x33\xd5\x2c\x90\x74\x59\xae\x22\x9c\x9b\x57\x49\x13\x00\xd8\x71\x4c\x5e\xc5\xe1\x72\x38\xa6\x36\x04\x8a\x9d\xc3\xb2\xbb\xfd\xb4\xc0\x51\x35\xbc\xa0\x4d\x3c\x54\xf1\x43\x13\x1d\x7e\x4c\x87\x03\x10\xe8\x21\xd3\x1d\xc7\xcc\x7c\xe0\xb3\xfd\xc0\xb2\xd6\x0d\xd5\x59\x32\xb2\x36\x65\x56\x45\xa6\x22\x1d\xa6\x73\xb1\xba\x0a\xeb\x73\xcc\xe1\x4b\xd1\xe8\xa5\xc4\xd9\x9d\x6e\xce\x0f\xfd\x8e\x3b\xec\x28\x7a\xda\xd4\x9c\xb2\x7e\xcc\x4b\xfa\xc0\xff\xb8\xbc\xb8\x78\x7e\x7f\x8f\xd5\x0a\xb3\x7e\x58\xaa\xb4\x2e\xb0\x5e\x0f\x17\x0d\x45\x0a\x83\x31\xfd\x73\x30\x67\x2c\x7b\xc6\x8f\xc4\x5f\xa3\xe1\x98\xb3\x18\x95\x8c\x6d\xd4\xd0\x1b\xb1\xf5\x65\x1b\xb9\xe8\xd5\xfa\xd0\x86\xf1\xab\x56\xbc\x39\x9c\x9c\x9e\x8c\x88\x75\x0d\x3b\x6c\x6b\x25\x1b\xb1\x48\x52\xb2\xc3\xea\xe2\x55\xcf\x77\x3d\xe2\x3b\xaa\x5b\x76\x78\xbb\x97\xb0\xd5\xe8\x70\xb5\xef\x9d\x78\xc7\xd1\x5b\xe7\x70\x4f\x51\xbc\x94\x03\x5a\xd3\x86\xeb\x51\x8b\xf9\xfc\x1b\x9b\xc3\xc7\xe5\x72\xfa\xe6\xf2\xf5\xd5\x01\xae\xbc\x94\x29\x33\xcd\x0a\x75\x38\xe9\x38\xed\xc9\x91\x88\x1a\x0d\xee\x1d\x3a\x14\x9c\xf2\xe8\x1b\x1b\x32\x30\x6e\x14\xf8\x5c\xf9\x84\x43\xa8\xe0\x13\x44\x0d\x5b\x2f\x3e\x55\x5c\x38\x54\x66\x4d\x72\xf3\xf9\x9d\x8f\x34\xeb\xcd\x52\xa1\x59\xae\x61\xbe\x89\x7a\x9f\x78\xfe\x82\x10\xb0\xd4\xc8\xb0\x8a\xe4\x8f\x15\x9d\x42\xb7\x7f\x89\x1b\x45\x7e\x8e\x1c\x54\xf0\xf8\x88\xff\xf6\xd2\x67\x89\xe3\xce\xe7\x9c\xf5\x5b\xc4\xd3\xd3\x38\xe1\xeb\xf5\xe2\xcb\x0d\xde\x0d\xa5\x83\xab\xc7\x97\x8b\xeb\x0f\x37\x8b\x5b\x77\xd4\xe7\xd7\xed\x8b\xcc\x9b\x2a\xfa\x5f\xe5\x37\x36\x67\x5e\x92\x91\xe4\xfc\x0f\x26\xe7\x93\x9f\x01\x00\x00\xff\xff\x92\x9b\x9a\x24\x89\x03\x00\x00")

func iatemplatesMysql_too_many_connectionsYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesMysql_too_many_connectionsYml,
		"iatemplates/mysql_too_many_connections.yml",
	)
}

func iatemplatesMysql_too_many_connectionsYml() (*asset, error) {
	bytes, err := iatemplatesMysql_too_many_connectionsYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/mysql_too_many_connections.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x78, 0x2c, 0x9e, 0x2d, 0x70, 0xa3, 0x90, 0xf7, 0xd2, 0x58, 0xbe, 0x7e, 0xa7, 0xdf, 0xf7, 0xd8, 0xe7, 0x93, 0x1c, 0xd1, 0x11, 0x9f, 0xec, 0xb7, 0xdd, 0xee, 0x3e, 0x57, 0x64, 0xe4, 0x78, 0x68}}
	return a, nil
}

var _iatemplatesNode_high_cpu_loadYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x52\xdf\x8b\xd3\x40\x10\x7e\xef\x5f\xf1\x51\x3c\x68\xc5\xf4\x07\x52\x94\x85\x13\xaa\xdc\x5b\x4f\x0e\xf4\x7c\x29\x25\x4c\xb3\xd3\x64\x35\xbb\x1b\x76\x27\xbd\x2b\xb5\xff\xbb\x24\x6d\xd7\x43\xd1\x3e\x75\xbf\x99\x7c\x3f\x66\x26\xcb\xb2\x81\xb0\x6d\x6a\x12\x8e\x6a\x00\x64\x70\x64\x59\xa1\xb1\x36\x77\x5e\x73\x5e\x99\xb2\xca\x8b\xa6\xcd\x6b\x4f\x7a\x00\x00\x7b\x0e\xd1\x78\xa7\x30\xef\x9f\xb1\xb5\x96\xc2\x41\xe1\xb3\xd7\x8c\xae\x1d\x9f\x1e\x1e\x91\xda\xf9\xb9\x09\x0a\x3f\xb3\xfe\x01\x8c\xe6\xc8\x40\xfb\x12\xdb\xc3\xa8\x17\xe8\xf4\xc6\x18\x05\x12\x3e\x03\x9d\x58\xe4\xc2\x3b\x1d\x73\xf1\x42\xf5\xd1\x7a\xcd\xb7\x43\xa3\x6b\x1e\x9e\xd6\x0b\xbb\x19\x8f\xc7\x17\xb6\xd7\x98\xcf\x66\x97\xff\x1f\xb0\x5e\x63\x22\x55\xe0\x58\xf9\x5a\x63\xb3\xe9\x0b\x0d\x05\xb2\x7d\xb6\xee\x77\xcd\x97\xda\x2e\xf8\x8b\x1c\x4b\x34\x1c\x0a\x76\x42\x25\x63\x17\xbc\x45\xe1\xdd\xce\x94\x6d\x60\x0d\x4b\xcf\xc6\xb6\x36\x7d\xd5\x3a\x23\x0a\xc3\x9b\x61\x42\xe4\xd0\xb0\xc2\xae\xf6\x24\x09\x0b\xe4\x4a\x56\x58\xcf\xde\x74\x7e\x37\x09\xdf\x53\xdd\xb2\xc2\xfb\x73\x84\x9d\x0f\x0a\x8b\x33\x77\xe4\x3d\x07\x23\x07\x85\x27\x0a\xce\xb8\xb2\x47\x6b\xda\x72\x9d\xb2\x88\x29\x7e\xb0\x28\x3c\xdc\xdf\x67\xef\xe6\x6f\x17\x57\xb8\x32\xae\xec\x46\x97\x6b\xaf\x30\x3c\x70\x3c\x9b\x23\xe7\xbc\x90\x18\xef\x12\x83\xe6\x58\x04\xd3\x48\xbf\xce\xb4\x22\xe0\x6b\x65\x22\xae\x77\x01\x13\xe1\xbc\x60\x67\x9c\x89\x15\x6b\x85\x4a\xa4\x89\x6a\x3a\xfd\x6e\x02\x4d\xba\x61\x79\x47\x93\xc2\xdb\xe9\x36\xf8\xa7\xc8\xd3\x3f\x0c\x01\xc7\x23\x5e\x9d\xbd\x4f\xd2\xce\x71\x3a\xa5\x43\xe9\x24\xac\x0f\x0c\xa9\xc8\xfd\xb5\xc6\x9b\xc4\xf3\x6d\xb9\x7a\xbc\xc3\x6d\xcf\xd7\xcf\x0e\xa7\x53\x2a\xae\x96\x1f\xef\x56\x5f\xd4\x0b\xb1\xdf\xd5\xff\x1c\x29\x46\xc6\x45\x21\x57\xf0\xbf\x6c\x8e\x07\xbf\x02\x00\x00\xff\xff\x42\xc9\x41\x20\x27\x03\x00\x00")

func iatemplatesNode_high_cpu_loadYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesNode_high_cpu_loadYml,
		"iatemplates/node_high_cpu_load.yml",
	)
}

func iatemplatesNode_high_cpu_loadYml() (*asset, error) {
	bytes, err := iatemplatesNode_high_cpu_loadYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/node_high_cpu_load.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xc, 0x7f, 0x34, 0x92, 0x77, 0x38, 0x6a, 0x20, 0x5c, 0xb2, 0x9e, 0x4a, 0xba, 0xfb, 0xe6, 0x83, 0xad, 0x96, 0x4c, 0x15, 0x2f, 0x7, 0x56, 0x38, 0xf7, 0x32, 0xfb, 0x3c, 0xa3, 0x7f, 0xb8, 0x52}}
	return a, nil
}

var _iatemplatesNode_low_free_memoryYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x52\x5d\x8b\xdb\x30\x10\x7c\xcf\xaf\x18\x42\x0f\xee\x4a\x9d\x8f\x96\xa3\x20\xda\x87\x14\xee\x2d\x29\x85\x5e\xfb\x12\x82\x51\xec\xb5\xad\x56\xd2\x1a\x69\x9d\xd4\x5c\xf3\xdf\x8b\xed\xc4\x0d\x57\x7a\x7e\xb2\x66\x56\xa3\xd9\x9d\x4d\x92\x64\x22\xe4\x6a\xab\x85\xa2\x9a\x00\x09\xbc\x76\xa4\x50\x3b\x97\x7a\xce\x29\xb5\x7c\x4c\x8b\x40\x94\x3a\x72\x1c\xda\x09\x00\x1c\x28\x44\xc3\x5e\x61\xd9\x1f\x63\xe3\x9c\x0e\xad\xc2\x67\xce\x09\xdc\x08\xb8\xc0\x55\x39\xfd\xaa\x83\xc2\xef\xa4\x3f\x00\xbd\xec\x40\xa7\x1b\x72\xab\x83\x36\x56\xef\x2d\xa5\xfb\x56\x28\x62\xfe\xbc\xe0\x91\x45\xdb\x81\x3c\x2b\xbc\xc6\x72\xb1\x38\xff\x7f\xc0\x76\x8b\x99\x54\x81\x62\xc5\x36\xc7\x6e\xd7\x13\xb5\x0e\xda\xf5\x1d\x75\xdf\xa5\xab\xb1\xec\x8c\x5f\x79\x5f\xa1\xa6\x90\x91\x17\x5d\x12\x8a\xc0\x0e\x19\xfb\xc2\x94\x4d\xa0\x1c\xce\x78\xe3\x1a\x37\xde\x6a\xbc\x11\x85\xe9\xcd\x74\x44\xa4\xad\x49\xa1\xb0\xac\x65\xc4\x82\xf6\x25\x29\x6c\x17\x6f\x3a\xbf\xbb\x11\x3f\x68\xdb\x90\xc2\xdb\xa1\x85\x82\x83\xc2\xfd\xa0\x1d\xe9\x40\xc1\x48\xab\x70\xd4\xc1\x1b\x5f\xf6\xa8\xd5\x7b\xb2\x63\x2f\x62\xb2\x9f\x24\x0a\x5f\x36\x9b\xe4\xfd\xf2\xdd\xfd\x05\xae\x8c\x2f\x63\x2a\x9c\xe6\xac\x30\x6d\x29\x0e\xe6\xb4\xf7\x2c\x5a\x0c\xfb\x51\x21\xa7\x98\x05\x53\x4b\x1f\xe1\x18\x0b\xf0\x58\x99\x88\xcb\x36\xc0\x44\x78\x16\x14\xc6\x9b\x58\x51\xae\x50\x89\xd4\x51\xcd\xe7\x3f\x4c\xd0\xb3\x6e\x58\xec\xf5\x2c\x63\x37\xdf\x07\x3e\x46\x9a\x3f\x33\x04\x3c\x3d\xe1\xd5\xe0\x7d\xd6\x47\xda\x65\x80\xd3\xe9\xbc\x1a\xdd\x03\x85\xb1\xd6\xf8\x12\x4d\x8d\xdb\x7f\x83\xbc\x81\xa5\x42\xee\x46\xbd\xef\xab\xf5\xb7\x07\x7c\xec\x75\xfb\x19\xe2\x74\x1a\xc9\xf5\xea\xd3\xc3\xfa\xab\xba\x7a\xf4\x2f\xfb\xc2\x82\xe2\xd6\xf8\x28\xda\x67\xf4\x3f\xbb\x77\x93\x3f\x01\x00\x00\xff\xff\xc6\x72\xf8\xef\x25\x03\x00\x00")

func iatemplatesNode_low_free_memoryYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesNode_low_free_memoryYml,
		"iatemplates/node_low_free_memory.yml",
	)
}

func iatemplatesNode_low_free_memoryYml() (*asset, error) {
	bytes, err := iatemplatesNode_low_free_memoryYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/node_low_free_memory.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x5c, 0xe4, 0x80, 0xd0, 0x52, 0x92, 0xfb, 0x7, 0xc9, 0xc5, 0x29, 0xa0, 0x8c, 0xdf, 0x26, 0xcb, 0x59, 0xdb, 0xd, 0x97, 0x90, 0x3c, 0x97, 0x6c, 0x2b, 0x16, 0x2, 0x9a, 0x90, 0xbe, 0x78, 0x72}}
	return a, nil
}

var _iatemplatesNode_swap_filled_upYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\x5d\x8b\xd3\x40\x14\x7d\xef\xaf\x38\x14\x17\x5a\x31\xfd\x40\x16\x65\x40\xa1\xc2\xfa\xd4\x95\x85\xae\xfa\x50\x4a\x98\x26\x37\xc9\x68\xe6\x83\x99\x9b\x76\xc3\xda\xff\x2e\x49\xda\xb1\xae\x98\xa7\xcc\xb9\x73\xcf\x3d\x67\xce\x4d\x92\x64\xc4\xa4\x5d\x2d\x99\x82\x18\x01\x09\x8c\xd4\x24\xe0\xb4\x4e\x8d\xcd\x29\x0d\x47\xe9\xd2\x42\xd5\x35\xe5\x69\xe3\x46\x00\x70\x20\x1f\x94\x35\x02\xcb\xfe\x18\x1a\xad\xa5\x6f\x05\xbe\xd8\x9c\x50\xa9\xb2\xc2\xe6\xfb\xea\x01\x5d\x93\x32\x25\xce\x5d\xf4\xe4\xbc\xc0\xaf\xa4\x3f\x00\x93\x25\x12\x4c\xfa\x11\x9a\xb4\xf5\x6d\xba\x39\x4a\xf7\xd9\x13\xa5\xfb\x96\x29\x60\x8e\x97\xc5\x47\xcb\xb2\x1e\xaa\xd3\xe9\x99\xe6\x35\x96\x8b\xc5\xf9\xff\x23\xb6\x5b\xcc\xb8\xf2\x14\x2a\x5b\xe7\xd8\xed\xfa\x82\x93\x5e\xea\xde\x5c\xf7\x5d\x0c\xc6\x6b\x67\xfc\xca\xc7\x0a\x8e\x7c\x46\x86\x65\x49\x28\xbc\xd5\xc8\xac\x29\x54\xd9\x78\xca\xa1\xe5\x93\xd2\x8d\x8e\x5d\x8d\x51\x2c\x30\xbe\x19\x47\x84\x5b\x47\x02\x45\x6d\x25\x47\xcc\x4b\x53\x92\xc0\x76\xf1\xa6\xd3\xbb\x8b\xf8\x41\xd6\x0d\x09\xbc\x1f\x2c\x14\xd6\x0b\xdc\x0e\xdc\x81\x0e\xe4\x15\xb7\x02\x47\xe9\x8d\x32\x65\x8f\xd6\x72\x4f\x75\xf4\xc2\x2a\xfb\x49\x2c\xf0\x70\x7f\x9f\xbc\x5b\xbe\xbd\xbd\xc0\x95\x32\x65\x48\xd9\xa6\xb9\x15\x18\xb7\x14\x06\x71\xd2\x18\xcb\x92\x95\x35\x91\x21\xa7\x90\x79\xe5\xb8\x8f\x33\x66\x03\x3c\x56\x2a\xe0\xb2\x18\x50\x01\xc6\x32\x0a\x65\x54\xa8\x28\x17\xa8\x98\x5d\x10\xf3\xf9\x0f\xe5\xe5\xac\x7b\x2c\x6b\xe4\x2c\xb3\x7a\xbe\xf7\xf6\x18\x68\xfe\x42\x10\xf0\xfc\x8c\x57\x83\xf6\x59\x9f\x6a\x97\x01\x4e\x27\x74\xb1\x76\xf4\xc3\x82\xa1\x71\xd0\xd6\x13\xb8\x92\xe6\x9f\x30\x6f\x22\xdb\xb7\xd5\xfa\xeb\x1d\x3e\xf4\xac\xfd\x0b\xe2\x74\x8a\xc5\xf5\xea\xd3\xdd\x7a\x23\xae\x46\xfe\xa9\xfe\xbd\xaa\xe1\x6a\xf8\xb0\xa8\x98\x28\x13\x58\x9a\x8c\xfe\xa7\x78\x3a\xfa\x1d\x00\x00\xff\xff\xd2\x81\xea\x23\x33\x03\x00\x00")

func iatemplatesNode_swap_filled_upYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesNode_swap_filled_upYml,
		"iatemplates/node_swap_filled_up.yml",
	)
}

func iatemplatesNode_swap_filled_upYml() (*asset, error) {
	bytes, err := iatemplatesNode_swap_filled_upYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/node_swap_filled_up.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xeb, 0xae, 0xc3, 0x7e, 0xb4, 0x7e, 0x67, 0xda, 0xfd, 0xbf, 0xad, 0xc0, 0x78, 0x5a, 0x2d, 0xc1, 0xfb, 0x7b, 0x9f, 0xfc, 0x72, 0xfe, 0x50, 0x9e, 0x6b, 0x44, 0x61, 0xfc, 0x39, 0x20, 0x1d, 0xd4}}
	return a, nil
}

var _iatemplatesPostgresql_downYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x51\xcf\xab\xd3\x40\x10\xbe\xe7\xaf\xf8\x78\x78\x68\xc1\x24\x3e\xe4\x21\x2c\xf4\xa0\xe8\xad\xd5\x4a\x05\x8f\x61\xbb\x99\x36\xab\xd9\x1f\xee\x4c\x5b\x43\xcd\xff\x2e\xdd\x9a\x5a\xc4\xc3\xbb\xed\xce\x7c\x33\xdf\x8f\x29\xcb\xb2\x10\x72\xb1\xd7\x42\xac\x0a\xa0\x84\xd7\x8e\x14\xa2\x73\x4d\x0c\x2c\xfb\x44\xfc\xa3\x6f\xda\x70\xf2\x05\x00\x1c\x29\xb1\x0d\x5e\xe1\x31\x7f\xf9\xe0\x9c\x4e\x83\xc2\xfa\x8a\xdd\x7c\x5e\xe2\xfd\xa7\xaf\x1f\x73\x93\x7e\xc6\xa4\xf0\xab\xcc\x9f\x8c\xc5\x76\xc0\x8c\x29\x1d\xad\xa1\xe6\x42\xf4\x12\x3e\xb4\xd3\x73\x6a\xc8\x10\x69\x8e\x59\xdc\x37\x87\x38\xc7\x62\x81\x57\x79\x43\xd4\x49\xbb\x2c\x12\xd8\x85\xa4\xf0\xc4\x57\x0d\x74\xa4\x64\x65\x50\x30\xc9\x8a\x35\xba\xcf\xe5\x5e\x6f\xa9\xff\x03\x07\xc4\x9a\xef\x24\x0a\xeb\xd5\xaa\x7c\xf3\xf8\xfa\x69\x2a\x77\xd6\xef\xb9\x91\xd0\xb4\x41\xe1\x61\x20\x7e\xc8\x1d\xed\x7d\x10\x2d\x36\xf8\xdb\x86\x96\xd8\x24\x1b\x25\xbb\xbf\x99\x02\xbe\x74\x96\x31\x65\x08\xcb\xf0\x41\xb0\xb3\xde\x72\x47\xad\x42\x27\x12\x59\xd5\xf5\x37\x9b\x74\x15\x29\x99\xe0\x75\x65\x82\xab\xb7\x29\x9c\x98\xea\x7f\x04\xe1\x3e\x49\xeb\x59\xb4\x37\x84\xf3\x19\x2f\xae\x7e\xaa\xfb\xf4\x30\x8e\x17\xc2\x4b\xe0\xd5\x6d\x7e\xf9\xf6\xdd\x87\xe5\x46\xdd\xcd\x60\x1c\xff\x9e\xe0\xbf\xe7\xc2\xec\x19\x54\xf3\xe2\x77\x00\x00\x00\xff\xff\x9f\x0f\xd0\x78\x2e\x02\x00\x00")

func iatemplatesPostgresql_downYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesPostgresql_downYml,
		"iatemplates/postgresql_down.yml",
	)
}

func iatemplatesPostgresql_downYml() (*asset, error) {
	bytes, err := iatemplatesPostgresql_downYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/postgresql_down.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd8, 0x2d, 0x2f, 0xf7, 0x84, 0xa0, 0xbf, 0xa8, 0x58, 0xc4, 0xae, 0x16, 0x1a, 0x84, 0x29, 0xde, 0x1c, 0x9, 0x6e, 0xbf, 0xaa, 0x6a, 0x2d, 0x36, 0xd3, 0xb7, 0x12, 0x30, 0x86, 0x7c, 0xb1, 0x25}}
	return a, nil
}

var _iatemplatesPostgresql_restartedYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\x4d\x6f\xd3\x50\x10\xbc\xfb\x57\x8c\x2a\x0e\x20\x61\xc7\x51\x85\x40\x4f\x5c\x40\xe2\xd6\x4a\x85\x70\x8b\x22\x6b\x6b\x6f\xec\x07\x7e\x1f\xbc\xdd\x24\x44\x25\xff\x1d\xe5\x25\x31\x69\x4f\xf8\xb6\x33\xf3\x46\xb3\x3b\x2e\xcb\xb2\x50\x76\x71\x24\x65\x31\x05\x50\xc2\x93\x63\x83\xe8\x5c\x13\x83\x68\x9f\x58\x7e\x8d\x4d\x62\x51\x4a\xca\x5d\x01\x00\x5b\x4e\x62\x83\x37\x98\xe7\x51\x36\xce\x51\xda\x1b\x3c\x9c\x1e\x2c\xbe\xde\xe1\xdb\xb3\x07\xfc\x3b\x26\x83\x3f\x65\x1e\x80\xd8\x67\x6f\x47\xa2\x9c\x9a\x4d\x54\xeb\xb8\x11\x6e\x83\xef\xe4\x2c\xf9\x88\xe5\x12\x95\x0e\x89\x65\x08\x63\x87\xd5\x2a\x13\x91\x12\xb9\x1c\xf4\xf8\x5d\xc2\x4e\xb2\x33\x7e\x15\x69\x71\x72\xc5\x3a\x05\x07\x1d\x18\x23\x89\xe2\xbc\xce\x24\xdf\x78\xab\x06\x32\xcd\xba\x8f\x6c\xb0\x1e\x03\xfd\xd3\x24\xf2\x3d\x1b\x2c\xeb\xb7\x98\x7f\xa8\xeb\xd5\x44\x6c\x69\xdc\xb0\xc1\x6d\x5d\x67\x68\x1d\x92\xc1\xbc\x3e\x99\x09\x6f\x39\x59\xdd\x1b\xec\x28\x79\xeb\xfb\x8c\x8e\xf4\xc8\xe3\xb4\x85\xda\xf6\x27\xab\xc1\xc3\xfd\x7d\xf9\x7e\x7e\xfb\xee\x02\x0f\xd6\xf7\xd2\x68\x68\xba\x60\x70\xb3\x67\xb9\xc9\x0c\x79\x1f\x94\xd4\x06\x3f\x39\x74\x2c\x6d\xb2\x51\x73\x27\xd3\x95\x81\xef\x83\x15\x5c\xea\x85\x15\xf8\xa0\x58\x5b\x6f\x65\xe0\xce\x60\x50\x8d\x62\x66\xb3\x1f\x36\x51\x15\x39\xb5\xc1\x53\xd5\x06\x37\x7b\x4c\x61\x27\x3c\x7b\x11\x08\xd7\xfd\x5a\x2f\x4a\xbe\x65\x3c\x3d\xe1\xd5\x69\x9f\x4a\x38\x6d\x6d\xcb\xcd\xb1\x14\x1c\x0e\xd8\x91\x60\xfa\x71\xb2\x30\x9f\xea\x48\x9d\xcb\x06\xf5\xa1\x9a\xec\xef\x3e\x7d\xfe\x72\xb7\x30\x57\x96\x38\x1c\x8a\x17\x85\x5e\x65\x78\xee\xff\xfa\x3f\x22\xbd\x29\xfe\x06\x00\x00\xff\xff\x6e\xb1\x20\x65\xf1\x02\x00\x00")

func iatemplatesPostgresql_restartedYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesPostgresql_restartedYml,
		"iatemplates/postgresql_restarted.yml",
	)
}

func iatemplatesPostgresql_restartedYml() (*asset, error) {
	bytes, err := iatemplatesPostgresql_restartedYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/postgresql_restarted.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x89, 0xba, 0xf1, 0x6c, 0x73, 0xbb, 0x8d, 0x49, 0x6b, 0x52, 0x53, 0x42, 0x6d, 0xc3, 0x7, 0xcd, 0x6a, 0x72, 0xa, 0x53, 0x1, 0x88, 0xec, 0x2d, 0xb6, 0x74, 0xc3, 0x51, 0x64, 0xfc, 0x59, 0x7c}}
	return a, nil
}

var _iatemplatesPostgresql_too_many_connectionsYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x93\x5d\x6b\xdb\x4c\x10\x85\xef\xfd\x2b\xce\x6b\xde\x40\x12\x2a\x3b\xa6\x84\x96\x85\x16\x5c\xc8\x9d\x0d\x29\x69\x7b\x63\x8c\x98\x48\x23\x69\x5b\xed\xae\xba\x3b\x72\x6c\x1c\xf7\xb7\x17\xad\x2d\xd5\xfd\xd4\xdd\x9e\x9d\x33\xfb\xcc\x87\x92\x24\x19\x09\x9b\xa6\x26\xe1\xa0\x46\x40\x02\x4b\x86\x15\x1a\x63\xd2\xc6\x05\x29\x3d\x87\xaf\x75\x2a\xce\xa5\x86\xec\x2e\xcd\x9c\xb5\x9c\x89\x76\x36\x8c\x00\x60\xc3\x3e\x68\x67\x15\x66\xf1\x18\x5a\x63\xc8\xef\x14\xee\x8f\xde\x87\xf7\x0b\x9c\x59\xa0\x2d\xda\xc0\x31\x94\xb7\x8d\x57\x78\x4e\xe2\x21\x3a\x2f\x9b\x32\x0d\x42\x92\x52\x26\x7a\xa3\xa5\x7b\xad\xb5\xb2\xcf\x49\x3a\xa8\xff\xbe\x8d\x7b\xd4\xc9\xf5\x73\x0f\x37\x3e\x5c\x9d\x32\xbc\x45\xe7\x67\x11\x6d\xcb\x90\x1a\xda\x9e\xc3\xe2\x1a\xab\x15\x26\x52\x79\x0e\x95\xab\x73\xac\xd7\x98\x62\x76\x73\x13\xcd\x0d\x79\x32\xb1\xfe\xee\xeb\x7b\x30\x04\x9f\xf4\xb3\xf2\xe6\x68\xd8\x67\x6c\x85\x4a\x46\xe1\x9d\xe9\xaa\x2c\x74\xd9\x7a\xce\x61\x68\xab\x4d\x6b\x06\x57\x6b\xb5\x28\x8c\x2f\xc6\x83\x22\xbb\x86\x15\x8a\xda\x91\x0c\x9a\x27\x5b\xb2\xc2\xea\xe6\x45\x87\xb5\x1e\xf4\x0d\xd5\x2d\x2b\xbc\x3e\x92\x16\xce\x2b\xdc\x1e\x73\x07\xde\xb0\xd7\xb2\x53\x78\x22\x6f\xb5\x2d\xa3\x5a\xd3\x23\xd7\x43\x2d\xa2\xb3\x2f\x2c\x0a\xf7\xcb\x65\xf2\x6a\xf6\xf2\xb6\x97\xab\xd8\x24\x71\x69\xee\x14\xc6\x3b\x0e\x47\x38\xb2\xd6\x09\xc5\x8e\xf5\x19\x72\x0e\x99\xd7\x8d\xc4\x29\x0f\xe3\x02\x3e\x54\x3a\xa0\x1f\x08\x74\x80\x75\x82\x42\x5b\x1d\x2a\xce\x15\x2a\x91\x26\xa8\xe9\xf4\xb3\xf6\x34\xe9\x9a\xe5\x2c\x4d\x32\x67\xa6\x8f\xde\x3d\x05\x9e\xfe\x02\x04\x2c\x9d\x67\x48\x45\xf6\xb7\x41\x5d\xc0\x15\x7f\x5b\x28\xf2\x7c\x5a\x2a\x38\x8b\xfd\x1e\xff\x1f\xeb\x9f\x04\xf6\x1b\x9d\x71\xda\x8d\x12\x87\xc3\xf0\xcc\xa7\xf9\xe2\xe3\x1d\xde\xc4\xd0\xd8\xda\xf3\xcb\xc5\xfc\xdd\xdd\xe2\x41\x9d\xe5\xf9\x71\xfb\xa7\xd5\x16\xe7\xd0\xfd\x16\x3f\x21\x5d\x6a\x1b\x84\x6c\xc6\xff\xc0\xb9\x1a\x7d\x0f\x00\x00\xff\xff\x7b\x94\x8f\x46\x7a\x03\x00\x00")

func iatemplatesPostgresql_too_many_connectionsYmlBytes() ([]byte, error) {
	return bindataRead(
		_iatemplatesPostgresql_too_many_connectionsYml,
		"iatemplates/postgresql_too_many_connections.yml",
	)
}

func iatemplatesPostgresql_too_many_connectionsYml() (*asset, error) {
	bytes, err := iatemplatesPostgresql_too_many_connectionsYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "iatemplates/postgresql_too_many_connections.yml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xf1, 0x3b, 0x1b, 0x85, 0xeb, 0x36, 0x50, 0x6d, 0xab, 0x1a, 0x9e, 0x3a, 0x13, 0x4d, 0xfe, 0x25, 0x72, 0x59, 0x72, 0xae, 0x3a, 0xc4, 0x27, 0x45, 0xec, 0xed, 0x40, 0x7c, 0xc7, 0x59, 0x17, 0x3d}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"iatemplates/mongodb_connections_memory_usage.yml": iatemplatesMongodb_connections_memory_usageYml,
	"iatemplates/mongodb_high_memory_usage.yml":        iatemplatesMongodb_high_memory_usageYml,
	"iatemplates/mongodb_restarted.yml":                iatemplatesMongodb_restartedYml,
	"iatemplates/mysql_down.yml":                       iatemplatesMysql_downYml,
	"iatemplates/mysql_restarted.yml":                  iatemplatesMysql_restartedYml,
	"iatemplates/mysql_too_many_connections.yml":       iatemplatesMysql_too_many_connectionsYml,
	"iatemplates/node_high_cpu_load.yml":               iatemplatesNode_high_cpu_loadYml,
	"iatemplates/node_low_free_memory.yml":             iatemplatesNode_low_free_memoryYml,
	"iatemplates/node_swap_filled_up.yml":              iatemplatesNode_swap_filled_upYml,
	"iatemplates/postgresql_down.yml":                  iatemplatesPostgresql_downYml,
	"iatemplates/postgresql_restarted.yml":             iatemplatesPostgresql_restartedYml,
	"iatemplates/postgresql_too_many_connections.yml":  iatemplatesPostgresql_too_many_connectionsYml,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"iatemplates": {nil, map[string]*bintree{
		"mongodb_connections_memory_usage.yml": {iatemplatesMongodb_connections_memory_usageYml, map[string]*bintree{}},
		"mongodb_high_memory_usage.yml":        {iatemplatesMongodb_high_memory_usageYml, map[string]*bintree{}},
		"mongodb_restarted.yml":                {iatemplatesMongodb_restartedYml, map[string]*bintree{}},
		"mysql_down.yml":                       {iatemplatesMysql_downYml, map[string]*bintree{}},
		"mysql_restarted.yml":                  {iatemplatesMysql_restartedYml, map[string]*bintree{}},
		"mysql_too_many_connections.yml":       {iatemplatesMysql_too_many_connectionsYml, map[string]*bintree{}},
		"node_high_cpu_load.yml":               {iatemplatesNode_high_cpu_loadYml, map[string]*bintree{}},
		"node_low_free_memory.yml":             {iatemplatesNode_low_free_memoryYml, map[string]*bintree{}},
		"node_swap_filled_up.yml":              {iatemplatesNode_swap_filled_upYml, map[string]*bintree{}},
		"postgresql_down.yml":                  {iatemplatesPostgresql_downYml, map[string]*bintree{}},
		"postgresql_restarted.yml":             {iatemplatesPostgresql_restartedYml, map[string]*bintree{}},
		"postgresql_too_many_connections.yml":  {iatemplatesPostgresql_too_many_connectionsYml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
