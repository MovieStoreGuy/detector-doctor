package client

import (
	"crypto/tls"
	"net/http"
)

// WithFlagSet is used when being used with evaluation arguments provided by a user or config
func WithFlagSet(cond bool, method func(*http.Client)) func(*http.Client) {
	if cond {
		return method
	}
	return noOpSetting
}

func noOpSetting(*http.Client) {}

// DisableHTTP2 follows the default standard on how to disable HTTP/2 with the default http client
func DisableHTTP2(cli *http.Client) {
	if trans, cast := cli.Transport.(*http.Transport); trans != nil && cast {
		trans.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
	}
}
