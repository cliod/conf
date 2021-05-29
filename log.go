package conf

import "log"

func eLog(err error, a ...interface{}) {
	if err != nil {
		iLog("[CONF ERROR]", a...)
	}
}

func wLog(err error, a ...interface{}) {
	if err != nil {
		iLog("[CONF WAINING]", a...)
	}
}

func iLog(msg string, a ...interface{}) {
	var arr []interface{}
	arr = append(arr, msg)
	arr = append(arr, a...)
	log.Println(arr...)
}
