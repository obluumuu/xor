package client

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/internal/server/utils"
)

var _ http.RoundTripper = (*RoundTripper)(nil)

type ProxySelector struct {
	Conn   *grpc.ClientConn
	Client proxy_pb.ProxyServiceClient
}

// Connect to proxy selector grpc server and return [ProxySelector] on success
func NewProxySelector(addr string) (*ProxySelector, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", addr, err)
	}
	client := proxy_pb.NewProxyServiceClient(conn)

	return &ProxySelector{Client: client, Conn: conn}, nil
}

// Get proxies by [proxyBlockId] and put them and [t] in [*RoundTripper]
//
// If [t] is nil, use [http.DefaultTransport]
func (c *ProxySelector) NewRoundTripper(t *http.Transport, proxyBlockId uuid.UUID) (*RoundTripper, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := c.Client.GetProxiesByProxyBlockId(ctx, &proxy_pb.GetProxiesByProxyBlockIdRequest{Id: proxyBlockId.String()})
	if err != nil {
		return nil, fmt.Errorf("get proxies by proxy block id: %w", err)
	}

	proxies := resp.Proxies
	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies matching this proxy block id")
	}

	if t == nil {
		t = http.DefaultTransport.(*http.Transport)
	}

	return &RoundTripper{T: t, Proxies: proxies}, nil
}

// Close conn to grpc server
func (c *ProxySelector) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return fmt.Errorf("close grpc conn: %w", err)
	}
	return nil
}

type RoundTripper struct {
	T       *http.Transport
	Proxies []*proxy_pb.GetProxiesByProxyBlockIdResponse_Proxy
}

// Select random proxy
func (r *RoundTripper) getRandProxyURL() *url.URL {
	proxy := r.Proxies[rand.Intn(len(r.Proxies))]

	var userInfo *url.Userinfo
	if proxy.Username != nil {
		if proxy.Password != nil {
			userInfo = url.UserPassword(*proxy.Username, *proxy.Password)
		} else {
			userInfo = url.User(*proxy.Username)
		}
	}

	return &url.URL{
		Scheme: utils.SchemaProtoToString(proxy.Schema),
		User:   userInfo,
		Host:   fmt.Sprintf("%s:%d", proxy.Host, proxy.Port),
	}
}

// Select random proxy and use it for [RoundTrip]
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t := r.T.Clone()
	t.Proxy = http.ProxyURL(r.getRandProxyURL())
	//nolint
	return t.RoundTrip(req)
}
