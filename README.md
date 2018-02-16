# GRedditNotifier
_Sends notifications to your phone when new posts are made in subreddit(s) you subscribe to_

[![Build Status](https://travis-ci.org/Jonatino/GRedditNotifier.svg?branch=master)](https://travis-ci.org/Jonatino/GRedditNotifier)
![license](https://img.shields.io/github/license/Jonatino/GRedditNotifier.svg)

This library is licensed under the MIT license.


GRedditNotifier is an open source Go application that provides new post push notifications straight to your mobile device from subreddit(s) you subscribe to.

# How Can I Use GRedditNotifier?

First you will need to install the Pushbullet App on your device.

One you have that app installed and an account created, generate an API key for Pushbullet here: https://www.pushbullet.com/#settings/account

Enter your API key and Reddit username into the GRedditNotifier.json config file and add any subreddit(s) you want to be notified of new posts in the GRedditNotifier.json config file.

Once you click run, the program will run in the background and check your subscribed subreddit(s) every X (number of configured seconds) seconds. If new posts are found, you will get a notification on your phone including the title, time, link, and the author of the post (See screenshot below).



---


# Screenshots

![Alt text](https://dl.dropboxusercontent.com/s/41m5b6tl3kjriyd/nomacs_2018-02-16_16-44-50.png "Notification Demo")
