package postgresdocker

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type options struct {
	tag      string
	username string
	password string
	dbDriver string
}

var defaultInstanceOptions = options{
	tag:      "14-alpine",
	username: "postgres",
	password: "password",
	dbDriver: "pgx",
}

type InstanceOption func(*options)

func Tag(s string) InstanceOption {
	return func(o *options) { o.tag = s }
}

func Username(s string) InstanceOption {
	return func(o *options) { o.username = s }
}

func Password(s string) InstanceOption {
	return func(o *options) { o.password = s }
}

func DbDriver(s string) InstanceOption {
	return func(o *options) { o.dbDriver = s }
}

type PostgresInstance struct {
	opts     options
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func New(opts ...InstanceOption) (*PostgresInstance, error) {
	options := defaultInstanceOptions
	for _, o := range opts {
		o(&options)
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        options.tag,
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", options.username),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", options.password),
		},
	}, func(hc *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, err
	}
	pgi := &PostgresInstance{
		opts:     options,
		resource: resource,
		pool:     pool,
	}
	return pgi, pgi.waitInit(10 * time.Second)
}

func (p *PostgresInstance) waitInit(timeout time.Duration) error {
	return p.pool.Retry(func() error {
		db, err := sql.Open(p.opts.dbDriver, p.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return db.PingContext(ctx)
	})
}

func (p *PostgresInstance) ConnectionString() string {
	hostPort := strings.Split(p.resource.GetHostPort("5432/tcp"), ":")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		hostPort[0], hostPort[1], p.opts.username, p.opts.password,
	)
}

func (p *PostgresInstance) Shutdown() {
	_ = p.pool.Purge(p.resource)
}
