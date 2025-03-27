package invoker

import (
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"office-demo/pkg/models/leveldb"
	"office-demo/pkg/services"
	"office-demo/ui"
)

var (
	Gin         *egin.Component
	Leveldb     *leveldb.LevelDB
	FileService *services.FileService
)

func Init() (err error) {
	Gin = egin.Load("server.demo.http").Build(egin.WithEmbedFs(ui.WebUI))
	Leveldb, err = leveldb.NewLevelDB()
	if err != nil {
		elog.Panic("Failed to initialize Leveldb")
	}
	FileService = services.NewFileService(Leveldb)
	FileService.InitCaseFile()
	return nil
}
