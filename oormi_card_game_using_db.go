package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbType string = "sqlite3"
var dbPath string = "./oormi_card_game.db"

//var dbPath string = ":memory:"

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
		} else if suite == "d" || suite == "D" {
			suite = "diamonds"
		} else if suite == "s" || suite == "S" {
			suite = "spades"
		} else if suite == "c" || suite == "C" {
			suite = "clubs"
		} else {
			fmt.Println("Wrong Suite Name. Must be of of : hearts-h/H, dianmonds-d/D, spades-s/S, clubs-c/C")
			continue
		}
		fmt.Print("Available Cards : ")
		for i := 0; i < cardsPerSuite; i++ {
			fmt.Print(cardNames[i])
			fmt.Print(" ")
		}
		fmt.Println()
		fmt.Println("Insert Card (ace-a/A. king-k/K, queen-q/A, jack-j/J, ten-t/T, nine-9, eight-8, seven-7, six-6, five-5, four-4, three-3. two-2)")
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
		} else if name == "10" || name == "t" || name == "T" {
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
		fmt.Print("# of cards in the first deal ? (4 or 3) : ")
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
				res, dbResponse = executeOnDB("UPDATE CARDS SET INMYHAND = true, PLCONF1 = 0,  PLCONF2 = 0,  PLCONF3 = 0,  PLCONF4 = 0,  PLCONF5 = 0,  PLCONF6 = 0,  PLCONF7 = 0 WHERE CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0003-UPDATE CARDS SET INMYHAND = true ...", true)
				if res == nil || !dbResponse {
					os.Exit(1)
				}
				cardNo++
				printMyHand()
				break
			} else {
				fmt.Println("!!!!! Cannot add this card:", inputCard)
				fmt.Println("It might be already added or not used in the game")
				printMyHand()
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
				res, dbResponse = executeOnDB("UPDATE CARDS SET INMYHAND = true, PLCONF1 = 0,  PLCONF2 = 0,  PLCONF3 = 0,  PLCONF4 = 0,  PLCONF5 = 0,  PLCONF6 = 0,  PLCONF7 = 0 WHERE CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0005-UPDATE CARDS SET INMYHAND = true ...", true)
				if res == nil || !dbResponse {
					os.Exit(1)
				}
				cardNo++
				printMyHand()
				break
			} else {
				fmt.Println("Cannot add this card:", inputCard)
				fmt.Println("It might be already added or not used in the game")
				printMyHand()
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
				res, dbResponse = executeOnDB("UPDATE CARDS SET INMYHAND = true, PLCONF1 = 0,  PLCONF2 = 0,  PLCONF3 = 0,  PLCONF4 = 0,  PLCONF5 = 0,  PLCONF6 = 0,  PLCONF7 = 0 WHERE CARDSUITE = '"+inputCard.Suite+"' AND CARDNAME = '"+inputCard.Name+"'", "SQL0007-UPDATE CARDS SET INMYHAND = true ...", true)
				if res == nil || !dbResponse {
					os.Exit(1)
				}
				cardNo++
				printMyHand()
				break
			} else {
				fmt.Println("Cannot add this card:", inputCard)
				fmt.Println("It might be already added or not used in the game")
				printMyHand()
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
	var err error
	var dbTryLimit int = 10
	var returnError bool = true
	var res sql.Result
	time.Sleep(50 * time.Millisecond)
	res, err = db.Exec(sqlString)
	for t := 0; ; t++ {
		if err == nil {
			break
		}
		if err != nil && t >= dbTryLimit {
			fmt.Println("EXEC: " + indetifier)
			log.Println(err)
			if exitOnErr {
				os.Exit(1)
			}
			returnError = false
		}
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Waiting for DB.... (executeOnDB:EXEC) ", t)
		res, err = db.Exec(sqlString)
	}
	time.Sleep(500 * time.Millisecond)

	return res, returnError
}

func queryFromDB(sqlString string, indetifier string, exitOnErr bool) (*sql.Rows, bool) {
	var err error
	var dbTryLimit int = 10
	var returnError bool = true
	var rows *sql.Rows

	time.Sleep(50 * time.Millisecond)
	rows, err = db.Query(sqlString)
	for t := 0; ; t++ {
		if err == nil {
			break
		}
		if err != nil && t >= dbTryLimit {
			fmt.Println("QUERY: " + indetifier)
			log.Println(err)
			if exitOnErr {
				os.Exit(1)
			}
			returnError = false
		}
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Waiting for DB.... (queryFromDB:QUERY) ", t)
		rows, err = db.Query(sqlString)
	}

	return rows, returnError
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
	res, dbResponse = executeOnDB("CREATE TABLE PLAYERS (PLAYERID INT,PLAYERNAME VARCHAR(32),FRIEND BOOLEAN,HEARTSPROB INT,SPADESPROB INT,DIAMONDSPROB INT,CLUBSPROB INT)", "SQL0011-CREATE TABLE PLAYERS", true)
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
	res, dbResponse = executeOnDB("CREATE TABLE ROUNDS (ROUND INT,ROUNDTURN INT,PLAYERID INT,FRIEND BOOLEAN,PLAYERNAME VARCHAR(32),CARDINDEX VARCHAR(2),CARDSUITE VARCHAR(8),CARDNAME VARCHAR(8),ROUNDPOINTS INT,WINNER INT,ROUNDSUITE VARCHAR(8),TRUMPSUITE VARCHAR(8),MYCARDPLAYCONDITION INT,NOTE VARCHAR(64))", "SQL0012-CREATE TABLE ROUNDS", true)
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
		res, dbResponse = executeOnDB("INSERT INTO PLAYERS (PLAYERID,PLAYERNAME,FRIEND,HEARTSPROB,SPADESPROB,DIAMONDSPROB,CLUBSPROB) VALUES('"+strconv.Itoa(p)+"','PLAYER"+strconv.Itoa(p)+"','"+strconv.FormatBool(friend)+"',1,1,1,1)", "SQL0013-INSERT INTO PLAYERS", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func initCardsTab(cardsPerSuite int, cardPoints [13]int, noOfPlayers int) {
	var selectedCards int = cardsPerSuite - 1 // since card point array is 1 less than # of cards
	var minCardPoints int = cardPoints[selectedCards]
	var res sql.Result
	var dbResponse bool
	res, dbResponse = executeOnDB("DELETE FROM CARDS WHERE ROUNDPOINTS < "+strconv.Itoa(minCardPoints), "SQL00140-DELETE FROM CARDS WHERE ROUNDPOINTS", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	if noOfPlayers == 4 {
		res, dbResponse = executeOnDB("UPDATE CARDS SET PLCONF4=0,PLCONF5=0,PLCONF6=0,PLCONF7=0", "SQL00141-UPDATE CARDS SET PLCONF4", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	} else if noOfPlayers == 6 {
		res, dbResponse = executeOnDB("UPDATE CARDS SET PLCONF6=0,PLCONF7=0", "SQL00142-UPDATE CARDS SET PLCONF6", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func initRoundsTab(cardsPerPlayer int, noOfPlayers int) {
	var res sql.Result
	var dbResponse bool
	for r := 0; r < cardsPerPlayer; r++ {
		for p := 0; p < noOfPlayers; p++ {
			res, dbResponse = executeOnDB("INSERT INTO ROUNDS (ROUND,ROUNDTURN,WINNER) VALUES('"+strconv.Itoa(r)+"','"+strconv.Itoa(p)+"',-1)", "SQL0015-INSERT INTO ROUNDS", true)
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

func printCardsPlayedInRound(roundNo int, roundWinnerCard CardValue, trumpSuite string, roundSuite string, friendPoints int, foePoints int, currentRoundWinnerID int) {
	var rows *sql.Rows
	var dbResponse bool
	var playerID, playerName, friend, cardSuite, cardName, roundPoints, winner string
	rows, dbResponse = queryFromDB("SELECT PLAYERID,PLAYERNAME,FRIEND,CARDSUITE,CARDNAME,ROUNDPOINTS FROM ROUNDS WHERE ROUND = "+strconv.Itoa(roundNo), "SQL0017-SELECT PLAYERID,PLAYERNAME,FRIEND,CARDSUITE,CARNAME FROM ROUNDS", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	fmt.Println("**** ROUND# : ", strconv.Itoa(roundNo), " INFO ******************************************")
	fmt.Println("Trump Suite --------- : " + trumpSuite)
	fmt.Println("Round Suite --------- : " + roundSuite)
	fmt.Println("Cards Played in Round : " + strconv.Itoa(roundNo) + " **** ")
	fmt.Println("ID\tNAME\tFRIEND\tSUITE\tCARD\tPTS\tWINNER?")
	for rows.Next() {
		rows.Scan(&playerID, &playerName, &friend, &cardSuite, &cardName, &roundPoints)
		if playerID == intToString(currentRoundWinnerID) {
			winner = " <-- CURRENT WINNER"
		} else {
			winner = ""
		}
		fmt.Println(playerID + "\t" + playerName + "\t" + friend + "\t" + cardSuite + "\t" + cardName + "\t" + roundPoints + "\t" + winner)
		playerID = "------"
		playerName = "------"
		friend = "------"
		cardSuite = "------"
		cardName = "------"
		roundPoints = "------"
		roundPoints = "------"
		winner = "------"
	}
	fmt.Println("--- Round Winner (so Far) -------------- : ", roundWinnerCard)
	fmt.Println("--- Rounds won by your team (so far) --- : ", friendPoints)
	fmt.Println("--- Rounds won by opposing team (so far) : ", foePoints)
	fmt.Println("*** END OF ROUND INFO *****************************************************************")
	fmt.Println()
	printMyHand()
	fmt.Println()
}

func printMyHand() {
	var rows *sql.Rows
	var dbResponse bool
	var myCard string
	rows, dbResponse = queryFromDB("SELECT CARDSUITE||'-'||CARDNAME FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE", "Cqy85VTxM4APfZ2lweDRZ0qwPr2cruUz", true)
	if rows == nil || !dbResponse {
		os.Exit(1)
	}
	fmt.Println("*** CARDS IN MY HAND ******************************************************************")
	for rows.Next() {
		rows.Scan(&myCard)
		fmt.Print(myCard)
		fmt.Print("  ")
	}
	fmt.Println()
	fmt.Println("*** END OF CARDS IN MY HAND **********************************************************")
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
		} else if suite == "d" || suite == "D" {
			suite = "diamonds"
		} else if suite == "s" || suite == "S" {
			suite = "spades"
		} else if suite == "c" || suite == "C" {
			suite = "clubs"
		} else {
			fmt.Println("Wrong Suite Name. Must be of of : hearts-h/H, dianmonds-d/D, spades-s/S, clubs-c/C")
			continue
		}
		fmt.Print("Available Cards : ")
		for i := 0; i < cardsPerSuite; i++ {
			fmt.Print(cardNames[i])
			fmt.Print(" ")
		}
		fmt.Println()
		fmt.Println("Insert Card (ace-a/A. king-k/K, queen-q/A, jack-j/J, ten-t/T, nine-9, eight-8, seven-7, six-6, five-5, four-4, three-3. two-2)")
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
		} else if name == "t" || name == "T" || name == "10" {
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

func updateCardsTabForPlayedCard(currentPlayerID int, currentCard CardValue, noOfPlayers int) (string, int) {
	var res sql.Result
	var dbResponse bool = false
	if noOfPlayers == 4 {
		res, dbResponse = executeOnDB("UPDATE CARDS SET INPLAY=false,INMYHAND=false,PLCONF1=0,PLCONF2=0,PLCONF3=0 WHERE CARDSUITE='"+currentCard.Suite+"' AND CARDNAME='"+currentCard.Name+"'", "SQL0033-UPDATE CARDS SET INPLAY=false", true)
	} else if noOfPlayers == 6 {
		res, dbResponse = executeOnDB("UPDATE CARDS SET INPLAY=false,INMYHAND=false,PLCONF1=0,PLCONF2=0,PLCONF3=0,PLCONF4=0,PLCONF5=0 WHERE CARDSUITE='"+currentCard.Suite+"' AND CARDNAME='"+currentCard.Name+"'", "SQL0033-UPDATE CARDS SET INPLAY=false", true)
	} else {
		res, dbResponse = executeOnDB("UPDATE CARDS SET INPLAY=false,INMYHAND=false,PLCONF1=0,PLCONF2=0,PLCONF3=0,PLCONF4=0,PLCONF5=0,PLCONF6=0,PLCONF7=0 WHERE CARDSUITE='"+currentCard.Suite+"' AND CARDNAME='"+currentCard.Name+"'", "SQL0033-UPDATE CARDS SET INPLAY=false", true)
	}
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

func updateCardsTabForTrumpCaller(trumpsCaller int, trumpSuite string, cardConfidenceAddForTrumpcCaller [13]float32, cardNames [13]string) {
	var res sql.Result
	var dbResponse bool = false
	for c := 0; c < len(cardNames); c++ {
		res, dbResponse = executeOnDB("UPDATE CARDS SET PLCONF"+strconv.Itoa(trumpsCaller)+"=PLCONF"+strconv.Itoa(trumpsCaller)+"*"+float32ToString(cardConfidenceAddForTrumpcCaller[c])+" WHERE INMYHAND=false AND CARDNAME='"+cardNames[c]+"' AND CARDSUITE='"+trumpSuite+"'", "SQL0039-UPDATE CARDS SET PLCONF", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func intToString(value int) string {
	return strconv.Itoa(value)
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

func checkIfFoeHasMoreConfidenceOnHigherTrumps(foeID int, currentCard CardValue, roundPoints int, trumpSuite string) bool {
	var dbResponse bool = false
	var rows *sql.Rows
	var higherTrumpCards int = 0
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE PLCONF"+intToString(foeID)+">=PLCONF1 AND PLCONF"+intToString(foeID)+">=PLCONF2 AND PLCONF"+intToString(foeID)+">=PLCONF3 AND PLCONF"+intToString(foeID)+">=PLCONF4 AND PLCONF"+intToString(foeID)+">=PLCONF5 AND PLCONF"+intToString(foeID)+">=PLCONF6 AND PLCONF"+intToString(foeID)+">=PLCONF7 AND ROUNDPOINTS>"+intToString(roundPoints), "SQL-me3pn6UsQ61bbdjVJEuzRlTQwzLK2XCl", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return false
	}
	for rows.Next() {
		rows.Scan(&higherTrumpCards)
		if higherTrumpCards > 0 {
			return true
		}
	}
	return false
}

func checkIfFoesDoNotHaveHigerTrumpConfidence(playedInRound [8]int, currentCard CardValue, roundPoints int, trumpSuite string) bool {
	for foeID := 0; foeID < len(playedInRound); foeID++ {
		if foeID == 1 || foeID == 3 || foeID == 5 || foeID == 7 {
			if playedInRound[foeID] == 0 {
				if checkIfFoeHasMoreConfidenceOnHigherTrumps(foeID, currentCard, roundPoints, trumpSuite) {
					return false
				}
			}
		}
	}
	return true
}

func getMidTrumpCardInMyHandWhereFoesHaveLessPosibilityForAHigherTrumpCard(trumpSuite string, playedInRound [8]int) CardValue {
	var currentCard, returnCard, nulCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	var roundPoints int = 0
	nulCard.Suite = "nul"
	returnCard.Suite = "nul"
	if countSuiteCardsInMyHand(trumpSuite) == 0 {
		return nulCard
	}
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,ROUNDPOINTS FROM CARDS WHERE INMYHAND=true AND INPLAY=true AND CARDSUITE='"+trumpSuite+"' ORDER BY ROUNDPOINTS ASC", "SQL0040-SELECT CARDSUITE,CARDNAME ", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return nulCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &roundPoints)
		if checkIfFoesDoNotHaveHigerTrumpConfidence(playedInRound, currentCard, roundPoints, trumpSuite) {
			returnCard.Suite = currentCard.Suite
			returnCard.Name = currentCard.Name
		} else {
			continue
		}
	}
	if returnCard.Suite != nulCard.Suite {
		return returnCard
	} else {
		return nulCard
	}
}

func getMaxRoundCardIfIHaveItInMyHandIfItsAboveCurrentWinner(roundSuite string, currentRoundWinnerPoints int) CardValue {
	var currentCard, nulCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	var inMyHand bool = false
	var roundPoints int = 0
	nulCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT INMYHAND,CARDSUITE,CARDNAME,ROUNDPOINTS FROM CARDS WHERE INPLAY=true AND CARDSUITE='"+roundSuite+"' ORDER BY ROUNDPOINTS DESC LIMIT 1", "SQL0040-SELECT CARDSUITE,CARDNAME ", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return nulCard
	}
	for rows.Next() {
		rows.Scan(&inMyHand, &currentCard.Suite, &currentCard.Name, &roundPoints)
	}
	if inMyHand && roundPoints > currentRoundWinnerPoints {
		return currentCard
	}
	return nulCard
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
	var garbageInt int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MAX(ROUNDPOINTS) FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+trumpSuite+"'", "SQL-iNNWre84D2ilFem4vh06LPKtyCiMxdJS", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &garbageInt)
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

func getFullFriendFoeConfidenceForMaxInSuite(cardSuite string) (friendConf int, foeConf int) {
	var currentCard, maxSuiteCard, nulCard CardValue
	var dbResponse = false
	var rows *sql.Rows
	nulCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT (PLCONF1 + PLCONF3 + PLCONF5 + PLCONF7) AS TOTAL_CONF, CARDNAME, CARDSUITE FROM CARDS WHERE CARDSUITE = '"+cardSuite+"' AND INMYHAND = FALSE AND INPLAY = TRUE ORDER BY ROUNDPOINTS DESC LIMIT 1", "SQL-80634037554310202519884252290862", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return 0, 0
	}
	for rows.Next() {
		rows.Scan(&foeConf, &maxSuiteCard.Name, &maxSuiteCard.Suite)
	}
	rows, dbResponse = queryFromDB("SELECT (PLCONF2 + PLCONF4 + PLCONF6) AS TOTAL_CONF, CARDNAME, CARDSUITE FROM CARDS WHERE CARDSUITE = '"+cardSuite+"' AND INMYHAND = FALSE AND INPLAY = TRUE ORDER BY ROUNDPOINTS DESC LIMIT 1", "SQL-80634037554310202519884252290862", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return 0, 0
	}
	for rows.Next() {
		rows.Scan(&friendConf, &currentCard.Name, &currentCard.Suite)
	}
	if currentCard != maxSuiteCard {
		return 0, 0
	} else {
		return friendConf, foeConf
	}
}

func checkIfPlayerHasMaxConfidenceOnTrump(playerNotHavingSuite int, trumpSuite string) bool {
	var dbResponse = false
	var rows *sql.Rows
	var plConfTotal, plCheck, plConf1, plConf2, plConf3, plConf4, plConf5, plConf6, plConf7 int = 0, 0, 0, 0, 0, 0, 0, 0, 0
	rows, dbResponse = queryFromDB("SELECT (PLCONF1+PLCONF2+PLCONF3+PLCONF4+PLCONF5+PLCONF6+PLCONF7), PLCONF"+intToString(playerNotHavingSuite)+", PLCONF1, PLCONF2, PLCONF3, PLCONF4, PLCONF5, PLCONF6, PLCONF7 FROM CARDS WHERE INPLAY=TRUE AND CARDSUITE='"+trumpSuite+"'", "SQL-80634037554310202519884252290862", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&plConfTotal, &plCheck, &plConf1, &plConf2, &plConf3, &plConf4, &plConf5, &plConf6, &plConf7)
		if plConfTotal == 0 {
			continue
		}
		if playerNotHavingSuite%2 == 0 {
			if plCheck < plConf1 || plCheck < plConf3 || plCheck < plConf5 || plCheck < plConf7 {
				return false
			}
		} else {
			if plCheck < plConf2 || plCheck < plConf4 || plCheck < plConf6 {
				return false
			}
		}
	}
	return true
}

func checkIfAnyOfTheFoesDontHaveASuiteAndDontHaveTrumpConf(cardSuite string, trumpSuite string) bool {
	var playerNotHavingSuite int = 0
	var dbResponse = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT PLAYERID FROM PLAYERS WHERE "+cardSuite+"PROB<>1 AND PLAYERID IN (1, 3, 5, 7)", "SQL-80634037554310202519884252290862", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&playerNotHavingSuite)
		if !checkIfPlayerHasMaxConfidenceOnTrump(playerNotHavingSuite, trumpSuite) {
			return false
		}
	}
	return true
}

func checkIfAnyOfTheFriendsDontHaveASuiteAndDontHaveTrumpConf(cardSuite string, trumpSuite string) bool {
	var playerNotHavingSuite int = 0
	var dbResponse = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT PLAYERID FROM PLAYERS WHERE "+cardSuite+"PROB<>1 AND PLAYERID IN (2, 4, 6)", "SQL-80634037554310202519884252290862", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&playerNotHavingSuite)
		if !checkIfPlayerHasMaxConfidenceOnTrump(playerNotHavingSuite, trumpSuite) {
			return false
		}
	}
	return true
}

func getMinFromANonTrumpSuiteIfFriendsHaveMaxOfThatSuiteAndFoesHaveLessPossibilityOfTrumps(trumpSuite string, noOfPlayers int, cardSuites [4]string, cardPoints [13]int, cardsPerSuite int) CardValue {
	var currentCard, nulCard, returnCard CardValue
	var dbResponse = false
	var garbageInt int = 0
	var rows *sql.Rows
	nulCard.Suite = "nul"
	returnCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MIN(ROUNDPOINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE AND CARDSUITE<>'"+trumpSuite+"'GROUP BY CARDSUITE ORDER BY ROUNDPOINTS ASC", "SQL-35AQbgWw1OEDt1qxBZwZpdenWlyo5D4M", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return nulCard
	}
	var friendConf, foeConf, maxFriendMinusFoe int = 0, 0, 0
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &garbageInt)
		friendConf, foeConf = getFullFriendFoeConfidenceForMaxInSuite(currentCard.Suite)
		if friendConf-foeConf > maxFriendMinusFoe {
			if checkIfAnyOfTheFoesDontHaveASuiteAndDontHaveTrumpConf(currentCard.Suite, trumpSuite) {
				maxFriendMinusFoe = friendConf - foeConf
				returnCard.Suite = currentCard.Suite
				returnCard.Name = currentCard.Name
			}
		}
	}
	if returnCard.Suite != nulCard.Suite {
		return returnCard
	} else {
		return nulCard
	}
}

func getMinFromANonTrumpSuiteIfFriendsHaveMorePossibilityOfMaxOfThatSuiteAndFoesHaveNoTrumps(trumpSuite string, noOfPlayers int, cardSuites [4]string, cardPoints [13]int, cardsPerSuite int) CardValue {
	var currentCard, nulCard, returnCard CardValue
	var dbResponse = false
	var garbageInt int = 0
	var rows *sql.Rows
	nulCard.Suite = "nul"
	returnCard.Suite = "nul"
	var foeConfidenceOnTrumps int = checkFoeConfidenceOnTrumps(trumpSuite, noOfPlayers)
	if foeConfidenceOnTrumps > 0 {
		return nulCard
	}
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MIN(ROUNDPOINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE AND CARDSUITE<>'"+trumpSuite+"'GROUP BY CARDSUITE ORDER BY ROUNDPOINTS ASC", "SQL-35AQbgWw1OEDt1qxBZwZpdenWlyo5D4M", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return nulCard
	}
	var friendConf, foeConf, maxFriendMinusFoe int = 0, 0, 0
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &garbageInt)
		friendConf, foeConf = getFullFriendFoeConfidenceForMaxInSuite(currentCard.Suite)
		if friendConf-foeConf > maxFriendMinusFoe {
			maxFriendMinusFoe = friendConf - foeConf
			returnCard.Suite = currentCard.Suite
			returnCard.Name = currentCard.Name
		}
	}
	if returnCard.Suite != nulCard.Suite {
		return returnCard
	} else {
		return nulCard
	}
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
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MIN(ROUNDPOINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE", "SQL00500-SELECT CARDSUITE,CARDNAME,MIN(POINTS)", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &garbageInt)
	}
	return currentCard
}

func getMinCardInMyHandFromSuite(cardSuite string) CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var garbageInt int
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MIN(POINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE AND CARDSUITE='"+cardSuite+"'", "SQL-3ZbrQRURHG3IFMFxXKxBnoElFh2gTUdg", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &garbageInt)
	}
	return currentCard
}

func getMaxCardInMyHandFromSuite(cardSuite string) CardValue {
	var currentCard CardValue
	var dbResponse bool = false
	var garbageInt int
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,MAX(POINTS) FROM CARDS WHERE INPLAY=TRUE AND INMYHAND=TRUE AND CARDSUITE='"+cardSuite+"'", "SQL-3ZbrQRURHG3IFMFxXKxBnoElFh2gTUdg", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return currentCard
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

func getMaxSuiteCardInMyHand(cardSuite string) int {
	var currentCard CardValue
	var cardPoints int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,POINTS FROM CARDS WHERE INPLAY=true and INMYHAND=true AND CARDSUITE='"+cardSuite+"'", "SQL-iNNWre84D2ilFem4vh06LPKtyCiMxdJS", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return cardPoints
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &cardPoints)
	}
	return cardPoints
}

func getNoOfSuiteCardsAboveGivenCard(cardSuite string, cardPoints int) int {
	var currentCard CardValue
	var dbResponse bool = false
	var rows *sql.Rows
	var noOfCards int = 0
	currentCard.Suite = "nul"
	rows, dbResponse = queryFromDB("SELECT CARDSUITE,CARDNAME,POINTS FROM CARDS WHERE INPLAY=true and INMYHAND=false AND CARDSUITE='"+cardSuite+"' AND POINTS>"+intToString(cardPoints)+" ORDER BY ROUNDPOINTS ASC", "SQL-iNNWre84D2ilFem4vh06LPKtyCiMxdJS", true)
	if !dbResponse {
		os.Exit(1)
	}
	if rows == nil {
		return noOfCards
	}
	for rows.Next() {
		rows.Scan(&currentCard.Suite, &currentCard.Name, &cardPoints)
		noOfCards++
	}
	return noOfCards
}

func findIfFoeHasMoreConfidenceInAnyCardOfThisSuiteThanFreinds(cardSuite string, foeID int, noOfPlayers int) int {
	var dbResponse bool = false
	var rows *sql.Rows
	var noOfSuiteCardFoesCouldHaveAtHand int = 0
	if noOfPlayers == 4 {
		rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+cardSuite+"' AND INPLAY=true and INMYHAND=false (PLCONF"+intToString(foeID)+">PLCONF2)", "SQL-WxLTxpcNmPm6MaQBF1sY7YG8pxHcWOWz", true)
	} else if noOfPlayers == 6 {
		rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+cardSuite+"' AND INPLAY=true and INMYHAND=false (PLCONF"+intToString(foeID)+">PLCONF2 OR PLCONF"+intToString(foeID)+">PLCONF4)", "SQL-WxLTxpcNmPm6MaQBF1sY7YG8pxHcWOWz", true)
	} else { // noOfPlayers == 8
		rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+cardSuite+"' AND INPLAY=true and INMYHAND=false (PLCONF"+intToString(foeID)+">PLCONF2 OR PLCONF"+intToString(foeID)+">PLCONF4 OR PLCONF"+intToString(foeID)+">PLCONF6)", "SQL-WxLTxpcNmPm6MaQBF1sY7YG8pxHcWOWz", true)
	}
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&noOfSuiteCardFoesCouldHaveAtHand)
	}
	return noOfSuiteCardFoesCouldHaveAtHand
}

func checkFoesForSuitesThatTheyDontHaveAtHand(cardSuites [4]string, noOfPlayers int) ([4]int, [4]int) {
	var cardsInEachSuiteThatFoesDontHaveAtHand [4]int
	var cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand [4]int
	for s := 0; s < len(cardSuites); s++ {
		var dbResponse bool = false
		var rows *sql.Rows
		cardsInEachSuiteThatFoesDontHaveAtHand[s] = 1
		cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand[s] = 0
		rows, dbResponse = queryFromDB("SELECT PLAYERID FROM PLAYERS WHERE "+cardSuites[s]+"PROB=1 AND PLAYERID%2=1", "SQL-cym4WPB08aLN5Duevwin42y30yPkZXzR", true)
		if !dbResponse {
			os.Exit(1)
		}
		if rows == nil {
			cardsInEachSuiteThatFoesDontHaveAtHand[s] = 0
			continue
		}
		var foeID int
		for rows.Next() {
			rows.Scan(&foeID)
			cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand[s] = cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand[s] + findIfFoeHasMoreConfidenceInAnyCardOfThisSuiteThanFreinds(cardSuites[s], foeID, noOfPlayers)
		}
	}
	return cardsInEachSuiteThatFoesDontHaveAtHand, cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand
}

func getNonTrumpCardIfIHaveSecondMaxTrumpAndCanConvinceFoesWhoHaveThatMaxTrumpOnlyToCutMyCard(trumpSuite string, noOfPlayers int, cardSuites [4]string, foePoints int, friendPoints int,
	cardsPerPlayer int) CardValue {
	var noOfTrumpsAboveMe, myMaxTrumpCardPoints int
	var currentCard, nulCard CardValue
	currentCard.Suite = "nul"
	nulCard.Suite = "nul"
	myMaxTrumpCardPoints = getMaxSuiteCardInMyHand(trumpSuite)
	if myMaxTrumpCardPoints == 0 {
		return currentCard
	}
	noOfTrumpsAboveMe = getNoOfSuiteCardsAboveGivenCard(trumpSuite, myMaxTrumpCardPoints)
	if noOfTrumpsAboveMe != 1 || foePoints+1 > (cardsPerPlayer/2) {
		return currentCard
	}
	var cardsInEachSuiteThatFoesDontHaveAtHand, cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand [4]int
	cardsInEachSuiteThatFoesDontHaveAtHand, cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand = checkFoesForSuitesThatTheyDontHaveAtHand(cardSuites, noOfPlayers)
	for s := 0; s < len(cardSuites); s++ {
		if cardsInEachSuiteThatFoesDontHaveAtHand[s] == 0 {
			currentCard = getMinCardInMyHandFromSuite(cardSuites[s])
			if currentCard.Suite != "nul" {
				return currentCard
			}
		}
		if cardsInEachSuiteThatFoeshaveLessProbToHaveAtHand[s] == 0 {
			currentCard = getMaxCardInMyHandFromSuite(cardSuites[s])
			if currentCard.Suite != "nul" {
				return currentCard
			}
		}
	}
	return nulCard
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

func resetPlayersTabPlayedInRound(noOfPlayers int, playedInRound [8]int) [8]int {
	for p := 0; p < noOfPlayers; p++ {
		playedInRound[p] = 0
	}
	return playedInRound
}

func updatePlayersTabPlayedInRound(currentPlayerID int, playedInRound [8]int) [8]int {
	playedInRound[currentPlayerID] = 1
	return playedInRound
}

func calcNewConfidence(oldConfidence int, factor float32) int {
	return int(float32(oldConfidence) * factor)
}

func calcConfidenceChange(oldConfidence int, factor float32) int {
	return (int(float32(oldConfidence) * factor)) - oldConfidence
}

func plConfUpdateString(currentPlayerID int, floatValue float32) string {
	return "PLCONF" + intToString(currentPlayerID) + "=PLCONF" + intToString(currentPlayerID) + "*" + float32ToString(floatValue)
}

func fn_adjustCONF_10_20(adjustCONF10A_dec_SuiteAce float32, adjustCONF20A_dec_NonSuiteAce float32, currentCard CardValue, currentCardPoints int, cardNames [13]string, cardsPerSuite int,
	trumpSuite string, currentPlayerID int) {
	// --- 10A) & 20A) player plays a non trump non ace card from lower half of the suite, when aces are still in game
	//			player expects to : lose, aces will be played agaist that card
	//          assumption 10A - player doesn't have that ace card
	//			assumption 20A - player doesn't have other non trump aces
	if currentCard.Suite == trumpSuite { // ignore if it's trump suite
		return
	}
	for c := 0; c < cardsPerSuite/2; c++ { // ignore if it's more upper half of the suite card in play
		if currentCard.Name == cardNames[c] {
			return
		}
	}
	var cardCount int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+currentCard.Suite+"' AND INPLAY=TRUE AND CARDNAME='ace'", "SQL10_20-1-fn_adjustCONF10_20", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&cardCount)
	}
	if cardCount == 0 { // return if suite ace is not in game
		return
	}
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF10A_dec_SuiteAce)+" WHERE CARDNAME='ace' AND CARDSUITE='"+currentCard.Suite+"'", "SQL10_20-2-fn_adjustCONF10_20", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF20A_dec_NonSuiteAce)+" WHERE CARDNAME='ace' AND CARDSUITE<>'"+currentCard.Suite+"' AND CARDSUITE<>'"+trumpSuite+"'", "SQL10_20-3-fn_adjustCONF10_20", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_30_40(adjustCONF30A_dec_LowerSuiteCards float32, adjustCONF40A_dec_SuiteCards float32, currentCard CardValue, currentCardPoints int, trumpSuite string, currentPlayerID int) {
	// --- 30A) & 40A) player plays a non trump card, when higher card from the same suite are still in game
	//		  	player expects to : lose, that there will be higher cards from the same suite challenging that card,
	//        	assumption 30A - player doesn't have lower cards from that suite
	//		 	assumption 40A - player wants to get rid if that suite
	if currentCard.Suite == trumpSuite { // ignore if it's trump suite
		return
	}
	var cardCount int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+currentCard.Suite+"' AND INPLAY=TRUE AND ROUNDPOINTS>"+intToString(currentCardPoints), "SQL30_40-1-fn_adjustCONF_30_40", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&cardCount)
	}
	if cardCount == 0 { // return if there are no card above this card in the suite
		return
	}
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF30A_dec_LowerSuiteCards)+" WHERE CARDSUITE='"+currentCard.Suite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL10_20-2-fn_adjustCONF10_20", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF40A_dec_SuiteCards)+" WHERE CARDSUITE='"+currentCard.Suite+"'", "SQL10_20-3-fn_adjustCONF10_20", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_50_60(adjustCONF50A_dec_UpperSuiteCards float32, adjustCONF60A_dec_LowerSuiteCards float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 50A) & 60A) PLAYER LOSE - roundsuite is not trump, Player played roundsuite, didn't win and roundwinner is roundsuite
	//			player expects to lose, and didn't have round cards above roundwinner
	//			assumption 50A - player didn't have higher cards than roundwinner - *if the roundwinner is opposing side
	//			assumption 60A - player didn't have lower cards than the round cardcard he played
	if roundSuite == trumpSuite || currentCard.Suite != roundSuite || currentPlayerID == currentRoundWinnerID || currentRoundWinnerCard.Suite != roundSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	if currentRoundWinnerID%2 != currentPlayerID%2 { //roundwinner and current player in opposing teams
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF50A_dec_UpperSuiteCards)+" WHERE CARDSUITE='"+currentCard.Suite+"' AND ROUNDPOINTS>"+intToString(currentRoundWinnerPoints), "SQL-A3nHVTXhlunEfI64zG3La8kudx09kTQu", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF60A_dec_LowerSuiteCards)+" WHERE CARDSUITE='"+currentCard.Suite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-A3nHVTXhlunEfI64zG3La8kudx09kTQu", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_70(adjustCONF70A_dec_LowerSuiteCards float32, currentPlayerID int, currentCard CardValue, currentCardPoints int, currentRoundWinnerID int, currentRoundWinnerCard CardValue,
	currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 70A) PLAYER LOSE - roundsuite is not trump, Player played roundsuite, didn't win and roundwinner is a trump
	//			player expects to lose, had roundcards but already round won by trump card
	//			assumption 60A - player didn't have lower cards than the round cardcard he played
	if roundSuite == trumpSuite || currentPlayerID == currentRoundWinnerID || currentCard.Suite != roundSuite || currentRoundWinnerCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF70A_dec_LowerSuiteCards)+" WHERE CARDSUITE='"+currentCard.Suite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-RoLgCn8NIsmOLhBcqRtELYyiQtLhEsOZ", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_80(adjustCONF80A_dec_NonTrumpCards float32, adjustCONF80B_dec_TrumpsBelowPlayed float32, adjustCONF80C_dec_TrumpsAboveWinner float32,
	currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 80A/B) PLAYER LOSE - roundsuite is not trump, Player played trumpCard, didn't win and roundwinner is a trump
	//			player expects to lose, didn't have any other cards except trumps
	//			assumption 80A - didn't have any other cards except trumps
	//			assumption 80B - player doesn't have lower trumps than played *if the roundwinner is same side
	// 			assumption 80C - player deosn't have higher trumps than roundwinner  *if the roundwinner is opposing side
	if currentPlayerID == currentRoundWinnerID || roundSuite == trumpSuite || currentCard.Suite != trumpSuite || currentRoundWinnerCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF80A_dec_NonTrumpCards)+" WHERE CARDSUITE<>'"+trumpSuite+"'", "SQL-boVihZoTxn9VR3ojKt7TwbplTGkCfvMP", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	if currentRoundWinnerID%2 == currentPlayerID%2 { // winner is same team
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF80B_dec_TrumpsBelowPlayed)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-G05tyXQR30hTOcHSrX25ZyJ6TRJSWp18", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	} else {
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF80C_dec_TrumpsAboveWinner)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS>"+intToString(currentRoundWinnerPoints), "SQL-hWdQs1wqtRtEcOlXGIAJ1ByRgRjniqfh", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func fn_adjustCONF_85(adjustCONF85A_dec_NonRoundCards float32, adjustCONF85B_dec_TrumpsAbovePlayed float32, adjustCONF85C_dec_OtherSuitesBelow float32,
	currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 85A/B) PLAYER LOSE - roundsuite is not trump, Player played non round non trump card
	//			player expects to lose, didn't have any roundcards or trumpcards
	//			assumption 85A - didn't have lesser cards in the suite played
	//			assumption 85B - player didn't have trump cards above the roundwinner *if the roundwinner is opposing side
	//			assumption 85C - player didn't have lesser points cards than the one played in other suites *if the roundwinner is same side
	if currentPlayerID == currentRoundWinnerID || roundSuite == trumpSuite || currentCard.Suite == roundSuite || currentCard.Suite == trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF85A_dec_NonRoundCards)+" WHERE CARDSUITE='"+currentCard.Suite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-ikSrkV910tNhz4coHq13byQVUNmtQq1y", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	if currentRoundWinnerID%2 == currentPlayerID%2 { // winner is same team 85C
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF85C_dec_OtherSuitesBelow)+" WHERE CARDSUITE<>'"+trumpSuite+"' AND CARDSUITE<>'"+currentCard.Suite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-qRHho9V2Oz2m5PrO7ILAd9jTcsvy1ZUm", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	} else { // 85B
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF85B_dec_TrumpsAbovePlayed)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS>"+intToString(currentRoundWinnerPoints), "SQL-G05tyXQR30hTOcHSrX25ZyJ6TRJSWp18", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func fn_adjustCONF_90(adjustCONF90A_dec_RoundCardBetween float32,
	currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string, prevRoundWinnerPoints int) {
	// --- 90A) PLAYER WINS -  roundsuite is not trump, Player played roundsuite
	//			player expects to win
	//			assumption 90A - player don't have round cards more than prev roundwinner and less than card player played
	if currentPlayerID != currentRoundWinnerID || roundSuite == trumpSuite || currentCard.Suite != roundSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF90A_dec_RoundCardBetween)+" WHERE CARDSUITE='"+currentCard.Suite+"' AND ROUNDPOINTS>"+intToString(prevRoundWinnerPoints)+" AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-r6HgwBC3hRL7CWcgLlluSjkyHa7kJvnR", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_100(adjustCONF100A_dec_LesserTrump float32, adjustCONF100B_dec_BetweenTrump float32,
	currentPlayerID int, currentCard CardValue, currentCardPoints int, currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int,
	roundSuite string, trumpSuite string, prevRoundWinnerPoints int, prevRoundWinnerCard CardValue) {
	// --- 100A/B) PLAYER WINS -  roundsuite is not trump, Player played trumpsuite
	//			player expects to win
	//			assumption 100A - prev roundwinner is roundcard, player don't have lesser trump cards than the one player played
	// 			assumption 100B - prev roundwinner is trumpcard, player don't have trump more than prev roundwinner and the less than the card that player played
	if currentPlayerID != currentRoundWinnerID || roundSuite == trumpSuite || currentCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	if prevRoundWinnerCard.Suite == roundSuite {
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF100A_dec_LesserTrump)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-F5p4GmJSrmX7TmKA8ZeFyh41slSfTx9N", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	} else if prevRoundWinnerCard.Suite == trumpSuite {
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF100A_dec_LesserTrump)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints)+" AND ROUNDPOINTS>"+intToString(prevRoundWinnerPoints), "SQL-ZacNjwLfLZ113Wx1JasfH1eQ4KRVDbha", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func fn_adjustCONF_105(adjustCONF105A_dec_LesserTrump float32, adjustCONF105B_dec_HigherTrump float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 105A) PLAYER LOSES -  roundsuite is trump, Player played trumpsuite below roundwinner
	//			player expects to lose
	//			assumption 105A - roundwinner is trump, player doesn't have lesser trump cards than the one player played *if the roundwinner is same side
	// 			assumption 105B - roundwinner is trump, player doesn't have higher trump cards than the the roundwinner *if the roundwinner is oppositing side
	if currentPlayerID == currentRoundWinnerID || trumpSuite != roundSuite || currentCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	if currentRoundWinnerID%2 == currentPlayerID%2 { // same team
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF105A_dec_LesserTrump)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-jDJkCs6SEeI9WDudJPJSLzfaGfdDJwmI", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	} else {
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF105B_dec_HigherTrump)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS>"+intToString(currentCardPoints), "SQL-TcVXZJ1e16YW0v2Ekmphvc9QCyegsOFg", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
}

func fn_adjustCONF_106(adjustCONF106A_dec_LesserTrump float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string,
	prevRoundWinnerPoints int, prevRoundWinnerCard CardValue) {
	// --- 106A) PLAYER WINS -  roundsuite is trump, Player played trumpsuite above roundwinner
	//			player expects to win, becuase there are no trump cards less than played card above roundwinner
	//			assumption 106A - roundwinner is trump, player don't have trump cards less than played card and above roundwinner
	if currentPlayerID != currentRoundWinnerID || roundSuite != trumpSuite || currentCard.Suite != trumpSuite || prevRoundWinnerCard.Suite != trumpSuite || currentCardPoints < prevRoundWinnerPoints {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF106A_dec_LesserTrump)+" WHERE CARDSUITE='"+trumpSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints)+" AND ROUNDPOINTS>"+intToString(prevRoundWinnerPoints), "SQL-TcVXZJ1e16YW0v2Ekmphvc9QCyegsOFg", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_107(adjustCONF107A_dec_NonTrumpCards float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 107A/B) PLAYER LOSES - roundsuite is trump, roundwinner is trump
	//			player expects to lose - player plays other suite
	//			assumption 107A: player doesn't have cards with points less than the one played
	//			assumption 107B - player didn't have lesser points cards than the one played in other suites
	if currentPlayerID == currentRoundWinnerID || roundSuite != trumpSuite || currentRoundWinnerCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF107A_dec_NonTrumpCards)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-TcVXZJ1e16YW0v2Ekmphvc9QCyegsOFg", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_110(adjustCONF110A_dec_LowerSuite float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// FOR PLAYERS IN MIDDLE OF THE ROUND ///////////////////////////////////////////////////////////////////
	// --- 110A) - roundsuite is not trump, roundwinner is not trump
	//			player expects to lose - player plays roundsuite above roundwinner when higher than played suite card are still in play
	//			assumption 110A - playter doesn't have round cards below the card played
	if roundSuite == trumpSuite || currentRoundWinnerID != currentPlayerID || currentCard.Suite == trumpSuite {
		return
	}
	var cardCount int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+currentCard.Suite+"' AND INPLAY=TRUE AND ROUNDPOINTS>"+intToString(currentCardPoints), "SQL-bh4TPM0qZ3z0yrtiMoOFk3zU5Sj8Oz5J", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&cardCount)
	}
	if cardCount == 0 { // return if there are no card above this card in the suite
		return
	}
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF110A_dec_LowerSuite)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_120(adjustCONF120A_dec_LesserSuite float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 120A) - roundsuite is not trump, roundwinner is not trump
	//			player expects to lose - player plays roundsutite below roundwinner
	//			assumption 120A - playter doesn't have round cards below the card played
	if roundSuite == trumpSuite || currentRoundWinnerID == currentPlayerID || currentCard.Suite == trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF120A_dec_LesserSuite)+" WHERE CARDSUITE='"+roundSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_140(adjustCONF140A_dec_LesserTrump float32, adjustCONF140B_dec_AboveTrump float32, adjustCONF140C_dec_OtherSuiteExcptTrump float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 140A) - roundsuite is not trump, roundwinner is trump
	//			player expects to lose - player plays trumpsuite below roundwinner
	//			assumption 140A - playter doesn't have trump cards below the card played
	//          assumption 140B - playter deosn't have trump cards above roundwinner *if roundwinner is opposing side
	//			assumption 140C - player doesn't have other suite cards
	if roundSuite == trumpSuite || currentRoundWinnerID == currentPlayerID || currentRoundWinnerCard.Suite != trumpSuite || currentCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF140A_dec_LesserTrump)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints)+" AND CARDSUITE='"+trumpSuite+"'", "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	if currentRoundWinnerID%2 != currentPlayerID%2 {
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF140B_dec_AboveTrump)+" WHERE ROUNDPOINTS>"+intToString(currentRoundWinnerPoints)+" AND CARDSUITE='"+trumpSuite+"'", "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF140C_dec_OtherSuiteExcptTrump)+" WHERE CARDSUITE<>'"+trumpSuite+"' AND CARDSUITE<>'"+roundSuite+"'", "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_160(adjustCONF160A_dec_NonRoundCards float32, adjustCONF160B_dec_TrumpsAbovePlayed float32, adjustCONF160C_dec_OtherSuitesBelow float32,
	currentPlayerID int, currentCard CardValue, currentCardPoints int, currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int,
	roundSuite string, trumpSuite string) {
	// -- 160A) - roundsuite is not trump, roundwinner is trump
	//			player expects to lose - player plays other suite
	//			assumption 160A: player doesn't have cards with points less than the one played
	//			assumption 160B - player didn't have trump cards above the roundwinner *if roundwinner is opposing side
	//			assumption 160C - player didn't have lesser points cards than the one played in other suites
	if roundSuite == trumpSuite || currentRoundWinnerCard.Suite != trumpSuite || currentCard.Suite == trumpSuite || currentCard.Suite == roundSuite || currentPlayerID == currentRoundWinnerID {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF160A_dec_NonRoundCards)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
	if currentRoundWinnerID%2 != currentPlayerID%2 {
		res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF160B_dec_TrumpsAbovePlayed)+" WHERE ROUNDPOINTS>"+intToString(currentRoundWinnerPoints)+" AND CARDSUITE='"+trumpSuite+"'", "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
		if res == nil || !dbResponse {
			os.Exit(1)
		}
	}
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF160C_dec_OtherSuitesBelow)+" WHERE CARDSUITE<>'"+trumpSuite+"' AND CARDSUITE<>'"+roundSuite+"' AND ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_170(adjustCONF170A_dec_LesserTrump float32, currentPlayerID int, currentCard CardValue, currentCardPoints int, currentRoundWinnerID int, currentRoundWinnerCard CardValue,
	currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 170A) PLAYER LOSES -  roundsuite is trump, Player played trumpsuite below roundwinner
	//			player expects to lose
	//			assumption 105A - roundwinner is trump, player don't have lesser trump cards than the one player played
	if currentPlayerID == currentRoundWinnerID || roundSuite != trumpSuite || currentCard.Suite != trumpSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF170A_dec_LesserTrump)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints)+" AND CARDSUITE='"+trumpSuite+"'", "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_180(adjustCONF180A_dec_LesserTrump float32, currentPlayerID int, currentCard CardValue, currentCardPoints int, currentRoundWinnerID int,
	currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 180A) PLAYER LOSES -  roundsuite is trump, Player played trumpsuite above roundwinner
	//			player expects to lose, there are higher trump cards than what player played
	//			assumption 106A - roundwinner is trump, player don't have lesser trump cards than the one player played
	if currentPlayerID != currentRoundWinnerID || roundSuite != trumpSuite || currentCard.Suite != trumpSuite {
		return
	}
	var cardCount int = 0
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT COUNT(*) FROM CARDS WHERE CARDSUITE='"+trumpSuite+"' AND INPLAY=TRUE AND ROUNDPOINTS>"+intToString(currentCardPoints), "SQL-bh4TPM0qZ3z0yrtiMoOFk3zU5Sj8Oz5J", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		rows.Scan(&cardCount)
	}
	if cardCount == 0 { // return if there are no card above this card in the suite
		return
	}
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF180A_dec_LesserTrump)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints)+" AND CARDSUITE='"+trumpSuite+"'", "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func fn_adjustCONF_190(adjustCONF190A_dec_NonRoundCards float32, currentPlayerID int, currentCard CardValue, currentCardPoints int,
	currentRoundWinnerID int, currentRoundWinnerCard CardValue, currentRoundWinnerPoints int, roundSuite string, trumpSuite string) {
	// --- 190A/B) PLAYER LOSES - roundsuite is trump, roundwinner is trump
	//			player expects to lose - player plays other suite
	//			assumption 190A: player doesn't have cards with points less than the one played
	if currentPlayerID == currentRoundWinnerID || roundSuite != trumpSuite || currentRoundWinnerCard.Suite != trumpSuite || currentCard.Suite == trumpSuite || currentCard.Suite == roundSuite {
		return
	}
	var dbResponse bool = false
	var res sql.Result
	res, dbResponse = executeOnDB("UPDATE CARDS SET "+plConfUpdateString(currentPlayerID, adjustCONF190A_dec_NonRoundCards)+" WHERE ROUNDPOINTS<"+intToString(currentCardPoints), "SQL-0sbERVJWu0KdmnRxa5quutANbqOYyol2", true)
	if res == nil || !dbResponse {
		os.Exit(1)
	}
}

func printGameRounds() {
	var dbResponse bool = false
	var rows *sql.Rows
	rows, dbResponse = queryFromDB("SELECT ROUND,ROUNDTURN,PLAYERID,FRIEND,PLAYERNAME,CARDSUITE,CARDNAME,WINNER,ROUNDSUITE,TRUMPSUITE from ROUNDS ORDER BY ROUND,ROUNDTURN", "SQL-ZPHkMHLpv2ijqceRlKyjOfcs79lHHlNa", true)
	if !dbResponse || rows == nil {
		os.Exit(1)
	}
	for rows.Next() {
		var round, roundTurn, playerID, winnerID int
		var friendB bool
		var teamS, playerName, cardSuite, cardName, winnerS, roundSuite, trumpSuite string
		cardSuite = "nul"
		playerID = -1
		rows.Scan(&round, &roundTurn, &playerID, &friendB, &playerName, &cardSuite, &cardName, &winnerID, &roundSuite, &trumpSuite)
		if friendB {
			teamS = "Friend"
		} else {
			teamS = "Foe  "
		}
		if winnerID == playerID {
			winnerS = "ROUND WINNER"
		} else {
			winnerS = "            "
		}
		if roundTurn == 0 {
			fmt.Println()
			fmt.Println("*** ROUND: ", intToString(round), ", TRUMP-SUITE: ", trumpSuite, ", ROUND-SUITE: ", roundSuite)
			fmt.Println("#\tID\tNAME\tTEAM\tPLAYEDCARD\tWINNER?---")
		}
		if cardSuite != "nul" {
			fmt.Println(intToString(roundTurn), "\t", intToString(playerID), "\t", playerName, "\t", teamS, "\t", cardSuite, "-", cardName, "\t", winnerS)
		}
	}
}

func nulPrint1(playedRoundSuite bool, playedNonRoundTrump bool, currentCardIsRoundWinner bool) {
	if playedRoundSuite && playedNonRoundTrump && currentCardIsRoundWinner {

	}
}

func nulPrint2(currentCardIsRoundWinner bool) {
	if currentCardIsRoundWinner {

	}
}

func main() {

	cardNames := [13]string{"ace", "king", "queen", "jack", "ten", "nine", "eight", "seven", "six", "five", "four", "three", "two"}
	cardNamesAb := [13]string{"a", "k", "q", "j", "t", "9", "8", "7", "6", "5", "4", "3", "2"}
	cardPoints := [13]int{45, 38, 29, 22, 17, 13, 10, 8, 6, 5, 4, 3, 2}
	// Points -------------A   K   Q   J  10   9   8  7  6  5  4  3  2
	cardSuites := [4]string{"hearts", "spades", "diamonds", "clubs"}
	cardSuitesAb := [4]string{"h", "s", "d", "c"}
	var pointsAddForRoundSuite int = 1000
	var pointsAddForTrumpSuite int = 2000
	var cardConfidenceAtStart int = 10000
	cardConfidenceAddForTrumpcCaller := [13]float32{1.150, 1.150, 1.150, 1.150, 1.100, 1.100, 1.100, 1.100, 1.100, 1.100, 1.100, 1.100, 1.100}
	// cardConfidenceAddForTrumpcCaller---------------A      K      Q      J      10     9      8      7      6      5      4      3      2
	playedInRound := [8]int{-1, -1, -1, -1, -1, -1, -1, -1}
	//************  roundwinner = the max card on round before player played
	// FOR ROUND STARTERS ////////////////////////////////////////////////////////////////////////////////////
	// --- 10A) & 20A) player plays a non trump non ace card from lower half of the suite, when aces are still in game
	//			player expects to : lose, aces will be played agaist that card
	//          assumption 10A - player doesn't have that ace card
	//			assumption 20A - player doesn't have other non trump aces
	var adjustCONF10A_dec_SuiteAce float32 = 0.350
	var adjustCONF20A_dec_NonSuiteAce float32 = 0.500
	// --- 30A) & 40A) player plays a non trump card, when higher card from the same suite are still in game
	//		  	player expects to : lose, that there will be higher cards from the same suite challenging that card,
	//        	assumption 30A - player doesn't have lower cards from that suite
	//		 	assumption 40A - player wants to get rid if that suite
	var adjustCONF30A_dec_LowerSuiteCards float32 = 0.500
	var adjustCONF40A_dec_SuiteCards float32 = 0.850
	// FOR ROUND ENDERS ////////////////////////////////////////////////////////////////////////////////////
	// --- 50A) & 60A) PLAYER LOSE - roundsuite is not trump, Player played roundsuite, didn't win and roundwinner is roundsuite
	//			player expects to lose, and didn't have round cards above roundwinner
	//			assumption 50A - player didn't have higher cards than roundwinner - *if the roundwinner is opposing side
	//			assumption 60A - player didn't have lower cards than the round cardcard he played
	var adjustCONF50A_dec_UpperSuiteCards float32 = 0.250
	var adjustCONF60A_dec_LowerSuiteCards float32 = 0.250
	// --- 70A) PLAYER LOSE - roundsuite is not trump, Player played roundsuite, didn't win and roundwinner is a trump
	//			player expects to lose, had roundcards but already round won by trump card
	//			assumption 60A - player didn't have lower cards than the round cardcard he played
	var adjustCONF70A_dec_LowerSuiteCards float32 = 0.250
	// --- 80A/B) PLAYER LOSE - roundsuite is not trump, Player played trumpCard, didn't win and roundwinner is a trump
	//			player expects to lose, didn't have any other cards except trumps
	//			assumption 80A - didn't have any other cards except trumps
	//			assumption 80B - player doesn't have lower trumps than played *if the roundwinner is same side
	// 			assumption 80C - player deosn't have higher trumps than roundwinner  *if the roundwinner is opposing side
	var adjustCONF80A_dec_NonTrumpCards float32 = 0.250
	var adjustCONF80B_dec_TrumpsBelowPlayed float32 = 0.250
	var adjustCONF80C_dec_TrumpsAboveWinner float32 = 0.250
	// --- 85A/B) PLAYER LOSE - roundsuite is not trump, Player played non round non trump card
	//			player expects to lose, didn't have any roundcards or trumpcards
	//			assumption 85A - didn't have lesser cards in the suite played
	//			assumption 85B - player didn't have trump cards above the roundwinner *if the roundwinner is opposing side
	//			assumption 85C - player didn't have lesser points cards than the one played in other suites *if the roundwinner is same side
	var adjustCONF85A_dec_NonRoundCards float32 = 0.200
	var adjustCONF85B_dec_TrumpsAbovePlayed float32 = 0.200
	var adjustCONF85C_dec_OtherSuitesBelow float32 = 0.200
	// --- 90A) PLAYER WINS -  roundsuite is not trump, Player played roundsuite
	//			player expects to win
	//			assumption 90A - player don't have round cards more than prev roundwinner and less than card player played
	var adjustCONF90A_dec_RoundCardBetween float32 = 0.2000
	// --- 100A/B) PLAYER WINS -  roundsuite is not trump, Player played trumpsuite
	//			player expects to win
	//			assumption 100A - prev roundwinner is roundcard, player don't have lesser trump cards that the one player played
	// 			assumption 100B - prev roundwinner is trumpcard, player don't have trump more than prev roundwinner and the less than the card that player played
	var adjustCONF100A_dec_LesserTrump float32 = 0.200
	var adjustCONF100B_dec_BetweenTrump float32 = 0.200
	// --- 105A) PLAYER LOSES -  roundsuite is trump, Player played trumpsuite below roundwinner
	//			player expects to lose
	//			assumption 105A - roundwinner is trump, player don't have lesser trump cards than the one player played *if the roundwinner is same side
	// 			assumption 105B - roundwinner is trump, player don't have higher trump cards than the the roundwinner *if the roundwinner is oppositing side
	var adjustCONF105A_dec_LesserTrump float32 = 0.200
	var adjustCONF105B_dec_HigherTrump float32 = 0.200
	// --- 106A) PLAYER WINS -  roundsuite is trump, Player played trumpsuite above roundwinner
	//			player expects to win, becuase there are no trump cards less than played card above roundwinner
	//			assumption 106A - roundwinner is trump, player don't have trump cards less than played card and above roundwinner
	var adjustCONF106A_dec_LesserTrump float32 = 0.200
	// --- 107A/B) PLAYER LOSES - roundsuite is trump, roundwinner is trump
	//			player expects to lose - player plays other suite
	//			assumption 107A: player doesn't have cards with points less than the one played in non trump Suites
	var adjustCONF107A_dec_NonTrumpCards float32 = 0.200
	// FOR PLAYERS IN MIDDLE OF THE ROUND ///////////////////////////////////////////////////////////////////
	// --- 110A) - roundsuite is not trump, roundwinner is not trump
	//			player expects to lose - player plays roundsuite above roundwinner when higher than played suite card are still in play
	//			assumption 110A - playter doesn't have round cards below the card played
	var adjustCONF110A_dec_LowerSuite float32 = 0.750
	// --- 120A) - roundsuite is not trump, roundwinner is not trump
	//			player expects to lose - player plays roundsutite below roundwinner
	//			assumption 120A - playter doesn't have round cards below the card played
	var adjustCONF120A_dec_LesserSuite float32 = 0.500
	// --- 130A) - roundsuite is not trump, roundwinner is trump
	//			player expects to lose - player plays roundsutite
	//			assumption - same as above two conditons

	// --- 140A) - roundsuite is not trump, roundwinner is trump
	//			player expects to lose - player plays trumpsuite below roundwinner
	//			assumption 140A - playter doesn't have trump cards below the card played
	//          assumption 140B - playter deosn't have trump cards above roundwinner *if roundwinner is opposing side
	//			assumption 140C - player doesn't have other suite cards
	var adjustCONF140A_dec_LesserTrump float32 = 0.200
	var adjustCONF140B_dec_AboveTrump float32 = 0.500
	var adjustCONF140C_dec_OtherSuiteExcptTrump float32 = 0.200
	// -- 150A) - roundsuite is not trump, roundwinner is trump
	//			player expects to win - player plays trumpsuite above roundwinner
	//			assumption: Can't predict on trump confidence
	//
	// -- 160A) - roundsuite is not trump, roundwinner is trump
	//			player expects to lose - player plays other suite
	//			assumption 160A: player doesn't have cards with points less than the one played
	//			assumption 160B - player didn't have trump cards above the roundwinner
	//			assumption 160C - player didn't have lesser points cards than the one played in other suites
	var adjustCONF160A_dec_NonRoundCards float32 = 0.200
	var adjustCONF160B_dec_TrumpsAbovePlayed float32 = 0.500
	var adjustCONF160C_dec_OtherSuitesBelow float32 = 0.200
	// --- 170A) PLAYER LOSES -  roundsuite is trump, Player played trumpsuite below roundwinner
	//			player expects to lose
	//			assumption 105A - roundwinner is trump, player don't have lesser trump cards than the one player played
	var adjustCONF170A_dec_LesserTrump float32 = 0.200
	// --- 180A) PLAYER LOSES -  roundsuite is trump, Player played trumpsuite above roundwinner
	//			player expects to lose, there are higher trump cards than what player played
	//			assumption 106A - roundwinner is trump, player don't have lesser trump cards than the one player played
	var adjustCONF180A_dec_LesserTrump float32 = 0.750
	// --- 190A/B) PLAYER LOSES - roundsuite is trump, roundwinner is trump
	//			player expects to lose - player plays other suite
	//			assumption 190A: player doesn't have cards with points less than the one played
	//			assumption 190B - player didn't have lesser points cards than the one played in other suites
	var adjustCONF190A_dec_NonRoundCards float32 = 0.200

	var err error
	db, err = sql.Open(dbType, dbPath)
	if err != nil {
		fmt.Println("SQLERR-main:Opening DB File")
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		fmt.Println("SQLERR-main:PRAGMA journal_mode=WAL")
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
	initPlayersTab(noOfPlayers)                          // insert player info plater tab
	initCardsTab(cardsPerSuite, cardPoints, noOfPlayers) // remove unused cards
	initRoundsTab(cardsPerPlayer, noOfPlayers)           // init records of rounds
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
		//		var prevRoundWinnerID int = -1
		//		var prevRoundWinnerName string = "nul-player"
		var prevRoundWinnerCard CardValue
		var prevRoundWinnerPoints int = 0
		playedInRound = resetPlayersTabPlayedInRound(noOfPlayers, playedInRound)
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
				if playerInRound == 0 { // round starter
					roundSuite = currentCard.Suite
					if roundSuite != trumpSuite {
						updateCardsTabRoundPointsForRounduite(pointsAddForRoundSuite, roundSuite)
					}
					updateRoundsTabWithRoundSuite(roundSuite, roundNo)
					playedInRound = updatePlayersTabPlayedInRound(currentPlayerID, playedInRound)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard, noOfPlayers)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex)
					//////////////////////////////////////// CRITERIA TO ADJUST PLCONF ///////////////////////
					fn_adjustCONF_10_20(adjustCONF10A_dec_SuiteAce, adjustCONF20A_dec_NonSuiteAce, currentCard, currentCardPoints, cardNames, cardsPerSuite, trumpSuite, currentPlayerID)
					fn_adjustCONF_30_40(adjustCONF30A_dec_LowerSuiteCards, adjustCONF40A_dec_SuiteCards, currentCard, currentCardPoints, trumpSuite, currentPlayerID)
					///////////////////////////////////////////////////////////////////////////////////////////////

				} else if playerInRound == (noOfPlayers - 1) { // round ender
					if roundSuite != currentCard.Suite {
						updatePlayersTabProbForRoundSuite(roundSuite, roundNo, currentPlayerID) // Players probability set to minus value for round
						updateCardsTabForRoundSuiteProb(roundSuite, currentPlayerID)            // Cards probabaility make 0 for roundsite cards for this player
						playedRoundSuite = false
						if currentCard.Suite == trumpSuite {
							playedNonRoundTrump = true
						}
					}
					playedInRound = updatePlayersTabPlayedInRound(currentPlayerID, playedInRound)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard, noOfPlayers)
					//					prevRoundWinnerID = currentRoundWinnerID
					//					prevRoundWinnerName = currentRoundWinnerName
					prevRoundWinnerCard = currentRoundWinnerCard
					prevRoundWinnerPoints = currentRoundWinnerPoints
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex)
					//////////////////////////////////////// CRITERIA TO ADJUST PLCONF ///////////////////////
					fn_adjustCONF_50_60(adjustCONF50A_dec_UpperSuiteCards, adjustCONF60A_dec_LowerSuiteCards, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_70(adjustCONF70A_dec_LowerSuiteCards, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_80(adjustCONF80A_dec_NonTrumpCards, adjustCONF80B_dec_TrumpsBelowPlayed, adjustCONF80C_dec_TrumpsAboveWinner, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_85(adjustCONF85A_dec_NonRoundCards, adjustCONF85B_dec_TrumpsAbovePlayed, adjustCONF85C_dec_OtherSuitesBelow, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_90(adjustCONF90A_dec_RoundCardBetween, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite, prevRoundWinnerPoints)
					fn_adjustCONF_100(adjustCONF100A_dec_LesserTrump, adjustCONF100B_dec_BetweenTrump, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite, prevRoundWinnerPoints, prevRoundWinnerCard)
					fn_adjustCONF_105(adjustCONF105A_dec_LesserTrump, adjustCONF105B_dec_HigherTrump, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_106(adjustCONF106A_dec_LesserTrump, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite, prevRoundWinnerPoints, prevRoundWinnerCard)
					fn_adjustCONF_107(adjustCONF107A_dec_NonTrumpCards, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					///////////////////////////////////////////////////////////////////////////////////////////////
				} else { // round mid
					if roundSuite != currentCard.Suite {
						updatePlayersTabProbForRoundSuite(roundSuite, roundNo, currentPlayerID) // Players probability set to minus value for round
						updateCardsTabForRoundSuiteProb(roundSuite, currentPlayerID)            // Cards probabaility make 0 for roundsite cards for this player
						playedRoundSuite = false
						if currentCard.Suite == trumpSuite {
							playedNonRoundTrump = true
						}
					}
					playedInRound = updatePlayersTabPlayedInRound(currentPlayerID, playedInRound)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard, noOfPlayers)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex)
					//////////////////////////////////////// CRITERIA TO ADJUST PLCONF ///////////////////////
					fn_adjustCONF_110(adjustCONF110A_dec_LowerSuite, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_120(adjustCONF120A_dec_LesserSuite, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_140(adjustCONF140A_dec_LesserTrump, adjustCONF140B_dec_AboveTrump, adjustCONF140C_dec_OtherSuiteExcptTrump, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_160(adjustCONF160A_dec_NonRoundCards, adjustCONF160B_dec_TrumpsAbovePlayed, adjustCONF160C_dec_OtherSuitesBelow, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_170(adjustCONF170A_dec_LesserTrump, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_180(adjustCONF180A_dec_LesserTrump, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					fn_adjustCONF_190(adjustCONF190A_dec_NonRoundCards, currentPlayerID, currentCard, currentCardPoints, currentRoundWinnerID, currentRoundWinnerCard, currentRoundWinnerPoints, roundSuite, trumpSuite)
					///////////////////////////////////////////////////////////////////////////////////////////
				}

				nulPrint1(playedRoundSuite, playedNonRoundTrump, currentCardIsRoundWinner) //////  REMOVE REMOVE REMOVE REMOVE ///////////////////////////

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
						currentCard = getMinFromANonTrumpSuiteIfFriendsHaveMorePossibilityOfMaxOfThatSuiteAndFoesHaveNoTrumps(trumpSuite, noOfPlayers, cardSuites, cardPoints, cardsPerSuite) // 40
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 44
							break
						}
						currentCard = getNonTrumpCardIfIHaveSecondMaxTrumpAndCanConvinceFoesWhoHaveThatMaxTrumpOnlyToCutMyCard(trumpSuite, noOfPlayers, cardSuites, foePoints, friendPoints, cardsPerPlayer) // 40
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 45
							break
						}
						currentCard = getMinFromANonTrumpSuiteIfFriendsHaveMaxOfThatSuiteAndFoesHaveLessPossibilityOfTrumps(trumpSuite, noOfPlayers, cardSuites, cardPoints, cardsPerSuite) // 40
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 46
							break
						}
						currentCard = getMinCardIfIHaveOneMinCardFromANonTrumpSuite(trumpSuite, cardSuites, cardPoints, cardsPerSuite) // 50
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 50
							break
						}
						currentCard = getMyMinPointsCard() // 60
						myCardPlayCondition = 60
						break
					}
					roundSuite = currentCard.Suite
					if roundSuite != trumpSuite {
						updateCardsTabRoundPointsForRounduite(pointsAddForRoundSuite, roundSuite)
					}
					updateRoundsTabWithRoundSuite(roundSuite, roundNo)
					playedInRound = updatePlayersTabPlayedInRound(currentPlayerID, playedInRound)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard, noOfPlayers)
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
					playedInRound = updatePlayersTabPlayedInRound(currentPlayerID, playedInRound)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard, noOfPlayers)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForMyPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex, myCardPlayCondition)
				} else { // i'm not the round starter or last player
					/////////// NOT OBVIOUS CHOICES //////////////////////
					currentCard.Suite = "nul"
					for currentCard.Suite == "nul" {
						currentCard = getMaxRoundCardIfIHaveItInMyHandIfItsAboveCurrentWinner(roundSuite, currentRoundWinnerPoints) ///////////////////////// get my Max RoundCard /////////////////////////
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 110
							break
						}
						/*currentCard = getMaxRoundCardInMyHand(roundSuite) ///////////////////////// get my Max RoundCard /////////////////////////
						if currentCard.Suite != "nul" {
							break
						}*/
						currentCard = getMinRoundCardInMyHand(roundSuite) ///////////////////////// get my Max RoundCard /////////////////////////
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 130
							break
						}
						currentCard = getMidTrumpCardInMyHandWhereFoesHaveLessPosibilityForAHigherTrumpCard(trumpSuite, playedInRound) //////////////////////////////get my Max trump Card ////////////////////
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 140
							break
						}
						currentCard = getMaxTrumpCardInMyHand(trumpSuite) ///////BEBUG CONDITION//////////
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 170
							break
						}
						//currentCard = getARandomCardInMyHand()
						currentCard = getMinCardIfIHaveOneMinCardFromANonTrumpSuite(trumpSuite, cardSuites, cardPoints, cardsPerSuite)
						if currentCard.Suite != "nul" {
							myCardPlayCondition = 190
							break
						}
						currentCard = getMyMinPointsCard()
						myCardPlayCondition = 200
						break
					}
					playedInRound = updatePlayersTabPlayedInRound(currentPlayerID, playedInRound)
					currentCardIndex, currentCardPoints = updateCardsTabForPlayedCard(currentPlayerID, currentCard, noOfPlayers)
					currentCardIsRoundWinner, currentRoundWinnerID, currentRoundWinnerName, currentRoundWinnerCard, currentRoundWinnerPoints = updateRoundsTabForMyPlayedCard(currentPlayerID, currentCard, currentCardPoints, roundNo, playerInRound, currentPlayerTeam, currentPlayerName, currentCardIndex, myCardPlayCondition)
				}
				fmt.Println("You Played : ", currentCard)
				nulPrint2(currentCardIsRoundWinner) //////  REMOVE REMOVE REMOVE REMOVE ///////////////////////////
			}
			printCardsPlayedInRound(roundNo, currentRoundWinnerCard, trumpSuite, roundSuite, friendPoints, foePoints, currentRoundWinnerID)
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
			friendPoints++
			if roundNo < cardsPerPlayer-1 {
				fmt.Println("**** YOUR TEAM WON THE ROUND. PlayerID : ", previousRoundWinner, " (", currentRoundWinnerName, ") STARTS THE NEXT ROUND ***")
				fmt.Println("So far, Your Team won     : ", friendPoints, " rounds")
				fmt.Println("So far, Opposing Team won : ", foePoints, " rounds")
			}
		} else {
			foePoints++
			if roundNo < cardsPerPlayer-1 {
				fmt.Println("**** YOUR TEAM LOST THE ROUND. PlayerID : ", previousRoundWinner, "(", currentRoundWinnerName, ") STARTS THE NEXT ROUND ***")
				fmt.Println("So far, Your Team won     : ", friendPoints, " rounds")
				fmt.Println("So far, Opposing Team won : ", foePoints, " rounds")
			}
		}
	} // rounds in a game end
	printGameRounds()
	if gameResult < 0 {
		fmt.Println("**** Your Team LOST THE GAME ******************")
	} else if gameResult > 0 {
		fmt.Println("**** Your Team WON THE GAME ******************")
	} else {
		fmt.Println("Your Team won     : ", friendPoints, " rounds")
		fmt.Println("Opposing Team won : ", foePoints, " rounds")
		if friendPoints > foePoints {
			fmt.Println("**** Your Team WON THE GAME ******************")
		} else if friendPoints < foePoints {
			fmt.Println("**** Your Team LOST THE GAME ******************")
		} else {
			fmt.Println("******** GAME DRAW ******************")
		}
	}
	db.Close()
}
