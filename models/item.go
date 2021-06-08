package models

type Item struct {
	Occurrences    int32
	AssociateItems map[string]int32
}

type PairItems struct {
	Item1         string
	Item2         string
	CoOccurrences int32
}