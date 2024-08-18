package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/tests/utils"
)

func getProxyBlockTagNames(tags []*proxy_pb.GetProxyBlockResponse_Tag) []string {
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}

func TestCreateGetProxyBlock(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "proxyblock"
	createRes, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name})
	require.NoError(t, err)

	id, err := uuid.Parse(createRes.Id)
	require.NoError(t, err)

	getRes, err := s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: id.String()})
	require.NoError(t, err)

	require.Equal(t, &proxy_pb.GetProxyBlockResponse{Id: id.String(), Name: name, Tags: []*proxy_pb.GetProxyBlockResponse_Tag{}}, getRes)
}

func TestCreateGetProxyBlockWithTags(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name, Description: description, Tags: tags1})
	require.NoError(t, err)

	id, err := uuid.Parse(createRes1.Id)
	require.NoError(t, err)

	getRes1, err := s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: id.String()})
	require.NoError(t, err)
	require.Equal(t, &proxy_pb.GetProxyBlockResponse{Id: id.String(), Name: name, Description: description, Tags: getRes1.Tags}, getRes1)
	require.ElementsMatch(t, tags1, getProxyBlockTagNames(getRes1.Tags))

	tags2 := []string{"amogus1", "amogus2", "kek1"}
	createRes2, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name, Description: description, Tags: tags2})
	require.NoError(t, err)

	getRes2, err := s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: createRes2.Id})
	require.NoError(t, err)
	require.ElementsMatch(t, tags2, getProxyBlockTagNames(getRes2.Tags))
}

func TestDeleteProxyBlock(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name, Description: description, Tags: tags1})
	require.NoError(t, err)

	tags2 := []string{"amogus1", "amogus2", "kek1"}
	createRes2, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name, Description: description, Tags: tags2})
	require.NoError(t, err)

	_, err = s.DeleteProxyBlock(context.Background(), &proxy_pb.DeleteProxyBlockRequest{Id: createRes1.Id})
	require.NoError(t, err)

	_, err = s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: createRes1.Id})
	require.ErrorContains(t, err, "NotFound")

	getRes2, err := s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: createRes2.Id})
	require.NoError(t, err)
	require.ElementsMatch(t, tags2, getProxyBlockTagNames(getRes2.Tags))
}

func TestUpdateProxyBlock(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name, Description: description, Tags: tags1})
	require.NoError(t, err)

	newName := "kek2"
	_, err = s.UpdateProxyBlock(context.Background(), &proxy_pb.UpdateProxyBlockRequest{Id: createRes1.Id, Name: newName, FieldMask: []string{"name"}})
	require.NoError(t, err)

	getRes1, err := s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: createRes1.Id})
	require.NoError(t, err)
	require.Equal(t, &proxy_pb.GetProxyBlockResponse{Id: createRes1.Id, Name: newName, Description: description, Tags: getRes1.Tags}, getRes1)
	require.ElementsMatch(t, tags1, getProxyBlockTagNames(getRes1.Tags))
}

func TestUpdateProxyBlockWithTags(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: name, Description: description, Tags: tags1})
	require.NoError(t, err)

	newTags1 := []string{"amogus1", "amogus2", "kek1"}
	_, err = s.UpdateProxyBlock(context.Background(), &proxy_pb.UpdateProxyBlockRequest{Id: createRes1.Id, Tags: newTags1, FieldMask: []string{"tags"}})
	require.NoError(t, err)

	getRes1, err := s.GetProxyBlock(context.Background(), &proxy_pb.GetProxyBlockRequest{Id: createRes1.Id})
	require.NoError(t, err)
	require.Equal(t, &proxy_pb.GetProxyBlockResponse{Id: createRes1.Id, Name: name, Description: description, Tags: getRes1.Tags}, getRes1)
	require.ElementsMatch(t, newTags1, getProxyBlockTagNames(getRes1.Tags))
}

func TestGetProxiesByProxyBlock(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	createProxyRes1, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy1", Tags: []string{"tag1", "tag2"}})
	require.NoError(t, err)

	createProxyRes2, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy2", Tags: []string{"tag1", "tag2", "tag3"}})
	require.NoError(t, err)

	_, err = s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy3", Tags: []string{"tag1"}})
	require.NoError(t, err)

	_, err = s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy4", Tags: []string{}})
	require.NoError(t, err)

	{
		createProxyBlockRes, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: "proxy_block1", Tags: []string{"tag1", "tag2"}})
		require.NoError(t, err)

		getProxiesRes, err := s.GetProxiesByProxyBlockId(context.Background(), &proxy_pb.GetProxiesByProxyBlockIdRequest{Id: createProxyBlockRes.Id})
		require.NoError(t, err)

		ids := make([]string, 0, len(getProxiesRes.Proxies))
		for _, proxy := range getProxiesRes.Proxies {
			ids = append(ids, proxy.Id)
		}
		require.ElementsMatch(t, []string{createProxyRes1.Id, createProxyRes2.Id}, ids)
	}

	{
		createProxyBlockRes, err := s.CreateProxyBlock(context.Background(), &proxy_pb.CreateProxyBlockRequest{Name: "proxy_block2", Tags: []string{}})
		require.NoError(t, err)

		getProxiesRes, err := s.GetProxiesByProxyBlockId(context.Background(), &proxy_pb.GetProxiesByProxyBlockIdRequest{Id: createProxyBlockRes.Id})
		require.NoError(t, err)
		require.Empty(t, getProxiesRes.Proxies)
	}
}
