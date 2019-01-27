package transport

import (
	"regexp"
	"strings"
)

const (
	hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-" +
		"([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"
)

func IsUUID(str string) bool {
	md := regexp.MustCompile(hexPattern).FindStringSubmatch(str)
	if md == nil {
		return false
	} else {
		return true
	}
}

func ConvertUriToQueueName(uriFull string) string {
	res := ""
	uri := uriFull
	if strings.Index(uri, "?") > 0 {
		uri = uri[:strings.Index(uri, "?")]
	}
	uriItems := strings.Split(uri, "/")
	if len(uriItems) == 1 {
		res = uri
	} else {
		for u := range uriItems {
			if uriItems[u] != "" {
				if IsUUID(uriItems[u]) {
					uriItems[u] = ":uuid"
				}
				res = res + "/" + uriItems[u]
			} else if u == len(uriItems)-1 {
				res = res + "/"
			}
		}
	}
	return res
}