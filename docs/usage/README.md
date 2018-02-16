# svc-profile usage

## Basic usage

- Run server

```shell
svc-profile serve --address <server-address>

# serve gRPC and HTTP multiplexed on localhost:3000
svc-profile serve --address localhost:3000

# serve gRPC on localhost:3000 and HTTP on localhost:3001
svc-profile serve --grpc-address localhost:3000 --http-address localhost:3001

# more info
svc-profile help serve
```

- Run cli client

  ```shell
  $ svc-profile cli version
  $ svc-profile cli services_status
  $ svc-profile cli create <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>
  $ svc-profile cli read <uuid [string]>
  $ svc-profile cli list <page_number [uint32]> <page_size [uint32]> <gender [UNKNOW|MALE|FEMALE]>
  $ svc-profile cli update <uuid [string]> <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>
  $ svc-profile cli soft_delete <uuid [string]>
  $ svc-profile cli hard_delete <uuid [string]>
  $ svc-profile cli --address localhost:42000 version

  # more info
  svc-profile help cli
  ```

- Run console client

```shell
$ svc-profile console --address=localhost:3000
INFO[0000] svc-profile console  Exit=exit HistoryFile="/tmp/svc-profile-62852.tmp" Interrupt="^C"
└─┤svc-profile-0.1.8+dev@localhost:13000├─$ help
INFO[0002] HELP :

	┌─ version
	└─ call version service

	┌─ services_status
	└─ call services_status service

	┌─ create <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>
	└─ call create service

	┌─ read <uuid [string]>
	└─ call read service

	┌─ list <page_number [uint32]> <page_size [uint32]> <gender [UNKNOW|MALE|FEMALE]>
	└─ call list service

	┌─ update <uuid [string]> <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>
	└─ call update service

	┌─ soft_delete <uuid [string]>
	└─ call soft_delete service

	┌─ hard_delete <uuid [string]>
	└─ call hard_delete service

	┌─ service_address
	└─ return service address

	┌─ jwt [<token>]
	└─ display current jwt or save none if it's set

	┌─ console_version
	└─ return console version

	┌─ tls_config
	└─ display TLS client configuration

	┌─ help
	└─ display this help

	┌─ exit
	└─ exit the console


└─┤svc-profile-0.1.8+dev@localhost:13000├─$ unknow
WARN[0003] Bad arguments : "unknow" unknow
└─┤svc-profile-0.1.8+dev@localhost:13000├─$ exit
```

- HTTP/1.1 usage (with curl):

  ```shell
  $ curl -X GET    http://localhost:13000/api/v1/version
  $ curl -X GET    http://localhost:13000/api/v1/services/status
  $ curl -X POST   http://localhost:13000/api/v1/create -d '{"gender": "UNKNOW|MALE|FEMALE", "email": "<string>", "name": "<string>", "birthday": "<string>"}'
  $ curl -X POST   http://localhost:13000/api/v1/read -d '{"uuid": "<string>"}'
  $ curl -X POST   http://localhost:13000/api/v1/list -d '{"page_number": <number>, "page_size": <number>, "gender": "UNKNOW|MALE|FEMALE"}'
  $ curl -X POST   http://localhost:13000/api/v1/update -d '{"uuid": "<string>", "gender": "UNKNOW|MALE|FEMALE", "email": "<string>", "name": "<string>", "birthday": "<string>"}'
  $ curl -X POST   http://localhost:13000/api/v1/soft_delete -d '{"uuid": "<string>"}'
  $ curl -X POST   http://localhost:13000/api/v1/hard_delete -d '{"uuid": "<string>"}'
  $ curl -X GET    http://localhost:13000/
  $ curl -X GET    http://localhost:13000/version
  $ curl -X GET    http://localhost:13000/metrics
  $ curl -X GET    http://localhost:13000/status
  $ curl -X GET    http://localhost:42000/version
  ```

- Get help

```shell
svc-profile help

# or get help directly for a command
svc-profile help <command[serve|cli|console]>
```

## Tests

- Use make directive

```shell
make test
```

- Unit tests

```shell
cd service
go test
```

- Functional tests (with an embedded server)

```shell
svc-profile functest -e
```

- Load tests

```shell
svc-profile loadtest --address <multiplexed server address> -n <number of sessions> -s <concurrency level>
```

## Mutual TLS authentication

- Create a Certificate Authority:

```shell
hack/gen-ca.sh gomeetexamples_ca
ls data/certs
```

- Create two key pairs with the common name "localhost":

```shell
hack/gen-cert.sh server gomeetexamples_ca
./gencert.sh client gomeetexamples_ca
ls data/certs
```

- Run the server with its TLS credentials:

```shell
svc-profile serve \
    --address localhost:3000 \
    --ca data/certs/gomeetexamples_ca.crt \
    --cert data/certs/server.crt \
    --key data/certs/server.key
```

- Run the clients with their TLS credentials:

```shell
svc-profile cli <grpc_service> <params...> \
    --address localhost:3000 \
    --ca data/certs/gomeetexamples_ca.crt \
    --cert data/certs/client.crt \
    --key data/certs/client.key

svc-profile console \
    --address localhost:3000 \
    --ca data/certs/gomeetexamples_ca.crt \
    --cert data/certs/client.crt \
    --key data/certs/client.key
```

## JSON Web Token support

JSON Web Token validation can be enabled on the server by providing a secret key:

```shell
svc-profile serve --jwt-secret foobar
```

The token subcommand is used to generate a JWT from the secret key:

```shell
svc-profile token --secret-key foobar
```

Then the cli and console subcommands can use the generated token for authentication against the JWT-enabled server:

```shell
svc-profile cli --jwt <generated token> <grpc_service> <params...>
svc-profile console --jwt <generated token>
```

JWT validation can be tested on the HTTP/1.1 endpoints by providing the bearer token in the "Authorization" HTTP header:

```shell
TOKEN=`svc-profile token --secret-key foobar`
curl -H "Authorization: Bearer $TOKEN" -X <HTTP_VERB> http://localhost:13000/api/v1/<grpc_service> -d '<HTTP_REQUEST_BODY json format>'
```


