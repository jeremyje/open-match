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
	"google.golang.org/grpc/metadata"
)

const (
	metadataNameTenantID = "Multitenancy-TenantID"
)

// GetTenantIDFromContext returns the Tenant ID from a RPC context.
func GetTenantIDFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(metadataNameTenantID)
	if len(values) == 1 {
		return values[0]
	}
	return ""
}

// AppendTenantID adds the Tenant ID to the request context.
func AppendTenantID(ctx context.Context, tenantID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metadataNameTenantID, tenantID)
}
