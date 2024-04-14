package controller

import "github.com/jeffthorne/tasky/database"

var db = database.CreateDBClientFromEnv()
