package main

import (
	cfg "github.com/conductorone/baton-procore/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("procore", cfg.Config)
}
