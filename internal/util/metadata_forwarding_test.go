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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	netlistenerTesting "open-match.dev/open-match/internal/util/netlistener/testing"
	"open-match.dev/open-match/pkg/pb"
)

type forwardingFrontendServer struct {
	forwardingAddress string
}

// CreateTicket forwards.
func (s *forwardingFrontendServer) CreateTicket(ctx context.Context, req *pb.CreateTicketRequest) (*pb.CreateTicketResponse, error) {
	conn, err := grpc.Dial(s.forwardingAddress, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewFrontendClient(conn)
	return c.CreateTicket(ctx, req)
}

// DeleteTicket is not implemented.
func (s *forwardingFrontendServer) DeleteTicket(ctx context.Context, req *pb.DeleteTicketRequest) (*pb.DeleteTicketResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetTicket is not implemented.
func (s *forwardingFrontendServer) GetTicket(ctx context.Context, req *pb.GetTicketRequest) (*pb.Ticket, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetAssignments forwards stream requests.
func (s *forwardingFrontendServer) GetAssignments(req *pb.GetAssignmentsRequest, stream pb.Frontend_GetAssignmentsServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

type namespaceFrontendServer struct {
}

// CreateTicket forwards.
func (s *namespaceFrontendServer) CreateTicket(ctx context.Context, req *pb.CreateTicketRequest) (*pb.CreateTicketResponse, error) {
	return &pb.CreateTicketResponse{
		Ticket: &pb.Ticket{
			Id: GetOpenMatchNamespace(ctx),
		},
	}, nil
}

// DeleteTicket is not implemented.
func (s *namespaceFrontendServer) DeleteTicket(ctx context.Context, req *pb.DeleteTicketRequest) (*pb.DeleteTicketResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetTicket is not implemented.
func (s *namespaceFrontendServer) GetTicket(ctx context.Context, req *pb.GetTicketRequest) (*pb.Ticket, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetAssignments forwards stream requests.
func (s *namespaceFrontendServer) GetAssignments(req *pb.GetAssignmentsRequest, stream pb.Frontend_GetAssignmentsServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func newGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(MetadataForwardStreamServerInterceptor(MetadataNameOpenMatchNamespace))),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(MetadataForwardUnaryServerInterceptor(MetadataNameOpenMatchNamespace))))
}

func serveForTest(t *testing.T) (func(), int) {
	nsNl := netlistenerTesting.MustListen()
	nsServer := newGRPCServer()
	pb.RegisterFrontendServer(nsServer, &namespaceFrontendServer{})
	go func() {
		lis, err := nsNl.Obtain()
		if err != nil {
			t.Errorf("cannot create grpc server: %s", err)
		}
		nsServer.Serve(lis)
	}()

	forwardNl := netlistenerTesting.MustListen()
	forwardServer := newGRPCServer()
	pb.RegisterFrontendServer(forwardServer, &forwardingFrontendServer{
		forwardingAddress: fmt.Sprintf("localhost:%d", nsNl.Number()),
	})
	go func() {
		lis, err := forwardNl.Obtain()
		if err != nil {
			t.Errorf("cannot create grpc server: %s", err)
		}
		forwardServer.Serve(lis)
	}()

	return func() {
		nsServer.Stop()
		nsNl.Close()

		forwardServer.Stop()
		forwardNl.Close()
	}, forwardNl.Number()
}
