package fs

import (
	//project based packages
	"time"

	"github.com/callmylist/callmylist-provapi/others/eventsocket"
)

//Event holds the struct of the event that is received from the server
type Event struct {
	Event        *eventsocket.Event
	ReceivedTime time.Time
}

//Message holds the struct of the complete message that needs to be sent back to the server
type Message struct {
	ID           int    //id of the message that was sent
	Host         string //IP of the server that a message is being sent to
	Message      string //message that needs to be sent
	Stop         bool   //in case a goroutine needs to be stopped, this along with host is used to confirm to stop its operation
	ReceivedTime time.Time
}

//Response holds the struct of the event that was a resonse to the message that was sent and anything else related to the message sent
type Response struct {
	ID          int               `json:"-"` //id of the message that was sent so it can be matched to the correct response
	Host        string            `json:"-"` //host that the response is going to or it came from
	ReturnEvent eventsocket.Event `json:"return_event"`
	// Header       map[string]interface{} `json:"header"` //the returned event header
	Body         string    `json:"body"`  //the returned event body
	Error        string    `json:"error"` //the returned error if there was one
	Event        []Event   `json:"event"` //the returned event or events that were tied to the original response
	ReceivedTime time.Time `json:"ommit"`
}

//Server is used to hold the info for freeswitch servers
type Server struct {
	Host      string `json:"host"`       //freeswitch hostname
	Port      string `json:"port"`       //freeswitch port
	Password  string `json:"password"`   //freeswitch password
	Timeout   string `json:"timeout"`    //freeswitch conneciton timeout in seconds
	EventType string `json:"event_type"` //freeswitch event type eg. json or plain
	EventList string `json:"event_list"` //freeswitch list of events eg. list of events to log or all
	Live      bool   `json:"live"`       //freeswitch server live monitoring system
}

//CurrentCall holds the current call info
type CurrentCall struct {
	UUID  string
	Agent string
}

//Call holds all the info about a call
type Call struct {
	ContactListId string `json:"ContactListId"`
}

//channel used for sending messages that are received from a server
var eventChannel = make(chan Event)

//channel used for sending events related to agents that are received from the server
var eventAgentResponseChannel = make(chan Event, 1000)

//channel used for sending events related to tiers that are received from the server
var eventTierResponseChannel = make(chan Event, 1000)

//channel used for sending events related to queues that are received from the server
var eventQueueResponseChannel = make(chan Event, 1000)

//channel used for sending events related to originate that are received from the server
var eventOriginateResponseChannel = make(chan Event, 1000)

//channel used for sending messages back to the server
// var messagesChannel = make(chan Message)
var messagesChannel = make(map[string]chan Message)

// channel to receive stop events channel
var CampaignStopChannel map[string]chan int = make(map[string]chan int, 10000)

// channel to receive campaign internal channels availability
var CampaignInternalLimitChannel map[string]chan int = make(map[string]chan int, 10000)

//channel used for sending to the response requester
var responseChannel = make(chan Response, 1000)

//channel used for sending messages back to the tier
var tierChannel = make(chan Response, 1000)

//channel used for a round robin queue of servers
var serverChannel = make(chan string, 1000)

//channel used for a round robin queue of dynamic agent servers
var dynamicServerChannel = make(chan string, 1000)

//channel used to stop goroutines
var done = make(map[string]chan Message)

//servers handler that holds all currently running servers
var servers []Server

//outCallStruct holds all the info about a call
type outCallStruct struct {
	UUID       string
	to_DID     string
	CampaignId string
	IP         string
}

//outboundCalls holds all current calls that are outbound
var outboundCalls []outCallStruct
