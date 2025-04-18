# Qoutey

A simple, scheduled email service that sends inspirational quotes to your inbox.

## Overview

Qoutey is a lightweight Go application that delivers inspirational quotes to configured email recipients at scheduled times throughout the day (7am, 12pm, and 7pm). It also sends a notification email when the service starts up, helping you verify that the application is running correctly. Qoutey intelligently manages quote selection to avoid repetition until a configurable number of other quotes have been used. This app is designed to run on a linux server

## Features

- Sends inspirational quotes via email at scheduled times
- Sends a startup notification email when the service begins
- Configurable SMTP settings for email delivery
- Customizable quote collection
- Avoids repeating quotes until a configurable number of other quotes have been used
- Runs as a system service with automatic restarts
- Comprehensive logging

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/ad-archer/qoutey.git
   cd qoutey
   ```

2. Create a configuration file by copying the example:
   ```
   cp config.json.example config.json
   ```

3. Edit the configuration file with your email settings:
   ```
   nano config.json
   ```
   
   Update the SMTP server details, email addresses, and customize your quotes collection.

4. Build the application:
   ```
   go build -o qoutey cmd/qoutey/main.go
   ```

## Running as a Service

### On Linux with systemd:

1. Copy the service file to systemd directory:
   ```
   sudo cp qoutey.service /etc/systemd/system/
   ```

2. Update the paths in the service file if necessary.

3. Enable and start the service:
   ```
   sudo systemctl enable qoutey
   sudo systemctl start qoutey
   ```

4. Check the status:
   ```
   sudo systemctl status qoutey
   ```

### On macOS:

1. Use the provided startup script:
   ```
   chmod +x start_qoutey.sh
   ./start_qoutey.sh
   ```

## Test Mode

You can test the email sending functionality without waiting for the scheduled times:

```
./qoutey test
```

This will immediately send a quote email to verify your configuration.

## Configuration Options

The `config.json` file contains several configurable options:

- SMTP server settings (server, port, username, password)
- Email details (from address, recipients, subject)
- List of inspirational quotes
- Maximum repetition setting (how many other quotes must be sent before repeating)

## Logging

Logs are saved to `qoutey.log` in the application directory. Check this file for any issues with email delivery or quote selection.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License
MIT