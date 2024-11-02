package main

import _ "github.com/lib/pq"

import (
    "os"
    "log"
    "database/sql"
    "github.com/ayushjaiswal22/gator/internal/config"
    "github.com/ayushjaiswal22/gator/internal/database"
    "github.com/ayushjaiswal22/gator/internal/commands"
)


func main() {
    if len(os.Args)<2 {
        log.Fatal("ERROR: Not enough arguments")
    }

    // Configuring...
    var conf config.Config
    conf, err := config.Read()
    if err!=nil {
        log.Fatal(err)
    }

    // Connecting to the DB
    db, err := sql.Open("postgres", conf.DbUrl)
    if err!=nil {
        log.Fatal(err)
    }
    defer db.Close()
    dbQueries := database.New(db)


    // Initialising...
    state := commands.State{Conf:&conf, Db: dbQueries}
    commandMap := commands.Commands{CmdMap: make(map[string]func(*commands.State, commands.Command) error)}
    
    // Registering commands...
    commandMap.Register("login", commands.LoginHandler)
    commandMap.Register("register", commands.RegisterUser)
    commandMap.Register("reset", commands.DeleteUsers)
    commandMap.Register("users", commands.GetUsers)
    commandMap.Register("agg", commands.ScrapeFeed)
    commandMap.Register("addfeed", commands.AddFeed)
    commandMap.Register("feeds", commands.GetFeeds)
    commandMap.Register("follow", commands.CreateFollow)
    commandMap.Register("following", commands.GetUserFeedFollow)
    commandMap.Register("unfollow", commands.UnfollowFeed)
    commandMap.Register("browse", commands.ListPosts)


    // Running command
    cmd := commands.Command{Name: os.Args[1], Args:os.Args[2:]}
    er := commandMap.Run(cmd.Name, &state, cmd)
    if er!=nil {
        log.Fatal(er)
    }
    
    
}
