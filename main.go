package main

import (
	"encoding/json"
	"fmt"
	"qawwali-syllabus/database"
	"qawwali-syllabus/process"
	"qawwali-syllabus/scrape"
	"qawwali-syllabus/translate"
	"qawwali-syllabus/utils"
)

var urlList = []string{"https://sufinama.org/poets/amir-khusrau/kalaam", "https://sufinama.org/poets/bulleh-shah/kaafi"}

// var urlList = []string{"https://sufinama.org/poets/amir-khusrau/kalaam"}

// var urlList = []string{"https://sufinama.org/sufi-kalam/best-10-kalams-of-shah-turab-ali-qalandar"}
var domain = "sufinama.org"

// var urlList = []string{"https://sufinama.org/kalaam/bahut-rahii-baabul-ghar-dulhan-amir-khusrau-kalaam-13"}

func main() {
	//collect words
	db := database.Connect()
	defer db.Close()
	usedUrls := database.QueryRows(db, "SELECT url FROM USED_LINKS", "")
	words, newUrls := scrape.Sufinama(domain, urlList, usedUrls)
	// fmt.Println(words)
	database.RecursiveEntries(db, newUrls, "USED_LINKS", "")

	//count words, remove duplicates
	countedWords := process.Words(words)

	//get translations on the collected words
	var translatedWords []interface{} //translate.TranslatedWord
	var untranslatedWords []interface{}
	for _, word := range countedWords {
		translatedWord, untranslatedWord := translate.Rekhta(word)
		if untranslatedWord.Word != "" {
			untranslatedWords = append(untranslatedWords, untranslatedWord)
		}
		for _, wordInstance := range translatedWord {

			translatedWords = append(translatedWords, wordInstance)
		}
	}
	p, _ := json.Marshal(translatedWords)
	utils.Print(p)
	database.RecursiveEntries(db, untranslatedWords, "UNTRANSLATED_WORDS", `ON CONFLICT(WORD)
DO UPDATE SET 
    word_count = word_count + excluded.word_count,
  occurrences = CASE 
        WHEN instr(occurrences, excluded.occurrences) = 0 THEN 
            occurrences || ', ' || excluded.occurrences
        ELSE 
            occurrences
    END;`)
	failedEntries := database.RecursiveEntries(db, translatedWords, "WORDS", `ON CONFLICT(Urdu_romanized, Translations)
DO UPDATE SET
    word_count = word_count + excluded.word_count;
occurrences = CASE 
        WHEN instr(occurrences, excluded.occurrences) = 0 THEN 
            occurrences || ', ' || excluded.occurrences
        ELSE 
            occurrences
    END;`)
	fmt.Println(failedEntries)
}
