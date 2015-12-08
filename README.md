# nlp4go


The idea of nlp4go is to provide a fast go-lang based nlp toolkit for researchers and developers which provides the most commonly used features of NLTK
and other NPL toolkits, but with production-ready computational performance. 
Things that can be executed in parallel mode (like POS tagging for independent sentences) should be processe in go-routines in parallel to utilize CPU cores efficiently.


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
* Immutable Article abstraction (should come up with good interface) 
 * Article abstraction itself (should come up with good interface) 
 * Handlers which can cache NLP algorithm resul in HandlerResulr (or something similar) (should come up with good interface) 
 * Article Meta Info storage (should come up with good interface) 
* Syntax Parsing 
 * We need a complete rules set for syntax parsing
 * FSM based parser
 * Investigate alternatives
* WordNet interface
* Corpuses interface (like Brown corpus in NLTK)
* Data Downloader

## Files and utils list

 
**errors.go** contains constatns for errors to return.  
**tagger.go** librarary src for set of differently implemented taggers  
**tokenizer.go** librarary src for set of differently implemented tokenizers  
**utils.go** set of helper functions  

**utils/parse** parses ontonotes dir and outputs sentences per line in (tag word) format to stdout  
**utils/train_stdin** trains model using input per line sentences (in (tag word) format) from stdin  
**utils/tags** reads regular sentences from stdin and ouputs words with POS and positions  

**tokenize/** experimental package for tokenizers   

POS tagger should be trained. For research we use ontonotes v.5 and provide a reader for it. Its output can be used to train a pos-tagger:

./parse -p /path/to/ontontoses/folder | go run ../train_stdin/main.go

We cannot share ontonotes, but you can use your own training data, just feed to train_stdin data in format:

`(IN In)(DT the)(NN summer)(IN of)(CD 2005)(, ,)(DT a)(NN picture)(WDT that)(NNS people)(VBP have)(RB long)(VBN been)(VBG looking)(RB forward)(IN to)(-NONE- *T*-1)(VBD started)(-NONE- *-2)(VBG emerging)(IN with)(NN frequency)(IN in)(JJ various)(JJ major)(NNP Hong)(NNP Kong)(NNS media)(. .)`

`(IN With)(PRP$ their)(JJ unique)(NN charm)(, ,)(DT these)(RB well)(HYPH -)(VBN known)(NN cartoon)(NNS images)(RB once)(RB again)(VBD caused)(NNP Hong)(NNP Kong)(TO to)(VB be)(DT a)(NN focus)(IN of)(JJ worldwide)(NN attention)(. .)`

`(DT The)(NN world)(POS 's)(JJ fifth)(NNP Disney)(NN park)(MD will)(RB soon)(VB open)(IN to)(DT the)(NN public)(RB here)(. .)`
