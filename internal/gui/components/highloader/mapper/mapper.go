package mapper

import (
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
)

const (
	defaultString = "qwerty"
	defaultInt    = 1
	defaultFloat  = 0.1
	defaultBool   = false
)

// ExampleMessage map for example message.
type ExampleMessage map[string]any

// Mapper struct for mapping data for GUI.
type Mapper struct {
	proto *entity.ParsedProto
}

// NewMapper create a new Mapper.
func NewMapper(proto *entity.ParsedProto) *Mapper {
	return &Mapper{proto: proto}
}

// MakeExampleMessage create example message.
func (m *Mapper) MakeExampleMessage(msg *entity.Message) *ExampleMessage {
	exampleMessage := make(ExampleMessage)
	for _, field := range msg.Fields {
		if field.IsMap {
			exampleMessage[field.Name] = make(map[string]any)
			continue
		}
		if field.Message != nil {
			exampleMessage[field.Name] = m.MakeExampleMessage(field.Message)
			continue
		}
		exampleMessage[field.Name] = m.getDefaultValueForScalar(&field)
	}

	return &exampleMessage
}

// getDefaultValueForScalar return value for scalar.
func (m *Mapper) getDefaultValueForScalar(field *entity.Field) any {
	switch field.Type {
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32.String(), descriptorpb.FieldDescriptorProto_TYPE_SINT64.String(),
		descriptorpb.FieldDescriptorProto_TYPE_INT32.String(), descriptorpb.FieldDescriptorProto_TYPE_INT64.String(),
		descriptorpb.FieldDescriptorProto_TYPE_UINT32.String(), descriptorpb.FieldDescriptorProto_TYPE_UINT64.String(),
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32.String(), descriptorpb.FieldDescriptorProto_TYPE_FIXED64.String():
		return defaultInt
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT.String(), descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.String():
		return defaultFloat
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL.String():
		return defaultBool
	case descriptorpb.FieldDescriptorProto_TYPE_STRING.String():
		return defaultString
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM.String():
		enum, ok := m.proto.FindEnumByName(field.EnumName)
		if !ok {
			return "EnumExampleNotFound"
		}
		return enum.RandomValue()
	default:
		return ""
	}
}

// MakeProtoMapUI create ProtoMapUI for GUI.
func (m *Mapper) MakeProtoMapUI(parsedProto *entity.ParsedProto) *ProtoMapUI {
	protoMapUI := ProtoMapUI{
		Services: make(map[string][]string),
		Methods:  make(map[string]*entity.Message),
	}
	for _, service := range parsedProto.Services {
		for _, method := range service.Methods {
			// TODO: remove after implementation other types.
			if method.Type == entity.MethodTypeUnaryRPC {
				protoMapUI.Services[service.Name] = append(protoMapUI.Services[service.Name], method.Name)
				protoMapUI.Methods[method.Name] = method.RequestMessage
			}
		}
	}

	return &protoMapUI
}

// ProtoMapUI struct with needed params for GUI.
type ProtoMapUI struct {
	Services map[string][]string
	Methods  map[string]*entity.Message
}

// GetServicesNames return services names.
func (p *ProtoMapUI) GetServicesNames() []string {
	names := make([]string, 0, len(p.Services))
	for name := range p.Services {
		names = append(names, name)
	}

	return names
}

// GetMethodsNamesByService return methods by service name.
func (p *ProtoMapUI) GetMethodsNamesByService(serviceName string) []string {
	methods := p.Services[serviceName]
	return methods
}

// GetMessageByMethodName return message by method name.
func (p *ProtoMapUI) GetMessageByMethodName(methodName string) (*entity.Message, bool) {
	m, ok := p.Methods[methodName]
	return m, ok
}
