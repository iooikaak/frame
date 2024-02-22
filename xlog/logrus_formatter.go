package xlog

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/iooikaak/frame/json"

	nested "github.com/antonfisher/nested-logrus-formatter"
)

var _team = "unknow"

// 应运维要求增加team字段, 来标识日志的不同使用方
func init() {
	if t := os.Getenv("TEAM"); t != "" {
		_team = t
	}
}

// LogrusFormatter formats logs into parsable json
type LogrusFormatter struct {
	Formatter *nested.Formatter
}

// Format renders a single log entry
func (f *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if _namespace == "" {
		return f.Formatter.Format(entry)
	}

	data := make(logrus.Fields, len(entry.Data)+6)

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data["message"] = entry.Message
	data["env"] = _namespace
	data["host"] = _host
	data["level"] = entry.Level.String()
	data["timestamp"] = entry.Time.Unix()
	data["team"] = _team
	serialized, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to NSQ JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
