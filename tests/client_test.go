package tests

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/armon/go-socks5"
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/pkg/client"
	"github.com/obluumuu/xor/tests/utils"
)

func splitHostPort(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, fmt.Errorf("split host port: %w", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("convert port to int: %w", err)
	}

	return host, port, nil
}

func TestClientWithHttpProxy(t *testing.T) {
	grpcLis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv, shutdown := utils.SetupServer(t, utils.WithListener(grpcLis))
	defer shutdown()

	proxyLis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	proxyUsername := "myusername"
	proxyPassword := ""

	proxy := goproxy.NewProxyHttpServer()
	auth.ProxyBasic(proxy, "realm", func(user, passwd string) bool {
		return user == proxyUsername && passwd == proxyPassword
	})

	srv := http.Server{Handler: proxy} //nolint:gosec
	defer func() {
		err := srv.Shutdown(ctx)
		require.NoError(t, err)
	}()

	go func() {
		err := srv.Serve(proxyLis)
		require.True(t, err == nil || errors.Is(err, http.ErrServerClosed))
	}()

	{
		host, port, err := splitHostPort(proxyLis.Addr().String())
		require.NoError(t, err)

		_, err = grpcSrv.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Tags: []string{"tag1"}, Schema: proxy_pb.Schema_SCHEMA_HTTP, Host: host, Port: uint32(port), Username: &proxyUsername, Password: &proxyPassword})
		require.NoError(t, err)

		_, err = grpcSrv.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Schema: proxy_pb.Schema_SCHEMA_HTTP, Host: "invalid_host", Port: 123, Tags: []string{"tag2"}})
		require.NoError(t, err)
	}

	createProxyBlockRes, err := grpcSrv.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: "proxy1", Tags: []string{"tag1"}})
	require.NoError(t, err)

	proxyBlockId, err := uuid.Parse(createProxyBlockRes.Id)
	require.NoError(t, err)

	proxySelector, err := client.NewProxySelector(grpcLis.Addr().String())
	require.NoError(t, err)

	roundTripper, err := proxySelector.NewRoundTripper(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, proxyBlockId) //nolint:gosec
	require.NoError(t, err)

	client := http.Client{
		Transport: roundTripper,
	}

	t.Run("http request", func(t *testing.T) {
		httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "somevalue", r.Header.Get("Test-Header"))
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`hello world`))
		}))
		defer httpSrv.Close()
		require.NotNil(t, httpSrv)

		url, err := url.Parse(httpSrv.URL)
		require.NoError(t, err)

		req := &http.Request{
			Method: http.MethodPatch,
			Header: map[string][]string{"Test-Header": {"somevalue"}},
			URL:    url,
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, resp.StatusCode)

		defer resp.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`hello world`), body)
	})
	t.Run("https request", func(t *testing.T) {
		httpsSrv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "somevalue", r.Header.Get("Test-Header"))
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`hello world`))
		}))
		defer httpsSrv.Close()
		require.NotNil(t, httpsSrv)

		url, err := url.Parse(httpsSrv.URL)
		require.NoError(t, err)

		req := &http.Request{
			Method: http.MethodPatch,
			Header: map[string][]string{"Test-Header": {"somevalue"}},
			URL:    url,
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, resp.StatusCode)

		defer resp.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`hello world`), body)
	})
}

func TestClientWithHttpsProxy(t *testing.T) {
	grpcLis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv, close := utils.SetupServer(t, utils.WithListener(grpcLis))
	defer close()

	proxyLis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	proxyUsername := "myusername"
	proxyPassword := "mypassword"

	proxy := goproxy.NewProxyHttpServer()

	auth.ProxyBasic(proxy, "realm", func(user, passwd string) bool {
		return user == proxyUsername && passwd == proxyPassword
	})

	srv := http.Server{Handler: proxy} //nolint:gosec
	defer func() {
		err := srv.Shutdown(ctx)
		require.NoError(t, err)
	}()

	go func() {
		err := srv.ServeTLS(proxyLis, "certs/server.crt", "certs/server.key")
		require.True(t, err == nil || errors.Is(err, http.ErrServerClosed))
	}()

	{
		host, port, err := splitHostPort(proxyLis.Addr().String())
		require.NoError(t, err)

		_, err = grpcSrv.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Schema: proxy_pb.Schema_SCHEMA_HTTPS, Host: host, Port: uint32(port), Username: &proxyUsername, Password: &proxyPassword, Tags: []string{"tag1"}})
		require.NoError(t, err)

		_, err = grpcSrv.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Host: "invalid_host", Port: 123, Tags: []string{"tag2"}})
		require.NoError(t, err)
	}

	createProxyBlockRes, err := grpcSrv.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: "proxy1", Tags: []string{"tag1"}})
	require.NoError(t, err)

	proxyBlockId, err := uuid.Parse(createProxyBlockRes.Id)
	require.NoError(t, err)

	proxySelector, err := client.NewProxySelector(grpcLis.Addr().String())
	require.NoError(t, err)

	roundTripper, err := proxySelector.NewRoundTripper(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, proxyBlockId) //nolint:gosec
	require.NoError(t, err)

	client := http.Client{
		Transport: roundTripper,
	}

	t.Run("http request", func(t *testing.T) {
		httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "somevalue", r.Header.Get("Test-Header"))
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`hello world`))
		}))
		defer httpSrv.Close()
		require.NotNil(t, httpSrv)

		url, err := url.Parse(httpSrv.URL)
		require.NoError(t, err)

		req := &http.Request{
			Method: http.MethodPatch,
			Header: map[string][]string{"Test-Header": {"somevalue"}},
			URL:    url,
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, resp.StatusCode)

		defer resp.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`hello world`), body)
	})
	t.Run("https request", func(t *testing.T) {
		httpsSrv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "somevalue", r.Header.Get("Test-Header"))
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`hello world`))
		}))
		defer httpsSrv.Close()
		require.NotNil(t, httpsSrv)

		url, err := url.Parse(httpsSrv.URL)
		require.NoError(t, err)

		req := &http.Request{
			Method: http.MethodPatch,
			Header: map[string][]string{"Test-Header": {"somevalue"}},
			URL:    url,
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, resp.StatusCode)

		defer resp.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`hello world`), body)
	})
}

func TestClientWithSocks5Proxy(t *testing.T) {
	grpcLis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv, close := utils.SetupServer(t, utils.WithListener(grpcLis))
	defer close()

	proxyLis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	proxyUsername := "myusername"
	proxyPassword := "mypassword"

	proxySrv, err := socks5.New(&socks5.Config{Credentials: socks5.StaticCredentials{proxyUsername: proxyPassword}})
	require.NoError(t, err)

	defer func() {
		err := proxyLis.Close()
		require.NoError(t, err)
	}()

	go func() {
		err := proxySrv.Serve(proxyLis)
		require.True(t, err != nil || errors.Is(err, net.ErrClosed))
	}()

	{
		host, port, err := splitHostPort(proxyLis.Addr().String())
		require.NoError(t, err)

		_, err = grpcSrv.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Schema: proxy_pb.Schema_SCHEMA_SOCKS5, Host: host, Port: uint32(port), Username: &proxyUsername, Password: &proxyPassword, Tags: []string{"tag1"}})
		require.NoError(t, err)

		_, err = grpcSrv.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Host: "invalid_host", Port: 123, Tags: []string{"tag2"}})
		require.NoError(t, err)
	}

	createProxyBlockRes, err := grpcSrv.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: "proxy1", Tags: []string{"tag1"}})
	require.NoError(t, err)

	proxyBlockId, err := uuid.Parse(createProxyBlockRes.Id)
	require.NoError(t, err)

	proxySelector, err := client.NewProxySelector(grpcLis.Addr().String())
	require.NoError(t, err)

	roundTripper, err := proxySelector.NewRoundTripper(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, proxyBlockId) //nolint:gosec
	require.NoError(t, err)

	client := http.Client{
		Transport: roundTripper,
	}

	t.Run("http request", func(t *testing.T) {
		httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "somevalue", r.Header.Get("Test-Header"))
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`hello world`))
		}))
		defer httpSrv.Close()
		require.NotNil(t, httpSrv)

		url, err := url.Parse(httpSrv.URL)
		require.NoError(t, err)

		req := &http.Request{
			Method: http.MethodPatch,
			Header: map[string][]string{"Test-Header": {"somevalue"}},
			URL:    url,
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, resp.StatusCode)

		defer resp.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`hello world`), body)
	})
	t.Run("https request", func(t *testing.T) {
		httpsSrv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "somevalue", r.Header.Get("Test-Header"))
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`hello world`))
		}))
		defer httpsSrv.Close()
		require.NotNil(t, httpsSrv)

		url, err := url.Parse(httpsSrv.URL)
		require.NoError(t, err)

		req := &http.Request{
			Method: http.MethodPatch,
			Header: map[string][]string{"Test-Header": {"somevalue"}},
			URL:    url,
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, resp.StatusCode)

		defer resp.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`hello world`), body)
	})
}
