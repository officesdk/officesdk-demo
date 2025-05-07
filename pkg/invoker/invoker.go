package invoker

import (
	"turbo-demo/pkg/models/leveldb"
	"turbo-demo/pkg/services"
	"turbo-demo/ui"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
)

var (
	Gin         *egin.Component
	Leveldb     *leveldb.LevelDB
	FileService *services.FileService
	AwsService  *services.AwsService
)

func Init() (err error) {
	Gin = egin.Load("server.demo.http").Build(egin.WithEmbedFs(ui.WebUI))
	Leveldb, err = leveldb.NewLevelDB()
	if err != nil {
		elog.Panic("Failed to initialize Leveldb")
	}
	AwsService, err = services.NewAwsService()
	if err != nil {
		elog.Panic("Failed to initialize AwsService")
	}

	FileService = services.NewFileService(Leveldb, AwsService)

	FileService.InitCaseFile()
	return nil
}
