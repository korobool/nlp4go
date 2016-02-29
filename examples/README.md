### Train model
```
go run tagger_train.go -corpus /home/user/ontonotes -model test-model.go
```
### Run STDIN tagger
```
go run tagger_tag.go -model test-model.go
```
