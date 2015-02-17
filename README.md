# go-nlp-utils 
(based on go-nlp-tools`)
 
**errors.go** contains constatns for errors to return.  
**tagger.go** librarary src for set of differently implemented taggers  
**tokenizer.go** librarary src for set of differently implemented tokenizers  
**utils.go** set of helper functions  

**utils/parse** parses ontonotes dir and outputs sentences per line in (tag word) format to stdout  
**utils/train_stdin** trains model using input per line sentences (in (tag word) format) from stdin  
**utils/tags** reads regular sentences from stdin and ouputs words with POS and positions  


POS tagger should be trained. For research we use ontonotes v.5 and provide a reader for it. Its output can be used to train a pos-tagger:

./parse -p /path/to/ontontoses/folder | go run ../train_stdin/main.go

We cannot share ontonotes, but you can use your own training data, just feed to train_stdin data in format:

`(IN In)(DT the)(NN summer)(IN of)(CD 2005)(, ,)(DT a)(NN picture)(WDT that)(NNS people)(VBP have)(RB long)(VBN been)(VBG looking)(RB forward)(IN to)(-NONE- *T*-1)(VBD started)(-NONE- *-2)(VBG emerging)(IN with)(NN frequency)(IN in)(JJ various)(JJ major)(NNP Hong)(NNP Kong)(NNS media)(. .)`

`(IN With)(PRP$ their)(JJ unique)(NN charm)(, ,)(DT these)(RB well)(HYPH -)(VBN known)(NN cartoon)(NNS images)(RB once)(RB again)(VBD caused)(NNP Hong)(NNP Kong)(TO to)(VB be)(DT a)(NN focus)(IN of)(JJ worldwide)(NN attention)(. .)`

`(DT The)(NN world)(POS 's)(JJ fifth)(NNP Disney)(NN park)(MD will)(RB soon)(VB open)(IN to)(DT the)(NN public)(RB here)(. .)`
