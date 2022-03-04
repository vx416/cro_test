package container

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
)

type DB struct {
	Host     string
	Port     int32
	Username string
	Password string
	DBName   string
	DBType   string
}

func (db *DB) GetDB() (*sql.DB, error) {
	var (
		dsn string
	)

	switch db.DBType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			db.Username, db.Password, db.Host, db.Port, db.DBName)
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			db.Username, db.Password, db.Host, db.Port, db.DBName)
	}
	client, err := sql.Open(db.DBType, dsn)
	if err != nil {
		return nil, err
	}
	err = client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (builder *Builder) RunDB(options *dockertest.RunOptions, driver, dbName string) (*DB, error) {
	container, err := builder.FindContainer(options.Name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		builder.containerIDs[container.ID] = true
		return &DB{
			Port:     int32(container.Ports[0].PublicPort),
			DBName:   dbName,
			Host:     "localhost",
			Username: "test",
			Password: "test",
			DBType:   driver,
		}, nil
	}

	resource, err := builder.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	builder.containerIDs[resource.Container.ID] = true
	dbPort := int64(0)

	switch driver {
	case "pg", "postgres":
		dbPort, err = strconv.ParseInt(resource.GetPort("5432/tcp"), 10, 64)
	case "mysql":
		dbPort, err = strconv.ParseInt(resource.GetPort("3306/tcp"), 10, 64)
	}

	if err != nil {
		return nil, err
	}

	db := &DB{
		Port:     int32(dbPort),
		DBName:   dbName,
		Host:     "localhost",
		Username: "test",
		Password: "test",
		DBType:   driver,
	}

	return db, nil
}

func (builder *Builder) RunPg(name string, dbName string, port ...string) (*DB, error) {
	options := &dockertest.RunOptions{
		Repository: "postgres", Tag: "12.3-alpine", Name: name,
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=" + dbName,
		},
	}
	if len(port) == 1 {
		options.PortBindings = make(map[dc.Port][]dc.PortBinding)
		options.PortBindings[dc.Port("5432/tcp")] = []dc.PortBinding{
			{
				HostPort: port[0],
			}}
	}

	pg, err := builder.RunDB(options, "pg", dbName)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pg.Username, pg.Password, pg.Host, pg.Port, pg.DBName)

	err = builder.checkConnection("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (builder *Builder) RunMysql(name string, dbName string, port ...string) (*DB, error) {
	options := &dockertest.RunOptions{
		Repository: "mysql/mysql-server", Tag: "5.7", Name: name,
		Env: []string{
			"MYSQL_USER=test",
			"MYSQL_PASSWORD=test",
			"MYSQL_ROOT_PASSWORD=test",
			"MYSQL_DATABASE=" + dbName,
		},
	}
	if len(port) == 1 {
		options.PortBindings = make(map[dc.Port][]dc.PortBinding)
		options.PortBindings[dc.Port("3306/tcp")] = []dc.PortBinding{
			{
				HostPort: port[0],
			}}
	}

	db, err := builder.RunDB(options, "mysql", dbName)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		db.Username, db.Password, db.Host, db.Port, db.DBName)

	err = builder.checkConnection("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (builder *Builder) checkConnection(driver, dsn string) error {
	builder.MaxWait = 8 * time.Second
	err := builder.Retry(func() error {
		db, err := sql.Open(driver, dsn)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	})
	return err
}

func RunMigration(db *sql.DB, driver, filesPath string) error {
	goose.SetDialect(driver)
	return goose.Up(db, filesPath)
}

func RunTestData(ctx context.Context, db *sql.DB, testDataDir string) error {
	return filepath.Walk(testDataDir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(file) == ".sql" {
				data, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}

				if string(data) != "" {
					dataStr := string(data)
					stmts := strings.Split(dataStr, ";")
					for _, stmt := range stmts {
						stmt = strings.TrimSpace(stmt)
						if stmt == "" {
							continue
						}
						_, err = db.ExecContext(ctx, stmt)
						if err != nil {
							fmt.Println(stmt)
							return err
						}
					}
				}
			}
		}
		return nil
	})
}

func DeleteTables(ctx context.Context, db *sql.DB, tableName ...string) error {
	for _, table := range tableName {
		_, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return err
		}
	}

	return nil
}
