package main

type Definition struct {
  Origin string
  Word string
  Explanation string
}

type Dictionary struct {
  Origin string
  Definitions []Definition
}

func (d *Dictionary) add(word, explanation string) {
  d.Definitions = append(d.Definitions, Definition{word, explanation, d.Origin})
}
