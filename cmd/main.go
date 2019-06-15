package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"

	"github.com/digital-idea/dilog"
)

var (
	// DBIP 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	DBIP = "127.0.0.1"

	// server setting
	regexpPort       = regexp.MustCompile(`:\d{2,5}$`)
	regexpID         = regexp.MustCompile(`\d{13}$`)
	flagHTTP         = flag.String("http", "", "dilog service port ex):8080")
	flagDBIP         = flag.String("dbip", DBIP, "MongoDB IP")
	flagPagenum      = flag.Int("pagenum", 10, "Number of items on page")
	flagProtocolPath = flag.String("protocolpath", "/show,/lustre,/project,/storage", "A path-aware string to associate with the protocol(dilink). Separate each character with a comma.")
	// add mode
	flagTool    = flag.String("tool", "", "tool name")
	flagProject = flag.String("project", "", "project name")
	flagSlug    = flag.String("slug", "", "shot or asset name")
	flagLog     = flag.String("log", "", "log strings")
	flagKeep    = flag.Int("keep", 180, "Days to keep")
	flagUser    = flag.String("user", "", "custom Username.")
	// find mode
	flagFind = flag.String("find", "", "search word")
	// remove mode
	flagRm   = flag.Bool("rm", false, "Delete data older than keep days")
	flagRmID = flag.String("rmid", "", "ID number to dalete")
	// flag help
	flagHelp = flag.Bool("help", false, "print help")
)

func username() string {
	user, err := user.Current()
	if err != nil {
		if runtime.GOOS == "darwin" {
			return os.ExpandEnv("$USER")
		} else if runtime.GOOS == "linux" {
			return os.ExpandEnv("$USER")
		} else {
			return user.Username
		}
	}
	return user.Username
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("dilog: ")
	flag.Parse()

	// webserver
	if regexpPort.MatchString(*flagHTTP) {
		Webserver()
	}

	// remove mode
	if *flagRm {
		itemlist, err := dilog.All(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range itemlist {
			isDelete, err := dilog.Timecheck(i.Time, i.Keep)
			if err != nil {
				log.Fatal(err)
			}
			if isDelete {
				err := dilog.Remove(*flagDBIP, i.ID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		return
	}
	// remove id mode
	if regexpID.MatchString(*flagRmID) {
		err := dilog.Remove(*flagDBIP, *flagRmID)
		if err != nil {
			log.Fatal(err)
		}
	}
	// add mode
	if *flagTool != "" && *flagLog != "" {
		if *flagUser == "" {
			*flagUser = username()
		}
		ip, err := serviceIP()
		if err != nil {
			log.Fatal(err)
		}
		err = dilog.Add(*flagDBIP, ip, *flagLog, *flagProject, *flagSlug, *flagTool, *flagUser, *flagKeep)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	flag.PrintDefaults()
}
