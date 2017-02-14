# aws-sign-proxy

HTTP proxy that signs requests for AWS service endpoints; e.g. Amazon Elasticsearch Service

This software is heavily inspired by [coreos/aws-auth-proxy](https://github.com/coreos/aws-auth-proxy).

## Usage

```bash
$ export AWS_ACCESS_KEY_ID=AKIAxxx
$ export AWS_SECRET_ACCESS_KEY=yyyyy
$ export AWS_REGION=ap-northeast-1
```

```bash
$ aws-auth-proxy --service-name es --upstream-host search-zzz.ap-northeast-1.es.amazonaws.com
$ open http://localhost:8080/_plugin/kibana
```

### Options

|Environment variable|Flag|Description|Required|Default|
|---|---|---|---|
|`AWS_ACCESS_KEY_ID`| |AWS access key ID|Required| |
|`AWS_SECRET_ACCESS_KEY`| |AWS secret access key|Required| |
|`AWS_REGION`|`--region`|AWS region|Required| |
|`AWS_SIGN_PROXY_SERVICE_NAME`|`--service-name`|AWS service name (e.g. `es`)|Required| |
|`AWS_SIGN_PROXY_UPSTREAM_HOST`|`--upstream-host`|Upstream endpoint|Required| |
|`AWS_SIGN_PROXY_UPSTREAM_SCHEME`|`--upstream-scheme`|Scheme for upstream endpoint| |`https`|
|`AWS_SIGN_PROXY_LISTEN_ADDRESS`|`--listen-address`|Address for proxy to listen on| |`:8080`|

## License

Original coreos/aws-auth-proxy is [released under Apache License Version 2.0](https://github.com/coreos/aws-auth-proxy/blob/9713146600f3aba055a5bfaf477af2a81dec272e/LICENSE).

This software is released under MIT License. [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
