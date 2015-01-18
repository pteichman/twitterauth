This is a very quick & dirty HTTP server for authorizing accounts
against a Twitter app. It's made to deploy to Heroku but would be easy
to run elsewhere.

    $ heroku create -b https://github.com/kr/heroku-buildpack-go.git

Create a new application at Twitter: https://apps.twitter.com/

When asked for a callback url, take the url from `heroku create` above
and append `/callback`, e.g. `https://ancient-temple-243.herokuapp.com/callback`.

    $ heroku config:set TWITTER_CONSUMER_KEY="your app consumer key"
    $ heroku config:set TWITTER_CONSUMER_SECRET="your app consumer secret"

    $ git push heroku master

Visit `https://[your app]/auth` to authorize. You'll get back your
access token and secret.
