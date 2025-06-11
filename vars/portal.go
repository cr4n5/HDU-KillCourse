package vars

import "fmt"

var portal = `
  _   _   ____    _   _           _  __  _   _   _    ____                                     
 | | | | |  _ \  | | | |         | |/ / (_) | | | |  / ___|   ___    _   _   _ __   ___    ___ 
 | |_| | | | | | | | | |  _____  | ' /  | | | | | | | |      / _ \  | | | | | '__| / __|  / _ \
 |  _  | | |_| | | |_| | |_____| | . \  | | | | | | | |___  | (_) | | |_| | | |    \__ \ |  __/
 |_| |_| |____/   \___/          |_|\_\ |_| |_| |_|  \____|  \___/   \__,_| |_|    |___/  \___|

HDU-KillCourse[https://github.com/cr4n5/HDU-KillCourse]               version: ` + Version + `
`

func ShowPortal() {
	fmt.Println(portal)
}
