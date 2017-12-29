package core

import (
	"bytes"
)

const(
	tokenDefaultCapSize = 100
)

var ignoreChars = [...]rune{' ','\t'}
var standTokenChars = [...]rune{':',';','/','[',']','{','}','*','+','-','>','(',')' }

type TokenInfo struct {
	Value 	string
	Line	int
}

type Scanner struct {

}

func (this*Scanner) IsStandChar(c rune) bool {
	for _,v := range standTokenChars{
		if v == c {
			return true
		}
	}
	return false
}

func (this*Scanner) IsIgnoreChar(c rune) bool {
	for _,v := range ignoreChars{
		if v == c {
			return true
		}
	}
	return false
}

func (this*Scanner) IsEOFChar(c rune) bool {
	 if c=='\n'{
	 	return true
	 }
	 return false
}

func (this*Scanner) GetTokens(text string) ([]*TokenInfo,error) {

	tokens := make([]*TokenInfo, 0, tokenDefaultCapSize)

	var buf bytes.Buffer
	i:=0
	for _, c := range text {
		//如果是换行，写入
		if this.IsEOFChar(c) {
			str := buf.String()
			//写入token
			if len(str) != 0 {
				tokens = append(tokens, &TokenInfo{
					Value: str,
					Line:  i,
				})
				buf.Reset()
			}
			//写入换行
			tokens = append(tokens, &TokenInfo{
				Value: "*EOF*",
				Line:  i,
			})
			i++
			continue
		}

		if this.IsStandChar(c) { //如果是独立tonken
			str := buf.String()
			//写入token
			if len(str) != 0 {
				tokens = append(tokens, &TokenInfo{
					Value: str,
					Line:  i,
				})
				buf.Reset()
			}
			//写入独立token
			tokens = append(tokens, &TokenInfo{
				Value: string(c),
				Line:  i,
			})

			continue
		}

		//如果是ignore字符，写入新的tonken
		if this.IsIgnoreChar(c) {
			str := buf.String()
			if len(str) != 0 {
				tokens = append(tokens, &TokenInfo{
					Value: str,
					Line:  i,
				})
				buf.Reset()
			}
		} else {
			buf.WriteRune(c)
		}
	}

	str := buf.String()
	if len(str) != 0 {
		tokens = append(tokens, &TokenInfo{
			Value: str,
			Line:  i,
		})
		buf.Reset()
	}

	return tokens, nil
}