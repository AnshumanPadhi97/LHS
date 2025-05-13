package main

import (
	"LHS/backend"
	docker "LHS/backend/Processors/Docker"
	"LHS/backend/Processors/db"
	"LHS/backend/models"
	"context"
	"log"
)

//WAILS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	if err := docker.InitDockerClient(); err != nil {
		log.Fatalf("Failed to init Docker client: %v", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	_ = docker.CloseClient()
}

// BACKEND MANAGEMENT

func (a *App) BuildStack(content string) error {
	tmpl, err := models.ParseStackYAML([]byte(content))
	if err != nil {
		return err
	}
	return backend.BuildStack(tmpl)
}

func (a *App) BuildStackByStackId(stackId int64) error {
	err := backend.BuildStackFromDB(stackId)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) RunStackById(stackID int64) error {
	return backend.RunStack(stackID)
}
