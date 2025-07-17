# whattime

A beautiful terminal UI application that displays your coworkers' current times by fetching timezone information from Slack.

## Features

- 🌍 Real-time timezone display for all Slack team members
- 🎨 Beautiful terminal UI using Charm's Bubble Tea framework
- 🔄 Auto-refreshing time display (updates every second)
- 📊 Organized table view with status indicators
- ⌨️ Simple keyboard controls

## Prerequisites

- Go 1.19 or higher
- Slack Bot Token or User Token with appropriate permissions

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd whattime
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your Slack token:
```bash
export SLACK_BOT_TOKEN=xoxb-your-bot-token-here
# OR
export SLACK_USER_TOKEN=xoxp-your-user-token-here
```

## Usage

Run the application:
```bash
go run .
```

### Controls

- `r` - Refresh data from Slack
- `q` or `Ctrl+C` - Quit the application

## Slack Setup

To use this application, you'll need a Slack token:

1. **Bot Token (Recommended)**: Create a Slack app and install it to your workspace
   - Required scopes: `users:read`
   - Token format: `xoxb-...`

2. **User Token**: Generate a legacy user token
   - Required scopes: `users:read`
   - Token format: `xoxp-...`

## Display Information

The application shows:
- **Name**: Real name from Slack profile
- **Username**: Slack username
- **Timezone**: User's configured timezone
- **Current Time**: Real-time clock for each user
- **Date**: Current date in user's timezone
- **Status**: Time-based status (Morning, Afternoon, Evening, Night)

## Building

To build a binary:
```bash
go build -o whattime .
```

## Fully vibe coded

Used claude code.
