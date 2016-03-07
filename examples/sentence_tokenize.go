package main

import (
	"fmt"
	"github.com/korobool/nlp4go/tokenize"
	"log"
)

var TEXT = "Julia with half an ear listened to the list Margery read out and, though she knew the room so well, idly looked about her. It was a very proper room for the manager of a first-class theatre. The walls had been panelled (at cost price) by a good decorator and on them hung engravings of theatrical pictures by Zoffany and de Wilde. The armchairs were large and comfortable. Michael sat in a heavily carved Chippendale* chair, a reproduction but made by a well-known firm, and his Chippendale table, with heavy ball and claw feet, was immensely solid. On it stood in a massive silver frame a photograph of herself and to balance it a photograph of Roger, their son. Between these was a magnificent silver ink-stand that she had herself given him on one of his birthdays and behind it a rack in red morocco, heavily gilt, in which he kept his private paper in case he wanted to write a letter in his own hand. The paper bore the address, Siddons Theatre, and the envelope his crest, a boar's head with the motto underneath: Nemo me impune lacessit.* A bunch of yellow tulips in a silver bowl, which he had got through winning the theatrical golf tournament three times running, showed Margery's care."

func main() {

	tokenizer, err := tokenize.NewDefaultTokenizer()
	if err != nil {
		log.Fatal("Failed to create Tokenizer")
	}

	for _, token := range tokenizer.Tokenize(TEXT) {
		fmt.Println(token)
	}
}
