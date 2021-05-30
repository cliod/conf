package conf

import "log"

func eLog(err error, a ...interface{}) {
	if err != nil {
		iLog("[CONF ERROR]: ", err, a...)
	}
}

func wLog(err error, a ...interface{}) {
	if err != nil {
		iLog("[CONF WAINING]: ", err, a...)
	}
}

func iLog(msg string, err error, a ...interface{}) {
	var arr []interface{}
	arr = append(arr, msg)
	arr = append(arr, err)
	arr = append(arr, a...)
	log.Println(arr...)
}
