// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package util provides utilities for net.Listener.
package util

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetadataForwardUnaryServerInterceptor forwards gRPC metadata fields from incoming context to outgoing context.
func MetadataForwardUnaryServerInterceptor(kList ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := incomingToOutgoingMetadata(ctx, kList)
		return handler(newCtx, req)
	}
}

// MetadataForwardStreamServerInterceptor forwards gRPC metadata fields from incoming context to outgoing context.
func MetadataForwardStreamServerInterceptor(kList ...string) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := incomingToOutgoingMetadata(stream.Context(), kList)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		err := handler(srv, wrapped)
		return err
	}
}

func incomingToOutgoingMetadata(ctx context.Context, kList []string) context.Context {
	mdOut, ok := metadata.FromOutgoingContext(ctx)
	logger.Infof("Outgoing Metadata: %+v, %t", mdOut, ok)
	md, ok := metadata.FromIncomingContext(ctx)
	logger.Infof("Incoming Metadata: %+v, %t", md, ok)
	attached := false
	if ok {
		for _, k := range kList {
			vals := md.Get(k)
			if len(vals) > 0 {
				ctx = metadata.AppendToOutgoingContext(ctx, k, vals[0])
				attached = true
			}
		}
	}
	if !attached {
		logger.Warningf("incomingToOutgoingMetadata has no metadata. %+v", ctx)
		panic("incomingToOutgoingMetadata has no metadata.")
	}
	return ctx
}
