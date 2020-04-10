package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/go-log/log"
	"github.com/micro/go-micro/v2/config/source/file"
	"honnef.co/go/tools/config"
)

var (
	err error
)

var (
	defaultRootPath         = "app"
	defaultConfigFilePrefix = "application-"
	etcdConfig              defaultEtcdConfig
	mysqlConig              defaultMysqlConfig
	profiles                defaultProfiles
	m                       sync.RWMutex
	inited                  bool
	sp                      = string(filepath.Separator)
)

func Init2qqq() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Logf("init 过了")
		return
	}

	appPath, _ := filepath.Abs(file.Dir(filepath.Join("."+sp, sp)))

	pt := filepath.Join(appPath, "conf")
	os.Chdir(appPath)

	if err = config.Load(file.NewSource(file.WithPath(pt + sp + "application.yml"))); err != nil {
		panic(err)
	}

	if err = config.Get(defaultRootPath, "profiles").Scan(&profiles); err != nil {
		panic(err)
	}

	log.Logf("加载文件:path:%s,%+v\n", pt+sp+"application.yml", profiles)
}
