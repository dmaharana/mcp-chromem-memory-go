package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	// Initialize the memory store
	store, err := NewMemoryStore("./memory.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize memory store")
	}
	defer store.Close()

	// Initialize MCP server
	server := NewMCPServer(store)
	
	log.Info().Msg("Starting MCP Memory Server")
	
	// Start stdio server
	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start MCP server")
	}
}