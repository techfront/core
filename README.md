# Techfront-core

## Gettting Started

The app requires postgresql just now to bootstrap locally (not Mysql). So make sure you have psql installed. The bootstrap process will create a database and settings for you, but you'll need to promote the first user to admin in order to use the site locally.

Go get this app:

    go get -u github.com/techfront/core

Then to build and run the server locally, as you'd expect:

    go run server.go

## Commands

Deploy production server:

    go get -u github.com/fragmenta/fragmenta
    fragmenta deploy production
    
Running generic server:
    
    systemctl start techfront-server
    
Running telegram-service:
    
    systemctl start techfront-telegram-bot-server
    
Running imageproxy-service with docker:
    
    docker run -d -p 8080:8080 techfront/imageproxy -addr 0.0.0.0:8080