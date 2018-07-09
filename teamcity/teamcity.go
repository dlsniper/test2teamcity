package teamcity

import (
	"strings"
)

const TeamCityTs = "2006-01-02T15:04:05.000"

func Escape(msg string) string {
	result := strings.Replace(msg, "|", "||", -1)
	result = strings.Replace(result, "\n", "|n", -1)
	result = strings.Replace(result, "[", "|[", -1)
	return strings.Replace(result, "]", "|]", -1)
}
