package database

import (
	"context"
	"fmt"
	"log"
	"mattermost-botpoll/config"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

type DB struct {
	Conn *tarantool.Connection
}

func InitConnectionDB(config *config.Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  fmt.Sprintf("%s:%s", config.TarantoolConf.TarantoolServerIp, config.TarantoolConf.TarantoolServerPort),
		User:     config.TarantoolConf.TarantoolUser,
		Password: config.TarantoolConf.TarantoolPassword,
	}

	opts := tarantool.Opts{
		Timeout: 5 * time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, err
	}

	_, err = conn.Do(tarantool.NewPingRequest()).Get()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("Успешно подключились к Tarantool")
	return &DB{Conn: conn}, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}

func (db *DB) InitSpaces() error {
	_, err := db.Conn.Do(
		tarantool.NewEvalRequest(`
            local space = box.schema.space.create('polls', {if_not_exists = true})
            space:format({
                {name = 'id', type = 'unsigned'},
                {name = 'title', type = 'string'},
                {name = 'description', type = 'string'},
                {name = 'variants', type = 'map'},
                {name = 'date_end', type = 'datetime'}
            })
            space:create_index('primary', {type = 'hash', parts = {'id'}, if_not_exists = true})
        `),
	).Get()
	if err != nil {
		return err
	}

	return nil
}
