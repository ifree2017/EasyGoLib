package redis

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/structs"

	"github.com/mitchellh/mapstructure"

	"github.com/go-redis/redis"
	"github.com/penggy/EasyGoLib/utils"
)

var Client *redis.Client
var cmd *exec.Cmd

func EXE() string {
	bin := utils.Conf().Section("redis").Key("bin").MustString("redis/redis-server")
	if !filepath.IsAbs(bin) {
		bin = filepath.Join(utils.CWD(), bin)
	}
	bin = strings.TrimSuffix(bin, ".exe")
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s.exe", bin)
	case "linux":
		return bin
	}
	return ""
}

func TestConnect() (ret bool) {
	if Client == nil {
		return
	}
	if _, err := Client.Ping().Result(); err == nil {
		ret = true
	}
	return
}

func HGetStruct(key string, out interface{}) (err error) {
	if Client == nil {
		err = fmt.Errorf("redis client not prepared")
		return
	}
	n, err := Client.Exists(key).Result()
	if err != nil {
		return
	}
	if n == 0 {
		err = fmt.Errorf("key[%s] not found", key)
		return
	}
	retMap, err := Client.HGetAll(key).Result()
	if err != nil {
		return
	}
	err = mapstructure.WeakDecode(retMap, out)
	return
}

func HSetStruct(key string, in interface{}, d time.Duration) (err error) {
	if Client == nil {
		err = fmt.Errorf("redis client not prepared")
		return
	}
	tx := Client.TxPipeline()
	for k, v := range structs.Map(in) {
		err = tx.HSet(key, k, v).Err()
		if err != nil {
			return
		}
	}
	if d > 0 {
		err = tx.Expire(key, d).Err()
		if err != nil {
			return
		}
	}
	_, err = tx.Exec()
	return
}

func HMSetStruct(key string, in interface{}, d time.Duration) (err error) {
	if Client == nil {
		err = fmt.Errorf("redis client not prepared")
		return
	}
	tx := Client.TxPipeline()
	err = tx.HMSet(key, structs.Map(in)).Err()
	if err != nil {
		return
	}
	if d > 0 {
		err = tx.Expire(key, d).Err()
		if err != nil {
			return
		}
	}
	_, err = tx.Exec()
	return
}

func Close() (err error) {
	log.Printf("redis close")
	if Client != nil {
		err = Client.Close()
		if err != nil {
			return
		}
		Client = nil
	}
	if cmd != nil {
		cmd.Process.Kill()
		os.RemoveAll(filepath.Join(filepath.Dir(EXE()), "redis.pid"))
		cmd = nil
	}
	return
}
