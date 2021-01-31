package pkg

import (
	"fmt"
	"github.com/briandowns/spinner"
)

func StartSpinner(s *spinner.Spinner, suffixMessage string, messages <-chan string) {
	s.Suffix = suffixMessage
	s.Color("green")
	s.Start()
	for v := range messages {
		msg := fmt.Sprintf("%v [ %v ]", suffixMessage, v)
		s.Suffix = msg
		s.Restart()
	}
}

