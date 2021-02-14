// Package adapter is the implementation of Newrelic LogSpout Adapter
package adapter

import (
    "log"
    "time"

    "github.com/gojektech/heimdall"
)

/*
    Definitions of All Structs
*/

// Configuration is Configuration Struct for Newrelic Adapter:
type Configuration struct {
    BackoffInterval     time.Duration
    FlushInterval       time.Duration
    Hostname            string
    HTTPTimeout         time.Duration
    JitterInterval      time.Duration
    NewrelicKey         string
    NewrelicURL         string
    MaxBufferSize       uint64
    RequestRetryCount   uint64
    Tags                string
}

// Adapter structure:
type Adapter struct {
    Config          Configuration
    HTTPClient      heimdall.Client
    Logger          *log.Logger
    Queue           chan Line
}

// Line structure for the queue of Adapter:
type Line struct {
    Timestamp   int64      `json:"timestamp"`
    Message     string     `json:"message"`
    Attributes  Attributes `json:"attributes"`
}

// Message structure:
type Attributes struct {
    Container   ContainerInfo `json:"container"`
    Source      string        `json:"source"`
    Hostname    string        `json:"hostname"`
}

// ContainerInfo structure for the Container of Message:
type ContainerInfo struct {
    Name    string          `json:"name"`
    ID      string          `json:"id"`
    Config  ContainerConfig `json:"config"`
}

// ContainerConfig structure for the Config of ContainerInfo:
type ContainerConfig struct {
    Image       string              `json:"image"`
    Hostname    string              `json:"hostname"`
    Labels      map[string]string   `json:"labels"`
}
