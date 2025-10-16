package proto

import (
	"errors"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
)

// ProtoParser struct for parse proto.
type ProtoParser struct {
	parser *protoparse.Parser
	fp     string
}

// NewProtoParser create a new ProtoParser.
func NewProtoParser() *ProtoParser {
	parser := &protoparse.Parser{}

	return &ProtoParser{
		parser: parser,
	}
}

// GetMethodDescriptor return desc.MethodDescriptor.
func (p *ProtoParser) GetMethodDescriptor(fp, methodName, serviceName string) (*desc.MethodDescriptor, error) {
	fd, err := p.parseFile(fp)
	if err != nil {
		return nil, err
	}

	svc := fd.FindSymbol(serviceName)
	fmt.Println(serviceName)
	if svc == nil {
		return nil, errors.New("service not found")
	}

	svcDesc, ok := svc.(*desc.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("cannot find service %q", serviceName)
	}

	methodDesc := svcDesc.FindMethodByName(methodName)
	if methodDesc == nil {
		return nil, fmt.Errorf("method %q not found for service %q", methodName, serviceName)
	}

	return methodDesc, nil
}

// ParseProto parse proto by filepath.
func (p *ProtoParser) ParseProto(fp string) (*entity.ParsedProto, error) {
	fd, err := p.parseFile(fp)
	if err != nil {
		return nil, err
	}

	parsedEntity := p.toEntity(fd, fp)

	return parsedEntity, nil
}

// parseFile parse file.
func (p *ProtoParser) parseFile(fp string) (*desc.FileDescriptor, error) {
	fds, err := p.parser.ParseFiles(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q: %w", p.fp, err)
	}

	fd := fds[0]
	return fd, nil
}

// toEntity convert internal struct to entity.ParsedProto.
func (p *ProtoParser) toEntity(fd *desc.FileDescriptor, fp string) *entity.ParsedProto {
	services := fd.GetServices()
	servicesEntity := make([]entity.Service, 0, len(services))
	for _, service := range services {
		methods := service.GetMethods()
		methodsEntity := make([]entity.Method, 0, len(methods))
		for _, method := range methods {
			m := entity.Method{
				Name:           method.GetName(),
				RequestMessage: p.makeMessage(method.GetInputType()),
				Type: p.getMethodType(method),
			}
			methodsEntity = append(methodsEntity, m)
		}

		s := entity.Service{
			Name:    service.GetName(),
			Methods: methodsEntity,
		}
		servicesEntity = append(servicesEntity, s)
	}

	enums := fd.GetEnumTypes()
	enumsEntity := make([]entity.Enum, 0, len(enums))
	for _, enum := range enums {
		values := enum.GetValues()
		valuesEnum := make([]string, 0, len(values))
		for _, enumValue := range values {
			valuesEnum = append(valuesEnum, enumValue.GetName())
		}

		e := entity.Enum{
			Name:   enum.GetName(),
			Values: valuesEnum,
		}
		enumsEntity = append(enumsEntity, e)
	}

	return &entity.ParsedProto{
		Services: servicesEntity,
		Enums:    enumsEntity,
		Package:  fd.GetPackage(),
		FilePath: fp,
	}
}

// getMethodType return method type.
func (p *ProtoParser) getMethodType(method *desc.MethodDescriptor) entity.MethodType {
	if method.IsClientStreaming() && !method.IsServerStreaming() {
		return entity.MethodTypeClientStreamingRPC
	} else if method.IsServerStreaming() && !method.IsClientStreaming() {
		return entity.MethodTypeServerStreamingRPC
	} else if method.IsClientStreaming() && method.IsServerStreaming() {
		return entity.MethodTypeBidirectionalStreamingRPC
	} else {
		return entity.MethodTypeUnaryRPC
	}
}

// makeMessage make entity.Message.
//
// Recursion.
func (p *ProtoParser) makeMessage(message *desc.MessageDescriptor) *entity.Message {
	fields := message.GetFields()
	fieldsEntity := make([]entity.Field, 0, len(fields))
	for _, field := range fields {
		f := entity.Field{
			Name:  field.GetJSONName(),
			Type:  field.GetType().String(),
			IsMap: field.IsMap(),
		}

		msg := field.GetMessageType()
		if msg != nil {
			// Calls itself.
			f.Message = p.makeMessage(msg)
		}

		enumType := field.GetEnumType()
		if enumType != nil {
			f.EnumName = enumType.GetName()
		}
		fieldsEntity = append(fieldsEntity, f)
	}
	m := entity.Message{
		Name:   message.GetName(),
		Fields: fieldsEntity,
	}

	return &m
}
