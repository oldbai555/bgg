package wordscheck

import (
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func TestDoCheckWords(t *testing.T) {
	_, err := DoCheckWords("http://www.wordscheck.com/wordcheck", "打打杀杀")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
