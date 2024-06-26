package crawler

import (
	"GameDB/internal/model"
	"regexp"
	"strings"
)

const FitgirlName string = "FitGirl-torrents"

func CrawlFitgirlMulti(pages []int) ([]*model.GameDownload, error) {
	return Crawl1337xMulti(FitgirlName, pages, DODIFormatter)
}

func CrawlFitgirlAll() ([]*model.GameDownload, error) {
	return Crawl1337xAll(FitgirlName, DODIFormatter)
}

var fitgirlRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\(.*\)`),
	regexp.MustCompile(`(?i)\[.*?\]`),
	regexp.MustCompile(`(?i)-.*?(Edition|Bundle|Pack|Set|Remake|Collection)`),
}

func FitgirlFormatter(name string) string {
	for _, re := range fitgirlRegexps {
		name = re.ReplaceAllString(name, "")
	}
	name = strings.ReplaceAll(name, "+ OST", "")
	name = strings.ReplaceAll(name, "- Digital Deluxe", "")
	return strings.TrimSpace(name)
}
