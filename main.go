package main

import (
	// "encoding/json"
	"fmt"
	"qawwali-glossary-db/database"
	"qawwali-glossary-db/process"
	"qawwali-glossary-db/scrape"
	"qawwali-glossary-db/translate"

	"sync"

	"github.com/schollz/progressbar/v3"
)

// var urlList = []string{"https://sufinama.org/poets/amir-khusrau/kalaam", "https://sufinama.org/poets/bulleh-shah/kaafi"}

// var urlList = []string{"https://sufinama.org/poets/amir-khusrau/kalaam"}
// var urlList = []string{"https://sufinama.org/sufi-kalam/best-10-qawwalis-of-rumi"}
// var urlList = []string{"https://sufinama.org/poets/kabir/pad"}
var urlList = []string{"https://sufinama.org/poets/sultan-bahu/kalaam"}

// var urlList = []string{
// 	"https://sufinama.org/poets/kabir/raga-based-poetries",
// 	"https://sufinama.org/poets/kabir/chaupaiyan",
// 	"https://sufinama.org/poets/kabir/shabad",
// 	"https://sufinama.org/poets/khwaja-ghulam-farid",
// 	"https://sufinama.org/poets/baba-farid",
// }

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

	bar := progressbar.New(len(countedWords))
	// for _, word := range countedWords {
	// 	translatedWord, untranslatedWord := translate.Rekhta(word)
	// 	if untranslatedWord.Word != "" {
	// 		untranslatedWords = append(untranslatedWords, untranslatedWord)
	// 	}
	// 	for _, wordInstance := range translatedWord {
	//
	// 		translatedWords = append(translatedWords, wordInstance)
	// 	}
	// }
	// p, _ := json.Marshal(translatedWords)
	// utils.Print(p)

	// Channels for concurrent processing
	translatedChan := make(chan translate.TranslatedWord, len(countedWords))
	untranslatedChan := make(chan translate.UntranslatedWord, len(countedWords))

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	for _, word := range countedWords {
		wg.Add(1)

		// Start a goroutine for each translation request
		go func(w translate.CountedWord) {
			defer wg.Done()

			translatedWord, untranslatedWord := translate.Rekhta(w)

			if untranslatedWord.Word != "" {
				untranslatedChan <- untranslatedWord
			}
			for _, wordInstance := range translatedWord {
				translatedChan <- wordInstance
			}
			// Update the progress bar after each word is processed
			bar.Add(1)
		}(word)
	}

	// Close the channels once all goroutines are done
	go func() {
		wg.Wait()
		close(translatedChan)
		close(untranslatedChan)
	}()

	// Collect results from the channels
	for translated := range translatedChan {
		translatedWords = append(translatedWords, translated)
	}

	for untranslated := range untranslatedChan {
		if untranslated.Word != "" {
			untranslatedWords = append(untranslatedWords, untranslated)
		}
	}

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
