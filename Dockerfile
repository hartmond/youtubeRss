FROM python:3

RUN pip install ttrss-python google-api-python-client oauth2client

ADD youtuberss.py /root/
ADD setup.py /root/

VOLUME /root/conf

WORKDIR /root/

CMD ["python", "youtuberss.py"]

