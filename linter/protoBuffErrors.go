package linter

import (
	"runtime"
	"sort"
	"sync"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type protoBufErrors []*protoBufError

func (p protoBufErrors) calculateLineCol(protoSource *descriptor.SourceCodeInfo) {
	if len(p) == 0 {
		return
	}

	numCPU := runtime.NumCPU()
	if len(p) < numCPU {
		numCPU = len(p)
	}

	ch := make(chan *protoBufError, numCPU*4)
	var wg sync.WaitGroup
	wg.Add(numCPU)

	for i := 0; i < numCPU; i++ {
		go func() {
			for v := range ch {
				v.setSourceLineNumber(protoSource)
			}
			wg.Done()
		}()
	}
	for _, v := range p {
		ch <- v
	}
	close(ch)
	wg.Wait()
}

func (p *protoBufErrors) addError(e *protoBufError) {
	*p = append(*p, e)
}

func (p *protoBufErrors) lintProtoMessage(
	pathIndex int32,
	pathType int32,
	parentPath []int32,
	protoMessage *descriptor.DescriptorProto,
) {

	path := append(
		parentPath,
		pathType,
		pathIndex,
	)

	if !isCamelCase(protoMessage.GetName()) {
		p.addError(&protoBufError{
			path:        path,
			errorCode:   errorMessageCase,
			errorString: protoMessage.GetName(),
		})
	}

	for i, v := range protoMessage.GetField() {
		p.lintProtoField(int32(i), path, v)
	}

	for i, v := range protoMessage.GetEnumType() {
		p.lintProtoEnumType(int32(i), pathMessageEnum, path, v)
	}

	for i, v := range protoMessage.GetNestedType() {
		p.lintProtoMessage(int32(i), pathMessageMessage, path, v)
	}
}

func (p *protoBufErrors) lintProtoField(
	pathIndex int32,
	parentPath []int32,
	messageField *descriptor.FieldDescriptorProto,
) {
	path := append(
		parentPath,
		pathMessageField,
		pathIndex,
	)
	if !isLowerUnderscore(messageField.GetName()) {
		p.addError(&protoBufError{
			path:        path,
			errorCode:   errorFieldCase,
			errorString: messageField.GetName(),
		})
	}
}

func (p *protoBufErrors) lintProtoEnumType(
	pathIndex int32,
	pathType int32,
	parentPath []int32,
	protoEnum *descriptor.EnumDescriptorProto,
) {
	path := append(
		parentPath,
		pathType,
		pathIndex,
	)

	if !isCamelCase(protoEnum.GetName()) {
		p.addError(&protoBufError{
			path:        path,
			errorCode:   errorEnumTypeCase,
			errorString: protoEnum.GetName(),
		})
	}

	for i, v := range protoEnum.GetValue() {
		p.lintProtoEnumValue(int32(i), path, v)
	}
}

func (p *protoBufErrors) lintProtoEnumValue(
	pathIndex int32,
	parentPath []int32,
	enumVal *descriptor.EnumValueDescriptorProto,
) {
	path := append(
		parentPath,
		pathEnumValue,
		pathIndex,
	)
	if !isUpperUnderscore(enumVal.GetName()) {
		p.addError(&protoBufError{
			path:        path,
			errorCode:   errorEnumValueCase,
			errorString: enumVal.GetName(),
		})
	}
}

func (p *protoBufErrors) lintProtoService(
	pathIndex int32,
	protoService *descriptor.ServiceDescriptorProto,
) {
	path := []int32{
		pathService,
		pathIndex,
	}
	if !isCamelCase(protoService.GetName()) {
		p.addError(&protoBufError{
			path:        path,
			errorCode:   errorServiceCase,
			errorString: protoService.GetName(),
		})
	}

	for i, v := range protoService.GetMethod() {
		p.lintProtoRPCMethod(int32(i), path, v)
	}
}

func (p *protoBufErrors) lintProtoRPCMethod(
	pathIndex int32,
	parentPath []int32,
	serviceMethod *descriptor.MethodDescriptorProto,
) {
	path := append(
		parentPath,
		pathRPCMethod,
		pathIndex,
	)
	if !isCamelCase(serviceMethod.GetName()) {
		p.addError(&protoBufError{
			path:        path,
			errorCode:   errorRPCMethodCase,
			errorString: serviceMethod.GetName(),
		})
	}
}

func (p *protoBufErrors) lintImportOrder(dependencies []string) {
	if !sort.StringsAreSorted(dependencies) {
		p.addError(&protoBufError{
			path:        []int32{},
			errorCode:   errorImportOrder,
			errorString: "import statements",
		})
	}
}
