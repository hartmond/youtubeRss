import sys
import configparser
from ttrss.client import TTRClient

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
    lst_ttrss.append({'title':((f.title).encode("utf8")), 'url':f.feed_url})
