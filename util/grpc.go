package util

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

type GrpcMeta struct {
	UserAgent string
	ClientIp  string
}

func ExtractMetaData(c context.Context) *GrpcMeta {
	mtd := GrpcMeta{}

	if md, ok := metadata.FromIncomingContext(c); ok {
		if uA, ok := md["user-agent"]; ok {
			mtd.UserAgent = strings.Join(uA, ",")
		}
		if uA, ok := md["grpcgateway-user-agent"]; ok {
			mtd.UserAgent = strings.Join(uA, ",")
		}
		if ip, ok := md[":authority"]; ok {
			mtd.ClientIp = strings.Join(ip, ",")
		}
		if ip, ok := md["x-forwarded-host"]; ok {
			mtd.ClientIp = strings.Join(ip, ",")
		}
	}
	return &mtd
}
