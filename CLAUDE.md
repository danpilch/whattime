# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go application that displays coworkers' current times by fetching timezone information from Slack. Features a beautiful terminal UI built with Charm's Bubble Tea framework.

## Development Commands

- `go run .` - Run the application
- `go build -o whattime .` - Build binary
- `go mod download` - Install dependencies

## Project Structure

- `main.go` - Main application with Bubble Tea TUI implementation
- `slack.go` - Slack API integration and timezone handling
- `go.mod` - Go module dependencies

## Key Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/slack-go/slack` - Slack API client

## Configuration

Requires Slack API token via environment variables:
- `SLACK_BOT_TOKEN` (recommended) - Bot token with `users:read` scope
- `SLACK_USER_TOKEN` - Alternative user token

## Architecture

The application uses the Bubble Tea pattern with:
- Model-View-Update architecture
- Real-time updates via ticker commands
- Async data loading from Slack API
- Table-based UI for displaying timezone information