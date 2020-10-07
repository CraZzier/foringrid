package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Backinfo model
type Backinfo struct {
	Routes    []RoutesStruct    `json:"routes"`
	Waypoints []WaypointsStruct `json:"waypoints"`
	Code      string            `json:"code"`
}

//RoutesStruct model
type RoutesStruct struct {
	Legs       []Legs  `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     float64 `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}

//Legs model
type Legs struct {
	Summary  string    `json:"summary"`
	Weight   float64   `json:"weight"`
	Duration float64   `json:"duration"`
	Steps    []float64 `json:"steps"`
	Distance float64   `json:"distance"`
}

//WaypointsStruct model
type WaypointsStruct struct {
	Hint     string    `json:"hint"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

//Feedback sending this back to client
type Feedback struct {
	Source string           `json:"source"`
	Routes []RoutesFeedback `json:"routes"`
}

//RoutesFeedback type inside feedback
type RoutesFeedback struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

//ErrorResponse servers informing client about error by message
type ErrorResponse struct {
	Message string `json:"message"`
}

//TestSort checks if sorting is valid
func TestSort() bool {
	var testgather []RoutesFeedback
	test1 := RoutesFeedback{
		Destination: "New York",
		Distance:    500,
		Duration:    30,
	}
	test2 := RoutesFeedback{
		Destination: "New Orlean",
		Distance:    500,
		Duration:    20,
	}
	test3 := RoutesFeedback{
		Destination: "New Jersey",
		Distance:    400,
		Duration:    10,
	}
	test4 := RoutesFeedback{
		Destination: "Washington",
		Distance:    500,
		Duration:    10,
	}
	testgather = append(testgather, test1)
	testgather = append(testgather, test2)
	testgather = append(testgather, test3)
	testgather = append(testgather, test4)
	testgather = SortRoads(testgather)
	if testgather[0].Destination == "New Jersey" && testgather[1].Destination == "Washington" && testgather[2].Destination == "New Orlean" && testgather[3].Destination == "New York" {
		return true
	}
	return false
}

//SortRoads sorts roads ascending
func SortRoads(feedbackroutes []RoutesFeedback) []RoutesFeedback {
	var temp RoutesFeedback
	for n := 0; n < len(feedbackroutes); n++ {
		for m := 1 + n; m < len(feedbackroutes); m++ {
			temp = feedbackroutes[n]
			if feedbackroutes[m].Duration < feedbackroutes[n].Duration {
				feedbackroutes[n] = feedbackroutes[m]
				feedbackroutes[m] = temp
			} else if feedbackroutes[m].Duration < feedbackroutes[n].Duration {
				if feedbackroutes[m].Distance < feedbackroutes[n].Distance {
					feedbackroutes[n] = feedbackroutes[m]
					feedbackroutes[m] = temp
				}
			} else {
				continue
			}
		}
	}
	return feedbackroutes
}

//CORS for task to get foreign-origin requests
func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
func routes(c *gin.Context) {
	//getting url get-data
	src, errSrc := c.Request.URL.Query()["src"]
	dst, errDst := c.Request.URL.Query()["dst"]

	//testing sort
	fmt.Println(TestSort())

	if !errSrc || !errDst {

		c.Writer.Header().Set("Content-type", "texthtml; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		fmt.Fprintln(c.Writer, "Variable dst,src or both are missing")

	} else {

		//declaring needed variables
		var feedbackroutes []RoutesFeedback
		var tempFeedback RoutesFeedback
		var info Backinfo
		var feed Feedback

		//requesting response for every destination
		for i := 0; i < len(dst); i++ {
			locReq := "http://router.project-osrm.org/route/v1/driving/" + src[0] + ";" + dst[i] + "?overview=false"
			resp, err := http.Get(locReq)
			if err != nil {
				var error ErrorResponse
				error.Message = "An error occured while getting service data"
				c.Writer.Header().Set("Content-type", "application/json")
				c.Writer.WriteHeader(http.StatusOK)
				json.NewEncoder(c.Writer).Encode(error)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			//getting json response into struct
			json.Unmarshal(body, &info)

			//cheking correctness of a reponse
			if info.Code != "Ok" {
				var error ErrorResponse
				error.Message = info.Code
				c.Writer.Header().Set("Content-type", "application/json")
				c.Writer.WriteHeader(http.StatusOK)
				json.NewEncoder(c.Writer).Encode(error)
				return
			} else {
				//adding all single request roads
				for p := 0; p < len(info.Routes); p++ {
					tempFeedback.Destination = dst[i]
					tempFeedback.Distance = info.Routes[p].Distance
					tempFeedback.Duration = info.Routes[p].Duration
					feedbackroutes = append(feedbackroutes, tempFeedback)
				}
			}

		}

		//sorting by bubble sort considering possible amount of roads
		feedbackroutes = SortRoads(feedbackroutes)

		//seting up final response
		feed.Source = src[0]
		feed.Routes = feedbackroutes
		c.Writer.Header().Set("Content-type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(feed)

	}
}
func hi(c *gin.Context) {
	fmt.Fprint(c.Writer, "hello!")
}
func main() {
	r := gin.Default()
	r.Use(CORS)
	r.GET("/", hi)
	r.GET("/routes", routes)
	r.Run(":8080")
}
