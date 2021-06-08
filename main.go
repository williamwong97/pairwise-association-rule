package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"pairwise-association-rule/models"
	"strings"
	"time"
)

func main() {
	fmt.Println("##### PAIRWISE ASSOCIATION RULE ####")

	var (
		args          = models.ParseArgsFromCommand()
		timer         time.Time
		trainedModels map[string]models.Item
		filteredRules map[string]models.Item
		listRF        map[string]float64
	)

	fmt.Println("---Read data file---")
	data := readCsv(args.DataFile)
	for _, row := range data {
		fmt.Println(row)
	}

	fmt.Println("---START: Training step---")
	timer = time.Now()
	trainedModels = training(data)
	fmt.Println("---FINISH--- Time(microseconds): ", time.Now().Sub(timer).Microseconds())
	for k, v := range trainedModels {
		fmt.Printf(" %s -> %v", k, v)
		fmt.Println()
	}

	fmt.Println("---START: Filter step---")
	timer = time.Now()
	filteredRules = filter(trainedModels, args.InputItems)
	fmt.Println("---FINISH--- Time(microseconds): ", time.Now().Sub(timer).Microseconds())
	for k, v := range filteredRules {
		fmt.Printf(" %s -> %v", k, v)
		fmt.Println()
	}

	fmt.Println("---Recommend step---")
	timer = time.Now()
	listRF = recommend(filteredRules)
	fmt.Println("---FINISH--- Time(microseconds): ", time.Now().Sub(timer).Microseconds())
	fmt.Println("!!!!!RECOMMENDATION RESULT!!!!!")
	fmt.Println(listRF)
}

func training(data [][]string) map[string]models.Item {
	var (
		trainedModels = make(map[string]models.Item)
		pairsCount    = make(map[string]models.PairItems)
		genKeyID      int32
	)

	//iterate over raw data to calculate od 
	for _, transaction := range data {
		var visitedList []string

		//Use 2 dimension slices to construct pairs of items.
		for _, x := range transaction {
			var (
				itemInfo  models.Item
				currentOD int32
			)
			//Assign value to temp variables if x is exist in trainedModels
			if _, ok := trainedModels[x]; ok {
				itemInfo = trainedModels[x]
				currentOD = itemInfo.Occurrences
			}

			for _, y := range transaction {
				if y == x || contain(visitedList, y) {
					continue
				}
				genKeyID++
				key := fmt.Sprintf("key_%d", genKeyID)
				value := models.PairItems{
					Item1:         x,
					Item2:         y,
					CoOccurrences: 1,
				}
				pairsCount[key] = value
			}

			//Increase OD of X
			currentOD++
			itemInfo.Occurrences = currentOD
			trainedModels[x] = itemInfo
			//add to visited list
			visitedList = append(visitedList, x)
		}
	}

	var deletedKey []string
	//Iterate over pairsCount to calculate CD of each pair items
	//Remove duplicate pairs, ex: pair [a,b] is the same with pair [b,a]
	for key, info := range pairsCount {
		if contain(deletedKey, key) {
			continue
		}
		for subKey, value := range pairsCount {
			if (((info.Item1 == value.Item1) && (info.Item2 == value.Item2)) ||
				((info.Item1 == value.Item2) && (info.Item2 == value.Item1))) && (key != subKey) {
				info.CoOccurrences++
				delete(pairsCount, subKey)
				deletedKey = append(deletedKey, subKey)
			}
		}
		pairsCount[key] = info
	}

	//add cd pairs to trainedModels
	for _, pair := range pairsCount {
		var (
			tempAscItem  = make(map[string]int32)
			tempAscItem2 = make(map[string]int32)
		)
		temp := trainedModels[pair.Item1]
		if temp.AssociateItems != nil {
			tempAscItem = temp.AssociateItems
		}
		tempAscItem[pair.Item2] = pair.CoOccurrences
		temp.AssociateItems = tempAscItem
		trainedModels[pair.Item1] = temp

		temp2 := trainedModels[pair.Item2]
		if temp2.AssociateItems != nil {
			tempAscItem2 = temp2.AssociateItems
		}
		tempAscItem2[pair.Item1] = pair.CoOccurrences
		temp2.AssociateItems = tempAscItem2
		trainedModels[pair.Item2] = temp2
	}

	return trainedModels
}

func filter(trainedModels map[string]models.Item, inputItems []string) map[string]models.Item {
	var (
		filteredRules = make(map[string]models.Item)
		itemInfo      models.Item
	)
	for _, item := range inputItems {
		if _, ok := trainedModels[item]; ok {
			itemInfo = trainedModels[item]
			for name, _ := range itemInfo.AssociateItems {
				//filter out associateItems which are input in inputItems.
				if contain(inputItems, name) {
					delete(itemInfo.AssociateItems, name)
				}
			}

			//Set filtered associateItems
			filteredRules[item] = models.Item{
				Occurrences:    itemInfo.Occurrences,
				AssociateItems: itemInfo.AssociateItems,
			}

		} else {
			fmt.Println("Input item not in db yet !!!")
		}
	}
	return filteredRules
}

func recommend(filteredRules map[string]models.Item) map[string]float64 {
	var (
		result             = make(map[string]float64)
		weightTable        = make(map[string]int32)
		probabilitiesTable = make(map[string]float64)
	)
	for _, rule := range filteredRules {
		for itemName, coOccurrences := range rule.AssociateItems {
			var (
				prob   float64
				weight int32
			)
			//Calculate probability conditional p
			prob = float64(coOccurrences) / float64(rule.Occurrences)
			probabilitiesTable[itemName] = probabilitiesTable[itemName] + prob

			//Calculate weight W
			weight = rule.Occurrences
			weightTable[itemName] = weightTable[itemName] + weight

			//Set recommendItems to result map
			result[itemName] = 0
		}
	}
	for item, recommendScore := range result {
		recommendScore = probabilitiesTable[item] * float64(weightTable[item])
		result[item] = math.Round(recommendScore*10) / 10 //rounding to 1 decimal places
	}
	return result
}

func readCsv(path string) [][]string {
	//Open file
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	//New scanner
	scanner := bufio.NewScanner(f)

	//Add data to [][]string
	var (
		numTransactions int32
		rows            [][]string
	)
	for scanner.Scan() {
		numTransactions++
		row := strings.Split(scanner.Text(), ",")
		rows = append(rows, row)
	}
	return rows
}

//func to check if input item is in input array
func contain(arr []string, item string) bool {
	for _, containedItem := range arr {
		if item == containedItem {
			return true
		}
	}
	return false
}
