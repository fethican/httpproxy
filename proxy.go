package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	Version = "0.1"
)

var (
	ProxyProto string
	ProxyTo    string
	BucketName string
	ServerPort string

	bytesProcessed = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "s3proxy_processed_bytes",
		Help: "Total of bytes proxied",
	})
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func (w *statusWriter) Log() {
	bytesProcessed.Observe(float64(w.length))
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {

	// Respond OK if request is to health endpoint
	if req.URL.Path == "/health" {
		res.WriteHeader(http.StatusOK)
		return
	}

	// Delegate to prometheus if request is to metrics endpoint
	if req.URL.Path == "/metrics" {
		promhttp.Handler().ServeHTTP(res, req)
		return
	}

	toURL := fmt.Sprintf("%s://%s/%s", ProxyProto, ProxyTo, BucketName)
	logrus.Infof("request %s from %s", req.URL.Path, toURL)

	_url, _ := url.Parse(toURL)

	proxy := httputil.NewSingleHostReverseProxy(_url)

	req.URL.Host = _url.Host
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = _url.Host

	sw := statusWriter{ResponseWriter: res}
	proxy.ServeHTTP(&sw, req)

	sw.Log()
}

func main() {
	ProxyProto = "https"
	ProxyTo = "s3.amazonaws.com"
	ServerPort = "8080"

	var ok bool
	BucketName, ok = os.LookupEnv("BUCKET_NAME")
	if !ok {
		logrus.Fatal("Provide a BUCKET_NAME")
	}

	if port, ok := os.LookupEnv("SERVER_PORT"); ok {
		ServerPort = port
	}

	if _proxyTo, ok := os.LookupEnv("PROXY_TO"); ok {
		ProxyTo = _proxyTo
	}

	if _proxyProto, ok := os.LookupEnv("PROXY_PROTO"); ok {
		ProxyProto = _proxyProto
	}

	logrus.Infof("Version: %s", Version)
	logrus.Infof("Proxying on port %s to %s://%s", ServerPort, ProxyProto, ProxyTo)
	logrus.Infof("Using bucket: %s", BucketName)

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", ServerPort), nil); err != nil {
		panic(err)
	}
}
