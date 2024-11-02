package commands

import(
    "fmt"
    "errors"
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/ayushjaiswal22/gator/internal/config"
    "github.com/ayushjaiswal22/gator/internal/database"
)

type State struct {
    Db *database.Queries
    Conf *config.Config
}

type Command struct {
    Name string
    Args []string
}

type Commands struct {
    CmdMap map[string]func(*State, Command) error
}

// Logs in a user and sets the user in the the config json
func LoginHandler(s *State, cmd Command) error {
    if len(cmd.Args) < 1 {
        return errors.New("ERROR: a username is required to login.")
    }

    _, er := s.Db.GetUser(context.Background(), cmd.Args[0])
    if er!=nil {
        return errors.New("User does not exist")
    }

    config.SetUser(*s.Conf, cmd.Args[0])
    fmt.Printf("%s has been logged in.\n", cmd.Args[0])
    return nil
}

// Creates a new user in the users table and logs in.
func RegisterUser(s *State, cmd Command) error {
    if len(cmd.Args) < 1 {
        return errors.New("ERROR: a username is required to login.")
    }
    ctx := context.Background()
    _, er := s.Db.GetUser(ctx, cmd.Args[0])
    if er==nil {
        return errors.New("User already exists")
    }
    params := database.CreateUserParams{ID:uuid.New(), CreatedAt:time.Now().UTC(), UpdatedAt:time.Now().UTC(), Name: cmd.Args[0]}
    u, err := s.Db.CreateUser(ctx, params)
    if err!=nil {
        return err
    }
    fmt.Printf("User Created:\nID: %v\nName: %s\nCreate at: %v\nUpdated at: %v\n", u.ID, u.Name, u.CreatedAt, u.UpdatedAt)

    config.SetUser(*s.Conf, cmd.Args[0])
    fmt.Printf("%s has been created.\n", cmd.Args[0])
    return nil
}


// Deletes the user from the users table
func DeleteUsers(s *State, cmd Command) error {
    ctx := context.Background()
    er := s.Db.DeleteUsers(ctx)
    if er!=nil {
        return er
    }
    fmt.Println("All users deleted!")
    return nil
}

// Gets all the users available in the users table
func GetUsers(s* State, cmd Command) error {
    ctx := context.Background()
    userList, er := s.Db.GetUsers(ctx)
    if er!=nil {
        return er
    }
    for _, u := range userList {
        if u.Name == s.Conf.CurrentUsername {
            fmt.Printf("* %s (current)\n", u.Name)
        } else {
            fmt.Printf("* %s\n", u.Name)
        }
    }
    return nil
}
