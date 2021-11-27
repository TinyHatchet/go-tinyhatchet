package mouseion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
}

// Logger sends log messages to MouseionHost using HTTPClient. MouseionHost
// is the only required field.
type Logger struct {
	HTTPClient          *http.Client
	LogErrors           bool
	MouseionHost        string
	APIToken, APISecret string
	DefaultTags         []string

	// AutoTagger is a function called when a log message is written. AutoTagger
	// parses the arg and determines if any tags should be added to the defaults
	AutoTagger func(defaultTags []string, arg interface{}) []string
}

func (logger *Logger) Print(v ...interface{}) {
	logger.send(fmt.Sprint(v...), logger.ArgsToTags(v))
}

func (logger *Logger) Printf(format string, v ...interface{}) {
	logger.send(fmt.Sprintf(format, v...), logger.ArgsToTags(v))
}

func (logger *Logger) Println(v ...interface{}) {
	logger.Print(v...)
}

func (logger *Logger) send(text string, tags []string) {
	err := send(logger.HTTPClient, logger.MouseionHost, logger.APIToken, logger.APISecret, text, tags)
	if err != nil && logger.LogErrors {
		log.Println(err)
	}
}

func (logger *Logger) ArgsToTags(args ...interface{}) []string {
	if logger.AutoTagger == nil {
		return logger.DefaultTags
	}
	tags := []string{}
	for _, arg := range args {
		autoTags := logger.AutoTagger(logger.DefaultTags, arg)
		if len(autoTags) > 0 {
			tags = append(tags, autoTags...)
		}
	}
	if len(tags) == 0 {
		tags = nil
	}
	return tags
}

func send(client *http.Client, host, username, password, text string, tags []string) error {
	if client != nil {
		client = http.DefaultClient
	}
	entry := LogEntry{Timestamp: time.Now(), Text: text, Tags: tags}
	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, jsonURL(host), bytes.NewBuffer(entryJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%d %s", resp.StatusCode, body)
	}
	return nil
}
func jsonURL(host string) string {
	return host + "/ingest.json"
}
