Slack Deploy Command
====================

Maintain deploy announcements in Slack channels.

Installation
------------

First you need to add a new slash command for your team:

1. Go to [Custom Integrations](https://api.slack.com/custom-integrations) and click on "New Command"
2. In the "Command" field type in <kbd>/deploy</kbd> — this will be your new Slack command to start, finish and list deploys in channel
3. Fill in "URL" field with an URL where `slack-deploy-command` is deployed
4. Set "Method" to `POST`
5. Copy the content of "Token" field, this will be needed to authenticate incoming requests

You may also like to customize name, icon and include this command into autocomplete list.

Now the server part. To execute a slash command Slack sends a request to an URL associated with it and outputs the response.
Slack requires this URL to be an HTTPS, so you will need some reverse-proxy (such as `nginx`, `caddy`, etc.) configured to serve HTTPS requests.
Self-signed certificates won't work for Slack, [Letsencrypt](https://letsencrypt.org) might be a good option to obtain a real one.

Now compile and run the server (assuming that you have `go` installed):

```
go get github.com/andrewslotin/slack-deploy-command
SLACK_TOKEN=<token you copied before> $GOPATH/bin/slack-deploy-command
```

This will run a server listening on `0.0.0.0:8081`. Check `$GOPATH/bin/slack-deploy-command --help` to see available options.

Usage
-----

Deploys are tracked per channel. This means that different channels can run different deploys at the same time.

* <kbd>/deploy status</kbd> — see if there is a deploy currently running.
* <kbd>/deploy &lt;subject&gt;</kbd> — initiate a deploy in the channel. <subject> is an arbitrary string describing what's being deployed.

    If there is already a deploy announced by another user in this channel, it needs to be finished first. However if you already initiated a deploy the channel, you can update its subject by executing this command again.
* <kbd>/deploy done</kbd> — finish current deploy, you can also finish a deploy started by another user.

License
-------

This software is distributed under LGPLv3 license. You can find the full text in [LICENSE](../master/LICENSE).
