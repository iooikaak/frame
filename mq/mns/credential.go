package mns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Credential mns credential
type Credential interface {
	Sign(http.Header, string, string, string) (string, error)
}

type credential struct {
}

// NewCredential new credential
func NewCredential() Credential {
	c := new(credential)
	return c
}

// Sign calc sign
func (c *credential) Sign(h http.Header, method, resource, accessSecret string) (string, error) {
	var date = time.Now().UTC().Format(http.TimeFormat)
	if v := h.Get("Date"); v != "" {
		date = v
	}

	var contentMD5 = h.Get("Content-MD5")
	var contentType = h.Get("Content-Type")

	var mnsHeaders []string
	for k := range h {
		if strings.HasPrefix(strings.ToLower(k), "x-mns-") {
			mnsHeaders = append(mnsHeaders, strings.ToLower(k)+":"+strings.TrimSpace(h.Get(k)))
		}
	}

	sort.Sort(sort.StringSlice(mnsHeaders))

	stringToSign := method + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		strings.Join(mnsHeaders, "\n") + "\n" +
		"/" + resource

	hmaced := hmac.New(sha1.New, []byte(accessSecret))
	if _, err := hmaced.Write([]byte(stringToSign)); err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString(hmaced.Sum(nil))
	return s, nil
}
