package utils

import (
	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
)

func SchemaStringToProto(schema string) proxy_pb.Schema {
	switch schema {
	case "http":
		return proxy_pb.Schema_SCHEMA_HTTP
	case "https":
		return proxy_pb.Schema_SCHEMA_HTTPS
	case "socks5":
		return proxy_pb.Schema_SCHEMA_SOCKS5
	default:
		return proxy_pb.Schema_SCHEMA_UNSPECIFIED
	}
}

func SchemaProtoToString(schema proxy_pb.Schema) string {
	switch schema {
	case proxy_pb.Schema_SCHEMA_HTTP:
		return "http"
	case proxy_pb.Schema_SCHEMA_HTTPS:
		return "https"
	case proxy_pb.Schema_SCHEMA_SOCKS5:
		return "socks5"
	case proxy_pb.Schema_SCHEMA_UNSPECIFIED:
		fallthrough
	default:
		return "unknown"
	}
}
