package route

import (
	"strings"
	"regexp"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type Matcher struct {}

func (m *Matcher) matchExactPath(curRoute *api.RouteEntry, path string) bool {
	return curRoute.Match.GetExactPath() != "" && curRoute.Match.GetExactPath() == path
}

func (m *Matcher) matchPrefixPath(curRoute *api.RouteEntry, path string) bool {
	return curRoute.Match.GetPrefix() != "" && strings.HasPrefix(path, curRoute.Match.GetPrefix())
}

func (m *Matcher) matchRegexPath(curRoute *api.RouteEntry, path string) bool {
	regexPattern := curRoute.Match.GetRegex()
	if regexPattern != "" {
		res, err := regexp.MatchString(regexPattern, path)
		return err == nil && res
	}
	return false
}

func (m *Matcher) matchExactHeader(match *api.HeaderMatch, v string) bool {
	return match.GetExactMatch() != "" && match.GetExactMatch() == v
}

func (m *Matcher) matchRegexHeader(match *api.HeaderMatch, v string) bool {
	regexPattern := match.GetRegexMatch()
	if regexPattern != "" {
		res, err := regexp.MatchString(regexPattern, v)
		return err == nil && res
	}
	return false
}
