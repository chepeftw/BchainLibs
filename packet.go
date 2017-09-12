package bchainlibs

import (
	"net"
	// "github.com/op/go-logging"
	"time"
	"strings"
	"crypto/sha256"
	"strconv"
	"encoding/hex"
)




// +++++++++ Constants
const (
	InternalUBlockType = iota // 0
	InternalVBlockType

	UBlockType
	VBlockType

	// Generated by the Blockchain layer, the idea is to remove the "stress" for the miner
	// to validate which is the last block and how to deal with it.
	// If we wanted to work TOTALLY separated, then this would not work, but for v1 it would.
	LastBlockType

	InternalQueryType
	QueryType

	InternalPing
	InternalPong
)


const (
	RouterPort        = ":10000"
	LocalPort         = ":0"
	BlockCPort 	      = ":10001"
	MinerPort 	      = ":10002"
	Protocol          = "udp"
	BroadcastAddr     = "255.255.255.255"
	LocalhostAddr     = "127.0.0.1"
)


// +++++++++ Packet structure
type Packet struct {
	TID          string     `json:"tid"` // Timed ID ... ip + created time
	Type         int        `json:"tp"`

	Source       net.IP     `json:"src,omitempty"`
	//Destination  net.IP     `json:"dst,omitempty"`

	PrID       string	`json:"prd,omitempty"` // Previous block ID
	Salt       string	`json:"slt,omitempty"` // Salt value
	BID        string	`json:"bid,omitempty"` // Block ID = hash(previous block ID + Timed ID + salt)

	Query        *Query     `json:"qry,omitempty"`
	Block    	 *Block 	`json:"blck,omitempty"`
}

type Query struct {
	Function  string    `json:"fct,omitempty"`
}

type Block struct {
	//Data       string	`json:"dat"`
	PacketID   string	`json:"pckt_id"`
	Protocol   string	`json:"pckt_proto"`
	Checksum   string	`json:"pckt_chcksm"`
	Source     net.IP	`json:"pckt_src"`
	Destination net.IP	`json:"pckt_dst"`
	ActualHop  net.IP	`json:"pckt_hop"`
	PreviousHop net.IP	`json:"pckt_prehop"`
	Timestamp  int64    `json:"pckt_tmstmp"`

	Function   string	`json:"bfct"`

	Created    int64    `json:"cts,omitempty"`
	Verified   int64    `json:"vts,omitempty"`
	Verifier   net.IP   `json:"vrfr,omitempty"`
}

func AssembleVerifiedBlock(payload Packet, prid string, salt string, puzzle string, me net.IP) Packet {
	payload.Type = InternalVBlockType

	h := sha256.New()
	h.Write([]byte(puzzle))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	payload.PrID = prid
	payload.Salt = salt
	payload.BID = sha1_hash
	payload.Block.Verified = time.Now().UnixNano()
	payload.Block.Verifier = me

	return payload
}

func AssembleUnverifiedBlock(me net.IP, data string, function string) Packet {
	now := time.Now().UnixNano()

	block := Block{
				Data: data,
				Function: function,
				Created: now,
			}

	payload := Packet{
		TID: generatePacketId( me, now ),
		Type: InternalUBlockType,
		Source: me,
		Block: &block,
	}

	return payload
}

func AssembleQuery(me net.IP, function string) Packet {
	now := time.Now().UnixNano()

	query := Query{
				Function: function,
			}

	payload := Packet{
		TID: generatePacketId( me, now ),
		Type: InternalQueryType,
		Source: me,
		Query: &query,
	}

	return payload
}

func AssemblePing(me net.IP ) Packet {
	return assembleInternal(me, InternalPing)
}

func AssemblePong(me net.IP ) Packet {
	return assembleInternal(me, InternalPong)
}

func assembleInternal(me net.IP, pingType int ) Packet {
	now := time.Now().UnixNano()

	payload := Packet{
		TID: generatePacketId( me, now ),
		Type: pingType,
		Source: me,
	}

	return payload
}



func generatePacketId(me net.IP, now int64) string {
	return me.String() + "_" + strconv.FormatInt(now, 10)
}

func (packet Packet) IsValid( piece string ) bool {
	valid := false

	h := sha256.New()
	puzzle := packet.PrID + packet.TID + packet.Salt
	h.Write([]byte( puzzle ))
	checksum := string(h.Sum(nil))

	sha256_hash := hex.EncodeToString(h.Sum(nil))

	if strings.Contains(checksum, piece) {
		if sha256_hash == packet.BID {
			valid = true
		}
	}

	return valid
}

func (packet Packet) Duplicate() Packet {
	clone := Packet{
		TID: packet.TID,
		Type: packet.Type,
		Source: packet.Source,

		PrID: packet.PrID,
		Salt: packet.Salt,
		BID: packet.BID,
	}

	if packet.Query != nil {
		query := Query{
			Function: packet.Query.Function,
		}

		clone.Query = &query
	}

	if packet.Block != nil {
		block := Block{
			Data: packet.Block.Data,
			Function: packet.Block.Function,
			Created: packet.Block.Created,
			Verifier: packet.Block.Verifier,
			Verified: packet.Block.Verified,
		}

		clone.Block = &block
	}

	return clone
}

func (packet Packet) String() string {
	val := "-> ( "
	val += packet.TID
	val += ", "
	val += packet.BID
	val += ", "
	val += packet.Salt
	val += " )"
	return val
}
