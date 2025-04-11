package leveldb

import (
	"encoding/json"
	"errors"
	"fmt"

	"turbo-demo/pkg/models/dto"

	"github.com/gotomicro/ego/core/econf"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	DB *leveldb.DB
}

func NewLevelDB() (*LevelDB, error) {
	db, err := leveldb.OpenFile(econf.GetString("leveldb.path"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %v", err)
	}
	return &LevelDB{
		DB: db,
	}, nil
}

func (l *LevelDB) GetFileMeta(fileID string) (file dto.FileMeta, err error) {
	data, err := l.DB.Get([]byte(fileID), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return dto.FileMeta{}, errors.New("record not found")
		}
		return dto.FileMeta{}, err
	}
	err = json.Unmarshal(data, &file)
	if err != nil {
		return dto.FileMeta{}, err
	}
	return
}

func (l *LevelDB) DeleteFileMeta(fileID string) error {
	return l.DB.Delete([]byte(fileID), nil)
}

func (l *LevelDB) SetFileMeta(key string, meta dto.FileMeta) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return l.DB.Put([]byte(key), data, nil)
}
