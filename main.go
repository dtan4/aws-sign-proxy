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

	flags.StringVar(&opts.awsRegion, "aws-region", os.Getenv("AWS_REGION"), "AWS region")
	flags.StringVar(&opts.serviceName, "service-name", "", "AWS service name")
	flags.StringVar(&opts.upstreamHost, "upstream-host", "", "Upstream endpoint")
	flags.StringVar(&opts.upstreamScheme, "upstream-scheme", defaultUpstreamScheme, "Scheme for upstream endpoint")
	flags.StringVar(&opts.listenAddress, "listen-address", defaultListenAddress, "Address for proxy to listen on")

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
