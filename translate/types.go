package translate

type rekhtaResponse struct {
	R []rekhtaData `json:"R"`
}

type rekhtaData struct {
	E           string        `json:"E"`
	H           string        `json:"H"`
	U           string        `json:"U"`
	WO          string        `json:"WO"`
	WP          string        `json:"WP"`
	AMF         string        `json:"AMF"`
	Translation []translation `json:"WM"`
}

type Word struct {
	Text, Data, Occurrence string
}
type CountedWord struct {
	Text, Data, Occurrences string
	Count                   int
}
type translation struct {
	Meaning string `json:"Meaning"`
}

type TranslatedWord struct {
	Urdu_romanized         string
	Urdu                   string
	Hindi                  string
	Origin                 string
	Word_type              string
	Word_count             int
	Audio_file             string
	Alternative_spelling_1 string
	Alternative_spelling_2 string
	Translations           string
	Occurrences            string
}
type WordTranslations struct {
	Original     string
	Translations []string
}
type NewUrl struct {
	Site string
	Url  string
}
