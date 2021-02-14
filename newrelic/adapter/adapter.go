// Package adapter is the implementation of Newrelic LogSpout Adapter
package adapter

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"

    "github.com/gliderlabs/logspout/router"
    "github.com/gojektech/heimdall"
    "github.com/gojektech/heimdall/httpclient"
)

// New method of Adapter:
func New(config Configuration) *Adapter {
    backoff := heimdall.NewConstantBackoff(config.BackoffInterval, config.JitterInterval)
    retrier := heimdall.NewRetrier(backoff)
    httpClient := httpclient.NewClient(
        httpclient.WithHTTPTimeout(config.HTTPTimeout),
        httpclient.WithRetrier(retrier),
        httpclient.WithRetryCount(int(config.RequestRetryCount)),
    )

    adapter := &Adapter{
        Config:         config,
        HTTPClient:     httpClient,
        Logger:         log.New(os.Stdout, config.Hostname + " ", log.LstdFlags),
        Queue:          make(chan Line),
    }

    go adapter.readQueue()
    return adapter
}

// getHost method is for deciding what to choose as a hostname:
func (adapter *Adapter) getHost(containerHostname string) string {
    host := containerHostname
    if (adapter.Config.Hostname != "") {
        host = adapter.Config.Hostname
    }
    return host
}

// Stream method is for streaming the messages:
func (adapter *Adapter) Stream(logstream chan *router.Message) {
    for m := range logstream {
        if m.Data == "" {
            continue
        }

        attributes := Attributes{
            Container:  ContainerInfo{
                Name:   strings.Trim(m.Container.Name, "/"),
                ID:     m.Container.ID,
                Config: ContainerConfig{
                    Image:      m.Container.Config.Image,
                    Hostname:   m.Container.Config.Hostname,
                    Labels:     m.Container.Config.Labels,
                },
            },
            Source:     m.Source,
            Hostname:   adapter.getHost(m.Container.Config.Hostname),
        }

        adapter.Queue <- Line{
            Message: m.Data,
            Timestamp:  time.Now().Unix(),
            Attributes: attributes,
        }
    }
}

// readQueue is a method for reading from queue:
func (adapter *Adapter) readQueue() {
    buffer := make([]Line, 0)
    bufferSize := 0

    timeout := time.NewTimer(adapter.Config.FlushInterval)

    for {
        select {
        case msg := <-adapter.Queue:
            if bufferSize >= int(adapter.Config.MaxBufferSize) {
                timeout.Stop()
                adapter.flushBuffer(buffer)
                buffer = make([]Line, 0)
                bufferSize = 0
            }

            buffer = append(buffer, msg)
            bufferSize += len(msg.Message)

        case <-timeout.C:
            if bufferSize > 0 {
                adapter.flushBuffer(buffer)
                buffer = make([]Line, 0)
                bufferSize = 0
            }
        }

        timeout.Reset(adapter.Config.FlushInterval)
    }
}

// flushBuffer is a method for flushing the lines:
func (adapter *Adapter) flushBuffer(buffer []Line) {
    var data bytes.Buffer

    body := struct {
        Lines []Line `json:"logs"`
    }{
        Lines: buffer,
    }

    if error := json.NewEncoder(&data).Encode(body); error != nil {
        adapter.Logger.Println(
            fmt.Errorf(
                "JSON Encoding Error: %s",
                error.Error(),
            ),
        )
        return
    }

    urlValues := url.Values{}
    urlValues.Add("hostname", "newrelic_logspout")
    url := "https://" + adapter.Config.NewrelicURL + "?" + urlValues.Encode()
    req, _ := http.NewRequest(http.MethodPost, url, &data)
    req.Header.Set("user-agent", "logspout-newrelic/" + os.Getenv("BUILD_VERSION"))
    req.Header.Set("Content-Type", "application/json; charset=UTF-8")
    req.Header.Set("X-License-Key", adapter.Config.NewrelicKey)

    resp, err := adapter.HTTPClient.Do(req)

    if err != nil {
        adapter.Logger.Println(
            fmt.Errorf(
                "HTTP Client Post Request Error: %s",
                err.Error(),
            ),
        )
        return
    }

    if resp != nil {
        if resp.StatusCode != http.StatusOK {
            adapter.Logger.Println(
                fmt.Errorf(
                    "Received Status Code: %s While Sending Message",
                    resp.StatusCode,
                ),
            )
        }
        defer resp.Body.Close()
    }
}
