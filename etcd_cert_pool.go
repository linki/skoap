package skoap

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/zalando/skipper/certpool"
)

const (
	poolPath = "/certpool"
	// interval in minutes
	interval = 1
)

type etcdCertPool struct {
	certPool  *certPool
	kapi      client.KeysAPI
	poolRoot  string
	etcdIndex uint64
}

func NewEtcdCertPool(endpoints []string, prefix string) (certpool.CertPool, error) {
	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	return &etcdCertPool{
		kapi:     client.NewKeysAPI(c),
		poolRoot: prefix + poolPath,
	}, nil
}

func (e *etcdCertPool) Set(pool **x509.CertPool) {
	e.certPool = newCertPool(pool)

	go e.run()
}

// run gets certpool from etcd and sets up a etcd watcher to watch for certpool
// updates.
func (e *etcdCertPool) run() {
	for {
		getOpts := &client.GetOptions{
			Recursive: true,
		}

		resp, err := e.kapi.Get(context.Background(), e.poolRoot, getOpts)
		if err != nil {
			logrus.Errorf("error getting certpool from etcd: %s", err)
			continue
		}

		for _, node := range resp.Node.Nodes {
			err := e.certPool.AppendCertsFromPEM(node.Key, []byte(node.Value))
			if err != nil {
				logrus.Errorf("failed to add Cert: %s", err)
				continue
			}
		}

		opts := &client.WatcherOptions{
			Recursive: true,
		}
		watcher := e.kapi.Watcher(e.poolRoot, opts)

		for {
			// TODO: timeout
			resp, err := watcher.Next(context.Background())
			if err != nil {
				logrus.Errorf("error watching: %s", err)
				continue
			}

			err = e.certPool.AppendCertsFromPEM(resp.Node.Key, []byte(resp.Node.Value))
			if err != nil {
				logrus.Errorf("failed to add certificate: %s", err)
				continue
			}
		}

		// sleep for x number of minutes
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

type certPool struct {
	state   map[string]*x509.Certificate
	poolPtr **x509.CertPool
}

func newCertPool(pool **x509.CertPool) *certPool {
	return &certPool{
		state:   make(map[string]*x509.Certificate),
		poolPtr: pool,
	}
}

// AppendCertsFromPEM appends a certificate defined as pem data and identified
// by a key to the cert pool.
func (p *certPool) AppendCertsFromPEM(key string, pemData []byte) error {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("failed to load pem data")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	p.state[key] = cert
	if *(p.poolPtr) == nil {
		*(p.poolPtr) = x509.NewCertPool()
	}
	pool := *(p.poolPtr)
	pool.AddCert(cert)

	return nil
}

// RemoveCert removes a certificate from the pool.
func (p *certPool) RemoveCert(key string) {
	if _, ok := p.state[key]; ok {
		delete(p.state, key)
		p.rebuildPool()
	}
}

// rebuildPool rebuilds the x509.CertPool based on the current state.
func (p *certPool) rebuildPool() {
	pool := x509.NewCertPool()
	for _, c := range p.state {
		pool.AddCert(c)
	}
	*(p.poolPtr) = pool
}
