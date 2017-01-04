import sys
import configparser
import httplib2
from ttrss.client import TTRClient
from apiclient.discovery import build
from oauth2client.file import Storage

#read config file
conf = configparser.ConfigParser()
conf.read('youtuberss.conf')

#fetch current Youtube Subscriptions from tt-rss

ttrss = TTRClient(conf['tt-rss']['url'], conf['tt-rss']['user'], conf['tt-rss']['password'])
ttrss.login()

if ttrss.logged_in() == False:
  print("Error logging in on TTRSS")
  sys.exit()

categories = ttrss.get_categories()

youtubeCatID = -1

for categorie in categories:
  if categorie.title == 'YouTube':
    youtubeCatID = categorie.id
    break

if youtubeCatID == -1:
  print('No YouTube Category')
  print('Please create Category with name YouTube')
  sys.exit()

lst_ttrss = []

ttrssfeeds = ttrss.get_feeds(cat_id=youtubeCatID)

for f in ttrssfeeds:
    lst_ttrss.append({'title':str(f.title), 'url':f.feed_url})

#fetch current Youtube Subscriptions from YouTube-API

yt_api_credentials = Storage(conf['yt']['credentials_file']).get()
yt = build('youtube', 'v3', http=(yt_api_credentials.authorize(httplib2.Http())))

#TODO lower page size and iterate through pages
ytdata = yt.subscriptions().list(part='snippet', mine=True, maxResults=50).execute()

lst_yt = []

for i in ytdata['items']:
  lst_yt.append({'title':i['snippet']['title'], 'url':('https://www.youtube.com/feeds/videos.xml?channel_id=' + i['snippet']['resourceId']['channelId'])})

