package crawler

import (
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"regexp"
	"strings"
)

const DODIName string = "DODI-torrents"

func CrawlDODIMulti(pages []int) ([]*model.GameDownload, error) {
	return Crawl1337xMulti(DODIName, pages, DODIFormatter)
}

func CrawlDODIAll() ([]*model.GameDownload, error) {
	return Crawl1337xAll(DODIName, DODIFormatter)
}

var dodiRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\s{2,}`),
	regexp.MustCompile(`(?i)[\-\+]\s?[^:\-]*?\s(Edition|Bundle|Pack|Set|Remake|Collection)`),
}

func DODIFormatter(name string) string {
	name = strings.Replace(name, "- [DODI Repack]", "", -1)
	name = strings.Replace(name, "- Campaign Remastered", "", -1)
	name = strings.Replace(name, "- Remastered", "", -1)
	if index := strings.Index(name, "+"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "–"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "("); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "["); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "- AiO"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "- All In One"); index != -1 {
		name = name[:index]
	}
	for _, re := range dodiRegexps {
		name = strings.TrimSpace(re.ReplaceAllString(name, ""))
	}
	name = strings.TrimSpace(name)
	name = strings.Replace(name, "- Portable", "", -1)
	name = strings.Replace(name, "- Remastered", "", -1)

	if index := strings.Index(name, "/"); index != -1 {
		names := strings.Split(name, "/")
		longestLength := 0
		longestName := ""
		for _, n := range names {
			if !utils.ContainsRussian(n) && len(n) > longestLength {
				longestLength = len(n)
				longestName = n
			}
		}
		name = longestName
	}

	return strings.TrimSpace(name)
}
