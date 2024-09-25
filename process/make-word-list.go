package process

import (
	"bufio"
	"fmt"
	// "fmt"
	"log"
	"os"
	"qawwali-syllabus/translate"
	"regexp"
	"sort"
	"strings"
)

func wordListFromFile() {
	var counts = make(map[string]float64)
	var countTotal float64

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		counts[re.ReplaceAllString(strings.ToLower(scanner.Text()), "")] += 1.0
		countTotal += 1.0
	}

	var words = make([]string, 0, len(counts))
	for word := range counts {
		words = append(words, word)
	}
	sort.SliceStable(words, func(i, j int) bool {
		return counts[words[i]] < counts[words[j]]
	})

	// for _, word := range words {
	// 	fmt.Println(word, ": ", counts[word], "/", countTotal)
	// }
}

//-------------//

func Words2(words []translate.Word) (map[translate.Word]int, []translate.Word) {
	counts := make(map[translate.Word]int)
	var countTotal int
	for _, word := range words {
		counts[word] += 1
		countTotal += 1
	}
	fmt.Println("counts: ", counts)

	keys := make([]translate.Word, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	fmt.Println("keys: ", keys)

	sort.SliceStable(keys, func(i, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})

	return counts, keys
}

func Words(words []translate.Word) []translate.CountedWord {
	wordMap := make(map[string]*translate.CountedWord)
	urlMap := make(map[string]map[string]bool) // To track URLs for each key

	// Aggregate words
	for _, word := range words {
		key := word.Text + "|" + word.Data
		if _, exists := wordMap[key]; !exists {
			wordMap[key] = &translate.CountedWord{
				Text:        word.Text,
				Data:        word.Data,
				Occurrences: "",
				Count:       1,
			}
			urlMap[key] = make(map[string]bool) // Initialize the URL tracking map
		} else {
			wordMap[key].Count++
		}

		// If the URL hasn't been added for this word, add it
		if !urlMap[key][word.Occurrence] {
			if wordMap[key].Occurrences == "" {
				wordMap[key].Occurrences = word.Occurrence
			} else {
				wordMap[key].Occurrences += ", " + word.Occurrence
			}
			urlMap[key][word.Occurrence] = true
		}
	}

	// Convert the map to a slice
	var result []translate.CountedWord
	for _, w := range wordMap {
		result = append(result, *w)
	}
	// Sort by count in descending order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}
