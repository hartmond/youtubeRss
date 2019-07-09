import os
import sys
import configparser
from oauth2client.client import OAuth2WebServerFlow
from oauth2client.file import Storage

conf_file = "conf/youtuberss.conf"

# check if config file exists

if not os.path.exists(conf_file):
    print('Error! No conifg file found.')
    sys.exit()

# read config file

conf = configparser.ConfigParser()
conf.read(conf_file)

# request access to account and receive credentials

flow = OAuth2WebServerFlow(client_id=conf['yt']['client_id'],
                           client_secret=conf['yt']['client_secret'],
                           scope='https://www.googleapis.com/auth/youtube',
                           redirect_uri='http://httpbin.org/get')

print('Copy the following link to a browser and allow access to your Account')
print('----------')
print(flow.step1_get_authorize_url())
print('----------')
print('Paste the displayed \'code\' parameter here')
auth_code = input()

credentials = flow.step2_exchange(auth_code)

# check if refresh token was received

if credentials.refresh_token is None:

    # inform user to rerun script
    print('!!!!!!!!!!')
    print('No refresh token received!')
    print('Please go to the Settings of your personal Google-Account (not the Developer account)'
          + 'and remove your application from connected applications and rerun the script')
    print('!!!!!!!!!!')

else:

    # save the credentials for the cronjob

    storage = Storage(conf['yt']['credentials_file'])
    storage.put(credentials)

    print('The OAuth credentials were saved.')
