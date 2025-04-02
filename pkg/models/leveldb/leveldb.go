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

func (l *LevelDB) GetFile(fileID string) (file dto.File, err error) {
	data, err := l.DB.Get([]byte(fileID), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return dto.File{}, errors.New("record not found")
		}
		return dto.File{}, err
	}
	err = json.Unmarshal(data, &file)
	if err != nil {
		return dto.File{}, err
	}
	file.Content, err = l.GetFileContent(fileID)
	return
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

func (l *LevelDB) SetFile(fileID string, f dto.File) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return l.DB.Put([]byte(fileID), data, nil)
}

func (l *LevelDB) DeleteFile(fileID string) error {
	return l.DB.Delete([]byte(fileID), nil)
}

func (l *LevelDB) SetFileMeta(key string, meta dto.FileMeta) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return l.DB.Put([]byte(key), data, nil)
}

func (l *LevelDB) SetFileContent(key string, content []byte) error {
	return l.DB.Put([]byte(key), content, nil)
}

func (l *LevelDB) GetFileContent(fileID string) (fileContent []byte, err error) {
	fileContent, err = l.DB.Get([]byte("content_"+fileID), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return []byte{}, errors.New("record not found")
		}
		return []byte{}, err
	}
	return
}

func (l *LevelDB) DeleteFileContent(fileID string) error {
	return l.DB.Delete([]byte("content_"+fileID), nil)
}
