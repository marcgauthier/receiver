/*

This package role is to receive UDP syslog message from the RELAY and extract the information from the syslog message and
search for a corresponding probe that have the same IP or hostname.  If probe is found the status of the probe is updated.
If the probe is not found a new probe is created.

Format that is required to be respected by the network monitoring tool.

	Example of InterMapper syslog message: tool='SSC IMap';tooltype='IMAP';host='<probe Name>';ip='<probe Address>';status='<probe Status>';network='DWAN'

	Example of WhatsUpGold syslog message: tool='WUG 76';tooltype='WUG';host='%probe.HostName';ip='%probe.Address';status='%probe.State';network='DWAN'

	Valid option for tooltype = WUG or IMAP


*/

package receiver

import (
	"github.com/antigloss/go/logger"
	"net"
	"strconv"
)

/* most packet coming from the broadcaster should be in the 100b-200bytes range, and
   should most of the time be under the MTU size.  But theoricly can go up to max
*/
const maxUDPSize = 65536

type RecvFunc func (data[]byte)

/* Start Main code executed at the start of the program.  This function will not
returned and must be called with the go prefix.
*/
func Start(listenIP string, listenPort, queueSize int, recvfunc RecvFunc) {

	logger.Info("Starting the receiver service listening for UDP on " + listenIP + ":" + strconv.Itoa(listenPort))

	/* Open UDP Socket and start listening.
	 */
	addr := net.UDPAddr{
		Port: listenPort,
		IP:   net.ParseIP(listenIP),
	}

	ServerConn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	/* start all the threads that will process the packets
	 */
	for i := 0; i < queueSize; i++ {
		go processingQueue(ServerConn, recvfunc)
	}

}

/* Threads that listen for UDP packet confirm there is a valid information and then
   send it to the probes.
*/

func processingQueue(connection *net.UDPConn, recvfunc RecvFunc) {

	// initialize all the required variables and buffer.

	n, err, buf := 0, error(nil), make([]byte, maxUDPSize)

	// infinite loop reading packet and processing

	for {
		// read UDP from the socket.
		n, _, err = connection.ReadFromUDP(buf)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		recvfunc(buf[:n])
		
	}

}
