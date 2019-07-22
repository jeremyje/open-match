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

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

const (
	// MetadataNameOpenMatchNamespace is the gRPC Metadata for Open Match namespace.
	MetadataNameOpenMatchNamespace = "om-namespace"
	// MetadataNameOpenMatchNamespaceHTTPHeader is the gRPC Metadata for Open Match namespace.
	MetadataNameOpenMatchNamespaceHTTPHeader = runtime.MetadataHeaderPrefix + MetadataNameOpenMatchNamespace
	// MetadataNameOpenMatchNamespaceHTTP is the gRPC Metadata for Open Match namespace.
	MetadataNameOpenMatchNamespaceHTTP = runtime.MetadataPrefix + MetadataNameOpenMatchNamespace
)

// RandomNamespace returns a randomly generated namespace based on UUID-4.
func RandomNamespace(baseName string) string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	if len(baseName) > 0 {
		return baseName + "-" + u.String()
	}
	return u.String()
}

// AppendOpenMatchNamespace attaches the namespace name to a request context.
func AppendOpenMatchNamespace(ctx context.Context, namespace string) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		values := getNamespaceFromMetadata(md)
		if len(values) == 1 {
			if values[0] != namespace {
				return ctx, errors.New("request already has a namespace")
			}
			return ctx, nil
		}
		if len(values) > 1 {
			return ctx, errors.New("context has more than 1 namespace value, it is corrupt")
		}
	}

	return metadata.AppendToOutgoingContext(ctx, MetadataNameOpenMatchNamespace, namespace), nil
}

// GetOpenMatchNamespace gets the namespace of the request.
func GetOpenMatchNamespace(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := getNamespaceFromMetadata(md)
	if len(values) == 1 {
		return values[0]
	}
	return ""
}

func getNamespaceFromMetadata(md metadata.MD) []string {
	values := md.Get(MetadataNameOpenMatchNamespace)
	if len(values) > 0 {
		return values
	}
	values = md.Get("grpcgateway-" + MetadataNameOpenMatchNamespace)
	if len(values) > 0 {
		return values
	}
	return []string{}
}

// NewNamespacedContext returns a context that holds an Open Match namespace.
func NewNamespacedContext(namespace string) context.Context {
	if len(namespace) > 0 {
		md := metadata.New(map[string]string{
			MetadataNameOpenMatchNamespace: namespace,
		})
		return metadata.NewIncomingContext(context.Background(), md)
	}
	return context.Background()
}

// DeferredContext returns a functor that creates the Open Match context outside of a request.
func DeferredContext(ctx context.Context) func() context.Context {
	namespace := GetOpenMatchNamespace(ctx)
	return func() context.Context {
		return NewNamespacedContext(namespace)
	}
}
