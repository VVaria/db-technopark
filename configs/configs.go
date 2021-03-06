package configs

import "fmt"

type Main struct {
	Host string
	Port int
}

type Postgres struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
}

type config struct {
	Main     Main
	Postgres Postgres
}

var Configs config

func (c *config) GetMainHost() string {
	return c.Main.Host
}

func (c *config) GetMainPort() int {
	return c.Main.Port
}

func (c *config) GetPostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName)
}

func init() {
	Configs = config{
		Main: Main{
			"localhost",
			5000,
		},
		Postgres: Postgres{
			"docker",
			"docker",
			"docker",
			"localhost",
			5432,
		},
	}
}
