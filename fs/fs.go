package fs

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"fmt"
	provModel "github.com/afzalabbasi/testrepo/api_call"
	"github.com/callmylist/callmylist-provapi/model"
	"github.com/callmylist/callmylist-provapi/manager/buffer"
	"github.com/callmylist/callmylist-provapi/model"
	"github.com/callmylist/callmylist-provapi/others/eventsocket"
	"github.com/callmylist/callmylist-provapi/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//StartFSClient starts the application with the messages handler listening
func StartFSClient() {

	systemThreshold := 1
	amChannels := 1

	//starts the go routine for handling of messages
	go handleMessage()

	//calls for all servers that are defined to be started, this can be done from the database later on, or they can be called by the api to be started
	startServers()

	pushInChannelBulk(buffer.CallsMaxCapacityChannel, systemThreshold, 1)
	pushInChannelBulk(buffer.AMMaxCapacityChannel, amChannels, 1)
}

func pushInChannelBulk(channel chan int, bufferSize int, value int) {
	for i := 0; i < bufferSize; i++ {
		channel <- value
	}
}

//startServers starts all the servers by using getServers
func startServers() {

	//gets all servers that need to be started
	fsServers := getServers()
	//starts a goroutine for each server
	for _, server := range fsServers {
		_, err := StartConnection(server)
		if err != nil {
			logrus.Println("Start connection => Server : ", server.Host, " => Error : ", err.Error())
			go ServerReconnectionPolling(server)
		}
	}
}

//getServers gets all the servers that need to be started
func getServers() []Server {

	//define servers, this can be passed in as json, a parameter through an API or read in from a database
	//for now only 2 servers are defined and they are defined manually
	var fsServers []Server
	//	fsServers = append(fsServers, Server{"208.76.55.72", "8021", "ClueCon", "10", "json", "ALL", true})
	fsServers = append(fsServers, Server{"167.99.51.83", "8021", "ClueCon", "10", "json", "CHANNEL_HANGUP_COMPLETE CHANNEL_HANGUP CUSTOM dnc::notify", true})
	//fsServers = append(fsServers, Server{"fs2.mycallblast.com", "8021", "ClueCon", "10", "json", "CHANNEL_HANGUP_COMPLETE CHANNEL_HANGUP CUSTOM dnc::notify", true})
	//fsServers = append(fsServers, Server{"fs1.mycallblast.com", "8021", "ClueCon", "10", "json", "ALL", true})

	return fsServers
}

//StartConnection starts a single server connection
func StartConnection(fsServer Server) ([]byte, error) {

	var conEst = make(chan int)
	go serverGoroutine(fsServer, conEst)

	value := <-conEst
	if value == 1 {
		servers = append(servers, fsServer)
		logrus.Println("Server reconnected : ", fsServer.Host)
		return json.Marshal(fsServer)
	} else {
		return nil, errors.New("Server connection not established")
	}
}

//StopConnection stops communicating with a single server
func StopConnection(host string) ([]byte, error, Server) {
	logrus.Println("stopping connect for host : ", host)

	//remove server from servers slice
	var i = 0
	var s Server

	for index, server := range servers {
		if server.Host == host {
			i = index
			s = server
			break
		}
	}

	servers = append(servers[:i], servers[i+1:]...)

	// remove server from server channel
	serverChannel = drainServers(serverChannel, host)
	data, err := json.Marshal(Server{host, "", "", "", "", "", true})

	return data, err, s
}

func drainServers(ch chan string, drain string) chan string {
	var newServers = make(chan string, 1000)
	for {
		select {
		case e := <-ch:
			if e == drain {
				continue
			}
			newServers <- e
			fmt.Printf("%s\n", e)
		default:
			return newServers
		}
	}

	return newServers
}

//ListConnections lists all servers the program is connected to currently
func ListConnections() ([]byte, error) {

	return json.Marshal(servers)
}

//handleMessage is started by StartFSClient to handle events and messages
func handleMessage() {

	//listens to eventChannel and done channel
	for {
		select {

		//case of event coming through
		case event := <-eventChannel:

			if event.Event.Get("Variable_uuid") != "" {
				//fmt.Printf("\nEVENT: %s\nUUID: %s\n\n", event.Event.Get("Event-Name"), event.Event.Get("Variable_uuid"))
				eventOriginateResponseChannel <- event
			}
		}
	}
}

//serverGoroutine is used to start a goroutine of a single server communicating its events through channels
func serverGoroutine(fsServer Server, ch chan int) {

	//connects to the server
	c, err := eventsocket.Dial(fsServer.Host+":"+fsServer.Port, fsServer.Password)
	if err != nil {
		logrus.Println("free switch connection not established")
		logrus.Println(err)
		ch <- 0
		return
	}

	ch <- 1

	logrus.Println("free switch connection established : " + fsServer.Host)

	//when the function ends it closes the connection to the server
	defer c.Close()

	messagesChannel[fsServer.Host] = make(chan Message)
	done[fsServer.Host] = make(chan Message)

	//starts a goroutine that will listen to all incoming events
	go serverEvents(fsServer, c)

	//this is a temporary solution to give time to run the event command on the server before running the cleanup and rest of it
	time.Sleep(time.Second * 5)

	serverChannel <- fsServer.Host

	//listens to channels
	for {

		//catches data from messages and done channels
		select {

		//messages come to the correct server because the channel is specifically made for this server
		case message := <-messagesChannel[fsServer.Host]:

			response, err := c.Send(message.Message)
			if err != nil {
				logrus.Println("!!!!!!!!!!!!!!! sending call failed !!!!!!!!!!!!!!")
				logrus.Println("Error : ", err)
				logrus.Println(message)
				logrus.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\nmessage host:\t" + message.Host + "\nserver host:\t" + fsServer.Host + "\nmessage sent:\t" + message.Message + "\nXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
				logrus.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				buffer.RefillOneCall()
				break
			}

			//assigns the returning message to response while checking variables
			if message.ID != 0 {
				var respond Response
				respond.ID = message.ID
				respond.Host = message.Host
				respond.ReceivedTime = time.Now()
				respond.ReturnEvent = *response

				if response.Body != ""{
					respond.Body = response.Body
				} else {
					respond.Body = ""
				}
				if err != nil {
					respond.Error = err.Error()
				} else {
					respond.Error = ""
				}
				responseChannel <- respond
			}

			//in case a message comes through the done channel, it checks if it is mean for this server and if so it shuts down the goroutine and the connection to the server
		case message := <-done[fsServer.Host]:

			logrus.Println("Done received : Server = ", fsServer.Host, " => Message Host = ", message.Host)
			if message.Host == fsServer.Host && message.Stop {
				_, ok := messagesChannel[fsServer.Host]
				if ok {
					delete(messagesChannel, fsServer.Host)
					messagesChannel[fsServer.Host] = make(chan Message)
				}
				logrus.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\nCLOSED server " + fsServer.Host + " goroutine\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

				_, err, server := StopConnection(fsServer.Host)
				if err != nil {
					logrus.Println("Error stopping connection : ", err)
				}
				// start polling of the free switch connection which just got disconnected
				go ServerReconnectionPolling(server)

				return
			} else {
				////send the message back to done and then sleep a bit so other goroutines can check the message
				//if time.Since(message.ReceivedTime) < 60*time.Second {
				//	//get rid of messages older than 60 seconds
				//	done <- message
				//	time.Sleep(time.Duration(rand.Intn(25)+5) * time.Millisecond)
				//}
			}
		}
	}
}

//serverEvents listens to all events from a specific channel and sends those messages to a channel
func serverEvents(fsServer Server, c *eventsocket.Connection) {

	defer close(messagesChannel[fsServer.Host])

	//set which event type of event the function should listen to
	eventsToListen := "events " + fsServer.EventType + " " + fsServer.EventList
	// eventsToListen := "events " + fsServer.EventType + " BACKGROUND_JOB API CHANNEL_CREATE CHANNEL_DESTROY CUSTOM callcenter::info"
	_, err := c.Send(eventsToListen)
	if err != nil {
		logrus.Println("------ Not able to connect with event socker => fs : ", fsServer.Host, " --------")
		log.Print(err)
		done[fsServer.Host] <- Message{0, fsServer.Host, "", true, time.Now()}
		return
	}

	for {

		//reads events as they come in and if they are not an error it sends them to the eventChannel
		ev, err := c.ReadEvent()
		if err != nil {
			done[fsServer.Host] <- Message{0, fsServer.Host, "", true, time.Now()}
			logrus.Println("Error for free switch : " + fsServer.Host)
			logrus.Println(err)
			return
		}

		if ev != nil {
			eventChannel <- Event{ev, time.Now()}
		}
	}

	logrus.Println("Finished server disconnection")
}

func ServerReconnectionPolling(server Server) {
	logrus.Println(server.Host, " trying to reconnect")
	// start polling
	<-time.After(10 * time.Second)
	_, err := StartConnection(server)
	if err != nil {
		ServerReconnectionPolling(server)
	}
}

//callStart sends a message to a server to start a call
func CallStart(campaign provModel.Call) {

	//get the next server in queue to be send a call to
	ip := <-serverChannel
	serverChannel <- ip

	logrus.Println("\n")
	logrus.Println("Started call fo campaign : " + campaign.Compaign_type + " :: on number : " + campaign.To_DID)
     			to_DID:=campaign.To_DID
     		//	from_DID:=campaign.From_DID
     		//	tx_DID:=campaign.Tx_DID
     		//	soundPath:=campaign.Human_audio

	go makeCall(ip, campaign, to_DID)
}

//makeCall sends a message to a server to start a call
func makeCall(ip string, campaign provModel.Call, to_DID string) {
	callUUID, err := utils.GetUUID()
	if err != nil {
		logrus.Println("UUID generation failed for campaign : " + campaign.Compaign_type + " :: On number : " + campaign.To_DID)
		return
	}

	//domain := "sip.mycallblast.com"
	//carrierToken := "MyCallBl5p0jo2i40h1h"

	dialString := campaignDialString(campaign, callUUID)

	gateway := "sofia/gateway/thinQ/"
	makeOriginate := "bgapi originate "
	makeOriginate += dialString + gateway
	makeOriginate += " "
	makeOriginate += to_DID
	makeOriginate += " "
	logrus.Println("")
	logrus.Println("Campaign : " + campaign.Compaign_type + " :: Dial string : " + makeOriginate)
	logrus.Println("call on ip : ", ip)

	if _, ok := messagesChannel[ip]; ok {
		//do something here
	} else {
		messagesChannel[ip] = make(chan Message)
	}

	select {
	case messagesChannel[ip] <- Message{0, ip, makeOriginate, false, time.Now()}:
		fmt.Println("Call initiating")
	case <-time.After(1 * time.Minute):
		messagesChannel[ip] = make(chan Message)
		messagesChannel[ip] <- Message{0, ip, makeOriginate, false, time.Now()}
		fmt.Println("Call send timeout : Campaign : " + campaign.Compaign_type + " => Number : " + to_DID)
	}

	messagesChannel[ip] <- Message{0, ip, makeOriginate, false, time.Now()}

	// Send API Event Trying to Let them know we are Trying to Create a New Call
	outboundCalls = append(outboundCalls, outCallStruct{UUID: callUUID,to_DID:to_DID, CampaignId: campaign.Compaign_type, IP: ip})

	// expiry channel
	var expiryChannel = make(chan int)

	go func() {
		<-time.After(120 * time.Second)
		expiryChannel <- 1
	}()

	for {
		select {
		case event := <-eventOriginateResponseChannel:

			if event.Event.Get("Variable_uuid") == callUUID {
				if event.Event.Get("Event-Name") == "CUSTOM" && event.Event.Get("Event-Subclass") == "dnc::notify" {
					campaignId := event.Event.Get("Campaign_Id")
					caller := event.Event.Get("Caller")
					userId := event.Event.Get("User_Id")
					// put number in dnc campaign
					// dnc info
					dncInfo := model.DncInfo{CampaignId: campaignId, UserId: userId, Caller: caller}
					logrus.Println("Number put in dnc")
					logrus.Println(dncInfo)
					// fire dnc event
					buffer.DncChannel <- dncInfo

				} else if event.Event.Get("Event-Name") == "CHANNEL_HANGUP_COMPLETE" {
					logrus.Println("*************Call hanged up***************")
					// refill buffer with once channel
					buffer.RefillOneCall()
					return
				}
			} else {
				////if it wasn't inteded for this function, send the return event back to the channel and wait to read a new response
				if time.Since(event.ReceivedTime) < 60*time.Second {
					//	//get rid of messages older than 60 seconds
					eventOriginateResponseChannel <- event
					time.Sleep(time.Duration(rand.Intn(25)+5) * time.Millisecond)
				}
			}
		case expire := <-expiryChannel:
			fmt.Println(expire)
			buffer.RefillOneCall()
			return
		}
	}
}
