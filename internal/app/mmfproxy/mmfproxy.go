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

package mmfproxy

import (
	"github.com/sirupsen/logrus"
	"open-match.dev/open-match/internal/config"
	"open-match.dev/open-match/internal/rpc"

	mmfHarness "open-match.dev/open-match/pkg/harness/function/golang"
	"open-match.dev/open-match/pkg/pb"
)

var (
	mmfProxyLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "mmfproxy",
	})
)

// RunApplication creates a server.
func RunApplication() {
	cfg, err := config.Read()
	if err != nil {
		mmfProxyLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot read configuration.")
	}
	clientCache := rpc.NewClientCache(cfg)
	// Invoke the harness to setup a GRPC service that handles requests to run the
	// match function. The harness itself queries open match for player pools for
	// the specified request and passes the pools to the match function to generate
	// proposals.
	mmfHarness.RunMatchFunction(&mmfHarness.FunctionSettings{
		Func: func(p *mmfHarness.MatchFunctionParams) ([]*pb.Match, error) {
			return proxyMatchMaker(clientCache, p)
		},
	})
}
