## Gator is rss feed colector. It's a guided projects from boot.dev

### Requirements
 - Go ver 1.23.2
 - Postgres installed localy or via docker

### Instalation
 - clone this repo
 - create DB for this project in Postgres
 - create a config file `.gatorconfig.json` inside your home directory
```
{
    "db_url":"",
    "current_user_name":""
}
```
 - fill `db_url` field with a connecting string for your DB
 - run `go install`

### Usage
Run `gator <command>`

List of commands:
 - `register`
 - `login`
 - `reset`
 - `users
 - `agg`
 - `addfeed`
 - `feeds`
 - `follow"
 - `following`
 - `unfollow`
 - `browse`

