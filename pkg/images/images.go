package images

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ListAppImages() (images []string) {
	imagesLength := 10
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	go startSpinner(s, "fetching images list...", nil)
	url := "https://registry-1.docker.io/"
	username := "cnvrghelm"
	password := "23e37770-0a2c-4111-b967-7e16e597a252"
	hub, err := registry.New(url, username, password)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error fetching images from docker hub")
	}
	tags, err := hub.Tags("cnvrg/app")
	tagRegex, _ := regexp.Compile("^master-\\d*-encode$")
	var filteredTags []int
	for _, tag := range tags {
		if tagRegex.MatchString(tag) {
			tagNumber, _ := strconv.Atoi(strings.Split(tag, "-")[1])
			filteredTags = append(filteredTags, tagNumber)
		}
	}
	logrus.Info(len(filteredTags) )
	if len(filteredTags) == 0 {
		logrus.Fatal("no images available for upgrade")
	}
	sort.Sort(sort.Reverse(sort.IntSlice(filteredTags)))
	for i := 0; i < imagesLength; i++ {
		images = append(images, "docker.io/cnvrg/app:master-"+strconv.Itoa(filteredTags[i])+"-encode")
	}
	s.Stop()
	return
}

func startSpinner(s *spinner.Spinner, suffixMessage string, messages <-chan string) {
	s.Suffix = suffixMessage
	s.Color("green")
	s.Start()
	for v := range messages {
		msg := fmt.Sprintf("%v [ %v ]", suffixMessage, v)
		s.Suffix = msg
		s.Restart()
	}
}
