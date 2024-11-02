package commands

import(
    "fmt"
    "errors"
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/ayushjaiswal22/gator/internal/database"
)

// Inserts a record in feed-follows table, this means user followed a feed
func CreateFollow(s *State, cmd Command) error {
    if len(cmd.Args)<1 {
        return errors.New("Not enough arguements")
    }
    ctx := context.Background()
    
    userUUID, err := s.Db.GetUserID(ctx, s.Conf.CurrentUsername)
    if err!=nil {
        return err
    }

    feedUUID, err := s.Db.GetFeedID(ctx, cmd.Args[0])
    if err!=nil {
        return err
    }

    params := database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID: userUUID,
        FeedID: feedUUID,
    }
    
    feedFollow, err := s.Db.CreateFeedFollow(ctx, params)
    if err!=nil{
        return err
    }
    fmt.Println(feedFollow.FeedName)
    fmt.Println(feedFollow.UserName)
    fmt.Println()
    return nil
}

// Gets all the feeds (feed names) that a user follows
func GetUserFeedFollow(s *State, cmd Command) error {
    ctx := context.Background()
    id, err := s.Db.GetUserID(ctx, s.Conf.CurrentUsername)
    if err!=nil {
        return err
    }
    feedList, err := s.Db.GetFeedFollowsForUser(ctx, id)
    if err!=nil {
        return err
    }
    for _, feed := range feedList {
        fmt.Println(feed)
    }
    return nil
}

func UnfollowFeed(s *State, cmd Command) error {
    if len(cmd.Args) < 1{
        errors.New("Not enough args.")
    }
    ctx := context.Background()

    feedUUID, err := s.Db.GetFeedID(ctx, cmd.Args[0])
    if err!=nil {
        return err
    }
    err = s.Db.DeleteFeedFollow(ctx, feedUUID)
    if err!=nil {
        return err
    }
    return nil
}
