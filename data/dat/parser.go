package dat

import (
	"encoding/json"
	"io"
	"io/fs"
	"strings"

	"github.com/oriath-net/pogo/dat"
	"github.com/pkg/errors"
)

var parser *dat.DataParser

func LoadParser() {
	LoadSchema()
	semverGameVersion := strings.Join(strings.Split(gameVersion, ".")[:3], ".")
	parser = dat.InitParser(semverGameVersion, &schemaFS{})
}

type schemaFS struct {
}

func (s *schemaFS) Open(name string) (fs.File, error) {
	data, err := json.Marshal(tableMap[strings.Split(name, ".")[0]].ToJSONFormat())
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal data")
	}
	return &schemaFSFile{Data: data}, nil
}

type schemaFSFile struct {
	Data   []byte
	Offset int
}

func (s *schemaFSFile) Stat() (fs.FileInfo, error) {
	// Do nothing
	return nil, nil
}

func (s *schemaFSFile) Read(bytes []byte) (int, error) {
	copied := copy(bytes, s.Data[s.Offset:])
	s.Offset += copied

	if s.Offset >= len(s.Data) {
		return copied, io.EOF
	}

	return copied, nil
}

func (s *schemaFSFile) Close() error {
	// Do nothing
	return nil
}

func ParseDat(data io.Reader, filename string) ([]interface{}, error) {
	parse, err := parser.Parse(data, filename)
	return parse, errors.Wrap(err, "failed to parse dat file")
}
