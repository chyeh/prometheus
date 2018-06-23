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
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/api/prometheus/v1"
)

func parseTime(s string) (time.Time, error) {
	// Default Value
	if s == "" {
		return time.Time{}, nil
	}

	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		return time.Unix(int64(s), int64(ns*float64(time.Second))), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q to a valid timestamp", s)
}

func timeToRange(start time.Time, end time.Time) (*v1.Range, error) {
	if end.IsZero() {
		end = time.Now()
	}
	if start.IsZero() {
		start = end.Add(-5 * time.Minute)
	}

	if start.After(end) {
		return nil, fmt.Errorf("start time is not before end time")
	}

	resolution := math.Max(math.Floor(end.Sub(start).Seconds()/250), 1)
	// Convert seconds to nanoseconds such that time.Duration parses correctly.
	step := time.Duration(resolution * 1e9)

	return &v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}, nil
}
