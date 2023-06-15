package main

import (
	"github.com/syusuke/ps-rtp-streams/packet"

	"github.com/nareix/joy4/format"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	log "github.com/sirupsen/logrus"
)

func init() {
	format.RegisterAll()
}

func main() {

	rtp := packet.NewRRtpTransfer("", packet.LocalCache)

	// send ip,port and recv ip,port
	rtp.Service("127.0.0.1", "172.20.25.2", 10086, 10087)

	f, err := avutil.Open("test.flv")
	if err != nil {
		log.Errorf("read file error(%v)", err)
		rtp.Exit()
		return
	}

	var pts uint64 = 10000
	streams, _ := f.Streams()
	var vindex int8
	for i, stream := range streams {
		if stream.Type() == av.H264 {
			vindex = int8(i)
			break
		}
	}

	for i := 0; i < 10; i++ {
		var pkt av.Packet
		var err error
		if pkt, err = f.ReadPacket(); err != nil {
			log.Errorf("read packet error(%v)", err)
			goto STOP
		}
		if pkt.Idx != vindex {
			continue
		}
		rtp.Send2data(pkt.Data, pkt.IsKeyFrame, pts)
		pts += 40
		//time.Sleep(time.Millisecond * 40)
	}
STOP:
	f.Close()
	rtp.Exit()
	return

}
