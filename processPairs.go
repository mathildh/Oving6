package main 
import (
	. "fmt"
	. "net"
	"time"
	"strings"
	"strconv"
	"os/exec"
)

func imAlive(port string){
	Println("Establishing iAmAlive connection")
	sendAddress, err := ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	checkIfError(err,"ERROR RESOLVING IMALIVE CONNECTION")
	connection, err := DialUDP("udp",nil,sendAddress)
	checkIfError(err,"ERROR DIALING THE UDP CONNECTION FOR IMALIVE")
	for {
		time.Sleep(time.Millisecond * 100)
		connection.Write([]byte("I'm alive!!!"))
	}
}

func checkIfError(err error, error_msg string){
	if err != nil{
		Println("Error of type: " + error_msg)
	}
}

func makeListnerConnection(port string) *UDPConn{
	udpAddress, err := ResolveUDPAddr("udp4", ":"+ port)
	Println("Establishing a lister to network..")
	checkIfError(err, "ERROR RESOLVING UDP ADDRESS ON PORT: " + port)
	connection, err := ListenUDP("udp", udpAddress)
	Println("Listening to the port: "+port)
	checkIfError(err, "ERROR LISTENING TO UDP ON PORT: " + port)
	return connection
}
func listenToNetworkTimeLimited(connection *UDPConn, outgoingMsg chan string, timeLimit int){
	data := make([]byte, 1024)
	for {
		connection.SetReadDeadline(time.Now().Add(time.Duration(timeLimit)* time.Millisecond))
		_,_,err := connection.ReadFromUDP(data)
		checkIfError(err, "ERROR WHILE READING FROM UDP")
		if err != nil{
			Println("Channeling data: connection is dead")
			go func(outgoingMsg chan string){
				outgoingMsg <- "connection is dead"
			}(outgoingMsg)
			Println("Backup: Breaking listen-loop")
			break
		}
		outgoingMsg <- string(data)
	}
	connection.Close()
}
func listenToNetwork(connection *UDPConn, incomingMsg chan string){
	data := make([]byte,1024)
	for{
		_,_,err := connection.ReadFromUDP(data)
		checkIfError(err, "ERROR READING FROM UDP IN listenToNetwork")
		if err != nil{
			data = []byte("connection is dead")
		}
		Println("listenToNetwork: channeling data: " + string(data))
		incomingMsg <-  string(data)
	}
	connection.Close()
}

func sendToNetwork(port string, msg string) {
	sendAddress, err := ResolveUDPAddr("udp4", "129.241.187.255:" + port)
	checkIfError(err, "ERROR WHILE RESOLVING UDP ADDRESS FOR SENDING")
	connection, err := DialUDP("udp4",nil,sendAddress)
	checkIfError(err, "ERROR WHILE DIALING TO UDP FOR SENDING")
	connection.Write([]byte(msg))
	connection.Close()

}

func main() {
  alivePort := "20013"
  countPort := "20014"
  imAliveConnection := makeListerConnection(alivePort)
	countConnection := makeListnerConnection(countPort)
	incomingMsg := make(chan string)
	//if there is a master we will detect imAlive messages
	var imAlive string
	listenToNetworkTimeLimited(imAliveConnection,incomingMsg,500)
	imAlive = <- incomingMsg
	
	if imAlive == ""I'm alive!!!"{
	// I will be backup 
	
	go listenToNetworkTimeLimited(imAliveConnection,incomingMsg,500) //recieving the alive msgs from master
	go listenToNetwork(countConnection,countChannel) //receive count variable 

	// updating and listening 
	count = func(incomingMsg chan string, countChannel chan string, count string) string{
		for{
			select{
			case update = <- incomingMsg:
				if update == "connection is dead"{
					Println("You have a dead main...")
					return count
				}
			case count = <- countChannel:
				Println("Backup recieved count: ", count)
			default:
			}

		}
	}(incomingMsg,countChannel,count)
	
	
	}
	
	else {
	//I will be master
	//Need to create a backupfile 
	cmd := exec.Command("gnome-terminal", "-e", "./main")
	cmd.Output()
	
	
	
	
	}
	
	
	
	


	incomingMsg := make(chan string)
	
	countPort := "20014"
	countConnection := makeListnerConnection(countPort)

	
	var continueToCount string 

	go imAlive(alivePort)
	go listenToNetworkTimeLimited(countConnection,incomingMsg,500)

	continueToCount = <- incomingMsg
	Println("ContinueToCount before Trim: " + string(continueToCount))

	continueToCount = strings.Trim(continueToCount,"\x00")
	Println("ContinueToCount after Trim: " + string(continueToCount))
	if continueToCount == "connection is dead"{
		continueToCount = "0"
	}

	i,_ := strconv.Atoi(continueToCount)
	Println("Printing i after conversion: ", i)
	

	//creates a new backup, spawns
	countConnection.Close()
	Println("Creating backup")
	cmd := exec.Command("gnome-terminal", "-e", "./backup")
	cmd.Output()

	//counts and update over UDP in a seperate thread 
	go func(countPort string, i int){
		for{
			sendToNetwork(countPort,strconv.Itoa(i))
			i += 1
			time.Sleep(time.Second*1)
			Println(i)
		}
	}(countPort,i)
	var exit string
	Scanln(&exit)

}
