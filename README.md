This project uses go and Postgres to run. You will need to have both installed.

To install gator you will need to run the following command: ```go install gator```

Once installed you will need to create a .gatorconfig.json file at the root of the command line. Then copy and paste the following in the file:
```{"db_url":"postgres://[username]:[password]@localhost:5432/gator?sslmode=disable","current_user_name":""}```
You will need to set up a username and password as you set up your postgres. Once that's done replace both fields in the db_url of the config file. This will tell the program where to find data and tell it that you are authorized to use it. Don't worry about the current_user_name field right not. That can be dealt with once the once gator starts running.

Once everything is running, you'll run commands by entering the following from the project root:
```go run . [command] [args]...```

Here is a list of commands you can run and some descriptions:
    -login: this is what will fill in the current_user_name of the config file. It will error if the user doesn't exist so you will need to       register one first
    -register: Creates a user and takes a string for a name argument.
    -reset: drops all data from tables including users. Made primarily for test purposes
    -users: returns a list of registered users
    -agg: Scrapes feeds at a user input rate based on the oldest viewed first. Takes a time duration argument in the for of #a. For example 5s  would be 5 seconds
    -addfeed: creates a feed by pulling information from a given url. The first arg is a title decided by the user. The second arg is the url. No two feeds can have the same url.
    -feeds: returns all feeds
    -follow: creates an association between a feed and the current user
    -following: returns all the feeds a current user is following
    -unfollow: deletes association between a feed and the current user
    -browse: Returns the oldest x posts from a url. x is an arg representing the number of posts the user wants to see