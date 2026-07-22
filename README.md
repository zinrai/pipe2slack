# pipe2slack

Reads standard input and posts it to a Slack Incoming Webhook.

Pipe the output of any command to it and the collected output is sent to Slack. Implemented with the Go standard library only (no external dependencies).

## Configuration

The webhook URL is read from a plain-text file containing only the URL.

By default the file is `$XDG_CONFIG_HOME/pipe2slack/webhook_url` (falling back to `~/.config/pipe2slack/webhook_url`):

```
$ mkdir -p ~/.config/pipe2slack
$ echo "https://hooks.slack.com/services/XXXX/YYYY/ZZZZ" > ~/.config/pipe2slack/webhook_url
```

Use `-config` to read the URL from any other path. This is the way to run under cron and other environments where `HOME` is not set, because the default location depends on it:

```
$ echo "hello" | pipe2slack -config /etc/pipe2slack/webhook_url
```

## Usage

Pipe the output of another command to it:

```
$ echo "Hello World" | pipe2slack
```

Empty input is skipped without posting.

## License

This project is licensed under the [MIT License](./LICENSE).
