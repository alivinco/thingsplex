package utils

import (
	"bufio"
	"bytes"
	"github.com/buger/jsonparser"
	"os"
	"strings"
)

func match(route []string, topic []string) bool {
	if len(route) == 0 {
		if len(topic) == 0 {
			return true
		}
		return false
	}

	if len(topic) == 0 {
		if route[0] == "#" {
			return true
		}
		return false
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return match(route[1:], topic[1:])
	}

	return false
}

func RouteIncludesTopic(route, topic string) bool {
	return match(strings.Split(route, "/"), strings.Split(topic, "/"))
}

type LogFilter struct {
	Component string
	FlowId string
}



func GetLogs(logFile string,filter *LogFilter,limit int64,isLogJsonFormated bool) []byte {
	file, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, limit)
	stat, err := os.Stat(logFile)
	start := stat.Size() - limit
	_, err = file.ReadAt(buf, start)
	//if err == nil {
	//	if filter.FlowId == "" {
	//		return buf
	//	}
	//}

	if err != nil {
		return nil
	}

	bytesReader := bytes.NewReader(buf)
	scanner := bufio.NewScanner(bytesReader)
	result := []byte{'['}
	scanner.Scan() // skipping first line in case if it's incomplete
	for scanner.Scan() {
		if filter.FlowId == "" {
			result = append(result, scanner.Bytes()...)
			result = append(result,',')
		}else {
			flowId ,err := jsonparser.GetString(scanner.Bytes(),"fid")
			if err == nil && flowId == filter.FlowId {
				result = append(result, scanner.Bytes()...)
				result = append(result,',')

			}
		}

	}
	result[len(result)-1]= ']'
	return result
}

