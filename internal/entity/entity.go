package entity

import (
	"fmt"
	"math/rand"
	"time"
)

// RequestParams params for request from form.
type RequestParams struct {
	Host            string
	Method          string
	MethodType      MethodType
	Service         string
	Message         string
	RPS             int
	Metadata        map[string]string
	RequestDeadline *time.Duration
	Proto           *ParsedProto
}

// Message from proto.
type Message struct {
	Name   string
	Fields []Field
}

// Field param for message.
type Field struct {
	Name     string
	Type     string
	Message  *Message
	IsMap    bool
	EnumName string
}

// MethodType type of method.
type MethodType int

// Available values for MethodType.
const (
	MethodTypeUnaryRPC                  MethodType = 0
	MethodTypeServerStreamingRPC        MethodType = 1
	MethodTypeClientStreamingRPC        MethodType = 2
	MethodTypeBidirectionalStreamingRPC MethodType = 3
)

// Method from proto.
type Method struct {
	Name           string
	RequestMessage *Message
	Type           MethodType
}

// Service from proto.
type Service struct {
	Name    string
	Methods []Method
}

// Enum from proto.
type Enum struct {
	Name   string
	Values []string
}

// RandomValue return random value for enum.
func (e Enum) RandomValue() string {
	randomIndex := rand.Intn(len(e.Values))
	return e.Values[randomIndex]
}

// ParsedProto proto that is serialized into a convenient structure.
type ParsedProto struct {
	Services []Service
	Enums    []Enum
	Package  string
	FilePath string
}

// FindMethodByName return enum by service and name.
func (p *ParsedProto) FindMethodByName(serviceName, methodName string) (Method, bool) {
	for _, s := range p.Services {
		fmt.Println(s.Name, serviceName)
		if s.Name == serviceName {
			for _, m := range s.Methods {
				fmt.Println(m.Name, methodName)
				if m.Name == methodName {
					return m, true
				}
			}
		}
	}

	return Method{}, false
}

// FindEnumByName return method by name.
func (p *ParsedProto) FindEnumByName(name string) (Enum, bool) {
	for _, enum := range p.Enums {
		if enum.Name == name {
			return enum, true
		}
	}

	return Enum{}, false
}
