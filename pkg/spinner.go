package pkg

import (
	"fmt"
	"github.com/briandowns/spinner"
	"strings"
)

func StartSpinner(s *spinner.Spinner, suffixMessage string, messages <-chan string) {
	s.Suffix = suffixMessage
	s.Color("green")
	s.Start()
	for v := range messages {
		msg := fmt.Sprintf("%v [ %v ]", strings.TrimSuffix(suffixMessage, "\n"), v)
		s.Suffix = msg
		s.Restart()
	}
}
