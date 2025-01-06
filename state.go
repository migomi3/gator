package main

import (
	"github.com/migomi3/gator/internal/config"
	"github.com/migomi3/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
