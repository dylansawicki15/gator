# gator

A multi-user RSS feed aggregator CLI built in Go. Add feeds, follow them, and browse posts — all from your terminal.

## Prerequisites

- [Go](https://go.dev/doc/install) 1.21+
- [PostgreSQL](https://www.postgresql.org/download/) 15+

## Installation

```bash
go install github.com/dylansawicki15/gator@latest
```

## Configuration

Gator reads its config from `~/.gatorconfig.json`. Create it with your Postgres connection string and (optionally) a starting username:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

### Database setup

Install [goose](https://github.com/pressly/goose) and run the migrations:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
```

## Usage

### User management

```bash
gator register <name>   # create a new user and log in
gator login <name>      # switch to an existing user
gator users             # list all users
```

### Feed management

```bash
gator addfeed "Hacker News" https://news.ycombinator.com/rss   # add a feed and auto-follow it
gator feeds                                                      # list all feeds
gator follow <url>                                               # follow an existing feed
gator unfollow <url>                                             # unfollow a feed
gator following                                                  # list feeds you follow
```

### Aggregation

```bash
gator agg 30s    # fetch feeds every 30 seconds (runs until Ctrl+C)
gator agg 1m     # fetch feeds every minute
```

Leave `agg` running in one terminal while you use other commands in another.

### Browsing posts

```bash
gator browse        # show the 2 most recent posts from your followed feeds
gator browse 10     # show the 10 most recent posts
```

### Other

```bash
gator reset    # delete all users and data (destructive!)
```