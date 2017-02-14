package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
	idLength = 8
)

// AWSProxy represents proxy object to sign request
type AWSProxy struct {
	region    string
	service   string
	signer    *v4.Signer
	targetURL *url.URL
}

// NewAWSProxy returns new AWSProxy object
func NewAWSProxy(region, service string, signer *v4.Signer, targetURL *url.URL) *AWSProxy {
	return &AWSProxy{
		region:    region,
		service:   service,
		signer:    signer,
		targetURL: targetURL,
	}
}

// ServeHTTP signs the given request and proceed real API request
func (h *AWSProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestID := generateID()

	log.Printf("[%s] > %s %s\n", requestID, r.Method, r.URL.Path)

	proxyURL := *r.URL
	proxyURL.Host = h.targetURL.Host
	proxyURL.Scheme = h.targetURL.Scheme

	req, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Printf("[%s] < %s %s %d\n", requestID, r.Method, r.URL.Path, http.StatusBadRequest)
		return
	}

	if _, err = h.signer.Sign(req, nil, h.service, h.region, time.Now()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Printf("[%s] < %s %s %d\n", requestID, r.Method, r.URL.Path, http.StatusBadRequest)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Printf("[%s] < %s %s %d\n", requestID, r.Method, r.URL.Path, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())

	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Printf("[%s] < %s %s %d\n", requestID, r.Method, r.URL.Path, http.StatusBadRequest)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(buf.Bytes())
	log.Printf("[%s] < %s %s %d\n", requestID, r.Method, r.URL.Path, resp.StatusCode)
}

func copyHeaders(src, dst http.Header) {
	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

func generateID() string {
	b := make([]byte, idLength)
	rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}
