package main

import (
	"context"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/philippgille/chromem-go"
	"github.com/rs/zerolog/log"
)

// StatisticalEmbedder implements a simple statistical embedding function
type StatisticalEmbedder struct {
	dimensions int
}

func NewStatisticalEmbedder() chromem.EmbeddingFunc {
	embedder := &StatisticalEmbedder{
		dimensions: 384, // Standard embedding dimension
	}
	
	return func(ctx context.Context, text string) ([]float32, error) {
		log.Debug().Str("text", text).Msg("Embedding text")
		return embedder.embedText(text), nil
	}
}

func (e *StatisticalEmbedder) embedText(text string) []float32 {
	// Normalize text
	text = strings.ToLower(text)
	text = regexp.MustCompile(`[^\p{L}\p{N}\s]+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	
	// Tokenize
	words := strings.Fields(text)
	if len(words) == 0 {
		return make([]float32, e.dimensions)
	}
	
	// Calculate word frequencies
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}
	
	// Create embedding vector
	embedding := make([]float32, e.dimensions)
	
	// Statistical features
	textLen := float32(len(text))
	wordCount := float32(len(words))
	uniqueWords := float32(len(wordFreq))
	
	// Basic statistical features (first 10 dimensions)
	embedding[0] = textLen / 1000.0                    // Text length (normalized)
	embedding[1] = wordCount / 100.0                   // Word count (normalized)
	embedding[2] = uniqueWords / wordCount             // Lexical diversity
	embedding[3] = e.averageWordLength(words)          // Average word length
	embedding[4] = e.calculateEntropy(wordFreq)        // Text entropy
	embedding[5] = e.countCapitalLetters(text)         // Capital letter ratio
	embedding[6] = e.countDigits(text)                 // Digit ratio
	embedding[7] = e.countPunctuation(text)            // Punctuation ratio
	embedding[8] = e.calculateReadability(words)       // Simple readability score
	embedding[9] = e.calculateSentenceComplexity(text) // Sentence complexity
	
	// Character n-gram features (dimensions 10-99)
	charNgrams := e.extractCharNgrams(text, 2, 3)
	for i, ngram := range charNgrams {
		if i+10 >= 100 {
			break
		}
		embedding[i+10] = float32(ngram.freq) / textLen
	}
	
	// Word-based features (dimensions 100-299)
	sortedWords := e.getSortedWords(wordFreq)
	for i, word := range sortedWords {
		if i+100 >= 300 {
			break
		}
		// Use a simple hash function to map words to dimensions
		hash := e.simpleHash(word.word) % 200
		embedding[100+hash] += float32(word.freq) / wordCount
	}
	
	// Positional features (dimensions 300-383)
	for i, word := range words {
		if i >= 84 {
			break
		}
		hash := e.simpleHash(word) % 84
		position := float32(i) / wordCount
		embedding[300+hash] += position
	}
	
	// Normalize the embedding vector
	return e.normalizeVector(embedding)
}

type ngramFreq struct {
	ngram string
	freq  int
}

type wordFreqType struct {
	word string
	freq int
}

func (e *StatisticalEmbedder) extractCharNgrams(text string, minN, maxN int) []ngramFreq {
	ngramMap := make(map[string]int)
	
	for n := minN; n <= maxN; n++ {
		for i := 0; i <= len(text)-n; i++ {
			ngram := text[i : i+n]
			ngramMap[ngram]++
		}
	}
	
	var ngrams []ngramFreq
	for ngram, freq := range ngramMap {
		ngrams = append(ngrams, ngramFreq{ngram, freq})
	}
	
	sort.Slice(ngrams, func(i, j int) bool {
		return ngrams[i].freq > ngrams[j].freq
	})
	
	return ngrams
}

func (e *StatisticalEmbedder) getSortedWords(wordFreq map[string]int) []wordFreqType {
	var words []wordFreqType
	for word, freq := range wordFreq {
		words = append(words, wordFreqType{word, freq})
	}
	
	sort.Slice(words, func(i, j int) bool {
		return words[i].freq > words[j].freq
	})
	
	return words
}

func (e *StatisticalEmbedder) averageWordLength(words []string) float32 {
	if len(words) == 0 {
		return 0
	}
	
	totalLen := 0
	for _, word := range words {
		totalLen += len(word)
	}
	
	return float32(totalLen) / float32(len(words)) / 10.0 // Normalized
}

func (e *StatisticalEmbedder) calculateEntropy(wordFreq map[string]int) float32 {
	total := 0
	for _, freq := range wordFreq {
		total += freq
	}
	
	if total == 0 {
		return 0
	}
	
	entropy := 0.0
	for _, freq := range wordFreq {
		p := float64(freq) / float64(total)
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	
	return float32(entropy) / 10.0 // Normalized
}

func (e *StatisticalEmbedder) countCapitalLetters(text string) float32 {
	capitals := 0
	total := 0
	
	for _, r := range text {
		if unicode.IsLetter(r) {
			total++
			if unicode.IsUpper(r) {
				capitals++
			}
		}
	}
	
	if total == 0 {
		return 0
	}
	
	return float32(capitals) / float32(total)
}

func (e *StatisticalEmbedder) countDigits(text string) float32 {
	digits := 0
	total := len(text)
	
	for _, r := range text {
		if unicode.IsDigit(r) {
			digits++
		}
	}
	
	if total == 0 {
		return 0
	}
	
	return float32(digits) / float32(total)
}

func (e *StatisticalEmbedder) countPunctuation(text string) float32 {
	punct := 0
	total := len(text)
	
	for _, r := range text {
		if unicode.IsPunct(r) {
			punct++
		}
	}
	
	if total == 0 {
		return 0
	}
	
	return float32(punct) / float32(total)
}

func (e *StatisticalEmbedder) calculateReadability(words []string) float32 {
	if len(words) == 0 {
		return 0
	}
	
	// Simple readability based on average word length and sentence count
	totalChars := 0
	sentences := 1 // Assume at least one sentence
	
	for _, word := range words {
		totalChars += len(word)
		if strings.HasSuffix(word, ".") || strings.HasSuffix(word, "!") || strings.HasSuffix(word, "?") {
			sentences++
		}
	}
	
	avgWordsPerSentence := float32(len(words)) / float32(sentences)
	avgCharsPerWord := float32(totalChars) / float32(len(words))
	
	// Simple readability score (lower is more readable)
	readability := (avgWordsPerSentence * 0.39) + (avgCharsPerWord * 11.8)
	
	return readability / 100.0 // Normalized
}

func (e *StatisticalEmbedder) calculateSentenceComplexity(text string) float32 {
	sentences := regexp.MustCompile(`[.!?]+`).Split(text, -1)
	if len(sentences) <= 1 {
		return 0
	}
	
	totalComplexity := 0.0
	validSentences := 0
	
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) == 0 {
			continue
		}
		
		words := strings.Fields(sentence)
		if len(words) == 0 {
			continue
		}
		
		// Complexity based on sentence length and word variety
		wordSet := make(map[string]bool)
		for _, word := range words {
			wordSet[word] = true
		}
		
		complexity := float64(len(words)) * (float64(len(wordSet)) / float64(len(words)))
		totalComplexity += complexity
		validSentences++
	}
	
	if validSentences == 0 {
		return 0
	}
	
	return float32(totalComplexity/float64(validSentences)) / 20.0 // Normalized
}

func (e *StatisticalEmbedder) simpleHash(s string) int {
	hash := 0
	for _, r := range s {
		hash = hash*31 + int(r)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

func (e *StatisticalEmbedder) normalizeVector(vec []float32) []float32 {
	// Calculate magnitude
	magnitude := float32(0)
	for _, v := range vec {
		magnitude += v * v
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))
	
	if magnitude == 0 {
		return vec
	}
	
	// Normalize
	normalized := make([]float32, len(vec))
	for i, v := range vec {
		normalized[i] = v / magnitude
	}
	
	return normalized
}