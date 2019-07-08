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

// Package mmfproxy provides a served version of the default MMF harness
// to be used for any languages. This service removes the need to talk to the Redis database
// directly.
package mmfproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"open-match.dev/open-match/internal/app/backend"
	"open-match.dev/open-match/internal/rpc"
	mmfHarness "open-match.dev/open-match/pkg/harness/function/golang"
	"open-match.dev/open-match/pkg/pb"
)

var (
	proxyFuncLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "app.mmfproxy.proxyfunc",
	})
)

func matchParamsToSimpleRunRequest(p *mmfHarness.MatchFunctionParams) *pb.SimpleRunRequest {
	/*
		tl := &pb.TicketList{
			Tickets: []*pb.Ticket{},
		}
	*/
	runReq := &pb.SimpleRunRequest{
		ProfileName: p.ProfileName,
		Properties:  p.Properties,
		Roster:      p.Rosters,
		//PoolNameToTickets: tl,
	}
	return runReq
}

func proxyMatchMaker(clientCache *rpc.ClientCache, p *mmfHarness.MatchFunctionParams) ([]*pb.Match, error) {
	ctx := p.Context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.Internal, "MMF Proxy requires metadata").Err()
	}
	configAsByteString, ok := md[backend.MMFProxyServiceConfigHeader]
	if !ok {
		status.Newf(codes.Internal, "%s was not found in request metadata", backend.MMFProxyServiceConfigHeader)
	}
	functionConfig := &pb.FunctionConfig{}
	err := proto.Unmarshal([]byte(configAsByteString[0]), functionConfig)

	configType := functionConfig.GetType()
	address := fmt.Sprintf("%s:%d", functionConfig.GetHost(), functionConfig.GetPort())

	runReq := matchParamsToSimpleRunRequest(p)
	switch configType {
	// MatchFunction Hosted as a GRPC service
	case pb.FunctionConfig_GRPC:
		var conn *grpc.ClientConn
		conn, err = clientCache.GetGRPC(address)
		if err != nil {
			proxyFuncLogger.WithFields(logrus.Fields{
				"error":    err.Error(),
				"function": functionConfig,
			}).Error("failed to establish grpc client connection to match function")
			return nil, status.Error(codes.InvalidArgument, "failed to connect to match function")
		}
		grpcClient := pb.NewSimpleMatchFunctionClient(conn)
		resp, err := grpcClient.SimpleRun(ctx, runReq)
		if err != nil {
			return nil, err
		}
		return resp.GetProposals(), nil
	// MatchFunction Hosted as a REST service
	case pb.FunctionConfig_REST:
		httpClient, baseURL, err := clientCache.GetHTTP(address)
		if err != nil {
			proxyFuncLogger.WithFields(logrus.Fields{
				"error":    err.Error(),
				"function": functionConfig,
			}).Error("failed to establish rest client connection to match function")
			return nil, status.Error(codes.InvalidArgument, "failed to connect to match function")
		}
		jsonReq, err := json.Marshal(runReq)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to marshal profile pb to string for profile %s: %s", runReq.GetProfileName(), err.Error())
		}

		reqBody, err := json.Marshal(map[string]json.RawMessage{"profile": jsonReq})
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to marshal request body for profile %s: %s", runReq.GetProfileName(), err.Error())
		}

		req, err := http.NewRequest("POST", baseURL+"/v1/matchfunction:run", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to create mmf http request for profile %s: %s", runReq.GetProfileName(), err.Error())
		}

		resp, err := httpClient.Do(req.WithContext(ctx))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get response from mmf run for proile %s: %s", runReq.GetProfileName(), err.Error())
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				proxyFuncLogger.WithError(err).Warning("failed to close response body read closer")
			}
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to read from response body for profile %s: %s", runReq.GetProfileName(), err.Error())
		}

		pbResp := &pb.RunResponse{}
		err = json.Unmarshal(body, pbResp)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to unmarshal response body to response pb for profile %s: %s", runReq.GetProfileName(), err.Error())
		}

		return pbResp.GetProposals(), nil
	default:
		return nil, status.Error(codes.InvalidArgument, "provided match function type is not supported")
	}
}
