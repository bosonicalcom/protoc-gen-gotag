package module

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	MinimumEdition = descriptorpb.Edition_EDITION_PROTO2
	MaximumEdition = descriptorpb.Edition_EDITION_2024
)

// SetEditionSupport takes a serialized CodeGeneratorResponse and sets the range of
// protobuf editions supported by this plugin, which protoc requires whenever
// FEATURE_SUPPORTS_EDITIONS is advertised.
func SetEditionSupport(raw []byte) ([]byte, error) {
	resp := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(raw, resp); err != nil {
		return nil, err
	}

	resp.MinimumEdition = proto.Int32(int32(MinimumEdition))
	resp.MaximumEdition = proto.Int32(int32(MaximumEdition))

	return proto.Marshal(resp)
}
