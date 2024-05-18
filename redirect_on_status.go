package traefik_plugin_redirect_on_status

import (
	"context"
	"fmt"
	"github.com/m-riedel/traefik-plugin-redirect-on-status/pkg/interceptor"
	"github.com/m-riedel/traefik-plugin-redirect-on-status/pkg/types"
	"net/http"
	"slices"
)

// Config holds configuration of the plugin
type Config struct {
	RedirectUri  string
	RedirectCode int
	Status       []string
	Method       []string
}

// CreateConfig creates Config with default attributes
func CreateConfig() *Config {
	return &Config{
		Status:       nil,
		RedirectUri:  "",
		RedirectCode: http.StatusTemporaryRedirect,
		Method:       make([]string, 0),
	}
}

// RedirectOnStatusPlugin holds the needed data for the plugin
type RedirectOnStatusPlugin struct {
	next           http.Handler
	name           string
	redirectUri    string
	httpCodeRanges types.HTTPCodeRanges
	redirectCode   int
	methods        []string
}

// New creates RedirectOnStatusPlugin from the given fields. Checks that the Config is correct
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.RedirectUri) == 0 {
		return nil, fmt.Errorf("uri cannot be empty")
	}
	if config.RedirectCode == 0 {
		return nil, fmt.Errorf("redirectCode cannot be empty")
	}
	// Only Allow Temporary Redirects, since this is an action that
	// happens due to a condition that may not be true every time
	if config.RedirectCode != 302 &&
		config.RedirectCode != 303 &&
		config.RedirectCode != 307 {
		return nil, fmt.Errorf("redirectCode must be a temporary redirection code")
	}

	if config.Status == nil || len(config.Status) == 0 {
		return nil, fmt.Errorf("code cannot be empty")
	}

	httpCodeRanges, err := types.NewHTTPCodeRanges(config.Status)

	if err != nil {
		return nil, err
	}

	return &RedirectOnStatusPlugin{
		next:           next,
		name:           name,
		httpCodeRanges: httpCodeRanges,
		redirectUri:    config.RedirectUri,
		redirectCode:   config.RedirectCode,
		methods:        config.Method,
	}, nil
}

func (e RedirectOnStatusPlugin) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// If no method is given or methods contains the request method,
	// use the interceptor
	if len(e.methods) == 0 ||
		slices.Contains(e.methods, request.Method) {
		rwInterceptor := interceptor.NewRedirectOnStatusInterceptor(rw, request, e.redirectUri, e.redirectCode, e.httpCodeRanges)
		e.next.ServeHTTP(rwInterceptor, request)
		return
	}

	e.next.ServeHTTP(rw, request)
}
