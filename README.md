# Quake log parser

This repository was developed to solve two tasks

 - **Task 1** : Build a parser capable of reading a log file "games.log" e expose an API to get that information.
- **Task 2** : Expose a consultation method that response an information related of each game.

## Problems found

 - Every user that reconnect to a match could change its ID in the log file.
- There is an identification called "*InitGame*" that sometimes happens before a "*ShutdownGame*" (which represents, in my opinion, the end of each game). For that reason I decided to consider  "*ShutdownGame*" being the end of each game. After my parser find that key in some line it resets the game detail counting.
- Game execution times with no pattern.

### Download and Install

#### Binary Distributions

Official binary distributions are available at https://golang.org/dl/.

After downloading a binary release, visit https://golang.org/doc/install in your web browser for installation
instructions.

### Dependencies

There was not used any dependency out of regular golang packages.

### How to use it?

Postman collection to this app: https://www.getpostman.com/collections/769f4e9d5c1c2f7092de

#### Task 1

In order to get all games in the game log you must do:

```
GET /get 
Host: {{url}}:9000
```
#### Task 2

If you want a specific game information

```
POST /post 
Host: {{url}}:9000
Content-Type: application/json
{
    "game":"game_8"
}
```





Any doubt please let me know.