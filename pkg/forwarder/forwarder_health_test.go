// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2021 Datadog, Inc.

package forwarder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasValidAPIKey(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts1.Close()
	defer ts2.Close()

	keysPerDomains := map[string][]string{
		ts1.URL: {"api_key1", "api_key2"},
		ts2.URL: {"key3"},
	}

	fh := forwarderHealth{keysPerDomains: keysPerDomains}
	fh.init()
	assert.True(t, fh.hasValidAPIKey())

	assert.Equal(t, &apiKeyValid, apiKeyStatus.Get("API key ending with _key1"))
	assert.Equal(t, &apiKeyValid, apiKeyStatus.Get("API key ending with _key2"))
	assert.Equal(t, &apiKeyValid, apiKeyStatus.Get("API key ending with key3"))
}

func TestComputeDomainsURL(t *testing.T) {
	keysPerDomains := map[string][]string{
		"https://app.datadoghq.com":              {"api_key1"},
		"https://custom.datadoghq.com":           {"api_key2"},
		"https://custom.agent.datadoghq.com":     {"api_key3"},
		"https://app.datadoghq.eu":               {"api_key4"},
		"https://app.us2.datadoghq.com":          {"api_key5"},
		"https://custom.agent.us2.datadoghq.com": {"api_key6"},
		// debatable whether the next one should be changed to `api.`, preserve pre-existing behavior for now
		"https://app.datadoghq.internal": {"api_key7"},
		"https://app.myproxy.com":        {"api_key8"},
	}

	expectedMap := map[string][]string{
		"https://api.datadoghq.com":      {"api_key1", "api_key2", "api_key3"},
		"https://api.datadoghq.eu":       {"api_key4"},
		"https://api.us2.datadoghq.com":  {"api_key5", "api_key6"},
		"https://api.datadoghq.internal": {"api_key7"},
		"https://app.myproxy.com":        {"api_key8"},
	}

	// just sort the expected map for easy comparison
	for _, keys := range expectedMap {
		sort.Strings(keys)
	}

	fh := forwarderHealth{keysPerDomains: keysPerDomains}
	fh.init()

	// lexicographical sort for assert
	for _, keys := range fh.keysPerAPIEndpoint {
		sort.Strings(keys)
	}

	assert.Equal(t, expectedMap, fh.keysPerAPIEndpoint)
}

func TestHasValidAPIKeyErrors(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("api_key") == "api_key1" {
			w.WriteHeader(http.StatusForbidden)
		} else if r.Form.Get("api_key") == "api_key2" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			assert.Fail(t, fmt.Sprintf("Unknown api key received: %v", r.Form.Get("api_key")))
		}
	}))
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts1.Close()
	defer ts2.Close()

	keysPerAPIEndpoint := map[string][]string{
		ts1.URL: {"api_key1", "api_key2"},
		ts2.URL: {"key3"},
	}

	fh := forwarderHealth{}
	fh.init()
	fh.keysPerAPIEndpoint = keysPerAPIEndpoint
	assert.True(t, fh.hasValidAPIKey())

	assert.Equal(t, &apiKeyInvalid, apiKeyStatus.Get("API key ending with _key1"))
	assert.Equal(t, &apiKeyStatusUnknown, apiKeyStatus.Get("API key ending with _key2"))
	assert.Equal(t, &apiKeyValid, apiKeyStatus.Get("API key ending with key3"))
}
