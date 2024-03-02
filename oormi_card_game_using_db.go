package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type CardValue struct {
	Suite string
	Name  string
}

func getNoOfPlayers() int {
	var noOfPlayers int = 0
	for {
		fmt.Print("Enter the number of players (4, 6, or 8): ")
		_, err := fmt.Scanf("%d", &noOfPlayers)
		if err != nil {
			fmt.Println("Input Error:", err)
		}
		if noOfPlayers == 4 || noOfPlayers == 6 || noOfPlayers == 8 {
			break
		} else {
			fmt.Println("# of Players cannot be :", noOfPlayers)
			fmt.Println("Enter the number of players (4, 6, or 8):")
		}
	}
	fmt.Println("No of Players in the Game :", noOfPlayers)
	return noOfPlayers
}

func getCardsPerPlayer(noOfPlayers int) int {
	var cardsPerPlayer int = 0
	for {
		fmt.Print("Enter the number of cards per player (between 5 and 8): ")
		_, err := fmt.Scanf("%d", &cardsPerPlayer)
		if err != nil {
			fmt.Println("Input Error:", err)
		}
		if cardsPerPlayer < 5 || cardsPerPlayer > 8 {
			fmt.Println("# of cards per player cannot be :", cardsPerPlayer)
			continue
		}
		if (cardsPerPlayer*noOfPlayers)%4 == 0 && (cardsPerPlayer*noOfPlayers) <= 52 {
			break
		} else {
			fmt.Println("# of cards per player cannot be :", cardsPerPlayer)
			fmt.Println("Condition 1 : Total number of cards in game must be dividable by 4.")
			fmt.Println("Condition 2 : Maximum number of cards in game is 52 (from two to ace).")
			fmt.Println("Consider : No of Players in the Game :", noOfPlayers)
			continue
		}
	}
	fmt.Println("No of Cards per player :", cardsPerPlayer)
	return cardsPerPlayer
}

func showCardsInGame(cardsPerSuite int, cardNames [13]string) {
	fmt.Println("Shuffle and deal " + strconv.Itoa(cardsPerSuite) + " cards from every suite -- CARD LIST --")
	for i := 0; i < cardsPerSuite; i++ {
		fmt.Println(cardNames[i])
	}
	fmt.Println("-- END OF CARD LIST --")
}

func printPlayerNameAndID() {
	var rows *sql.Rows
	var dbResponse bool
	rows, dbResponse = queryFromDB("SELECT PLAYERNAME,PLAYERID FROM PLAYERS", "SQL0001-SELECT PLAYERNAME,PLAYERID FROM PLAYERS", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	fmt.Println("[PlayerName] : [PlayerID] -- PLAYER LIST --")
	for rows.Next() {
		var playerName string
		var PlayerID int
		rows.Scan(&playerName, &PlayerID)
		fmt.Println(playerName + " : " + strconv.Itoa(PlayerID))
	}
	fmt.Println("-- END OF PLAYER LIST --")
}

func getTrumpsCaller(noOfPlayers int, ourTrumps *bool) int {
	var trumpsCaller int = -1
	for {
		printPlayerNameAndID()
		var confrim string
		fmt.Print("Who calls Trumps? (Player ID) : ")
		_, err := fmt.Scanf("%d", &trumpsCaller)

		if err != nil {
			fmt.Println("Input Error:", err)
			continue
		}

		if trumpsCaller < 0 || trumpsCaller >= noOfPlayers {
			fmt.Println("Trump caller cannot be :", trumpsCaller)
			continue
		} else {
			fmt.Println("Trumps called by (PlayerID):    ", trumpsCaller)
			//			fmt.Println("Trumps called by (Player Name): ", playerList[trumpsCaller])
			if trumpsCaller%2 == 0 {
				*ourTrumps = true
			} else {
				*ourTrumps = false
			}
			fmt.Print("Confrim Trump Caller? (y/Y) : ")
			fmt.Scan(&confrim)
			if confrim == "y" || confrim == "Y" {
				fmt.Println("input Confrimed")
				break
			} else {
				fmt.Println("input discarded")
				continue
			}
		}
	}
	return trumpsCaller
}

func inputCardValue(cardNames []string, cardsPerSuite int, cardNo int) CardValue {
	var suite string
	var name string
	fmt.Println("Insert Card # : ", cardNo)
	for {
		fmt.Println("Insert Suite (hearts-h/H, diamonds-d/D, spades-s/S, clubs-c/C): ")
		fmt.Print("Insert Suite : ")
		fmt.Scanln(&suite)
		if suite == "h" || suite == "H" {
			suite = "hearts"
			break
		} else if suite == "d" || suite == "D" {
			suite = "diamonds"
			break
		} else if suite == "s" || suite == "S" {
			suite = "spades"
			break
		} else if suite == "c" || suite == "C" {
			suite = "clubs"
			break
		} else {
			fmt.Println("Wrong Suite Name. Must be of of : hearts-h/H, dianmonds-d/D, spades-s/S, clubs-c/C")
			continue
		}
	}
	for {
		fmt.Print("Available Cards : ")
		for i := 0; i < cardsPerSuite; i++ {
			fmt.Print(cardNames[i])
			fmt.Print(" ")
		}
		fmt.Println()
		fmt.Println("Insert Card (ace-a/A. king-k/K, queen-q/A, jack-j/J, ten-10, nine-9, eight-8, seven-7, six-6, five-5, four-4, three-3. two-2)")
		fmt.Print("Insert Card : ")
		fmt.Scanln(&name)
		if name == "a" || name == "A" {
			name = "ace"
			break
		} else if name == "k" || name == "K" {
			name = "king"
			break
		} else if name == "q" || name == "Q" {
			name = "queen"
			break
		} else if name == "j" || name == "J" {
			name = "jack"
			break
		} else if name == "10" {
			name = "ten"
			break
		} else if name == "9" {
			name = "nine"
			break
		} else if name == "8" {
			name = "eight"
			break
		} else if name == "7" {
			name = "seven"
			break
		} else if name == "6" {
			name = "six"
			break
		} else if name == "5" {
			name = "five"
			break
		} else if name == "4" {
			name = "four"
			break
		} else if name == "3" {
			name = "three"
			break
		} else if name == "2" {
			name = "two"
			break
		} else {
			fmt.Println("Wrong Card Name.!")
			continue
		}
	}

	return CardValue{
		Suite: suite,
		Name:  name,
	}
}

func addCardToHandAndFindMyTrump(cardNames [13]string, cardsPerSuite int, cardsPerPlayer int) string {
	var res sql.Result
	var rows *sql.Rows
	var dbResponse bool
	var cardNo, trumpCheck, heartsPoints, diamondsPoints, spadesPoints, clubsPoints int
	cardNo = 0
	trumpCheck = 0
	heartsPoints = 0
	diamondsPoints = 0
	spadesPoints = 0
	clubsPoints = 0
	for {
		fmt.Println("# of cards in the first deal ? (4 or 3) : ")
		_, err := fmt.Scanf("%d", &trumpCheck)
		if err != nil {
			fmt.Println("Input Error:", err)
			continue
		}
		if trumpCheck == 3 || trumpCheck == 4 {
			break
		} else {
			continue
		}
	}
	fmt.Println("Add cards from 1st deal to decide Trump Suite : ")
	for { // while cards need to be entered
		if cardNo >= trumpCheck {
			break
		}
		for { // for each card
			inputCard := inputCardValue(cardNames[:], cardsPerSuite, cardNo+1)
			rows, dbResponse = queryFromDB("SELECT POINTS FROM CARDS WHERE INPLAY = true AND INMYHAND = false AND CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0002-SELECT POINTS FROM CARDS WHERE INPLAY", true)
			if rows == nil || !dbResponse {
				os.Exit(1)
			}
			var cardPoints int = 0
			for rows.Next() {
				rows.Scan(&cardPoints)
			}
			if cardPoints != 0 {
				if inputCard.Suite == "hearts" {
					heartsPoints = heartsPoints + cardPoints
				} else if inputCard.Suite == "diamonds" {
					diamondsPoints = diamondsPoints + cardPoints
				} else if inputCard.Suite == "spades" {
					spadesPoints = spadesPoints + cardPoints
				} else {
					clubsPoints = clubsPoints + cardPoints
				}
				fmt.Println("**** Card Added to your hand:", inputCard)
				//showCardsInHand(cards)
				res, dbResponse = executeOnDB("UPDATE CARDS SET INMYHAND = true, PLCONF1 = 0,  PLCONF2 = 0,  PLCONF3 = 0,  PLCONF4 = 0,  PLCONF5 = 0,  PLCONF6 = 0,  PLCONF7 = 0 WHERE CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0003-UPDATE CARDS SET INMYHAND = true ...", true)
				if res == nil || !dbResponse {
					os.Exit(1)
				}
				cardNo++
				break
			} else {
				fmt.Println("Cannot add this card:", inputCard)
				fmt.Println("It might be already added or not used in the game")
				//showCardsInHand(cards)
				continue
			}
		}
	}

	var suitePoints int = 0

	fmt.Println("Points - hearts   : ", heartsPoints)
	fmt.Println("Points - diamonds : ", diamondsPoints)
	fmt.Println("Points - spades   : ", spadesPoints)
	fmt.Println("Points - clubs    : ", clubsPoints)

	suitePoints = heartsPoints
	var trumpSuite string = "hearts"
	if diamondsPoints > suitePoints {
		trumpSuite = "diamonds"
		suitePoints = diamondsPoints
	}
	if spadesPoints > suitePoints {
		trumpSuite = "spades"
		suitePoints = spadesPoints
	}
	if clubsPoints > suitePoints {
		trumpSuite = "clubs"
	}

	fmt.Println("Announce to players ******* Trumps is : ", trumpSuite)
	fmt.Println("Add cards from 2nd deal : ")

	for { // while cards need to be entered
		if cardNo >= cardsPerPlayer {
			break
		}
		for { // for each card
			inputCard := inputCardValue(cardNames[:], cardsPerSuite, cardNo+1)
			rows, dbResponse = queryFromDB("SELECT POINTS FROM CARDS WHERE INPLAY = true AND INMYHAND = false AND CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0004-SELECT POINTS FROM CARDS WHERE INPLAY", true)
			if rows == nil || !dbResponse {
				os.Exit(1)
			}
			var cardPoints int = 0
			for rows.Next() {
				rows.Scan(&cardPoints)
			}
			if cardPoints != 0 {
				fmt.Println("**** Card Added to your hand:", inputCard)
				//showCardsInHand(cards)
				res, dbResponse = executeOnDB("UPDATE CARDS SET INMYHAND = true, PLCONF1 = 0,  PLCONF2 = 0,  PLCONF3 = 0,  PLCONF4 = 0,  PLCONF5 = 0,  PLCONF6 = 0,  PLCONF7 = 0 WHERE CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0005-UPDATE CARDS SET INMYHAND = true ...", true)
				if res == nil || !dbResponse {
					os.Exit(1)
				}
				cardNo++
				break
			} else {
				fmt.Println("Cannot add this card:", inputCard)
				fmt.Println("It might be already added or not used in the game")
				//showCardsInHand(cards)
				continue
			}
		}
	}
	return trumpSuite
}

func addCardsToHandAndGetTrump(cardNames [13]string, cardsPerSuite int, cardsPerPlayer int) string {
	var res sql.Result
	var rows *sql.Rows
	var dbResponse bool
	var cardNo int = 0
	var trumpSuite string
	fmt.Println("Add cards You got (1st and 2nd deals) : ")
	for { // while cards need to be entered
		if cardNo >= cardsPerPlayer {
			break
		}
		for { // for each card
			inputCard := inputCardValue(cardNames[:], cardsPerSuite, cardNo+1)
			rows, dbResponse = queryFromDB("SELECT POINTS FROM CARDS WHERE INPLAY = true AND INMYHAND = false AND CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0006-SELECT POINTS FROM CARDS WHERE INPLAY", true)
			if rows == nil || !dbResponse {
				os.Exit(1)
			}
			var cardPoints int = 0
			for rows.Next() {
				rows.Scan(&cardPoints)
			}
			if cardPoints != 0 {
				fmt.Println("**** Card Added to your hand:", inputCard)
				//showCardsInHand(cards)
				res, dbResponse = executeOnDB("UPDATE CARDS SET INMYHAND = true, PLCONF1 = 0,  PLCONF2 = 0,  PLCONF3 = 0,  PLCONF4 = 0,  PLCONF5 = 0,  PLCONF6 = 0,  PLCONF7 = 0 WHERE CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0007-UPDATE CARDS SET INMYHAND = true ...", true)
				if res == nil || !dbResponse {
					os.Exit(1)
				}
				cardNo++
				break
			} else {
				fmt.Println("Cannot add this card:", inputCard)
				fmt.Println("It might be already added or not used in the game")
				//showCardsInHand(cards)
				continue
			}
		}
	}
	for { //while Trump is given
		var confirm string
		fmt.Println("Input Trump Suite: hearts-h/H, dianmonds-d/D, spades-s/S, clubs-c/C")
		fmt.Print("Input Trump Suite: ")
		fmt.Scanln(&trumpSuite)
		if trumpSuite == "h" || trumpSuite == "H" {
			trumpSuite = "hearts"
		} else if trumpSuite == "d" || trumpSuite == "D" {
			trumpSuite = "diamonds"
		} else if trumpSuite == "s" || trumpSuite == "S" {
			trumpSuite = "spades"
		} else if trumpSuite == "c" || trumpSuite == "C" {
			trumpSuite = "clubs"
		} else {
			fmt.Println("Trump Suite has to be one of : hearts-h/H, dianmonds-d/D, spades-s/S, clubs-c/C")
			continue
		}
		fmt.Println("Confirm Trump is : ", trumpSuite)
		fmt.Print("Confirm ? (y/Y) : ")
		fmt.Scanln(&confirm)
		if confirm == "y" || confirm == "Y" {
			fmt.Println("**** Trump Suite : ", trumpSuite)
			break
		} else {
			continue
		}
	}
	return trumpSuite
}

func executeOnDB(sqlString string, indetifier string, exitOnErr bool) (sql.Result, bool) {
	var res sql.Result
	statement, err := db.Prepare(sqlString)
	if err != nil {
		fmt.Println("PREPARE: " + indetifier)
		log.Println(err)
		if exitOnErr {
			os.Exit(1)
		}
		return res, false
	}
	res, err = statement.Exec()
	if err != nil {
		fmt.Println("EXEC: " + indetifier)
		log.Println(err)
		if exitOnErr {
			os.Exit(1)
		}
		return res, false
	}
	return res, true
}

func queryFromDB(sqlString string, indetifier string, exitOnErr bool) (*sql.Rows, bool) {
	var rows *sql.Rows
	statement, err := db.Prepare(sqlString)
	if err != nil {
		fmt.Println("PREPARE: " + indetifier)
		log.Println(err)
		if exitOnErr {
			os.Exit(1)
		}
		return rows, false
	}
	rows, err = statement.Query()
	if err != nil {
		fmt.Println("QUERY: " + indetifier)
		log.Println(err)
		if exitOnErr {
			os.Exit(1)
		}
		return rows, false
	}
	return rows, true
}

func createCardsTab(cardSuites [4]string, cardSuitesAb [4]string, cardNames [13]string, cardNamesAb [13]string, cardPoints [13]int, cardConfidenceAtStart int) {
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("DROP TABLE IF EXISTS CARDS", "SQL0008-DROP TABLE IF EXISTS CARDS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	res, dbResponse = executeOnDB("CREATE TABLE CARDS (CARDINDEX VARCHAR(2), CARDSUITE VARCHAR(8), CARDNAME VARCHAR(8), INPLAY BOOLEAN, INMYHAND BOOLEAN, POINTS INT, ROUNDPOINTS INT, PLCONF1 INT, PLCONF2 INT, PLCONF3 INT, PLCONF4 INT, PLCONF5 INT, PLCONF6 INT, PLCONF7 INT)", "SQL0009-CREATE TABLE CARDS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	for s := 0; s < len(cardSuites); s++ {
		for n := 0; n < len(cardNames); n++ {
			res, dbResponse = executeOnDB("INSERT INTO CARDS (CARDINDEX, CARDSUITE, CARDNAME, INPLAY, INMYHAND, POINTS, ROUNDPOINTS, PLCONF1, PLCONF2, PLCONF3, PLCONF4, PLCONF5, PLCONF6, PLCONF7) VALUES ('"+cardSuitesAb[s]+cardNamesAb[n]+"','"+cardSuites[s]+"','"+cardNames[n]+"', true, false,'"+strconv.Itoa(cardPoints[n])+"','"+strconv.Itoa(cardPoints[n])+"', "+strconv.Itoa(cardConfidenceAtStart)+", "+strconv.Itoa(cardConfidenceAtStart)+", "+strconv.Itoa(cardConfidenceAtStart)+", "+strconv.Itoa(cardConfidenceAtStart)+", "+strconv.Itoa(cardConfidenceAtStart)+", "+strconv.Itoa(cardConfidenceAtStart)+", "+strconv.Itoa(cardConfidenceAtStart)+")", "SQL0010-INSERT INTO CARDS", true)
			if res == nil || !dbResponse {
				os.Exit(1)
			}
		}
	}
}

func createPlayersTab() {
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("DROP TABLE IF EXISTS PLAYERS", "DROP TABLE IF EXISTS PLAYERS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	res, dbResponse = executeOnDB("CREATE TABLE PLAYERS (PLAYERID INT,PLAYERNAME VARCHAR(32),FRIEND BOOLEAN,PLAYEDINROUND BOOLEAN,HEARTSPROB INT,SPADESPROB INT,DIAMONDSPROB INT,CLUBSPROB INT)", "SQL0011-CREATE TABLE PLAYERS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func createRoundsTab() {
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("DROP TABLE IF EXISTS ROUNDS", "DROP TABLE IF EXISTS ROUNDS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	res, dbResponse = executeOnDB("CREATE TABLE ROUNDS (ROUND INT,ROUNDTURN INT,PLAYERID INT,FRIEND BOOLEAN,PLAYERNAME VARCHAR(32),CARDINDEX VARCHAR(2),CARDSUITE VARCHAR(8),CARDNAME VARCHAR(8),ROUNDPOINTS INT,WINNER BOOLEAN,ROUNDSUITE VARCHAR(8),TRUMPSUITE VARCHAR(8),MYCARDPLAYCONDITION INT,NOTE VARCHAR(64))", "SQL0012-CREATE TABLE ROUNDS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func initPlayersTab(noOfPlayers int) {
	var res sql.Result
	var dbResponse bool
	for p := 0; p < (noOfPlayers); p++ {
		// get player name in future :-)
		var friend bool
		if p%2 == 0 {
			friend = true
		} else {
			friend = false
		}
		res, dbResponse = executeOnDB("INSERT INTO PLAYERS (PLAYERID,PLAYERNAME,FRIEND,PLAYEDINROUND,HEARTSPROB,SPADESPROB,DIAMONDSPROB,CLUBSPROB) VALUES('"+strconv.Itoa(p)+"','PLAYER"+strconv.Itoa(p)+"','"+strconv.FormatBool(friend)+"',false,1,1,1,1)", "SQL0013-INSERT INTO PLAYERS", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func initCardsTab(cardsPerSuite int, cardPoints [13]int) {
	var selectedCards int = cardsPerSuite - 1 // since card point array is 1 less than # of cards
	var minCardPoints int = cardPoints[selectedCards]
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("DELETE FROM CARDS WHERE ROUNDPOINTS < "+strconv.Itoa(minCardPoints), "SQL0014-DELETE FROM CARDS WHERE ROUNDPOINTS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func initRoundsTab(cardsPerPlayer int, noOfPlayers int) {
	var res sql.Result
	var dbResponse bool
	for r := 0; r < cardsPerPlayer; r++ {
		for p := 0; p < noOfPlayers; p++ {
			res, dbResponse = executeOnDB("INSERT INTO ROUNDS (ROUND,ROUNDTURN,WINNER) VALUES('"+strconv.Itoa(r)+"','"+strconv.Itoa(p)+"',false)", "SQL0015-INSERT INTO ROUNDS", true)
			if res == nil || !dbResponse {
				os.Exit(1)
			}
		}
	}
}

func updateCardsTabRoundPointsForTrumpSuite(trumpSuite string, pointsAddForTrumpSuite int) {
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("UPDATE CARDS SET ROUNDPOINTS = (POINTS + "+strconv.Itoa(pointsAddForTrumpSuite)+") WHERE CARDSUITE = '"+trumpSuite+"'", "SQL0016-UPDATE CARDS SET ROUNDPOINTS... FOR TRUMPSUITE", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func printCardsPlayedInRound(roundNo int, roundWinnerCard CardValue) {
	var rows *sql.Rows
	var dbResponse bool
	var playerID, playerName, friend, cardSuite, cardName, roundPoints, winner string
	rows, dbResponse = queryFromDB("SELECT PLAYERID,PLAYERNAME,FRIEND,CARDSUITE,CARDNAME,ROUNDPOINTS,WINNER FROM ROUNDS WHERE ROUND = "+strconv.Itoa(roundNo), "SQL0017-SELECT PLAYERID,PLAYERNAME,FRIEND,CARDSUITE,CARNAME FROM ROUNDS", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	fmt.Println("Cards Played in Round : " + strconv.Itoa(roundNo) + " --- ")
	fmt.Println("--PLAYERID PLAYERNAME FRIEND CARDSUITE CARDNAME ROUNDPOINTS WINNER")
	for rows.Next() {
		rows.Scan(&playerID, &playerName, &friend, &cardSuite, &cardName, &roundPoints, &winner)
		fmt.Println(playerID + " " + playerName + " " + friend + " " + cardSuite + " " + cardName + " " + roundPoints + " " + winner)
		playerID = "------"
		playerName = "------"
		friend = "------"
		cardSuite = "------"
		cardName = "------"
		roundPoints = "------"
		roundPoints = "------"
		winner = "------"
	}
	fmt.Println("--- Round Winner (So Far) : ", roundWinnerCard)
	fmt.Println("--- END OF Cards Played in Round --- ")
}

func getCurrentPlayerName(currentPlayerID int) string {
	var rows *sql.Rows
	var dbResponse bool
	var playerName string = "nul-player"
	rows, dbResponse = queryFromDB("SELECT PLAYERNAME FROM PLAYERS WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0018-SELECT PLAYERNAME ROM PLAYERS WHERE", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&playerName)
	}
	return playerName
}

func updateRoundsTabTrumpSuiteColumn(trumpSuite string) {
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("UPDATE ROUNDS SET TRUMPSUITE = '"+trumpSuite+"'", "SQL0019-UPDATE ROUNDS SET TRUMPSUITE", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updateRoundsTabWithPlayerIDs(roundNo int, currentPlayerID int, noOfPlayers int) {
	var res sql.Result
	var dbResponse bool
	for p := 0; p < noOfPlayers; p++ {
		var playerIDForRound int = (currentPlayerID + p) % noOfPlayers
		res, dbResponse = executeOnDB("UPDATE ROUNDS SET PLAYERID = "+strconv.Itoa(playerIDForRound)+" WHERE ROUND = "+strconv.Itoa(roundNo)+" AND ROUNDTURN = "+strconv.Itoa(p), "SQL0020-UPDATE ROUNDS SET PLAYERID ... ", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func getPlayCard(cardNames []string, cardsPerSuite int) CardValue {
	var suite string
	var name string
	for {
		fmt.Println("Insert Suite (hearts-h/H, diamonds-d/D, spades-s/S, clubs-c/C): ")
		fmt.Print("Insert Suite : ")
		fmt.Scanln(&suite)
		if suite == "h" || suite == "H" {
			suite = "hearts"
			break
		} else if suite == "d" || suite == "D" {
			suite = "diamonds"
			break
		} else if suite == "s" || suite == "S" {
			suite = "spades"
			break
		} else if suite == "c" || suite == "C" {
			suite = "clubs"
			break
		} else {
			fmt.Println("Wrong Suite Name. Must be of of : hearts-h/H, dianmonds-d/D, spades-s/S, clubs-c/C")
			continue
		}
	}
	for {
		fmt.Print("Available Cards : ")
		for i := 0; i < cardsPerSuite; i++ {
			fmt.Print(cardNames[i])
			fmt.Print(" ")
		}
		fmt.Println()
		fmt.Println("Insert Card (ace-a/A. king-k/K, queen-q/A, jack-j/J, ten-10, nine-9, eight-8, seven-7, six-6, five-5, four-4, three-3. two-2)")
		fmt.Print("Insert Card : ")
		fmt.Scanln(&name)
		if name == "a" || name == "A" {
			name = "ace"
			break
		} else if name == "k" || name == "K" {
			name = "king"
			break
		} else if name == "q" || name == "Q" {
			name = "queen"
			break
		} else if name == "j" || name == "J" {
			name = "jack"
			break
		} else if name == "10" {
			name = "ten"
			break
		} else if name == "9" {
			name = "nine"
			break
		} else if name == "8" {
			name = "eight"
			break
		} else if name == "7" {
			name = "seven"
			break
		} else if name == "6" {
			name = "six"
			break
		} else if name == "5" {
			name = "five"
			break
		} else if name == "4" {
			name = "four"
			break
		} else if name == "3" {
			name = "three"
			break
		} else if name == "2" {
			name = "two"
			break
		} else {
			fmt.Println("Wrong Card Name.!")
			continue
		}
	}
	return CardValue{
		Suite: suite,
		Name:  name,
	}
}

func checkIfCurrentCardIsValid(currentCard CardValue) bool {
	var rows *sql.Rows
	var dbResponse bool
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE = '"+currentCard.Suite+"' AND CARDNAME ='"+currentCard.Name+"' AND INPLAY = true and INMYHAND = false", "SQL0021-SELECT PLAYERNAME ROM PLAYERS WHERE", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	var rowCount int = 0
	for rows.Next() {
		rows.Scan(&rowCount)
	}
	if rowCount == 0 {
		return false
	} else {
		return true
	}
}

func checkIfSurrentCardSuiteIsLegalForPlayer(currentPlayerID int, currentCardSuite string) bool {
	var rows *sql.Rows
	var dbResponse bool = false
	if currentCardSuite == "hearts" {
		rows, dbResponse = queryFromDB("SELECT HEARTSPROB FROM PLAYERS WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0022-SELECT HEARTSPROB FROM PLAYERS WHERE ", true)
	} else if currentCardSuite == "spades" {
		rows, dbResponse = queryFromDB("SELECT SPADESPROB FROM PLAYERS WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0023-SELECT SPADESPROB FROM PLAYERS WHERE ", true)
	} else if currentCardSuite == "diamonds" {
		rows, dbResponse = queryFromDB("SELECT DIAMONDSPROB FROM PLAYERS WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0024-SELECT DIAMONDSPROB FROM PLAYERS WHERE ", true)
	} else if currentCardSuite == "clubs" {
		rows, dbResponse = queryFromDB("SELECT CLUBSPROB FROM PLAYERS WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0025-SELECT CLUBSPROB FROM PLAYERS WHERE ", true)
	}
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	var cardProb int = 0
	for rows.Next() {
		rows.Scan(&cardProb)
	}
	if cardProb < 1 {
		if cardProb < 0 {
			cardProb = -(cardProb)
		}
		fmt.Println("Game lost for Player : ", currentPlayerID, ". Should not have this suite : ", currentCardSuite, ". This Player Didn't use the Suite in Round : ", strconv.Itoa(cardProb))
		return false
	} else {
		return true
	}
}

func updateCardsTabRoundPointsForRounduite(pointsAddForRoundSuite int, roundSuite string) {
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE CARDS SET ROUNDPOINTS = (POINTS + "+strconv.Itoa(pointsAddForRoundSuite)+") WHERE CARDSUITE = '"+roundSuite+"'", "SQL0026-UPDATE CARDS SET ROUNDPOINTS... for ROUNDSUITE", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updateRoundsTabWithRoundSuite(roundSuite string, roundNo int) {
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE ROUNDS SET ROUNDSUITE = '"+roundSuite+"' WHERE ROUND = "+strconv.Itoa(roundNo), "SQL0027-UPDATE ROUNDS SET ROUNDSUITE", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updatePlayersTabProbForRoundSuite(roundSuite string, roundNo int, currentPlayerID int) { // Players probability set to minus value for round
	var res sql.Result
	var dbResponse bool = false
	if roundNo > 0 {
		roundNo = -(roundNo)
	}
	if roundSuite == "hearts" {
		res, dbResponse = executeOnDB("UPDATE PLAYERS SET HEARTSPROB = "+strconv.Itoa(roundNo)+" WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0028-UPDATE PLAYERS SET HEARTSPROB", true)
	} else if roundSuite == "spades" {
		res, dbResponse = executeOnDB("UPDATE PLAYERS SET SPADESPROB = "+strconv.Itoa(roundNo)+" WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0029-UPDATE PLAYERS SET SPADESPROB", true)
	} else if roundSuite == "diamonds" {
		res, dbResponse = executeOnDB("UPDATE PLAYERS SET DIAMONDSPROB = "+strconv.Itoa(roundNo)+" WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0030-UPDATE PLAYERS SET DIAMONDSPROB", true)
	} else if roundSuite == "clubs" {
		res, dbResponse = executeOnDB("UPDATE PLAYERS SET CLUBSPROB = "+strconv.Itoa(roundNo)+" WHERE PLAYERID = "+strconv.Itoa(currentPlayerID), "SQL0031-UPDATE PLAYERS SET CLUBSPROB", true)
	}
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updateCardsTabForRoundSuiteProb(roundSuite string, currentPlayerID int) { // Cards probabaility make 0 for roundsite cards for this player
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE CARDS SET PLCONF"+strconv.Itoa(currentPlayerID)+" = 0 WHERE CARDSUITE = '"+roundSuite+"'", "SQL0032-UPDATE CARDS SET PLCONF", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updateCardsTabForPlayedCard(currentPlayerID int, currentCard CardValue) (string, int) {
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE CARDS SET INPLAY=false,INMYHAND=false,PLCONF1=0,PLCONF2=0,PLCONF3=0,PLCONF4=0,PLCONF5= 0,PLCONF6=0,PLCONF7=0 WHERE CARDSUITE='"+currentCard.Suite+"' AND CARDNAME='"+currentCard.Name+"'", "SQL0033-UPDATE CARDS SET INPLAY=false", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT CARDINDEX,ROUNDPOINTS FROM CARDS WHERE CARDSUITE = '"+currentCard.Suite+"' AND CARDNAME ='"+currentCard.Name+"'", "SQL0034-SELECT ROUNDPOINTS FROM CARDS", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	var currentCardPoints int = 0
	var currentCardIndex string = "nul"
	for rows.Next() {
		rows.Scan(&currentCardIndex, &currentCardPoints)
	}
	return currentCardIndex, currentCardPoints
}

func updateRoundsTabForPlayedCard(currentPlayerID int, currentCard CardValue, currentCardPoints int, roundNo int, playerInRound int, currentPlayerTeam string, currentPlayerName string, currentCardIndex string) (bool, int, string, CardValue, int) {
	var res sql.Result
	var dbResponse bool = false
	var friend string
	var winnerCard CardValue
	if currentPlayerTeam == "friend" {
		friend = "true"
	} else {
		friend = "false"
	}
	res, dbResponse = executeOnDB("UPDATE ROUNDS SET PLAYERID="+strconv.Itoa(currentPlayerID)+",FRIEND="+friend+",PLAYERNAME='"+currentPlayerName+"',CARDINDEX='"+currentCardIndex+"',CARDSUITE='"+currentCard.Suite+"',CARDNAME='"+currentCard.Name+"',ROUNDPOINTS="+strconv.Itoa(currentCardPoints)+" WHERE ROUND="+strconv.Itoa(roundNo)+" AND ROUNDTURN="+strconv.Itoa(playerInRound), "SQL0035-UPDATE ROUNDS SET PLAYERID", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT PLAYERID,PLAYERNAME,MAX(ROUNDPOINTS),CARDNAME,CARDSUITE FROM ROUNDS WHERE ROUND = "+strconv.Itoa(roundNo), "SQL0037-SELECT PLAYERID FROM ROUNDS WHERE ROUNDPOINTS ", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	var currentRoundWinnerID int = -1
	var currentRoundWinnerName string = "nul-player"
	var maxRoundPoints int = 0
	for rows.Next() {
		rows.Scan(&currentRoundWinnerID, &currentRoundWinnerName, &maxRoundPoints, &winnerCard.Name, &winnerCard.Suite)
	}
	if maxRoundPoints == currentCardPoints {
		return true, currentRoundWinnerID, currentRoundWinnerName, winnerCard, maxRoundPoints
	} else {
		return false, currentRoundWinnerID, currentRoundWinnerName, winnerCard, maxRoundPoints
	}
}

func updateRoundsTabForMyPlayedCard(currentPlayerID int, currentCard CardValue, currentCardPoints int, roundNo int, playerInRound int, currentPlayerTeam string, currentPlayerName string, currentCardIndex string, myCardPlayCondition int) (bool, int, string, CardValue, int) {
	var res sql.Result
	var dbResponse bool = false
	var friend string
	var winnerCard CardValue
	if currentPlayerTeam == "friend" {
		friend = "true"
	} else {
		friend = "false"
	}
	res, dbResponse = executeOnDB("UPDATE ROUNDS SET PLAYERID="+strconv.Itoa(currentPlayerID)+",FRIEND="+friend+",PLAYERNAME='"+currentPlayerName+"',CARDINDEX='"+currentCardIndex+"',CARDSUITE='"+currentCard.Suite+"',CARDNAME='"+currentCard.Name+"',ROUNDPOINTS="+strconv.Itoa(currentCardPoints)+",MYCARDPLAYCONDITION='"+strconv.Itoa(myCardPlayCondition)+"' WHERE ROUND="+strconv.Itoa(roundNo)+" AND ROUNDTURN="+strconv.Itoa(playerInRound), "SQL0035-UPDATE ROUNDS SET PLAYERID", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT PLAYERID,PLAYERNAME,MAX(ROUNDPOINTS),CARDNAME,CARDSUITE FROM ROUNDS WHERE ROUND = "+strconv.Itoa(roundNo), "SQL0037-SELECT PLAYERID FROM ROUNDS WHERE ROUNDPOINTS ", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	var currentRoundWinnerID int = -1
	var currentRoundWinnerName string = "nul-player"
	var maxRoundPoints int = 0
	for rows.Next() {
		rows.Scan(&currentRoundWinnerID, &currentRoundWinnerName, &maxRoundPoints, &winnerCard.Name, &winnerCard.Suite)
	}
	if maxRoundPoints == currentCardPoints {
		return true, currentRoundWinnerID, currentRoundWinnerName, winnerCard, maxRoundPoints
	} else {
		return false, currentRoundWinnerID, currentRoundWinnerName, winnerCard, maxRoundPoints
	}
}

func updateRoundsTabWithWinner(currentRoundWinnerID int, roundNo int) {
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE ROUNDS SET WINNER="+strconv.Itoa(currentRoundWinnerID)+" WHERE ROUND="+strconv.Itoa(roundNo), "SQL0038-UPDATE ROUNDS SET PLAYERID", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updateCardsTabForTrumpCaller(trumpsCaller int, trumpSuite string, cardConfidenceAddForTrumpcCaller [13]int, cardNames [13]string) {
	var res sql.Result
	var dbResponse bool = false
	for c := 0; c < len(cardNames); c++ {
		res, dbResponse = executeOnDB("UPDATE CARDS SET PLCONF"+strconv.Itoa(trumpsCaller)+"=PLCONF"+strconv.Itoa(trumpsCaller)+"+"+strconv.Itoa(cardConfidenceAddForTrumpcCaller[c])+" WHERE INMYHAND=false AND CARDNAME='"+cardNames[c]+"' AND CARDSUITE='"+trumpSuite+"'", "SQL0039-UPDATE CARDS SET PLCONF", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func getARandomCardInMyHand() CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME FROM CARDS WHERE INPLAY=true and INMYHAND=true ORDER BY RANDOM() LIMIT 1", "SQL0039-SELECT CARDSUITE,CARDNAME FROM CARDS", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name)
	}
	return currentCard
}

func getMaxRoundCardInMyHand(roundSuite string) CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+roundSuite+"' ORDER BY ROUNDPOINTS DESC LIMIT 1", "SQL0040-SELECT CARDSUITE,CARDNAME ", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name)
	}
	return currentCard
}

func getMinRoundCardInMyHand(roundSuite string) CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+roundSuite+"' ORDER BY ROUNDPOINTS ASC LIMIT 1", "SQL0041-SELECT CARDSUITE,CARDNAME ", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name)
	}
	return currentCard
}

func getMinRoundCardInMyHandAbovePoints(roundSuite string, currentRoundWinnerPoints int) CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+roundSuite+"' AND ROUNDPOINTS>"+strconv.Itoa(currentRoundWinnerPoints)+" ORDER BY ROUNDPOINTS ASC LIMIT 1", "SQL0042-SELECT CARDSUITE,CARDNAME ", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name)
	}
	return currentCard
}

func getMaxRoundCardPointsInMyHand(roundSuite string) int {
	var cardPoints int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT ROUNDPOINTS FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+roundSuite+"' ORDER BY ROUNDPOINTS DESC LIMIT 1", "SQL0043-SELECT ROUNDPOINTS FROM CARDS", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return 0
	}
	for rows.Next() {
		rows.Scan(&cardPoints)
	}
	return cardPoints
}

func getMaxTrumpCardInMyHand(trumpSuite string) CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+trumpSuite+"' ORDER BY ROUNDPOINTS DESC LIMIT 1", "SQL0044-SELECT CARDSUITE,CARDNAME FROM CARDS", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name)
	}
	return currentCard
}

func getANonTrumpAceInMyHand(trumpSuite string) CardValue {
	var currentCard, tempCard CardValue
	var suiteRounds, minSuiteRounds int = 0, 99
	var dbResponse, dbResponse1 bool = false, false
	var rows, rows1 *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE<>'"+trumpSuite+"' AND CARDNAME='ace'", "SQL00410-SELECT CARDSUITE,CARDNAME ", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&tempCard.Suite, &tempCard.Name)
		rows1, dbResponse1 = queryFromDB("SELECT COUNT(DISTINCT ROUNDSUITE || ' ' || ROUND) from ROUNDS WHERE ROUNDSUITE='"+tempCard.Suite+"'", "SQL00420-SELECT COUNT(DISTINCT ROUNDSUITE ", true)
		if !dbResponse1 {
			os.Exit(1)
		}
		for rows1.Next() {
			rows1.Scan(&suiteRounds)
		}
		if suiteRounds <= minSuiteRounds {
			currentCard.Suite = tempCard.Suite
			currentCard.Name = tempCard.Name
		}
	}
	return currentCard
}

func getMaxTrumpIfIHaveMaxTrumpAndMoreThanHalfOfTrumpCardsInPlay(trumpSuite string) CardValue {
	var currentCard, nulCard CardValue
	var dbResponse, iHaveMaxTrump bool = false, false
	var rows *sql.Rows
	var garbageInt, trumpsInOtherHands, trumpsInMyHand int = 0, 0, 0
	nulCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT INMYHAND,MAX(POINTS),CARDSUITE,CARDNAME from CARDS WHERE INPLAY=TRUE AND CARDSUITE = '"+trumpSuite+"'", "SQL00430-SELECT INMYHAND,MAX(POINTS) from", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&iHaveMaxTrump, &garbageInt, &currentCard.Suite, &currentCard.Name)
	}
	if !iHaveMaxTrump {
		return nulCard
	} else {
		rows, dbResponse = queryFromDB("SELECT COUNT(CARDINDEX) from CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE AND CARDSUITE='"+trumpSuite+"'", "SQL00440-SELECT COUNT(CARDINDEX) from CARDS", true)
		if rows == nil || !dbResponse {
			os.Exit(1)
		}
		for rows.Next() {
			rows.Scan(&trumpsInMyHand)
		}
		rows, dbResponse = queryFromDB("SELECT COUNT(CARDINDEX) from CARDS WHERE INPLAY=TRUE AND INMYHAND=FALSE AND CARDSUITE='"+trumpSuite+"'", "SQL00450-SELECT COUNT(CARDINDEX) from CARDS", true)
		if rows == nil || !dbResponse {
			os.Exit(1)
		}
		for rows.Next() {
			rows.Scan(&trumpsInOtherHands)
		}
		if trumpsInMyHand < trumpsInOtherHands {
			return nulCard
		} else {
			return currentCard
		}
	}
}

func getMaxTrumpIfIHaveMaxTrumpAndNextRoundWinWinsTheGame(trumpSuite string, friendPoints int, cardsPerPlayer int) CardValue {
	var currentCard, nulCard CardValue
	var dbResponse, iHaveMaxTrump bool = false, false
	var rows *sql.Rows
	var garbageInt int = 0
	nulCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT INMYHAND,MAX(POINTS),CARDSUITE,CARDNAME from CARDS WHERE INPLAY=TRUE AND CARDSUITE = '"+trumpSuite+"'", "SQL00460-SELECT INMYHAND,MAX(POINTS) from", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&iHaveMaxTrump, &garbageInt, &currentCard.Suite, &currentCard.Name)
	}
	if !iHaveMaxTrump {
		return nulCard
	} else {
		if (friendPoints + 1) > (cardsPerPlayer / 2) {
			return currentCard
		} else {
			return nulCard
		}
	}
}

func checkFoeConfidenceOnTrumps(trumpSuite string, noOfPlayers int) int {
	var foeConfidenceOnTrumps int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	for p := 0; p < noOfPlayers; p++ {
		if p%2 == 1 {
			var trumpConfForPlayer int = 0
			rows, dbResponse = queryFromDB("SELECT SUM(PLCONF"+strconv.Itoa(p)+") FROM CARDS WHERE CARDSUITE='"+trumpSuite+"' AND INMYHAND=FALSE AND INPLAY=TRUE", "SQL00470-SELECT CARDSUITE,CARDNAME ", true)
			if !dbResponse || rows == nil {
				os.Exit(1)
			}
			for rows.Next() {
				rows.Scan(&trumpConfForPlayer)
			}
			foeConfidenceOnTrumps = foeConfidenceOnTrumps + trumpConfForPlayer
		}
	}
	return foeConfidenceOnTrumps
}

func checkFriendConfidenceOnTrumps(trumpSuite string, noOfPlayers int) int {
	var foeConfidenceOnTrumps int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	for p := 0; p < noOfPlayers; p++ {
		if p%2 == 0 {
			var trumpConfForPlayer int = 0
			rows, dbResponse = queryFromDB("SELECT SUM(PLCONF"+strconv.Itoa(p)+") FROM CARDS WHERE CARDSUITE='"+trumpSuite+"' AND INMYHAND=FALSE AND INPLAY=TRUE", "SQL00470-SELECT CARDSUITE,CARDNAME ", true)
			if !dbResponse || rows == nil {
				os.Exit(1)
			}
			for rows.Next() {
				rows.Scan(&trumpConfForPlayer)
			}
			foeConfidenceOnTrumps = foeConfidenceOnTrumps + trumpConfForPlayer
		}
	}
	return foeConfidenceOnTrumps
}

func getMaxFromANonTrumpSuiteIfIHaveMaxOfThatSuiteAndFoesHaveNoTrumps(trumpSuite string, noOfPlayers int, cardSuites [4]string) CardValue {
	var currentCard, nulCard CardValue
	var dbResponse = false
	var rows *sql.Rows
	var garbageInt int = 0
	nulCard.Suite = "nul"
	var foeConfidenceOnTrumps int = checkFoeConfidenceOnTrumps(trumpSuite, noOfPlayers)
	if foeConfidenceOnTrumps > 0 {
		return nulCard
	} else {
		for s := 0; s < len(cardSuites); s++ {
			if cardSuites[s] == trumpSuite {
				continue
			}
			rows, dbResponse = queryFromDB("SELECT INMYHAND,MAX(POINTS),CARDSUITE,CARDNAME FROM CARDS WHERE CARDSUITE='"+cardSuites[s]+"' AND INPLAY=TRUE", "SQL00480-SELECT INMYHAND,MAX(POINTS)", true)
			if !dbResponse {
				os.Exit(1)
			}
			if rows == nil {
				continue
			} else {
				var inMyHand bool
				for rows.Next() {
					rows.Scan(&inMyHand, &garbageInt, &currentCard.Suite, &currentCard.Name)
				}
				if inMyHand {
					return currentCard
				}
			}
		}
		return nulCard
	}
}

func getMinCardIfIHaveOneMinCardFromANonTrumpSuite(trumpSuite string, cardSuites [4]string, cardPoints [13]int, cardsPerSuite int) CardValue {
	var maxCutOffCardPointsForLeastValueCards int = cardPoints[cardsPerSuite/2] - 1
	var currentCard, nulCard, tempCard CardValue
	var dbResponse = false
	var rows *sql.Rows
	var suiteCount, minSuitePoints int = 0, cardPoints[0] // minSuitePoints = ace points
	nulCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,COUNT(DISTINCT CARDNAME),MIN(POINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE GROUP BY CARDSUITE;", "SQL00490-SELECT INMYHAND,MAX(POINTS)", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return nulCard
	} else if countSuiteCardsInMyHand(trumpSuite) == 0 {
		return nulCard
	} else {
		var minCardFound bool = false
		for rows.Next() {
			var minPointForThisSuite int
			rows.Scan(&tempCard.Suite, &tempCard.Name, &suiteCount, &minPointForThisSuite)
			if tempCard.Suite != trumpSuite && suiteCount == 1 && minPointForThisSuite < minSuitePoints && minPointForThisSuite < maxCutOffCardPointsForLeastValueCards {
				currentCard.Suite = tempCard.Suite
				currentCard.Name = tempCard.Name
				minSuitePoints = minPointForThisSuite
				minCardFound = true
			}
		}
		if minCardFound {
			return currentCard
		} else {
			return nulCard
		}
	}
}

func getMyMinPointsCard() CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var garbageInt int
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MIN(POINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE", "SQL00500-SELECT CARDSUITE,CARDNAME,MIN(POINTS)", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &garbageInt)
	}
	return currentCard
}

func countSuiteCardsInMyHand(cardCuite string) int {
	var suiteCount int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+cardCuite+"' AND INPLAY=TRUE AND INMYHAND=TRUE", "SQL00510-SELECT COUNT(CARDNAMES)", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&suiteCount)
	}
	return suiteCount
}

func getARoundCardInMyHandToCloseRound(roundSuite string, currentRoundWinnerPoints int, currentRoundWinnerID int) CardValue {
	var nulCard CardValue
	nulCard.Suite = "nul"
	if countSuiteCardsInMyHand(roundSuite) > 0 { // have to return a card
		if currentRoundWinnerID%2 == 0 { // friend is the winner
			return getMinRoundCardInMyHand(roundSuite)
		} else { // foe is the winner
			if getMaxRoundCardPointsInMyHand(roundSuite) > currentRoundWinnerPoints {
				return getMinRoundCardInMyHandAbovePoints(roundSuite, currentRoundWinnerPoints)
			} else {
				return getMinRoundCardInMyHand(roundSuite)
			}
		}
	} else { // no suite cards
		return nulCard
	}
}

func getATrumpCardInMyHandToCloseRound(trumpSuite string, currentRoundWinnerPoints int, currentRoundWinnerID int) CardValue {
	var nulCard CardValue
	nulCard.Suite = "nul"
	if currentRoundWinnerID%2 == 0 { // friend is the winner
		return nulCard
	} else {
		if getMaxRoundCardPointsInMyHand(trumpSuite) > currentRoundWinnerPoints {
			return getMinRoundCardInMyHandAbovePoints(trumpSuite, currentRoundWinnerPoints)
		} else {
			return nulCard
		}
	}
}

func resetPlayersTabPlayedInRound() {
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE PLAYERS SET PLAYEDINROUND=false", "SQL00511-UPDATE PLAYERS SET PLAYEDINROUND=false", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func updatePlayersTabPlayedInRound(currentPlayerID int) {
	var res sql.Result
	var dbResponse bool = false
	res, dbResponse = executeOnDB("UPDATE PLAYERS SET PLAYEDINROUND=true WHERE PLAYERID="+strconv.Itoa(currentPlayerID), "SQL00512-UPDATE PLAYERS SET PLAYEDINROUND=false", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func main() {

	cardNames := [13]string{"ace", "king", "queen", "jack", "ten", "nine", "eight", "seven", "six", "five", "four", "three", "two"}
	cardNamesAb := [13]string{"a", "k", "q", "j", "t", "9", "8", "7", "6", "5", "4", "3", "2"}
	cardPoints := [13]int{378, 234, 145, 90, 56, 35, 22, 14, 9, 6, 4, 3, 2}
	cardSuites := [4]string{"hearts", "spades", "diamonds", "clubs"}
	cardSuitesAb := [4]string{"h", "s", "d", "c"}
	var pointsAddForRoundSuite int = 1000
	var pointsAddForTrumpSuite int = 2000
	var cardConfidenceAtStart int = 128
	cardConfidenceAddForTrumpcCaller := [13]int{32, 24, 18, 12, 10, 8, 6, 4, 2, 1, 1, 1, 1}

	var err error
	db, err = sql.Open("sqlite3", "./oormi_card_game.db") // disk - for debugging
	// db, err = sql.Open("sqlite3", ":memory:") // memory - faster
	if err != nil {
		fmt.Println("SQLERR - func main : Opening DB File")
		log.Println(err)
		os.Exit(1)
	}

	createCardsTab(cardSuites, cardSuitesAb, cardNames, cardNamesAb, cardPoints, cardConfidenceAtStart)
	createPlayersTab()
	createRoundsTab()

	var noOfPlayers int = getNoOfPlayers()
	var cardsPerPlayer int = getCardsPerPlayer(noOfPlayers)
	var cardsPerSuite int = (cardsPerPlayer * noOfPlayers) / 4
	//	var totalCardsInGame int = cardsPerPlayer * noOfPlayers
	initPlayersTab(noOfPlayers)                // insert player info plater tab
	initCardsTab(cardsPerSuite, cardPoints)    // remove unused cards
	initRoundsTab(cardsPerPlayer, noOfPlayers) // init records of rounds
	showCardsInGame(cardsPerSuite, cardNames)
	var ourTrumps bool
	var trumpsCaller int = getTrumpsCaller(noOfPlayers, &ourTrumps)
	var trumpSuite string

	if trumpsCaller == 0 {
		trumpSuite = addCardToHandAndFindMyTrump(cardNames, cardsPerSuite, cardsPerPlayer)
	} else {
		trumpSuite = addCardsToHandAndGetTrump(cardNames, cardsPerSuite, cardsPerPlayer)
		updateCardsTabForTrumpCaller(trumpsCaller, trumpSuite, cardConfidenceAddForTrumpcCaller, cardNames) // increase confidence for trumpscaller on trumps suite
	}
	updateCardsTabRoundPointsForTrumpSuite(trumpSuite, pointsAddForTrumpSuite) // update cards table with incrasing round points to trump suite cards
	updateRoundsTabTrumpSuiteColumn(trumpSuite)                                // update rounds table with trumpsuite

	var friendPoints int = 0
	var foePoints int = 0
	var gameResult int = 0
	var previousRoundWinner int = trumpsCaller
	for roundNo := 0; roundNo < cardsPerPlayer; roundNo++ { // for every round
		var roundSuite string = "notknown"
		var currentRoundWinnerID int = -1
		var currentRoundWinnerName string = "nul-player"
		var currentRoundWinnerCard CardValue
		var currentRoundWinnerPoints int = 0
		resetPlayersTabPlayedInRound()
		for playerInRound := 0; playerInRound < noOfPlayers; playerInRound++ { // for every player
			var playedRoundSuite bool = true     // if the palyer played round suite
			var playedNonRoundTrump bool = false // if the player played trump when he played a non round card
			var currentPlayerID int = (previousRoundWinner + playerInRound) % noOfPlayers
			var currentPlayerTeam string
			if currentPlayerID%2 == 0 {
				currentPlayerTeam = "friend"
			} else {
				currentPlayerTeam = "foe"
			}
			var currentPlayerName string = getCurrentPlayerName(currentPlayerID)
			var currentCard CardValue
			if playerInRound == 0 {
				updateRoundsTabWithPlayerIDs(roundNo, currentPlayerID, noOfPlayers)
			}

			if currentPlayerID != 0 { // other player play
				fmt.Print("Round: ", roundNo, " Insert Card for Player - ", currentPlayerID, " ", currentPlayerName, " ", currentPlayerTeam, " : ")
				for { // until the user gives a card available for play
					currentCard = getPlayCard(cardNames[:], cardsPerSuite)
					if checkIfCurrentCardIsValid(currentCard) { // check if current card is in the game (and not used already)
						break
					} else {
						fmt.Println("Invalids card : Player can't play this card (already played/not in game/in your hand) : ", currentCard)
					}
				}
				if checkIfSurrentCardSuiteIsLegalForPlayer(currentPlayerID, currentCard.Suite) { // check if player should have that card (used another card before in s round asking for this suite)
				} else {
					if currentPlayerTeam == "friend" {
						gameResult = -1
						break
					} else {
						gameResult = 1
						break
					}
				}
				var currentCardPoints int = 0
				var currentCardIndex string = "nul"
				var currentCardIsRoundWinner bool = false
				if playerInRound == 0 {
					roundSuite = currentCard.Suite
					if roundSuite != trumpSuite {
						updateCardsTabRoundPointsForRounduite(pointsAddForRoundSuite, roundSuite)
					}
					updateRoundsTabWithRoundSuite(roundSuite, roundNo)
					updatePlayersTabPlayedInRound(currentPlayerID)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex)
					//////////////////////////////////////// NEED MORE CRITERIA ///////////////////////
				} else {
					if roundSuite != currentCard.Suite {
						updatePlayersTabProbForRoundSuite(roundSuite, roundNo, currentPlayerID) // Players probability set to minus value for round
						updateCardsTabForRoundSuiteProb(roundSuite, currentPlayerID)            // Cards probabaility make 0 for roundsite cards for this player
						playedRoundSuite = false
						if currentCard.Suite == trumpSuite {
							playedNonRoundTrump = true
						}
					}
					updatePlayersTabPlayedInRound(currentPlayerID)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex)
					//////////////////////////////////////// NEED MORE CRITERIA ///////////////////////
				}

				fmt.Println(playedRoundSuite, playedNonRoundTrump, currentCardIsRoundWinner)

			} else { // my card play
				var currentCardPoints int = 0
				var currentCardIndex string = "nul"
				var currentCardIsRoundWinner bool = false
				var myCardPlayCondition int = 0
				if playerInRound == 0 { // i'm the round starter
					currentCard.Suite = "nul"
					for currentCard.Suite == "nul" {
						// currentCard = getARandomCardInMyHand() ///////////////////// needs to be changed for better selection ///////////////////////////////
						currentCard = getANonTrumpAceInMyHand(trumpSuite) // 10
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 10
							break
						}
						currentCard = getMaxTrumpIfIHaveMaxTrumpAndMoreThanHalfOfTrumpCardsInPlay(trumpSuite) // 20
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 20
							break
						}
						currentCard = getMaxTrumpIfIHaveMaxTrumpAndNextRoundWinWinsTheGame(trumpSuite, friendPoints, cardsPerPlayer) // 30
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 30
							break
						}
						currentCard = getMaxFromANonTrumpSuiteIfIHaveMaxOfThatSuiteAndFoesHaveNoTrumps(trumpSuite, noOfPlayers, cardSuites) // 40
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 40
							break
						}
						currentCard = getMinCardIfIHaveOneMinCardFromANonTrumpSuite(trumpSuite, cardSuites, cardPoints, cardsPerSuite) // 50
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 50
							break
						}
						/////////////////// more methods ?? ////////////////////////////////////////////////////////
						currentCard = getMyMinPointsCard() // 60
						myCardPlayCondition = 60
						break
					}
					roundSuite = currentCard.Suite
					if roundSuite != trumpSuite {
						updateCardsTabRoundPointsForRounduite(pointsAddForRoundSuite, roundSuite)
					}
					updateRoundsTabWithRoundSuite(roundSuite, roundNo)
					updatePlayersTabPlayedInRound(currentPlayerID)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForMyPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex, myCardPlayCondition)

				} else if playerInRound == (noOfPlayers - 1) { // i'm the last player in round
					currentCard.Suite = "nul"
					for currentCard.Suite == "nul" {
						//currentCard = getMaxRoundCardInMyHand(roundSuite) ///////////////////////// get my Max RoundCard /////////////////////////
						currentCard = getARoundCardInMyHandToCloseRound(roundSuite, currentRoundWinnerPoints, currentRoundWinnerID)
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 70
							break
						}
						//currentCard = getMaxTrumpCardInMyHand(trumpSuite) //////////////////////////////get my Max trump Card ////////////////////
						currentCard = getATrumpCardInMyHandToCloseRound(trumpSuite, currentRoundWinnerPoints, currentRoundWinnerID)
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 80
							break
						}
						//currentCard = getARandomCardInMyHand() ////////////////////////getMaxNonRoundNonTrumpCardInMyHand //////////////////
						currentCard = getMinCardIfIHaveOneMinCardFromANonTrumpSuite(trumpSuite, cardSuites, cardPoints, cardsPerSuite)
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 90
							break
						}
						currentCard = getMyMinPointsCard()
						myCardPlayCondition = 100
						break
					}
					updatePlayersTabPlayedInRound(currentPlayerID)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForMyPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex, myCardPlayCondition)
				} else { // i'm not the round starter or last player
					/////////// NOT OBVIOUS CHOISES //////////////////////
					currentCard.Suite = "nul"
					for currentCard.Suite == "nul" {
						currentCard = getMaxRoundCardInMyHand(roundSuite) ///////////////////////// get my Max RoundCard /////////////////////////
						if currentCard.Suite != "nul" {
							break
						}
						currentCard = getMaxTrumpCardInMyHand(trumpSuite) //////////////////////////////get my Max trump Card ////////////////////
						if currentCard.Suite != "nul" {
							break
						}
						currentCard = getARandomCardInMyHand() ////////////////////////getMaxNonRoundNonTrumpCardInMyHand //////////////////
						break
					}
					updatePlayersTabPlayedInRound(currentPlayerID)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForMyPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex, myCardPlayCondition)
				}
				fmt.Println("You Played : ", currentCard)
				fmt.Println(currentCardIsRoundWinner) // REMOVE
			}
			printCardsPlayedInRound(roundNo, currentRoundWinnerCard)
			if gameResult != 0 {
				break
			}
		} // player truns in a round end
		if gameResult != 0 {
			break
		}
		previousRoundWinner = currentRoundWinnerID
		updateRoundsTabWithWinner(currentRoundWinnerID, roundNo)
		if currentRoundWinnerID%2 == 0 {
			if roundNo < cardsPerPlayer-1 {
				fmt.Println("**** YOUR TEAM WON THE ROUND. PlayerID : ", previousRoundWinner, " (", currentRoundWinnerName, ") STARTS THE NEXT ROUND ***")
			}
			friendPoints++
		} else {
			if roundNo < cardsPerPlayer-1 {
				fmt.Println("**** YOUR TEAM LOST THE ROUND. PlayerID : ", previousRoundWinner, "(", currentRoundWinnerName, ") STARTS THE NEXT ROUND ***")
			}
			foePoints++
		}
	} // rounds in a game end
	if gameResult < 0 {
		fmt.Println("**** Your Team LOST THE GAME ******************")
	} else if gameResult > 0 {
		fmt.Println("**** Your Team WON THE GAME ******************")
	} else {
		fmt.Println("Your Team won ", friendPoints, " rounds")
		fmt.Println("Opposing Team won ", foePoints, " rounds")
		if friendPoints > foePoints {
			fmt.Println("**** Your Team WON THE GAME ******************")
		} else if friendPoints < foePoints {
			fmt.Println("**** Your Team LOST THE GAME ******************")
		} else {
			fmt.Println("******** GAME DRAW ******************")
		}
	}
}