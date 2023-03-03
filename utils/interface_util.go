package utils

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func StructToString(v interface{}) string {
	ops := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	if protoMessage, ok := v.(proto.Message); ok {
		if b, err := ops.Marshal(protoMessage); err == nil {
			return string(b)
		}
	}
	if b, err := json.Marshal(v); err == nil {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}
