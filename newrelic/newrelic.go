package newrelic

import (
    "errors"
    "log"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/gliderlabs/logspout/router"
    "github.com/mgreenw/logspout-newrelic/newrelic/adapter"
)

/*
    Common Functions
*/

// Getting Uint Variable from Environment:
func getUintOpt(name string, dfault uint64) uint64 {
    if result, err := strconv.ParseUint(os.Getenv(name), 10, 64); err == nil {
        return result
    }
    return dfault
}

// Getting Duration Variable from Environment:
func getDurationOpt(name string, dfault time.Duration) time.Duration {
    if result, err := strconv.ParseInt(os.Getenv(name), 10, 64); err == nil {
        return time.Duration(result)
    }
    return dfault
}

// Getting String Variable from Environment:
func getStringOpt(name, dfault string) string {
    if value := os.Getenv(name); value != "" {
        return value
    }
    return dfault
}

func init() {
    router.AdapterFactories.Register(NewNewrelicRouter, "newrelic")

    filterLabels := make([]string, 0)
    if filterLabelsValue := os.Getenv("FILTER_LABELS"); filterLabelsValue != "" {
        filterLabels = strings.Split(filterLabelsValue, ",")
    }

    filterSources := make([]string, 0)
    if filterSourcesValue := os.Getenv("FILTER_SOURCES"); filterSourcesValue != "" {
        filterSources = strings.Split(filterSourcesValue, ",")
    }

    r := &router.Route{
        Adapter:        "newrelic",
        FilterName:     getStringOpt("FILTER_NAME", ""),
        FilterID:       getStringOpt("FILTER_ID", ""),
        FilterLabels:   filterLabels,
        FilterSources:  filterSources,
    }

    if err := router.Routes.Add(r); err != nil {
        log.Fatal("Cannot Add New Route: ", err.Error())
    }
}

// NewNewrelicRouter creates adapter:
func NewNewrelicRouter(route *router.Route) (router.LogAdapter, error) {
    newrelicKey := os.Getenv("NEWRELIC_KEY")
    if newrelicKey == "" {
        return nil, errors.New("Cannot Find Environment Variable \"NEWRELIC_KEY\"")
    }

    if os.Getenv("INACTIVITY_TIMEOUT") == "" {
        os.Setenv("INACTIVITY_TIMEOUT", "1m")
    }

    config := adapter.Configuration{
        BackoffInterval:    getDurationOpt("HTTP_CLIENT_BACKOFF", 2) * time.Millisecond,
        FlushInterval:      getDurationOpt("FLUSH_INTERVAL", 250) * time.Millisecond,
        Hostname:           os.Getenv("HOSTNAME"),
        HTTPTimeout:        getDurationOpt("HTTP_CLIENT_TIMEOUT", 30) * time.Second,
        JitterInterval:     getDurationOpt("HTTP_CLIENT_JITTER", 5) * time.Millisecond,
        NewrelicKey:        newrelicKey,
        NewrelicURL:          getStringOpt("NEWRELIC_URL", "log-api.newrelic.com/log/v1"),
        MaxBufferSize:      getUintOpt("MAX_BUFFER_SIZE", 2) * 1024 * 1024,
        RequestRetryCount:  getUintOpt("MAX_REQUEST_RETRY", 5),
        Tags:               os.Getenv("TAGS"),
    }

    return adapter.New(config), nil
}
