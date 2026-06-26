package grpcclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnConfig struct {
	Addr               string
	TLSEnabled         bool
	CACertPath         string
	InsecureSkipVerify bool
}

func NewConnection(cfg ConnConfig) (*grpc.ClientConn, error) {
	creds, err := transportCredentials(cfg)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("create grpc client for %s: %w", cfg.Addr, err)
	}
	return conn, nil
}

func transportCredentials(cfg ConnConfig) (credentials.TransportCredentials, error) {
	if !cfg.TLSEnabled {
		return insecure.NewCredentials(), nil
	}

	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}
	if cfg.InsecureSkipVerify {
		tlsCfg.InsecureSkipVerify = true
		return credentials.NewTLS(tlsCfg), nil
	}
	if cfg.CACertPath != "" {
		ca, err := os.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("read ca cert %s: %w", cfg.CACertPath, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("append ca cert %s", cfg.CACertPath)
		}
		tlsCfg.RootCAs = pool
	}

	return credentials.NewTLS(tlsCfg), nil
}
