# Under a Rock

## Client HTTP/REST API for UnderARock


### The Point

You are to write a command interpreter using the provided SimpleShell class. You're going to create a way for commands to be typed
into your shell, read the typed commands and arguments, send them off to the Under-A-Rock server using a REST API, read the JSON data returned from the URL call, and print it out nicely formatted for your user. 

Undar-A-Rock acts a little (very little) like a twitter server or chat server.

The Under-A-Rock Server can be reached at `http://zipcode.rocks:8085` 

There are two segments to the API, the ID segment and the Messages segment.

## IDs

#### ID commands in shell
In the shell, 
`ids` should return a formatted list of the IDs available to you.

`ids your_name your_github_id` command should post your Name and your GithubId to the server.
If you do this twice with two different Names, but the name GithubId, the name on the server gets changed.

the IDs API is:

#### URL: /ids/
* `GET` : Get all github ids registered
* `POST` : add your github id / name to be registered
* `PUT` : change the name linked to your github id

json payload for /ids/ - this is a sample
```json
{
    "userid": "-", // gets filled w id
    "name": "Zipcode",
    "githubid": "zipcoder"
}
```
 
## Messages

#### Message comands in shell

in the shell, 
* `messages` should return the last 20 messages, nicely formatted.
* `messages your_github_id` should return the last 20 messages sent to you.
* `send your_github_id 'Hello World' ` should post a new message in the timeline
* `send your_github_id 'my string message' to some_friend_githubid` should post a message to your friend from you on the timeline.

the Messages API is:

#### URL: /messages/
* `GET` : Get last 20 msgs - returns an JSON array of message objects

#### URL: /ids/:mygithubid/messages/ 
* `GET` : Get last 20 msgs for myid  - returns an JSON array of message objects, E.G. `/ids/xt0fer/messages/`
* `POST` : Create a new message in timeline - need to POST a new message object, and will get back one with a message sequence number and timestamp of the server inserted. In one case of `send`, you don't have a `toid`, just a `fromid`. In the other send/to case, you have both `toid` and `fromid`.

#### URL: /ids/:mygithubid/messages/:sequence
* `GET` : Get msg with a sequence  - returns a JSON message object for a sequence number, E.G. `/ids/xt0fer/messages/472830272bac8d`

#### URL: /ids/:mygithubid/from/:friendgithubid
* `GET` : Get last 20 msgs for myid from friendid, E.G. `/ids/xt0fer/from/zipcoder`

json payload for /meassages/ these are samples, one to a specific friend, one to the timeline.
```json
{
    "sequence": "-",
    "timestamp": "_",
    "fromid": "xt0fer",
    "toid": "zipcoder",
    "message": "Hello, ZipCode!"
}

{
    "sequence": "-",
    "timestamp": "_",
    "fromid": "xt0fer",
    "toid": "",
    "message": "Hello, World!"
}
```
