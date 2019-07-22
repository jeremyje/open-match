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

package util

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"

	"github.com/stretchr/testify/assert"
	"open-match.dev/open-match/pkg/pb"
)

func TestEmptyOpenMatchNamespace(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Equal("", getNamespaceFromIncomingRequest(ctx, t))
}

func TestMultiAttachNamespace(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	ctx, err := AppendOpenMatchNamespace(ctx, "one")
	assert.Equal("one", getNamespaceFromIncomingRequest(ctx, t))
	assert.Nil(err)
	ctx, err = AppendOpenMatchNamespace(ctx, "two")
	assert.Equal("one", getNamespaceFromIncomingRequest(ctx, t))
	assert.NotNil(err)
}

func TestAppendOpenMatchNamespace(t *testing.T) {
	testCases := []string{
		"",
		"open-match",
		"open-match-0123456789-0123456789-0123456789-0123456789-0123456789-0123456789-0123456789-0123456789-0123456789",
		"😀😁😂🤣😃😄😅😆😉😊😋😎😍😘",
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			assert := assert.New(t)
			ctx := context.Background()
			ctx, err := AppendOpenMatchNamespace(ctx, tc)

			assert.Equal(tc, getNamespaceFromIncomingRequest(ctx, t))
			assert.Nil(err)
		})
	}
}

func getNamespaceFromIncomingRequest(ctx context.Context, t *testing.T) string {
	close, port := serveForTest(t)
	defer close()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("cannot connect to gRPC server, %s", err)
	}
	defer conn.Close()
	c := pb.NewFrontendClient(conn)
	r, err := c.CreateTicket(ctx, &pb.CreateTicketRequest{})
	if err != nil {
		t.Errorf("cannot call CreateTicket %s", err)
	}
	return r.Ticket.Id
}
