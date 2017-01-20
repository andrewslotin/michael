![Let's Get Ready To Rumble!](../master/docs/michael-buffer.jpg)

Michael
=======

[![Build Status](https://travis-ci.org/andrewslotin/michael.svg?branch=master)](https://travis-ci.org/andrewslotin/michael)

Announce deploys in Slack channels.

Installation
------------

First you need to add a new slash command for your team:

1. Go to [Custom Integrations](https://api.slack.com/custom-integrations) and click on "New Command"
2. In the "Command" field type in <kbd>/deploy</kbd> — this will be your new Slack command to start, finish and list deploys in channel
3. Fill in "URL" field with an URL where `michael` is deployed
4. Set "Method" to `POST`
5. Copy the content of "Token" field, this will be needed to authenticate incoming requests

You may also like to customize name, icon and include this command into autocomplete list.

Now the server part. To execute a slash command Slack sends a request to an URL associated with it and outputs the response.
Slack requires this URL to be an HTTPS, so you will need some reverse-proxy (such as `nginx`, `caddy`, etc.) configured to serve HTTPS requests.
Self-signed certificates won't work for Slack, [Letsencrypt](https://letsencrypt.org) might be a good option to obtain a real one.

Now compile and run the server (assuming that you have `go` installed):

```
go get github.com/andrewslotin/michael
SLACK_TOKEN=<token you copied before> $GOPATH/bin/michael
```

This will run a server listening on `0.0.0.0:8081`. Check `$GOPATH/bin/michael --help` to see available options.

Optionally you may provide your [GitHub personal access token](https://github.com/settings/tokens) with `repo` permissions by
setting `GITHUB_TOKEN` environment variable. This token is used to get PR details (title, description and author) and attach them to an announcement.
If no token is provided only public pull requests will be have detailed information, and others will only contain a link to GitHub.

Usage
-----

Deploys are tracked per channel. This means that different channels can run different deploys at the same time.

* <kbd>/deploy status</kbd> — see if there is a deploy currently running.
    <img src="../master/docs/deploy-status.png" alt="Deploy status response" height="50">
* <kbd>/deploy &lt;subject&gt;</kbd> — initiate a deploy in the channel. <subject> is an arbitrary string describing what's being deployed.
    <img src="../master/docs/deploy-start.png" alt="Deploy announcement" height="134">

    If there is already a deploy announced by another user in this channel, it needs to be finished first.
    <img src="../master/docs/deploy-running.png" alt="Deploy already started message" height="54">
    
    However if you already initiated a deploy the channel, you can update its subject by executing this command again.
* <kbd>/deploy done</kbd> — finish current deploy.

    <img src="../master/docs/deploy-done.png" alt="Deploy completion announcement" height="44">
    
    You can also finish a deploy started by another user.
    
    <img src="../master/docs/deploy-finish-other.png" alt="Complete unfinished deploy" height="42">

### Deploy status in channel topic

In addition to announcing deploys in channel you may find it useful to have a small sign in the channel topic. This way you can quickly check
if it's safe to deploy. Slack deploy command uses :white_check_mark: and :no_entry: to mark channel as clear for deployment and show that there
is a deploy in progress. To use this feature you need to provide [Slack Web API token](https://api.slack.com/docs/oauth-test-tokens) in
`SLACK_WEBAPI_TOKEN` environment variable and add either `:white_check_mark:` or `:no_entry:` to the channel topic. Whenever the deploy changes
the deploy bot will swap these emojis.

<img src="../master/docs/topic-deploy.png" alt="Channel topic notification" height="270">

To disable this feature without re-deploying the whole service simply remove emojis from channel topic.

### Persistent deploy statuses

To keep the deploy status between service restarts you might want to use built-in BoltDB database. To do this you need to specify the path to
your BoltDB file in `BOLTDB_PATH` environment variable.

```go
BOLTDB_PATH=/path/to/your/bolt.db $GOPATH/bin/michael
```

### Deploy history

To see the history of deploys in channel run <kbd>/deploy history</kbd> in this channel and click the link returned by bot.

<img src="../master/docs/deploy-history.png" alt="Channel history link" height="52">

This will open a page in your browser with all deploys that were ever announced in this channel.

```
* suddendef was deploying https://github.com/andrewslotin/michael/pull/15 since 24 Aug 16 20:54 UTC until 24 Aug 16 20:54 UTC
* suddendef was deploying https://github.com/andrewslotin/michael/pull/15 https://github.com/andrewslotin/michael/pull/11 since 24 Aug 16 20:54 UTC until 24 Aug 16 20:55 UTC
* suddendef was deploying history since 25 Aug 16 08:35 UTC until 25 Aug 16 08:35 UTC
* suddendef was deploying https://github.com/andrewslotin/michael/pull/19 since 25 Aug 16 08:35 UTC until 25 Aug 16 08:35 UTC
```

#### Authorization and authentication

While handling the <kbd>/deploy history</kbd> command deploy bot generates a one-time token that grants access to current channel
deploy history. This access is being granted for the next 30 days and can be renewed at any time by requesting and opening a link
to history in the same channel.

Deploy bot uses JSON Web Tokens (JWT) to store channel access lists. A secret key to sign JWT can be set via `HISTORY_AUTH_SECRET`
environment variable. If there was no secret provided, deploy bot generates a random string and writes it into the log. On next
start you should use this string as a value for `HISTORY_AUTH_SECRET`, otherwise all issued authorizations will be revoked.

Why Michael?
------------

Because Buffer might be not the best name for such tool.

License
-------

This software is distributed under LGPLv3 license. You can find the full text in [LICENSE](../master/LICENSE).
