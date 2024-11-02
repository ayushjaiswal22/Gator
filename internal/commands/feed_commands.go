package commands

import(
    "fmt"
    "errors"
    "encoding/xml"
    "io"
    "log"
    "strings"
    "strconv"
    "database/sql"
    "html"
    "net/http"
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/ayushjaiswal22/gator/internal/database"
)

/*
func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
    
}
*/

const RSSUrl = "https://www.wagslane.dev/index.xml"

// RSS Feed Struct
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Fetches All RSS Feeds from the RSSUrl
func FetchFeed(s *State, url string, feed_id uuid.UUID) error {
    ctx := context.Background()
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err!=nil {
        return err
    }
    client := http.Client{Timeout: time.Second * 10}
    req.Header.Add("User-Agent", "gator")
    res, err := client.Do(req)
    if err!=nil {
        return err
    }
    defer res.Body.Close()
    data, err := io.ReadAll(res.Body)
    var feed RSSFeed
    er := xml.Unmarshal(data, &feed)
    if er!=nil {
        return er
    }

    feedTitle := html.UnescapeString(feed.Channel.Title)
    feedDesc := html.UnescapeString(feed.Channel.Description)
    fmt.Printf("\nTitle: %s\n\n", feedTitle)

    fmt.Println(feed.Channel.Link)
    fmt.Println(feedDesc)
    for _, item := range feed.Channel.Item {
        itemTitle := html.UnescapeString(item.Title)
        itemDesc := html.UnescapeString(item.Description)
        /*
        fmt.Printf(" * %s\n", itemTitle)
        fmt.Println(item.Link)
        fmt.Println(itemDesc)
        fmt.Println(item.PubDate)
        */
        publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

        postParams := database.CreatePostParams{
            ID: uuid.New(),
            CreatedAt: time.Now().UTC(),
            Title: itemTitle,
            Url: item.Link,
            Description: itemDesc,
            PublishedAt: publishedAt,
            FeedID: feed_id,
        }
        _, err = s.Db.CreatePost(ctx, postParams)
        if err!=nil {
            if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
        }
    }
    log.Printf("Feed %s collected, %v posts found", feed.Channel.Title, len(feed.Channel.Item))
    return nil

}


// Adds or Posts an RSS Feed, and inserts it into the feeds table in DB
func AddFeed(s *State, cmd Command) error {
    if len(cmd.Args) < 2 {
        return errors.New("Not enough arguements to create a feed.")
    }
    ctx := context.Background()
    id, err := s.Db.GetUserID(ctx, s.Conf.CurrentUsername)
    if err!=nil {
        return err
    }


    params := database.CreateFeedParams{
        ID: uuid.New(), 
        CreatedAt: time.Now().UTC(), 
        UpdatedAt: time.Now().UTC(), 
        Name: cmd.Args[0], 
        Url: cmd.Args[1],
        UserID: id, 
    }
    feed, er := s.Db.CreateFeed(ctx, params);
    if er!=nil {
        return er
    }

    followParams := database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID: id,
        FeedID: feed.ID,
    }
    
    _, err = s.Db.CreateFeedFollow(ctx, followParams)
    if err!=nil{
        return err
    }
    return nil

}

// Gets all the user defined feeds from the feeds table
func GetFeeds(s *State, cmd Command) error {
    ctx := context.Background()
    feeds, err := s.Db.GetAllFeeds(ctx)
    if err!=nil {
        return err
    }
    for _, feed := range feeds {
        fmt.Println(feed.Name)
        fmt.Println(feed.Url)
        name, err := s.Db.GetUsername(ctx, feed.UserID)
        if err !=nil {
            return err
        }
        fmt.Println(name)
        fmt.Println()
    }
    return nil
}

func ScrapeFeed(s* State, cmd Command) error {
    if len(cmd.Args) < 1 {
        return errors.New("Not enough arguements to fetch feed.")
    }
    ctx := context.Background()
    duration, err := time.ParseDuration(cmd.Args[0])
    if err!=nil {
        return err
    }
    ticker := time.NewTicker(duration)
    fmt.Printf("Collecting feeds every %s...\n", cmd.Args[0])
    for ; ; <-ticker.C {
        feed,err := s.Db.GetNextFeedToFetch(ctx)
        if err!=nil {
            return err
        }
        t := sql.NullTime{Time:time.Now().UTC(), Valid:true}
        params := database.MarkFeedFetchedParams{LastFetchedAt:t, ID:feed.ID}
        err = s.Db.MarkFeedFetched(ctx, params)
        if err!=nil{
            return err
        }
        err = FetchFeed(s, feed.Url, feed.ID)
        if err!=nil{
            return err
        }
    }
    return nil
}

func ListPosts(s* State, cmd Command) error {
    limit := 2
    var err error
    if len(cmd.Args) == 1 {
        limit, err = strconv.Atoi(cmd.Args[0])
        if err!=nil {
            return err
        }
    }
    ctx := context.Background()
    posts, err := s.Db.GetPostsUser(ctx, int32(limit))
    if err!=nil {
        return err
    }
    for _, post := range posts {
        fmt.Println("===============================================================")
        fmt.Printf("-----%s-----\n", post.Title)
        fmt.Println("===============================================================")
        fmt.Printf("Description: %s\n", post.Description)
        fmt.Printf("Link: %s\n", post.Url)
        fmt.Printf("Published: %s\n\n", post.PublishedAt)
        fmt.Println("=====================================")

    }
    
    return nil
}
