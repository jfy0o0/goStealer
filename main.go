package main

import (
	"fmt"
	"goHero/net/ghssh"
	"log"
	"runtime"
)

func test(skip int) {
	call(skip)
}

func call(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	pcName := runtime.FuncForPC(pc).Name()
	fmt.Println(fmt.Sprintf("%v   %s   %d   %t   %s", pc, file, line, ok, pcName))
}

func main() {
	sf, err := ghssh.NewSftpSSH("root", "_shang@0107", "192.168.12.32:822")
	if err != nil {
		log.Fatalln(err)
	}
	defer sf.Close()
	if err := sf.PutFile("/opt/TM8K.yaml", "/opt/TM8K.yaml.bak2"); err != nil {
		log.Println(err)
	}
	//c := ghcode.New(1, "hi", "detailed")
	//fmt.Println(c.Code(), c.Message(), c.Detail())
	//st := ghssh.NewTunnelSSH("root", "_shang@0107", "1190", "192.168.12.32:822", "127.0.0.1:30000")
	//st.Run()
}
