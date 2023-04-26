package mariadb

import (
	"fmt"
	"time"

	"github.com/plzapsys/go-pkg/utils"

	// "database/sql"
	"context"
	// "microservices/pkg/tracing"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"

	// "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbName"`
}

type Sqlx struct {
	SqlxDB *sqlx.DB
	// DB  *sql.DB
	config *Config
}

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// NewSqlxConn func for connection to Mariadb database.
func NewSqlxConn(cfg *Config) (*Sqlx, error) {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	sqlx := &Sqlx{SqlxDB: db, config: cfg}

	return sqlx, nil
}

func (db *Sqlx) Close() {
	// _ = db.DB.Close()
	_ = db.SqlxDB.Close()
}

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, db *sqlx.DB, countQuery string, itemsQuery string) (*utils.ListResult[T], error) {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "mariadbdb.Paginate")

	var count int64
	err := db.GetContext(ctx, &count, countQuery)
	if err != nil {
		// tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "Paginate.CountItems")
	}

	itemsList := make([]T, 0, listQuery.GetSize())

	rows, err := db.QueryxContext(ctx, itemsQuery, listQuery.GetOffset(), listQuery.GetLimit())
	if err != nil {
		// tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "Paginate.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {

		var itm T
		// var itm []T
		// itm := &T{}
		if err = rows.StructScan(itm); err != nil {
			// fmt.Println(err)
			// tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "Paginate.StructScan")
		}
		fmt.Println(itm)

		itemsList = append(itemsList, itm)

	}

	if err = rows.Err(); err != nil {
		// tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "Paginate.rows.Err")
	}

	return utils.NewListResult[T](itemsList, listQuery.GetSize(), listQuery.GetPage(), count), nil
}
