// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2021 Datadog, Inc.

// +build !kubeapiserver

package orchestrator

// GetStatus returns status info for the orchestrator explorer.
func GetStatus(apiCl interface{}) map[string]interface{} {
	status := make(map[string]interface{})
	status["Disabled"] = "The orchestrator explorer is not compiled-in"
	return status
}
