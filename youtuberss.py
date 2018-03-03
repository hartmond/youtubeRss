#!/usr/bin/env python3
# coding: utf-8

"""
connector for watching YouTube subscriptions in Tiny Tiny RSS
uses the APIs of tt-rss and YouTube find new subscribed/unsubscribed
channels and updated subscribed feeds accordingly
"""

import sys
import os
import configparser
import httplib2
from ttrss.client import TTRClient
from apiclient.discovery import build
from oauth2client.file import Storage


class TTRssClient:
    """class to handle the tt-rss api"""
    def __init__(self, url, user, password):
        self.client = TTRClient(url, user, password)

        self.client.login()
        if not self.client.logged_in():
            print("Error logging in on TTRSS")
            sys.exit()

        self.youtube_cat_id = None
        for category in self.client.get_categories():
            if category.title == 'YouTube':
                self.youtube_cat_id = category.id
                break

        if self.youtube_cat_id is None:
            print('No YouTube Category')
            print('Please create Category with name YouTube')
            sys.exit()

        self.id_lookup = None

    def get_feeds(self):
        """read currently subscribed feed from tt-rss"""
        raw_feeds = self.client.get_feeds(cat_id=self.youtube_cat_id)
        self.id_lookup = dict((f.feed_url, f.id) for f in raw_feeds)
        subscriptions = set()
        for feed in raw_feeds:
            subscriptions.add((feed.feed_url, feed.title))
        return subscriptions

    def subscribe(self, feed):
        """subscribe to feed in tt-rss"""
        self.client.subscribe(feed[0], self.youtube_cat_id)

    def unsubscribe(self, feed):
        """unsubscribe from feed in tt-rss"""
        self.client.unsubscribe(self.id_lookup[feed[0]])


class YTConnector:
    """class to handle the youtube api"""
    def __init__(self, credentials_file):
        if not os.path.exists(credentials_file):
            print('Error! No OAuth credentials found. Run setup.py first')
            sys.exit()
        yt_api_credentials = Storage(credentials_file).get()
        authorize = yt_api_credentials.authorize(httplib2.Http())
        self.client = build('youtube', 'v3',
                            http=authorize)

    def get_subscriptions(self):
        """read current subscriptions from youtube api"""
        lst_yt = set()

        yt_req = self.client.subscriptions().list(part='snippet', mine=True,
                                                  maxResults=5)
        while yt_req is not None:
            yt_data = yt_req.execute()
            base = 'https://www.youtube.com/feeds/videos.xml?channel_id='

            for feed in yt_data['items']:
                lst_yt.add((base + feed['snippet']['resourceId']['channelId'],
                            feed['snippet']['title']))

            yt_req = self.client.subscriptions().list_next(yt_req, yt_data)

        return lst_yt


if __name__ == '__main__':
    # read config file
    if not os.path.exists('youtuberss.conf'):
        print('Error! No conifg file found.')
        sys.exit()
    CONF = configparser.ConfigParser()
    CONF.read('youtuberss.conf')

    # initalize tt-rss connector
    TTRSS_URL, TTRSS_USER, TTRSS_PW = CONF['tt-rss'].values()
    TTRSS = TTRssClient(TTRSS_URL, TTRSS_USER, TTRSS_PW)

    # initalize yt connecot
    OAUTH_FILE = CONF['yt']['credentials_file']
    YT = YTConnector(OAUTH_FILE)

    # read acutal state
    FEEDS_TTRSS = TTRSS.get_feeds()
    FEEDS_YT = YT.get_subscriptions()

    # update feed subscriptions in tt-rss
    for subscribe_feed in FEEDS_YT.difference(FEEDS_TTRSS):
        TTRSS.subscribe(subscribe_feed)
    for unsubscribe_feed in FEEDS_TTRSS.difference(FEEDS_YT):
        TTRSS.unsubscribe(unsubscribe_feed)
