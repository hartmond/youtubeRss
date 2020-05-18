# youtubeRss

**The old version written in Python and working with Tiny Tiny RSS can still be found under the tag v1.0.0 in this repository**

Even if there are tons of crap on YouTube nowadays there are a few channels which I like and want to watch.

Accessing my subscriptions through the YouTube homepage does not work for me because I want to see every video and I am not sure if the subscription page shows my every video, especially when I do not watch videos for over a week.

Because of that I access my subscriptions through my Feed-Reader. Currently, I use Miniflux.

Fortunately, YouTube provides RSS feeds for channels. But I do not want to search an address and add this manually to my feedreader when I subscribe to a new channel. I use the subscribe buttons from YouTube. This program updates the subscriptions in miniflux each hour.

## Setup
- Create an Project in the [Google Developers Console](https://console.developers.google.com/) to interact with YouTube
- Enable the YouTube Data Api v3 for this project
- Create Credentials for an Desktop Application
- Download the Credential-JSON and place it in the working directory with the name ``youtube_secret.json``
- Create the file ``miniflux_secret.json`` and enter URL and API key (example file is in this repo)
- Run the application. The Application will interactively setup oauth permissions for a user and the user token to the file ``youtube_user.json``
- If all three json files are present the appliation can fron now on run non-interactively.
