# Google Sheets Updater

## Description

The purpose of this script is to automate some of the mundane tasks that I had to do periodically(weekly basis). 
Say a part of your job on a daily basis was to report a certain piece of information from software `a` to a particular row and column in `Google Sheets`. 
In this case -
1. The software is called `Teamwork Projects`. 
2. The description of the row is `Number of bugs generated in last 7 days`
3. The column title had to be the current week number of the year

###### How this works

1. The script queries the  `Teamwork Projects` API and retrieves info about `Number of bugs generated in the last 7 days`
2. Accesses the Google sheet of the user and finds the row with the description `Number of bugs generated in the last 7 days` 
3. Creates a column that is the current week number from the start of the year, if not already present. 
4. Updates that row and column with the `Number of bugs generated in the last 7 days`

###### Modifications 

For security reasons, I've removed the API token to access `Teamwork Projects` and my personal `client configuration file`. 

## Download size

~15Mb

## Prerequisite

First and foremost you need to have your Go Development environment setup i.e., have your __GOROOT__ and __GOPATH__ set up. If you have this sorted,please move to the section _Configuring the authentication for the script_. If not, please follow the steps below.
1. Download Go from [here](https://golang.org/dl/) and install it
2. Set your environment variable __GOROOT__ to point to _Go_ folder
3. Create a new environment variable __GOPATH__ to point to the folder where you want to store your Go Projects. For the sake of this     tutorial, lets name this folder _Go Projects_. So, your __GOPATH__ should point to the folder _Go Projects_.
4. Navigate into _Go Projects_ and create a folder _src_. Inside _src_ is where all your projects will be stored.
5. Navigate into _src_ and create a folder _github.com_. The reason for this is _go_ requires you to set up your development environment categorically depending upon the source of the project

## What is GOROOT and GOPATH
###### GOROOT
The `GOROOT` environment variable lists the place where to look for the Go binary distribution. The Go binary distributions asume they will be installed in `usr/local/go` (or `C:\Go`in windows) but it is possible to install the `Go Tools`in a different location. In this case you must set the  GOROOT environment variable to point to the directory in which it was installed.

For example, if you installed Go to your home directory you should add the following commands to `$HOME/.profile`:
```
export GOROOT=$HOME/go
export PATH=$PATH:$GOROOT/bin
```
###### GOPATH
The `GOPATH` environment variable lists places to look for Go code. On Unix, the value is a colon-separated string. On Windows, the value is a semicolon-separated string. On Plan 9, the value is a list.

`GOPATH` must be set to get, build and install packages outside the standard Go tree.
More information on `GOPATH` can be found [here](https://golang.org/doc/install#tarball_non_standard)

## Configuring the essential authentication information for the script

1. In order to be able to access the spreadsheet, youre google account should have already been given access to this Google Sheet. 
2. Navigate to the following [link](https://developers.google.com/sheets/api/quickstart/go) and enable _Google Sheets API_ 
3. Download `Client secret configuration` from the pop up and save to the directory where the script is present
4. Place the API token of `Teamwork Projects` in the file `accesstoken.txt` 
5. Once this is done follow the instructions from the script after running it.


## Running The Script

1. Navigate into _github.com_ and clone this repository __gsheets-metrics-updater__
2. Now open a terminal inside __gsheets-metrics-updater__ and type `brew install dep` if on mac or `go get -u github.com/golang/dep/cmd/dep` for windows
3. `dep`is a dependency management tool similar to `npm` but for Go
4. If you're on windows, add `dep` to your global path
5. Open a terminal inside __gsheets-metrics-updater__ and type `dep ensure -v`
6. Now type `go run *.go`