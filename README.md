# proxy

## Service

По дефолту берет данные о проксях из бд, в standalone mode берет данные из заданного конфига
- CRUD for `Proxy`
- CRUD for `ProxyBlock`

### Proxy
- uuid uuid
- name string
- description string
- tags []string
- ip string
- port uint16
- username string
- password string

### ProxyBlock
- uuid uuid
- name string
- description string
- filters []Tag


## ClientLibrary
- MakeClient(c *http.Client, proxyBlockId uuid) http.Client
- client.MakeRequest(...)

`MakeClient` оборачивает клиент так, чтобы все запросы (`client.MakeRequest`) ходили через рандомные прокси из заданного `ProxyBlock`

## ClientBinary

Поднимает прокси без авторизации, которую можно указать в хроме, и которая проксирует все запросы на указанные прокси
