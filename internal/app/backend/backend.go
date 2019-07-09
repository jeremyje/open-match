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

package backend

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"open-match.dev/open-match/internal/config"
	"open-match.dev/open-match/internal/rpc"
	"open-match.dev/open-match/internal/statestore"
	"open-match.dev/open-match/pkg/pb"
)

var (
	backendLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "backend",
	})
)

// RunApplication creates a server.
func RunApplication() {
	cfg, err := config.Read()
	if err != nil {
		backendLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot read configuration.")
	}
	p, err := rpc.NewServerParamsFromConfig(cfg, "api.backend")
	if err != nil {
		backendLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot construct server.")
	}

	if err := BindService(p, cfg); err != nil {
		backendLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("failed to bind backend service.")
	}

	rpc.MustServeForever(p)
}

// BindService creates the backend service and binds it to the serving harness.
func BindService(p *rpc.ServerParams, cfg config.View) error {
	service := &backendService{
		cfg:          cfg,
		synchronizer: &synchronizerClient{cfg: cfg},
		store:        statestore.New(cfg),
		mmfClients:   rpc.NewClientCache(cfg),
	}

	p.AddHealthCheckFunc(service.store.HealthCheck)
	p.AddHandleFunc(func(s *grpc.Server) {
		pb.RegisterBackendServer(s, service)
	}, pb.RegisterBackendHandlerFromEndpoint)

	return nil
}
