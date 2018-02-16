package functest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gomeet-examples/svc-profile/client"
	gomeetContext "github.com/gomeet/gomeet/utils/context"
)

func grpcClient(config FunctionalTestConfig) (*client.GomeetClient, context.Context, error) {
	serverAddr := config.ServerAddress
	if config.GrpcServerAddress != "" {
		serverAddr = config.GrpcServerAddress
	}

	c, err := client.NewGomeetClient(
		serverAddr,
		config.TimeoutSeconds,
		config.CaCertificate,
		config.ClientCertificate,
		config.ClientPrivateKey,
	)

	if err != nil {
		return nil, nil, err
	}

	// prepare the context
	ctx := gomeetContext.AuthContextFromJWT(context.Background(), config.JsonWebToken)

	return c, ctx, nil
}

func httpClient(config FunctionalTestConfig) (*http.Client, string, string, error) {
	proto := "http"
	client := &http.Client{Timeout: time.Duration(config.TimeoutSeconds) * time.Second}

	// TLS only if gRPC and HTTP are served on different addresses
	if config.CaCertificate != "" && config.ClientCertificate != "" && config.ClientPrivateKey != "" && config.GrpcServerAddress != "" && config.HttpServerAddress != "" {
		proto = "https"

		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(config.CaCertificate)
		if err != nil {
			return client, "", "", fmt.Errorf("failed to read CA certificate: %v", err)
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return client, "", "", fmt.Errorf("failed to build certificate pool")
		}

		serverHost, _, err := net.SplitHostPort(config.HttpServerAddress)
		if err != nil {
			return client, "", "", fmt.Errorf("failed to parse server hostname in %s: %v", config.HttpServerAddress, err)
		}

		tlsConfig := &tls.Config{
			ServerName: serverHost,
			RootCAs:    certPool,
		}

		transport := &http.Transport{TLSClientConfig: tlsConfig}

		client = &http.Client{Transport: transport, Timeout: time.Duration(config.TimeoutSeconds) * time.Second}
	}

	serverAddr := config.ServerAddress
	if config.HttpServerAddress != "" {
		serverAddr = config.HttpServerAddress
	}

	return client, serverAddr, proto, nil
}
