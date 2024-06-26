package crawler

import (
	"GameDB/internal/model"
	"regexp"
	"strings"
)

const KaOsKrewName string = "KaOsKrew-torrents"

func CrawlKaOsKrewMulti(pages []int) ([]*model.GameDownload, error) {
	return Crawl1337xMulti(DODIName, pages, DODIFormatter)
}

func CrawlKaOsKrewAll() ([]*model.GameDownload, error) {
	return Crawl1337xAll(DODIName, DODIFormatter)
}

var kaOsKrewRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\.REPACK2?-KaOs`),
	regexp.MustCompile(`(?i)\.UPDATE-KaOs`),
	regexp.MustCompile(`(?i)v\.?\d+(\.\d+)*|Build\.\d+`),
	regexp.MustCompile(`(?i)\.MULTi\d+`),
}

func KaOsKrewFormatter(name string) string {
	if index := kaOsKrewRegexps[2].FindIndex([]byte(name)); index != nil {
		name = name[:index[0]]
	}
	for _, re := range kaOsKrewRegexps {
		name = re.ReplaceAllString(name, "")
	}
	name = strings.Replace(name, ".", " ", -1)
	name = regexp.MustCompile(`(?i)\sgoty`).ReplaceAllString(": Game Of The Year", name)
	return strings.TrimSpace(name)
}
