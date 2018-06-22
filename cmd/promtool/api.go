// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type API struct {
	requestTimeout time.Duration
	v1             v1.API
}

func newAPI(serverURL string) (*API, error) {
	c, err := api.NewClient(api.Config{Address: serverURL})
	if err != nil {
		return nil, err
	}
	api := v1.NewAPI(c)
	return &API{
		requestTimeout: defaultTimeout,
		v1:             api,
	}, nil
}

func (i *API) QueryInstant(query string, ts time.Time) (model.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.requestTimeout)
	defer cancel()
	return i.v1.Query(ctx, query, ts)
}

func (i *API) QueryRange(query string, r v1.Range) (model.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.requestTimeout)
	defer cancel()
	return i.v1.QueryRange(ctx, query, r)
}
