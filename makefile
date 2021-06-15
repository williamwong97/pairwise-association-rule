build:
	clear
	go build pairwise-association-rule

run: build
	./pairwise-association-rule --dataFile data/groceries.csv --inputItems yogurt,chocolate

run-1: build
	./pairwise-association-rule --dataFile data/groceries_2.csv --inputItems yogurt,chocolate

compare: build
	./pairwise-association-rule --dataFile data/groceries_2.csv --inputItems yogurt
	./pairwise-association-rule --dataFile data/groceries_2.csv --inputItems yogurt,chocolate
	./pairwise-association-rule --dataFile data/groceries_2.csv --inputItems yogurt,chocolate,sugar,
	./pairwise-association-rule --dataFile data/groceries_2.csv --inputItems yogurt,chocolate,sugar,bottled beer
	./pairwise-association-rule --dataFile data/groceries_2.csv --inputItems yogurt,chocolate,sugar,bottled beer,cream cheese

