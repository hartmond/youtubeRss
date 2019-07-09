# youtubeRss

Even if there are tons of crap on YouTube nowadays there are a few channels which I like and want to watch.

Accessing my subscriptions through the YouTube homepage does not work for me because I want to see every video and I am not sure if the subscription page shows my every video, especially when I do not watch videos for over a week.

Because of that I access my subscriptions through my Feed-Reader. I use Tiny Tiny RSS.

Fortunately, YouTube provides RSS feeds for channels. But I do not want to search an address and add this manually to tt-rss when I subscribe to a new channel. I use the subscribe buttons from YouTube. This script updates the subscriptions in tt-rss.

## Usage with docker
Place the config file in the mounted volume and refence a youtube-credentials file in the relative folder conf/.
