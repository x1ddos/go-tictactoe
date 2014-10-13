appengine-endpoints-tictactoe-go
====================================

This application implements a simple backend for a Tic Tac Toe game using
Google Cloud Endpoints, App Engine, and Go.

## Products
- [App Engine][1]

## Language
- [Go][2]

## APIs
- [Cloud Endpoints for Go][3]

## Setup Instructions

1. Make sure to have the [App Engine SDK for Go][4] installed, version
   1.8.0 or higher.
2. Change `'YOUR-CLIENT-ID'` in [`static/js/render.js`][5] and 
   [`api.go`][6] to the respective client ID(s) you have registered 
   in the [APIs Console][7].
3. Update the value of `application` in [`app.yaml`][8] from `go-endpoints` 
   to the app ID you have registered in the App Engine admin console and would 
   like to use to host your instance of this sample.
4. Run the application, and ensure it's running by visiting your local server's
   admin console (by default [localhost:8000][9].)
5. Test your Endpoints by visiting the Google APIs Explorer: 
  [localhost:8080/_ah/api/explorer][10]

Once setup, the app should be accessible at [localhost:8080](http://localhost:8080).

[1]: https://developers.google.com/appengine
[2]: http://golang.org/
[3]: https://github.com/crhym3/go-endpoints#cloud-endpoints-for-go
[4]: https://developers.google.com/appengine/downloads
[5]: https://github.com/crhym3/go-tictactoe/blob/master/app/static/js/render.js
[6]: https://github.com/crhym3/go-tictactoe/blob/master/tictactoe/api.go
[7]: https://console.developers.google.com
[8]: https://github.com/crhym3/go-tictactoe/blob/master/app/app.yaml
[9]: http://localhost:8000/
[10]: http://localhost:8080/_ah/api/explorer
