# nlp4go

[![Join the chat at https://gitter.im/korobool/nlp4go](https://badges.gitter.im/korobool/nlp4go.svg)](https://gitter.im/korobool/nlp4go?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


The idea of nlp4go is to provide a fast go-lang based nlp toolkit for researchers and developers which provides the most commonly used features of NLTK
and other NPL toolkits, but with production-ready computational performance. 
Things that can be executed in parallel mode (like POS tagging for independent sentences) should be processed in go-routines in parallel to utilize CPU cores efficiently.


### Currently supported languages
* English
* Russian


### General plan for implementation
* Tokenizer(s)
 * Tree bank
 * regex
 * split
 * investigate alternatives
* POS tagger(s)
 * Percepton 
* String abstraction to imporove performance on unicode
 * Compatibility with regex
 * O(1) len() operation for unicode strings
 * Slises in bytes and characters
* Immutable Article abstraction (should come up with good interface) 
* NER support
* Parsing 
 * We need a complete rules set for syntax parsing
 * Syntax Parsing
 * Dependency Parsing
* WordNet interface

## Repository structure should be idiomatically similar to following tree:

```
legacy  // to be removed later
core
    strings.go
    exregex.go
ml
    perceptron.go 
    ...
tokenize
    wordsplit.go
    wordregex.go
    sentencesplit.go
    sentencesregex.go
    ...
tagg
    perceptron-pos
parse
    syntax
    dependency
    ...
utils
   train_pos_tagger.go
   read_ontonotes.go
```

### Train model
```
go run tagger_train.go -corpus /home/user/ontonotes -model test-model.go
```
### Run STDIN tagger
```
go run tagger_tag.go -model test-model.go
```

We cannot share ontonotes, but you can use your own training data, just feed to train_stdin data in format:
```
`(IN In)(DT the)(NN summer)(IN of)(CD 2005)(, ,)(DT a)(NN picture)(WDT that)(NNS people)(VBP have)(RB long)(VBN been)(VBG looking)(RB forward)(IN to)(-NONE- *T*-1)(VBD started)(-NONE- *-2)(VBG emerging)(IN with)(NN frequency)(IN in)(JJ various)(JJ major)(NNP Hong)(NNP Kong)(NNS media)(. .)`

`(IN With)(PRP$ their)(JJ unique)(NN charm)(, ,)(DT these)(RB well)(HYPH -)(VBN known)(NN cartoon)(NNS images)(RB once)(RB again)(VBD caused)(NNP Hong)(NNP Kong)(TO to)(VB be)(DT a)(NN focus)(IN of)(JJ worldwide)(NN attention)(. .)`

`(DT The)(NN world)(POS 's)(JJ fifth)(NNP Disney)(NN park)(MD will)(RB soon)(VB open)(IN to)(DT the)(NN public)(RB here)(. .)`
```
