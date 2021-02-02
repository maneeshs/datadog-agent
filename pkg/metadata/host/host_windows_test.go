// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2021 Datadog, Inc.

package host

import (
	"testing"
)

func TestFillOsVersion(t *testing.T) {
	stats := &systemStats{}
	info := getHostInfo()
	fillOsVersion(stats, info)
}
