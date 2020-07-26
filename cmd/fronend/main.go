package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const page = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Chat Example</title>
<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js" integrity="sha384-OgVRvuATP1z7JjHLkuOU7Xw704+h835Lr+6QL9UvYjZE3Ipu6Tp75j7Bh/kR0JKI" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css" integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous">
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js" integrity="sha384-OgVRvuATP1z7JjHLkuOU7Xw704+h835Lr+6QL9UvYjZE3Ipu6Tp75j7Bh/kR0JKI" crossorigin="anonymous"></script>
<script type="text/javascript">
var mainLoopId 

function buildMsg(message) {
	var msg = document.createElement("div")
	msg.innerHTML = "<span> <u>" + message.sender + "</u></span>" + "<p>" + message.body + "</p>" + "<span>" + message.time + "</span>"
	msg.classList.add("msg")
	
	var body = msg.getElementsByTagName("p")[0]
	body.classList.add("msg-body")
	return msg
}

function buildMyMsg(message){
	var msg = buildMsg(message)
	msg.classList.add("darker")
	
	var time = getMsgTime(msg)
	time.classList.add("time-left")
	
	var sender = getMsgSender(msg)
	sender.classList.add("sender-left")
	return msg
}

function buildOtherMsg(message){
	var msg = buildMsg(message)

	var time = getMsgTime(msg)
	time.classList.add("time-right")

	var sender = getMsgSender(msg)
	sender.classList.add("sender-right")

	return msg
}

function getMsgTime(msg) {
	return msg.getElementsByTagName("span")[1]
}


function getMsgSender(msg) {
	return msg.getElementsByTagName("span")[0]
}

function logOthersMsgs(messages) {
    for (var i = 0; i < messages.length; i++) {
		var item = buildOtherMsg(messages[i])
        log.appendChild(item);
	}
}

function logMyMsg(message) {
	var item = buildMyMsg(message)
	log.appendChild(item)
}

function getMsgs() {
	fetch('http://localhost:8001/chat', {method: 'POST', credentials: 'include'}).then(function(response) {
		if (response.status == 200) {
			console.log("gooooood")
		  response.json().then(data => logOthersMsgs(data));
		}
		else if (response.status != 204){ 
			moveToLogin(mainLoopId)(response)
		}
	})
}

function setPageState(state){
    switch(state) {
		case "login":
			document.getElementById("chat").style.display = "none";
			document.getElementById("login").style.display = "block";
            break
		case "chat":
			document.getElementById("login").style.display = "none";
			document.getElementById("chat").style.display = "block";
            break
		}
}


function getCookie(cname) {
  var name = cname + "=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var ca = decodedCookie.split(';');
  for(var i = 0; i <ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

function signin(){
	  const form = document.getElementById("form-login");
	  var credentials = {username: form.userName.value, password: form.userPassword.value}
	  
	  // Post data using the Fetch API
	   fetch('http://localhost:8002/signin', {
		method: 'POST',
		body: JSON.stringify(credentials),
        credentials: 'include',
		}
	  ).then(function(response){
			if (response.status == 200){
					setPageState("chat")
					mainLoopId = setInterval(getMsgs, 1000)
				}
			else {
			form.userName.value = ""
			form.userPassword.value = ""
		}
		}
		)
		return false
	}


function moveToLogin(loopId) {
	return function(response) {
		if (response.status == 403) {
			setPageState("login")
			clearInterval(loopId)

		}
	}
}

function sendMsg(){
	  const form = document.getElementById("form-msg")
	  var msg = form.msg.value
	  form.msg.value = ""
	  logMyMsg(JSON.parse(msg))
	  
	  // Post data using the Fetch API
	  fetch('http://localhost:8000/send', {
		method: 'POST',
		body: msg,
		mode: 'no-cors',
		credentials: 'include'
		}
	  ).then(moveToLogin(mainLoopId))
	return false

}

window.onload = function () {
	setPageState("login")
	
	
}
</script>
<style type="text/css">
html {
    overflow: hidden;
}

body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}

#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}

#form-msg {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}

.msg {
  border: 2px solid #dedede;
  background-color: #f1f1f1;
  border-radius: 5px;
  padding: 10px;
  margin: 10px 0;
}

/* Darker chat container */
.darker {
  border-color: #ccc;
  background-color: #ddd;
}

/* Style time text */
.time-right {
  float: right;
  color: #aaa;
}

/* Style time text */
.time-left {
  float: left;
  color: #999;
}


.msg-body{
  margin-top: 1rem;
}    

.sender-left{
  color: cadetblue;
  float: left;
}

.sender-right{
  color: red;
  float: right;
}


#overlay {
  position: absolute;
  display: block;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0,0,0,0.5);
}

#form-login {
    background-color: #EDEDED;
    padding-top: 10px;
    padding-bottom: 20px;
    padding-left: 20px;
    padding-right: 20px;
    border-radius: 15px;
    border-color:#d2d2d2;
    border-width: 5px;
    box-shadow:0 1px 0 #cfcfcf;
}

h4 { 
 border:0 solid #fff; 
 border-bottom-width:1px;
 padding-bottom:10px;
 text-align: center;
}

.form-control {
    border-radius: 10px;
}

.wrapper {
    text-align: center;
}
</style>
</head>
<body>
<div id="login" class="container">
    <div class="row">
        <div class="col-md-offset-5 col-md-3">
            <form id="form-login">
            <h4>Welcome back.</h4>
            <input type="text" id="userName" class="form-control input-sm chat-input" placeholder="username" />
            </br>
            <input type="text" id="userPassword" class="form-control input-sm chat-input" placeholder="password" />
            </br>
            <div class="wrapper">
            <button type="submit" onclick="return signin()">send</button>

            </div>
            </form>
        
        </div>
    </div>
</div>
<div id="chat">
<div id="log"></div>
<form id="form-msg">
    <button type="submit" onclick="return sendMsg()">send</button>
    <input type="text" id="msg" size="64" autofocus />
</form>
</div> 
</body>
</html>
`

func serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, page)
}

func main() {
	http.HandleFunc("/", serve)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
