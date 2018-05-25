package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/teris-io/shortid"

	"github.com/go-ini/ini"
)

func LocalIP() string {
	ip := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func MD5(str string) string {
	encoder := md5.New()
	encoder.Write([]byte(str))
	return hex.EncodeToString(encoder.Sum(nil))
}

func CWD() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}

func EXEName() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func HomeDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}

func DataDir() string {
	dir := CWD()
	_dir := Conf().Section("").Key("data_dir").Value()
	if _dir != "" {
		dir = _dir
	}
	dir = ExpandHomeDir(dir)
	EnsureDir(dir)
	// if _dir == "" {
	// 	conf := ReloadConf()
	// 	conf.Section("").Key("data_dir").SetValue(dir)
	// 	conf.SaveTo(ConfFile())
	// }
	return dir
}

func VODDir() string {
	dir := filepath.Join(DataDir(), "vod")
	dir = ExpandHomeDir(Conf().Section("").Key("vod_dir").MustString(dir))
	EnsureDir(dir)
	return dir
}

func LogDir() string {
	dir := filepath.Join(DataDir(), "logs")
	dir = ExpandHomeDir(Conf().Section("").Key("log_dir").MustString(dir))
	EnsureDir(dir)
	return dir
}

func ConfFile() string {
	return filepath.Join(CWD(), strings.ToLower(EXEName())+".ini")
}

func ConfFileDev() string {
	return filepath.Join(CWD(), strings.ToLower(EXEName())+".dev.ini")
}

var conf *ini.File

func Conf() *ini.File {
	if conf != nil {
		return conf
	}
	if _conf, err := ini.InsensitiveLoad(ConfFile()); err != nil {
		log.Println("load empty conf")
		_conf, _ = ini.LoadSources(ini.LoadOptions{Insensitive: true}, []byte(""))
		conf = _conf
	} else {
		conf = _conf
	}
	return conf
}

func ReloadConf() *ini.File {
	if _conf, err := ini.InsensitiveLoad(ConfFile()); err != nil {
		log.Println("load empty conf")
		_conf, _ = ini.LoadSources(ini.LoadOptions{Insensitive: true}, []byte(""))
		conf = _conf
	} else {
		conf = _conf
	}
	return conf
}

func SaveToConf(section string, kvmap map[string]string) {
	ReloadConf()
	sec := conf.Section(section)
	for k, v := range kvmap {
		sec.Key(k).SetValue(v)
	}
	conf.SaveTo(ConfFile())
}

func ExpandHomeDir(path string) string {
	if len(path) == 0 {
		return path
	}
	if path[0] != '~' {
		return path
	}
	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return path
	}
	return filepath.Join(HomeDir(), path[1:])
}

func EnsureDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func Exisit(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func ShortID() string {
	return shortid.MustGenerate()
}

func PauseExit() {
	log.Println("Press any to exit")
	keyboard.GetSingleKey()
	os.Exit(0)
}

func PauseGo(msg ...interface{}) {
	log.Println(msg...)
	keyboard.GetSingleKey()
}

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(StringArray(""))
	ini.PrettyFormat = false
}
