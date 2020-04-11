package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-log/log"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/file"
)

var (
	err error
)

var (
	defaultRootPath         = "app"
	defaultConfigFilePrefix = "application-"
	etcdConfig              defaultEtcdConfig
	profiles                defaultProfiles
	m                       sync.RWMutex
	inited                  bool
	sp                      = string(filepath.Separator)
)

func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Logf("init 过了")
		return
	}

	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))

	pt := filepath.Join(appPath, "conf")
	os.Chdir(appPath)

	if err = config.Load(file.NewSource(file.WithPath(pt + sp + "application.yml"))); err != nil {
		panic(err)
	}

	if err = config.Get(defaultRootPath, "profiles").Scan(&profiles); err != nil {
		panic(err)
	}

	log.Logf("加载文件:path:%s,%+v\n", pt+sp+"application.yml", profiles)

	if len(profiles.GetInclude()) > 0 {
		include := strings.Split(profiles.GetInclude(), ",")

		sources := make([]source.Source, len(include))

		for i := 0; i < len(include); i++ {
			filePath := pt + string(filepath.Separator) + defaultConfigFilePrefix + strings.TrimSpace(include[i]) + ".yml"
			log.Logf("init path : %s\n", filePath)
			sources[i] = file.NewSource(file.WithPath(filePath))
		}

		if err = config.Load(sources...); err != nil {
			panic(err)
		}
	}

	config.Get(defaultRootPath, "etcd").Scan(&etcdConfig)

	inited = true
}

func GetEtcdConfig() (ret EtcdConfig) {
	return etcdConfig
}
