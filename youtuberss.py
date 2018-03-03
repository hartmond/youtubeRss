import sys
import os
import configparser
import httplib2
from ttrss.client import TTRClient
from apiclient.discovery import build
from oauth2client.file import Storage

# check if config file exists

if not os.path.exists('youtuberss.conf'):
    print('Error! No conifg file found.')
    sys.exit()

# read config file

conf = configparser.ConfigParser()
conf.read('youtuberss.conf')

# check if oauth credentials exists

if not os.path.exists(conf['yt']['credentials_file']):
    print('Error! No OAuth credentials found. Run setup.py first')
    sys.exit()

# fetch current Youtube Subscriptions from tt-rss

ttrss = TTRClient(conf['tt-rss']['url'], conf['tt-rss']['user'],
                  conf['tt-rss']['password'])
ttrss.login()

if not ttrss.logged_in():
    print("Error logging in on TTRSS")
    sys.exit()

categories = ttrss.get_categories()

youtubeCatID = -1

for category in categories:
    if category.title == 'YouTube':
        youtubeCatID = category.id
        break

if youtubeCatID == -1:
    print('No YouTube Category')
    print('Please create Category with name YouTube')
    sys.exit()

lst_ttrss = set()

ttrssfeeds = ttrss.get_feeds(cat_id=youtubeCatID)

for feed in ttrssfeeds:
    lst_ttrss.add((feed.feed_url, feed.title))


# fetch current Youtube Subscriptions from YouTube-API

yt_api_credentials = Storage(conf['yt']['credentials_file']).get()
yt = build('youtube', 'v3',
           http=(yt_api_credentials.authorize(httplib2.Http())))

lst_yt = set()

yt_req = yt.subscriptions().list(part='snippet', mine=True, maxResults=5)

while yt_req is not None:
    yt_data = yt_req.execute()

    for feed in yt_data['items']:
        lst_yt.add(('https://www.youtube.com/feeds/videos.xml?channel_id='
                    + feed['snippet']['resourceId']['channelId'],
                    feed['snippet']['title']))

    yt_req = yt.subscriptions().list_next(yt_req, yt_data)

# subscribe

for feed in lst_yt.difference(lst_ttrss):
    ttrss.subscribe(feed[0], youtubeCatID)

# unsubscribe
if lst_ttrss:
    id_lookup = dict((f.feed_url, f.id) for f in ttrssfeeds)

for feed in lst_ttrss.difference(lst_yt):
    ttrss.unsubscribe(id_lookup[feed[0]])
