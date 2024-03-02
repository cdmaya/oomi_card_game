# oomi_card_game
This program sketch is designed to facilitate gameplay of the Sri Lankan card game, Oomi, accommodating 4, 6, or 8 players.
The flexibility extends to the variable number of cards utilized in the game, ranging from 2 to ace from each suite, as per the requirements. Each player can hold a maximum of 8 cards.

![oomi_card_game_setup](https://github.com/cdmaya/oomi_card_game/assets/51185952/b08229b3-5167-4948-81df-bd79f680d496)


Currently, the program operates in a text-based format, offering detailed prompts during runtime to guide users through the gameplay. Developed in the Go language, it utilizes an SQLite database to temporarily store, sort, and retrieve the necessary values during gameplay.

For Debian-based Linux hosts, the following commands can be used to install the required packages:

    sudo apt install sqlite3 libsqlite3-dev gcc golang-go git -y
    go get github.com/mattn/go-sqlite3
    
Key features of the program include:

The program excels in rounds where it starts the round or ends the round, providing status updates and concluding rounds appropriately.
In other scenarios, the program currently plays cards with the objective of winning the round, without taking into account the strategies of other players on the same team or the opposing team. Future enhancements may involve implementing a more sophisticated logical flow or integrating a Neural Network to enhance strategic decision-making.

Feel free to explore the code and adapt it to your needs. If running on a different operating system, ensure compatibility by checking the Golang and SQLite drivers accordingly. Enjoy your Oomi gameplay!




