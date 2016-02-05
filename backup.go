package main 
import (
	. "fmt"
	. "net"
	"time"
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
func main(){
	var (
		update string
		count string
	)
	alivePort := "20013"
	countPort := "20014"
	aliveConnection := makeListnerConnection(alivePort)
	countConnection := makeListnerConnection(countPort)
	incomingMsg := make(chan string)
	countChannel := make(chan string)
	
	go listenToNetworkTimeLimited(aliveConnection,incomingMsg,500) //recieving the alive msgs
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

	
	countConnection.Close()
	aliveConnection.Close()
	Println("Creating a new main")
	cmd := exec.Command("gnome-terminal", "-e", "./main")
	cmd.Run()

	Println("Count =",count)
	for i := 0; i<500; i++{
		go sendToNetwork(countPort,count)
		time.Sleep(time.Millisecond)
	}

}
