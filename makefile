build:
	clear
	go build pairwise-association-rule

run: build
	./pairwise-association-rule --dataFile data/sample.csv --inputItems ab


