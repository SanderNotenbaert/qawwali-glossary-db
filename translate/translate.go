package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "qawwali-syllabus/utils"
	"github.com/avast/retry-go/v4"
	"strings"
)

type UntranslatedWord struct {
	Word        string
	Word_count  int
	Occurrences string
}

func Rekhta(word CountedWord) ([]TranslatedWord, UntranslatedWord) {
	count := word.Count
	data := map[string]string{
		"Word":         word.Data,
		"SelectedWord": word.Text,
	}
	dataType := "application/json"
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Make POST request with JSON data
	r, err := postWithRetry("https://world.rekhta.org/api/v2/shayari/GetGroupWordMeaning?lang=1", dataType, bytes.NewBuffer(jsonData))
	if (err == nil) && (len(r.R) == 0) {
		fmt.Println("No translation found")
		return []TranslatedWord{}, UntranslatedWord{Word: word.Text, Word_count: count, Occurrences: word.Occurrences}
	}

	var words []TranslatedWord

	for _, t := range r.R {
		var translations []string
		for _, translation := range t.Translation {
			translation.Meaning = strings.ReplaceAll(translation.Meaning, `"`, "'")
			translations = append(translations, fmt.Sprintf(`"%s"`, translation.Meaning))
		}
		words = append(words, TranslatedWord{
			Urdu_romanized:         word.Text,
			Urdu:                   t.U,
			Hindi:                  t.H,
			Word_count:             count,
			Audio_file:             t.AMF,
			Alternative_spelling_1: t.E,
			Origin:                 t.WO,
			Word_type:              t.WP,
			Translations:           fmt.Sprintf("[%s]", strings.Join(translations, ", ")),
			Occurrences:            word.Occurrences,
		})

	}
	return words, UntranslatedWord{}
}

func postWithRetry(url string, dataType string, data io.Reader) (rekhtaResponse, error) {
	body, err := retry.DoWithData(
		func() (rekhtaResponse, error) {
			resp, err := http.Post(url, dataType, data)
			if err != nil {
				return rekhtaResponse{}, err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return rekhtaResponse{}, err
			}
			// utils.Print(body)

			var r rekhtaResponse
			err = json.Unmarshal(body, &r)
			if err != nil {
				panic(err)
			}
			if len(r.R) < 1 {
				// err = fmt.Errorf("no translation found")
				// fmt.Println(err)
				return rekhtaResponse{}, nil
			}
			// retry.RetryIf(func(err error) bool {
			// 	return err == fmt.Errorf("no translation found")
			// })
			return r, nil
		},
	)

	return body, err
}
