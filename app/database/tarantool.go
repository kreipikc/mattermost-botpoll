package database

import (
	"context"
	"fmt"
	"log"
	"mattermost-botpoll/config"
	"mattermost-botpoll/models"
	"mattermost-botpoll/utils"
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

func (db *DB) CreatePoll(poll *models.PollBody) (int, error) {
	if (poll.DateEnd).Before(time.Now()) {
		return 0, fmt.Errorf("время окончания голосования не может быть раньше завтрашнего дня")
	}

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

	// log.Printf("Ответ от Tarantool при создании poll: %v", resp)

	if len(resp) == 0 {
		return 0, fmt.Errorf("не удалось создать опрос: пустой ответ от Tarantool")
	}

	tuple := resp[0].([]interface{})
	if len(tuple) != 6 {
		return 0, fmt.Errorf("неверный формат кортежа: ожидается 6 полей, получено %d", len(tuple))
	}

	generatedID, err := utils.ConvertAllIntToInt(tuple[0])
	if err != nil {
		return 0, fmt.Errorf("ошибка при конвертации к int: %v", err)
	}

	poll.Id = generatedID

	return generatedID, nil
}

func (db *DB) GetPollByID(id int) (*models.PollBody, error) {
	req := tarantool.NewSelectRequest("polls").
		Index("primary").
		Limit(1).
		Iterator(tarantool.IterEq).
		Key([]interface{}{id})

	resp, err := db.Conn.Do(req).Get()
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении SELECT: %w", err)
	}

	// log.Printf("Ответ от Tarantool для id %d: %v", id, resp)

	if len(resp) == 0 {
		return nil, fmt.Errorf("опрос с ID %d не найден", id)
	}

	tuple, ok := resp[0].([]interface{})
	if !ok || len(tuple) != 6 {
		return nil, fmt.Errorf("неверный формат кортежа: ожидается 6 полей, получено %v", resp[0])
	}

	pollID, err := utils.ConvertAllIntToInt(tuple[0])
	if err != nil {
		return nil, fmt.Errorf("ошибка при конвертации к int: %v", err)
	}

	authorID, ok := tuple[1].(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип для author_id: %T", tuple[1])
	}

	title, ok := tuple[2].(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип для title: %T", tuple[2])
	}

	description, ok := tuple[3].(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип для description: %T", tuple[3])
	}

	variantsRaw, ok := tuple[4].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("неверный тип для variants: %T", tuple[4])
	}

	variants := make(map[string]int)
	for key, value := range variantsRaw {
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("неверный тип ключа в variants: %T", key)
		}
		valueInt, err := utils.ConvertAllIntToInt(value)
		if err != nil {
			return nil, fmt.Errorf("ошибка при конвертации к int: %v", err)
		}
		variants[keyStr] = int(valueInt)
	}

	dateEndStr, ok := tuple[5].(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип для date_end: %T", tuple[5])
	}
	dateEnd, err := time.Parse(time.RFC3339, dateEndStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга date_end: %w", err)
	}

	poll := &models.PollBody{
		Id:          pollID,
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Variants:    variants,
		DateEnd:     dateEnd,
	}

	return poll, nil
}

func (db *DB) UpdatePollVote(idPoll int, variant string) error {
	poll, err := db.GetPollByID(idPoll)
	if err != nil {
		return fmt.Errorf("ошибка получения опроса: %w", err)
	}

	if (poll.DateEnd).Before(time.Now()) {
		return fmt.Errorf("голосование завершено")
	}

	if _, exists := poll.Variants[variant]; !exists {
		return fmt.Errorf("вариант %s не найден в опросе с ID %d", variant, idPoll)
	}

	poll.Variants[variant]++

	req := tarantool.NewUpdateRequest("polls").
		Index("primary").
		Key([]interface{}{idPoll}).
		Operations(tarantool.NewOperations().Assign(4, poll.Variants))

	resp, err := db.Conn.Do(req).Get()
	if err != nil {
		return fmt.Errorf("ошибка обновления опроса: %w", err)
	}

	if len(resp) == 0 {
		return fmt.Errorf("не удалось обновить опрос: пустой ответ от Tarantool")
	}

	log.Printf("Успешно обновлён опрос с ID %d: вариант %s теперь имеет %d голосов", idPoll, variant, poll.Variants[variant])
	return nil
}

func (db *DB) UpdatePollEnd(idPoll int, idAuthor string) error {
	poll, err := db.GetPollByID(idPoll)
	if err != nil {
		return fmt.Errorf("ошибка получения опроса: %w", err)
	}

	if idAuthor != poll.AuthorID {
		return fmt.Errorf("завершать голосование досрочно может только автор")
	}

	req := tarantool.NewUpdateRequest("polls").
		Index("primary").
		Key([]interface{}{idPoll}).
		Operations(tarantool.NewOperations().Assign(5, time.Now().Format(time.RFC3339)))

	resp, err := db.Conn.Do(req).Get()
	if err != nil {
		return fmt.Errorf("ошибка обновления опроса: %w", err)
	}

	if len(resp) == 0 {
		return fmt.Errorf("не удалось обновить опрос: пустой ответ от Tarantool")
	}

	return nil
}

func (db *DB) DeletePoll(idPoll int, idAuthor string) error {
	poll, err := db.GetPollByID(idPoll)
	if err != nil {
		return fmt.Errorf("ошибка получения опроса: %w", err)
	}

	if idAuthor != poll.AuthorID {
		return fmt.Errorf("удалить голосование может только автор")
	}

	req := tarantool.NewDeleteRequest("polls").
		Index("primary").
		Key([]interface{}{idPoll})

	resp, err := db.Conn.Do(req).Get()
	if err != nil {
		return fmt.Errorf("ошибка при выполнении DELETE: %w", err)
	}

	// log.Printf("Ответ от Tarantool при удалении ID %d: %v", idPoll, resp)

	if len(resp) == 0 {
		return fmt.Errorf("опрос с ID %d не найден", idPoll)
	}

	return nil
}
