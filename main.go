package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var errorsList = [...]int{
	http.StatusInternalServerError,
	http.StatusServiceUnavailable,
	http.StatusBadRequest,
	http.StatusBadGateway,
	http.StatusConflict}

func main() {
	rand.Seed(time.Now().UnixNano())
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/kill", killHandler).Methods("POST")
	r.HandleFunc("/.well-known/health-check", healthCheckHander).Methods("GET")
	r.HandleFunc("/errors/random", randomHandler).Methods("GET")
	r.HandleFunc("/codes/{code}", codeHandler).Methods("GET")

	log.Println("serving port http://localhost:9000")

	err := http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	out := `
		<html>
			<head>
				<title>Error Service</title>
				<style>
					%s
				</style>
			</head>
			<body>
				<div class="header">
					<h1>Error Service</h1>
				</div>
				<div class="content">
					<h2>Shutdown</h2>
					<form method="post" action="/kill">
						<input class="danger-button" type="submit" value="shutdown" />
					</form>
					<h2>Quick Errors</h2>
					<a target="_new" href="/codes/random">random error</a>, <a target="_new" href="/codes/404">error 404</a>, <a target="_new" href="/codes/500">error 500</a>, <a target="_new" href="/codes/503">error 503</a>
					<h2>Error Generator</h2>
					<div class="buttons">
						<input id="button-stop" type="button" value="stop" /> <input id="button-start" type="button" value="start" />
					</div>
				</div>
				<div id="result">
				</div>
				<script>
					%s
				</script>
			</body>
		</html>
	`
	fmt.Fprintln(w, fmt.Sprintf(out, style, script))
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	log.Fatalln("bye bye 😥")
	http.Error(w, "bye bye 😥", http.StatusInternalServerError)
}

func healthCheckHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, "{\"status\":\"ok\",\"icon\":\"👌\"}")
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	value := rand.Intn(10)
	if value < 4 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, fmt.Sprintf("🙌 got: %d", value))
	} else {
		i := rand.Intn(4)
		code := errorsList[i]
		sendError(w, code)
	}
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	i, err := strconv.Atoi(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if i > 399 {
		sendError(w, i)
	} else {
		w.WriteHeader(i)
		fmt.Fprintln(w, fmt.Sprintf("🙌 got: %d", i))
	}
}

func sendError(w http.ResponseWriter, code int) {
	log.Printf("error %d", code)
	switch code {
	case http.StatusInternalServerError:
		http.Error(w, "🤯 StatusInternalServerError!", http.StatusInternalServerError)
		return
	case http.StatusServiceUnavailable:
		http.Error(w, "🤡 StatusServiceUnavailable!", http.StatusServiceUnavailable)
		return
	case http.StatusNotFound:
		http.Error(w, "🙄 StatusNotFound!", http.StatusNotFound)
		return
	case http.StatusBadGateway:
		http.Error(w, "👻 StatusBadGateway!", http.StatusBadGateway)
		return
	case http.StatusBadRequest:
		http.Error(w, "😕 StatusBadRequest!", http.StatusBadRequest)
		return
	case http.StatusConflict:
		http.Error(w, "😭 StatusConflict!", http.StatusConflict)
		return
	default:
		http.Error(w, fmt.Sprintf("👽 error code %d not found!", code), http.StatusNotFound)
		return
	}
}

var style = `
body {
	margin: auto;
	padding: 25px;
	text-align: center;
}

.content {
	justify-content: space-between;
	background-color: #dedede;
	max-width: 800px;
	width: 50%;
	margin: 0 auto;
	padding: 10px;
}

.buttons input, .danger-button {
	border:1px solid black;
	padding:10px;
	margin: 0px 15px;
	width: 45%;
}

.danger-button {
	color: white;
	font-weight: bold;
	background-color: red;
}

.result {
	margin: 10px auto;
}

.result-item {
    text-align: left;
    width: 48%;
    margin: 0px auto 0px auto;
    border: solid 1px lightgray;
    padding: 5px;
}
`

var script = `
	console.log("script starting!");
	var startButton = document.getElementById("button-start");
	var stopButton = document.getElementById("button-stop");
	var resultList = document.getElementById("result");
	var oldDiv = document.createElement("div"); 
	var oldBGColor;
	var toggle = false;
	var running = false;
	bootstrap();
	startButton.onclick = function () { 
		if (!running) {
			oldBGColor = startButton.style.backgroundColor;
			startButton.style.backgroundColor = 'lime';
			running = true;
		}
		callErrors();
	};
	stopButton.onclick = function () { 
		if (running) {
			startButton.style.backgroundColor = oldBGColor;
			running = false;
		}
	};

	function callErrors() {
		var url = '/errors/random';
		var request = new XMLHttpRequest();
		request.open("GET",url);
		request.addEventListener('load', function(event) {
			console.log(request.responseText);
			addNode(request.status+" "+request.statusText+" "+request.responseText);
		});
		request.send();
		if(running) {
			toggleColor();
			setTimeout(callErrors, 500);
		}
	}

	function toggleColor() {
		if(toggle) {
			startButton.style.backgroundColor = oldBGColor;
			toggle = false;
		} else {
			startButton.style.backgroundColor = 'lime';
			toggle = true;
		}
	}

	function addNode(text) {
		var newDiv = document.createElement("div"); 
		newDiv.classList.add("result-item");
		var newContent = document.createTextNode(text); 
		newDiv.appendChild(newContent);
		resultList.insertBefore(newDiv, oldDiv); 
		oldDiv = newDiv;
	}

	function bootstrap() {
		oldDiv.classList.add("result-item");
		resultList.appendChild(oldDiv);
	}
`
