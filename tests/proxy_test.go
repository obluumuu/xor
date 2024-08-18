package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/tests/utils"
)

func getProxyTagNames(tags []*proxy_pb.GetProxyResponse_Tag) []string {
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}

func TestCreateGetProxy(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	t.Run("not_existed_id", func(t *testing.T) {
		createRes, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: "proxy_check_id"})
		require.NoError(t, err)

		id, err := uuid.Parse(createRes.Id)
		require.NoError(t, err)

		wrongId := uuid.New()
		require.NotEqual(t, id.String(), wrongId.String())

		_, err = s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: wrongId.String()})
		require.ErrorContains(t, err, "NotFound")
	})
	t.Run("without_tags", func(t *testing.T) {
		name := "proxy_without_tags"
		host := "127.0.0.1"
		port := uint32(80)
		username := "username"
		createRes, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Schema: proxy_pb.Schema_SCHEMA_HTTP, Host: host, Port: port, Username: &username})
		require.NoError(t, err)

		getRes, err := s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes.Id})
		require.NoError(t, err)

		require.Equal(t, &proxy_pb.GetProxyResponse{
			Id:       createRes.Id,
			Name:     name,
			Host:     host,
			Port:     port,
			Schema:   proxy_pb.Schema_SCHEMA_HTTP,
			Username: &username,
			Tags:     []*proxy_pb.GetProxyResponse_Tag{},
		}, getRes)
	})
	t.Run("with_tags", func(t *testing.T) {
		name := "proxy_with_tags"
		description := "amogus"
		host := "aaaaaa"
		tags1 := []string{"kek1", "kek2", "kek3"}
		createRes1, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Description: description, Host: host, Tags: tags1})
		require.NoError(t, err)

		getRes1, err := s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes1.Id})
		require.NoError(t, err)
		require.Equal(t, &proxy_pb.GetProxyResponse{Id: createRes1.Id, Name: name, Description: description, Host: host, Tags: getRes1.Tags}, getRes1)
		require.ElementsMatch(t, tags1, getProxyTagNames(getRes1.Tags))

		tags2 := []string{"amogus1", "amogus2", "kek1"}
		createRes2, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Description: description, Host: host, Tags: tags2})
		require.NoError(t, err)

		getRes2, err := s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes2.Id})
		require.NoError(t, err)
		require.ElementsMatch(t, tags2, getProxyTagNames(getRes2.Tags))
	})
}

func TestDeleteProxy(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	host := "aaaaaa"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Description: description, Host: host, Tags: tags1})
	require.NoError(t, err)

	tags2 := []string{"amogus1", "amogus2", "kek1"}
	createRes2, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Description: description, Host: host, Tags: tags2})
	require.NoError(t, err)

	_, err = s.DeleteProxy(context.Background(), &proxy_pb.DeleteProxyRequest{Id: createRes1.Id})
	require.NoError(t, err)

	_, err = s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes1.Id})
	require.ErrorContains(t, err, "NotFound")

	getRes2, err := s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes2.Id})
	require.NoError(t, err)
	require.ElementsMatch(t, tags2, getProxyTagNames(getRes2.Tags))
}

func TestUpdateProxy(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	host := "aaaaaa"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Description: description, Schema: proxy_pb.Schema_SCHEMA_HTTP, Host: host, Tags: tags1})
	require.NoError(t, err)

	newName := "kek2"
	username1 := "qwerty"
	_, err = s.UpdateProxy(context.Background(), &proxy_pb.UpdateProxyRequest{Id: createRes1.Id, Name: newName, Schema: proxy_pb.Schema_SCHEMA_HTTPS, Username: username1, FieldMask: []string{"name", "schema", "username"}})
	require.NoError(t, err)

	getRes1, err := s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes1.Id})
	require.NoError(t, err)
	require.Equal(t, &proxy_pb.GetProxyResponse{Id: createRes1.Id, Name: newName, Description: description, Schema: proxy_pb.Schema_SCHEMA_HTTPS, Host: host, Tags: getRes1.Tags, Username: &username1}, getRes1)
	require.ElementsMatch(t, tags1, getProxyTagNames(getRes1.Tags))
}

func TestUpdateProxyWithTags(t *testing.T) {
	s, close := utils.SetupServer(t)
	defer close()

	name := "kek"
	description := "amogus"
	host := "aaaaaa"
	tags1 := []string{"kek1", "kek2", "kek3"}
	createRes1, err := s.CreateProxy(context.Background(), &proxy_pb.CreateProxyRequest{Name: name, Description: description, Host: host, Tags: tags1})
	require.NoError(t, err)

	newTags1 := []string{"amogus1", "amogus2", "kek1"}
	_, err = s.UpdateProxy(context.Background(), &proxy_pb.UpdateProxyRequest{Id: createRes1.Id, Tags: newTags1, FieldMask: []string{"tags"}})
	require.NoError(t, err)

	getRes1, err := s.GetProxy(context.Background(), &proxy_pb.GetProxyRequest{Id: createRes1.Id})
	require.NoError(t, err)
	require.Equal(t, &proxy_pb.GetProxyResponse{Id: createRes1.Id, Name: name, Description: description, Host: host, Tags: getRes1.Tags}, getRes1)
	require.ElementsMatch(t, newTags1, getProxyTagNames(getRes1.Tags))
}
