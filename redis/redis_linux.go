package redis

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/penggy/EasyGoLib/utils"
)

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
		cmd.Dir = filepath.Dir(EXE())
		err = cmd.Start()
		if err != nil {
			return
		}
		pidPath := filepath.Join(filepath.Dir(EXE()), "redis.pid")
		ioutil.WriteFile(pidPath, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	}

	log.Printf("redis server --> redis://%s:%d/db%d", host, port, db)
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: auth,
		DB:       db,
	})
	return
}
