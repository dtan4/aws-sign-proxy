package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/spf13/pflag"
)

const (
	defaultUpstreamScheme = "https"
	defaultListenAddress  = ":8080"
)

type Options struct {
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsRegion          string
	serviceName        string
	upstreamHost       string
	upstreamScheme     string
	listenAddress      string
}

func main() {
	opts := optsFromEnv()

	flags := pflag.NewFlagSet("aws-es-auth-proxy", pflag.ExitOnError)

	flags.StringVar(&opts.serviceName, "service-name", "", "AWS service name")
	flags.StringVar(&opts.upstreamHost, "upstream-host", "", "Upstream endpoint")
	flags.StringVar(&opts.upstreamScheme, "upstream-scheme", defaultUpstreamScheme, "Scheme for upstream endpoint")
	flags.StringVar(&opts.listenAddress, "listen-address", defaultListenAddress, "Address for proxy to listen on")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	creds := credentials.NewStaticCredentials(opts.awsAccessKeyID, opts.awsSecretAccessKey, "")
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

func optsFromEnv() *Options {
	opts := &Options{}

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		opts.awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		opts.awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	if os.Getenv("AWS_REGION") != "" {
		opts.awsRegion = os.Getenv("AWS_REGION")
	}

	return opts
}
