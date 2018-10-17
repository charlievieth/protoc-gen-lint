package linter

import (
	"bufio"
	"io"
	"sort"
	"strconv"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

const (
	//Path Types
	pathMessage        = 4
	pathEnumType       = 5
	pathEnumValue      = 2
	pathService        = 6
	pathRPCMethod      = 2
	pathMessageName    = 1
	pathMessageField   = 2
	pathMessageMessage = 3
	pathMessageEnum    = 4
)

type errorCode int

const (
	//Error Types
	errorImportOrder errorCode = iota
	errorMessageCase
	errorFieldCase
	errorEnumTypeCase
	errorEnumValueCase
	errorServiceCase
	errorRPCMethodCase
)

var linterErrors = [...]string{
	"Sort import statements alphabetically.",
	"Use CamelCase (with an initial capital) for message names.",
	"Use underscore_separated_names for field names.",
	"Use CamelCase (with an initial capital) for enum type names.",
	"Use CAPITALS_WITH_UNDERSCORES  for enum value names.",
	"Use CamelCase (with an initial capital) for service names.",
	"Use CamelCase (with an initial capital) for RPC method names.",
}

type Config struct {
	ProtoFile   *descriptor.FileDescriptorProto
	OutFile     io.WriteCloser
	SortImports bool
}

// LintProtoFile takes a file name, proto file description, and a file.
// It checks the file for errors and writes them to the output file
func LintProtoFile(conf Config) (int, error) {
	var errors protoBufErrors

	if conf.SortImports {
		errors.lintImportOrder(conf.ProtoFile.GetDependency())
	}

	for i, v := range conf.ProtoFile.GetMessageType() {
		errors.lintProtoMessage(int32(i), pathMessage, []int32{}, v)
	}

	for i, v := range conf.ProtoFile.GetEnumType() {
		errors.lintProtoEnumType(int32(i), pathEnumType, []int32{}, v)
	}

	for i, v := range conf.ProtoFile.GetService() {
		errors.lintProtoService(int32(i), v)
	}

	// filter errors
	a := errors[:0]
	for _, v := range errors {
		if v.errorCode != errorMessageCase {
			a = append(a, v)
		}
	}
	errors = a

	// calculate line and column in parallel
	errors.calculateLineCol(conf.ProtoFile.GetSourceCodeInfo())
	// sort output by line and column
	sort.Sort(byLineCol(errors))

	buf := make([]byte, 0, 128)
	w := bufio.NewWriter(conf.OutFile)

	prefix := *conf.ProtoFile.Name + ":"
	buf = append(buf, prefix...)

	for _, v := range errors {
		buf = buf[:len(prefix)]
		buf = strconv.AppendInt(buf, int64(v.line), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(v.col), 10)
		buf = append(buf, ": '"...)
		buf = append(buf, v.errorString...)
		buf = append(buf, "' - "...)
		buf = append(buf, linterErrors[v.errorCode]...)
		buf = append(buf, '\n')
		if _, err := w.Write(buf); err != nil {
			return 0, err
		}
	}
	if err := w.Flush(); err != nil {
		return len(errors), err
	}

	return len(errors), nil
}
