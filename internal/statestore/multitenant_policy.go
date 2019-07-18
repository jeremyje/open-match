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

package statestore

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"open-match.dev/open-match/internal/config"
)

// MultitenantPolicy controls the logic for how multi-tenancy works.
type MultitenantPolicy struct {
	cfg config.View
}

func (mtp *MultitenantPolicy) isRequired() {
	return mtp.cfg.GetBool("multitenancy.required")
}

// Verify the tenant ID
func (mtp *MultitenantPolicy) Verify(ctx context.Context, tenantID string) error {
	if mtp.isRequired() && tenantID == "" {
		return status.Error(codes.Internal, "tenant_id is required for this request")
	}
	return nil
}
