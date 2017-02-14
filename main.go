package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/spf13/pflag"
)

const (
	defaultUpstreamScheme = "https"
	defaultListenAddress  = ":8080"
)

// Options represents options of aws-sign-proxy
type Options struct {
	awsRegion      string
	serviceName    string
	upstreamHost   string
	upstreamScheme string
	listenAddress  string
}

func main() {
	opts := Options{}

	flags := pflag.NewFlagSet("aws-sign-proxy", pflag.ExitOnError)

	var upstreamScheme, listenAddress string

	if os.Getenv("AWS_SIGN_PROXY_UPSTREAM_SCHEME") == "" {
		upstreamScheme = defaultUpstreamScheme
	} else {
		upstreamScheme = os.Getenv("AWS_SIGN_PROXY_UPSTREAM_SCHEME")
	}

	if os.Getenv("AWS_SIGN_PROXY_LISTEN_ADDRESS") == "" {
		listenAddress = defaultListenAddress
	} else {
		listenAddress = os.Getenv("AWS_SIGN_PROXY_LISTEN_ADDRESS")
	}

	flags.StringVar(&opts.awsRegion, "aws-region", os.Getenv("AWS_REGION"), "AWS region")
	flags.StringVar(&opts.serviceName, "service-name", os.Getenv("AWS_SIGN_PROXY_SERVICE_NAME"), "AWS service name")
	flags.StringVar(&opts.upstreamHost, "upstream-host", os.Getenv("AWS_SIGN_PROXY_UPSTREAM_HOST"), "Upstream endpoint")
	flags.StringVar(&opts.upstreamScheme, "upstream-scheme", upstreamScheme, "Scheme for upstream endpoint")
	flags.StringVar(&opts.listenAddress, "listen-address", listenAddress, "Address for proxy to listen on")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	creds := defaults.CredChain(defaults.Config(), defaults.Handlers())
	if _, err := creds.Get(); err != nil {
		log.Fatal(err)
	}

	signer := v4.NewSigner(creds)
	proxyHandler := NewAWSProxy(opts.awsRegion, opts.serviceName, signer, &url.URL{
		Host:   opts.upstreamHost,
		Scheme: opts.upstreamScheme,
	})

	fmt.Printf("Listening on %s\n", opts.listenAddress)

	if err := http.ListenAndServe(opts.listenAddress, proxyHandler); err != nil {
		log.Fatal(err)
	}
}
