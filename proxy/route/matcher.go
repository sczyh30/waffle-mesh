package route

import (
	"strings"
	"regexp"
	"github.com/sczyh30/waffle-mesh/api/gen"
)

func matchExactPath(curRoute *api.RouteEntry, path string) bool {
	return curRoute.Match.GetExactPath() != "" && curRoute.Match.GetExactPath() == path
}

func matchPrefixPath(curRoute *api.RouteEntry, path string) bool {
	return curRoute.Match.GetPrefix() != "" && strings.HasPrefix(path, curRoute.Match.GetPrefix())
}

func matchRegexPath(curRoute *api.RouteEntry, path string) bool {
	regexPattern := curRoute.Match.GetRegex()
	if regexPattern != "" {
		res, err := regexp.MatchString(regexPattern, path)
		return err == nil && res
	}
	return false
}

func matchExactHeader(match *api.HeaderMatch, v string) bool {
	return match.GetExactMatch() != "" && match.GetExactMatch() == v
}

func matchRegexHeader(match *api.HeaderMatch, v string) bool {
	regexPattern := match.GetRegexMatch()
	if regexPattern != "" {
		res, err := regexp.MatchString(regexPattern, v)
		return err == nil && res
	}
	return false
}
