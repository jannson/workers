package fetch

import (
	"net/http"
	"syscall/js"
)

// Client is an HTTP client.
type Client struct {
	// namespace - Objects that Fetch API belongs to. Default is Global
	namespace js.Value
}

// applyOptions applies client options.
func (c *Client) applyOptions(opts []ClientOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// HTTPClient returns *http.Client.
func (c *Client) HTTPClient(redirect RedirectMode) *http.Client {
	return &http.Client{
		Transport: &transport{
			namespace: c.namespace,
			redirect:  redirect,
		},
	}
}

// ClientOption is a type that represents an optional function.
type ClientOption func(*Client)

// WithBinding changes the objects that Fetch API belongs to.
// This is useful for service bindings, mTLS, etc.
func WithBinding(bind js.Value) ClientOption {
	return func(c *Client) {
		c.namespace = bind
	}
}

// NewClient returns new Client
func NewClient(opts ...ClientOption) *Client {
	// In TinyGo's wasm_exec.js, js.Global() can be a Proxy used to inject a per-request `context`.
	// Some Cloudflare runtime APIs (notably `fetch`) are strict about the receiver (`this`) and may
	// throw "Illegal invocation" when called with the Proxy as `this`.
	//
	// Use the real global object when available.
	namespace := js.Global()
	if globalThis := namespace.Get("globalThis"); !globalThis.IsUndefined() && !globalThis.IsNull() {
		namespace = globalThis
	}
	c := &Client{
		namespace: namespace,
	}
	c.applyOptions(opts)

	return c
}
