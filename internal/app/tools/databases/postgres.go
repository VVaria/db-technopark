package databases

import (
	"github.com/jackc/pgx"
)

type Postgres struct {
	postgresDatabase *pgx.ConnPool
}

func NewPostgres(dataSourceName string) (*Postgres, error) {
	pgxConn, err := pgx.ParseConnectionString(dataSourceName)
	if err != nil {
		return nil, err
	}

	pgxConn.PreferSimpleProtocol = true

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgxConn,
		MaxConnections: 200,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		postgresDatabase: pool,
	}, nil
}

func (p *Postgres) GetDatabase() *pgx.ConnPool{
	return p.postgresDatabase
}

func (p *Postgres) Close() {
	p.postgresDatabase.Close()
}
