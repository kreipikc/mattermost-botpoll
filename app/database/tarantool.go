package database

import (
	"context"
	"fmt"
	"log"
	"mattermost-botpoll/config"
	"mattermost-botpoll/models"
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
			box.schema.sequence.create('poll_id_seq', {start = 1, min = 1, if_not_exists = true})

            local space = box.schema.space.create('polls', {if_not_exists = true})
            space:format({
                {name = 'id', type = 'unsigned'},
				{name = 'author_id', type = 'string'},
                {name = 'title', type = 'string'},
                {name = 'description', type = 'string'},
                {name = 'variants', type = 'map'},
                {name = 'date_end', type = 'string'}
            })
            space:create_index('primary', {type = 'tree', parts = {'id'}, sequence = 'poll_id_seq', if_not_exists = true})
        `),
	).Get()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreatePoll(poll *models.PollBody) (uint32, error) {
	req := tarantool.NewCallRequest("box.space.polls:auto_increment").
		Args([]interface{}{
			[]interface{}{
				poll.AuthorID,
				poll.Title,
				poll.Description,
				poll.Variants,
				poll.DateEnd.Format(time.RFC3339),
			},
		})

	resp, err := db.Conn.Do(req).Get()
	if err != nil {
		return 0, err
	}

	if len(resp) == 0 {
		return 0, fmt.Errorf("не удалось создать опрос: пустой ответ от Tarantool")
	}

	tuple := resp[0].([]interface{})
	if len(tuple) != 6 {
		return 0, fmt.Errorf("неверный формат кортежа: ожидается 6 полей, получено %d", len(tuple))
	}

	var generatedID uint32
	switch v := tuple[0].(type) {
	case int8:
		generatedID = uint32(v)
	case int16:
		generatedID = uint32(v)
	case int32:
		generatedID = uint32(v)
	case int64:
		generatedID = uint32(v)
	case uint8:
		generatedID = uint32(v)
	case uint16:
		generatedID = uint32(v)
	case uint32:
		generatedID = v
	case uint64:
		generatedID = uint32(v)
	case float64:
		generatedID = uint32(v)
	default:
		return 0, fmt.Errorf("неожиданный тип для Id: %T", v)
	}

	poll.Id = generatedID

	return generatedID, nil
}

func (db *DB) GetPollByID(id uint32) (*models.PollBody, error) {
	return &models.PollBody{}, nil
}

func (db *DB) UpdatePollVote(idPoll uint32, variant string) error {
	return nil
}
