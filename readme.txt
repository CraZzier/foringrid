Application for INRID

1.Introduction
2.How to use
3.Code

1.Introduction
The App has been created in order to provide information about the fastest roads between 
location. It is written in Go with help of a 'gin-gonic'. The entire app is contained in one file
- main.go. In app.yaml file it is defined what runtime should application run on. It is "must have" due to the fact that the
app has been deployed to Google Cloud Service which requires the file and use it as settings provider.
Routes to be checked and sorted are typed into the parameter section of the url with following names(of the parameter)
 - src - start point(should be only one)
 - dst - destination point(can be many)
Example:
https://new-go-routes-07.oa.r.appspot.com/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219
The final data served by app back to client has a format of json(Javascript Object Notation). The import data are in "routes" scrap of the json
Object. They are served ascendingly as it is asked to designate the fastest road.
Except gin-gonic package, app uses also the following packages:
"encoding/json" - managing json inside app
"fmt" - access to I/O functions
"io/ioutil" - reading and bufferes
"net/http" - http solutions

2. How to use
In order to run application locally you must:
- have Go installed on your computer(https://golang.org/dl/)
- download gin using command go get -u github.com/gin-gonic/gin
- run application using go run main.go in the application folder
- then you need to open your browser and open localhost:8080/routes
The other way is to run binary(only when you are using linux with amd64 architecture)
and then open localhost:8080/routes
The last option is to open:
https://new-go-routes-07.oa.r.appspot.com/routes?
and give suitable parameters ine the url

3. Code
Code has a few model classes:
Type made for unmarshalling json response from road distance provider
type Backinfo struct {
	Routes    []RoutesStruct    `json:"routes"`
	Waypoints []WaypointsStruct `json:"waypoints"`
	Code      string            `json:"code"`
}
Subobject of a json returned by road server - unused in the concept of a
type RoutesStruct struct {
	Legs       []Legs  `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     float64 `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}

Subobject of a json returned by road server
type Legs struct {
	Summary  string    `json:"summary"`
	Weight   float64   `json:"weight"`
	Duration float64   `json:"duration"`
	Steps    []float64 `json:"steps"`
	Distance float64   `json:"distance"`
}

Subobject of a json returned by road server
type WaypointsStruct struct {
	Hint     string    `json:"hint"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

Final object class that is given as a final reponse to API
type Feedback struct {
	Source string           `json:"source"`
	Routes []RoutesFeedback `json:"routes"`
}

Struct used in sorting function but also part of a road service
type RoutesFeedback struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

Struct that is used to give message when an error occurs
type ErrorResponse struct {
	Message string `json:"message"`
}

TestSort function is created to validate sort function (Sort Roads)
Rest of helping information is within code
