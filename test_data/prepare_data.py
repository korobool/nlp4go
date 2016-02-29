import json
from nltk.tokenize import TreebankWordTokenizer


def read_referrence_data(filename):
    with open(filename, "r") as f:
        lines = f.readlines()
        return [l.strip('\n\r') for l in lines]

def span_tokenize(text):

    tokens = TreebankWordTokenizer().tokenize(text)
    
    dt = 0
    end = 0
    spaned_tokens = []

    for token in tokens:

        raw_token = token
        is_quote_end = False
        is_quote_start = False
        is_ellipsis = False

        if token == '``':
            raw_token = '"'
            is_quote_start = True
        elif token == "''":
            raw_token = '"'
            is_quote_end = True
        elif token == "...":
            is_ellipsis = True

        start = text[dt:].find(raw_token)
        if start != -1:
            end = start + len(raw_token)
            spaned_tokens.append({
                "word": token,
                "runes": [ord(c) for c in token],
                "pos": dt+start,
                "pos_end": dt+end,
                "is_quote_start": is_quote_start,
                "is_quote_end": is_quote_end,
                "is_ellipsis": is_ellipsis,
            })
            dt += end

    return spaned_tokens


if __name__ == "__main__":

    ref_sentences = read_referrence_data("sentences.en.txt")

    with open('sentences.en.json', 'w') as outfile:
        sent =[]
        for sentence in ref_sentences:
            sent.append({ 
                "sentence": sentence,
                "tokens":  span_tokenize(sentence),
            })
        json.dump(sent, outfile, indent=4)
