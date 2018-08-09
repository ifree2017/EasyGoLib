package redis

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

func Init() (err error) {
	sec := utils.Conf().Section("redis")
	host := sec.Key("host").MustString("localhost")
	port := sec.Key("port").MustInt(6379)
	auth := sec.Key("auth").MustString("")
	db := sec.Key("db").MustInt(0)

	if host == "localhost" && utils.Exist(EXE()) {
		if buf, _ := ioutil.ReadFile(filepath.Join(filepath.Dir(EXE()), "redis.pid")); len(buf) > 0 {
			pid, _ := strconv.Atoi(string(buf))
			if p, err := os.FindProcess(pid); err == nil {
				p.Kill()
				for i := 0; i < 10; i++ {
					time.Sleep(1 * time.Second)
					if p, _ := os.FindProcess(pid); p == nil {
						break
					}
				}
			}
		}
		if utils.IsPortInUse(port) {
			err = fmt.Errorf("Port[%d] In Use", port)
			return
		}
		args := []string{"--port", strconv.Itoa(port)}
		if auth != "" {
			args = append(args, "--requirepass", auth)
		}
		cmd = exec.Command(EXE(), args...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		cmd.Dir = filepath.Dir(EXE())
		err = cmd.Start()
		if err != nil {
			return
		}
		pidPath := filepath.Join(filepath.Dir(EXE()), "redis.pid")
		ioutil.WriteFile(pidPath, []byte(strconv.Itoa(cmd.Process.Pid)), os.ModeAppend)
	}

	log.Printf("redis server --> redis://%s:%d/db%d", host, port, db)

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: auth,
		DB:       db,
	})
	if _, e := Client.Ping().Result(); e != nil {
		err = fmt.Errorf("redis connect failed, %v", e)
		return
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
