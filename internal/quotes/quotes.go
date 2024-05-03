package quotes

import (
	"math/rand"
)

var quotes = []string{
	"Do not go where the path may lead, go instead where there is no path and leave a trail.",
	"What you do not want done to yourself, do not do to others.",
	"The only true wisdom is in knowing you know nothing.",
	"The only limit to our realization of tomorrow will be our doubts of today.",
	"It does not matter how slowly you go as long as you do not stop.",
}

func GetRandomQuote() string {
	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}
