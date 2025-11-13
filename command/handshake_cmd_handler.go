package command

import (
	"net"
	"strings"

	"com.github.redisgo/config"
	"com.github.redisgo/replication"
	"com.github.redisgo/util"
)

func (cmd *Cmd) handleREPLCONFCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'replconf' command\r\n")
		return
	}
	key := strings.ToLower(cmd.Args[1])
	switch key {
	case "listening-port":
		writeResponse(conn, util.ReturnOkResponse())
	case "capa":
		writeResponse(conn, util.ReturnOkResponse())
	case "getack":
		writeResponse(conn, util.ReturnReplconfGetAckResponse())
	default:
		writeResponse(conn, "-ERR unknown key: "+key+"\r\n")
	}
}

// todo: make all the handshake commands internal. should not be exposed to the user or redis clients.
func (cmd *Cmd) handlePsyncCommand(conn *net.Conn) {
	if len(cmd.Args) < 3 {
		writeResponse(conn, "-ERR wrong number of arguments for 'psync' command\r\n")
		return
	}
	master_repl_offset := cmd.Args[2]
	if master_repl_offset == "-1" {
		master_repl_offset = "0"
	}
	writeResponse(conn, "+FULLRESYNC "+config.ServerConfig.Master_replid+" "+master_repl_offset+"\r\n")
	replication.AddSlaveConnection(conn)
	//connection.AddSlaveConnection(conn)
	// Send empty RDB file after FULLRESYNC response
	/*if cmd.Args[2] == "-1" {
		fmt.Println("Sending RDB snapshot")
		sendRDBSnapshot(conn, database.GetDB())
	}*/
}

/*func sendRDBSnapshot(conn *net.Conn, db *sync.Map) error {
	rdbBytes, err := createRDBFromSyncMap(db)
	if err != nil {
		return fmt.Errorf("failed to create RDB: %v", err)
	}

	// RESP bulk string format: $<length>\r\n<rdb_bytes>\r\n
	// First write the length prefix
	header := fmt.Sprintf("$%d\r\n", len(rdbBytes))
	if _, err := (*conn).Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write RDB header: %v", err)
	}

	// Then write the RDB bytes directly (not as string to preserve binary data)
	if _, err := (*conn).Write(rdbBytes); err != nil {
		return fmt.Errorf("failed to write RDB bytes: %v", err)
	}

	// Finally write the CRLF terminator
	if _, err := (*conn).Write([]byte("\r\n")); err != nil {
		return fmt.Errorf("failed to write RDB terminator: %v", err)
	}

	fmt.Printf("Sent RDB snapshot (%d bytes) to slave\n", len(rdbBytes))
	return nil
}

func createRDBFromSyncMap(db *sync.Map) ([]byte, error) {
	buf := new(bytes.Buffer)

	// 1️⃣ Write RDB header: "REDIS" + version (9 bytes total)
	buf.WriteString("REDIS")
	buf.WriteByte(0x00)
	buf.WriteByte(0x30)
	buf.WriteByte(0x30)
	buf.WriteByte(0x30)
	buf.WriteByte(0x39) // Version 0009

	// 2️⃣ Write DB selector (0)
	buf.WriteByte(0xFE) // FE = Select DB opcode
	buf.WriteByte(0x00) // DB 0

	// 3️⃣ Write key-value pairs
	db.Range(func(key, value any) bool {
		k := key.(string)

		// Skip internal keys (MULTI-* keys used for transactions)
		if len(k) > 6 && k[:6] == "MULTI-" {
			return true
		}

		// Determine type and serialize accordingly
		switch v := value.(type) {
		case *object:
			// Check if expired before writing anything
			now := time.Now()
			if v.ExpireInMillies > 0 && now.After(v.Time) {
				// Skip expired keys
				return true
			}

			// String type (0x00)
			buf.WriteByte(0x00)
			writeLength(buf, uint64(len(k)))
			buf.WriteString(k)

			// Write value
			writeLength(buf, uint64(len(v.Value)))
			buf.WriteString(v.Value)

		case *objectList:
			// List type (0x0B = list with quicklist encoding)
			buf.WriteByte(0x0B)
			writeLength(buf, uint64(len(k)))
			buf.WriteString(k)

			// Write list length
			writeLength(buf, uint64(len(v.Value)))

			// Write each list element
			for _, item := range v.Value {
				writeLength(buf, uint64(len(item)))
				buf.WriteString(item)
			}

		case *streamList:
			// Stream type (0x15 = stream)
			buf.WriteByte(0x15)
			writeLength(buf, uint64(len(k)))
			buf.WriteString(k)

			// Write stream length
			writeLength(buf, uint64(len(v.Value)))

			// Write each stream entry
			for _, streamObj := range v.Value {
				// Write ID
				writeLength(buf, uint64(len(streamObj.Id)))
				buf.WriteString(streamObj.Id)

				// Write field count (key-value pairs)
				writeLength(buf, uint64(len(streamObj.Value)*2))

				// Write key-value pairs
				for field, val := range streamObj.Value {
					writeLength(buf, uint64(len(field)))
					buf.WriteString(field)
					writeLength(buf, uint64(len(val)))
					buf.WriteString(val)
				}
			}

		default:
			// Unknown type, skip
			return true
		}

		return true
	})

	// 4️⃣ Write EOF opcode
	buf.WriteByte(0xFF)

	// 5️⃣ Write checksum (8 bytes, zero for now)
	buf.Write(make([]byte, 8))

	return buf.Bytes(), nil
}

// Helper: write variable length integers like Redis RDB format
// Redis RDB uses a special encoding for lengths:
// - 6-bit: 0xxxxxxx (0-63)
// - 14-bit: 01xxxxxx xxxxxxxx (64-16383)
// - 32-bit: 10xxxxxx xxxxxxxx xxxxxxxx xxxxxxxx (16384+)
// - 64-bit: 11xxxxxx ... (for very large values)
func writeLength(buf *bytes.Buffer, length uint64) {
	if length < 0x40 {
		// 6-bit encoding (0xxxxxxx)
		buf.WriteByte(byte(length))
	} else if length < 0x4000 {
		// 14-bit encoding (01xxxxxx xxxxxxxx)
		b1 := byte((length >> 8) | 0x40)
		b2 := byte(length & 0xFF)
		buf.WriteByte(b1)
		buf.WriteByte(b2)
	} else if length < 0x40000000 {
		// 32-bit encoding (10xxxxxx xxxxxxxx xxxxxxxx xxxxxxxx)
		buf.WriteByte(0x80)
		tmp := make([]byte, 4)
		binary.BigEndian.PutUint32(tmp, uint32(length))
		buf.Write(tmp)
	} else {
		// 64-bit encoding (11xxxxxx ...)
		buf.WriteByte(0x81)
		tmp := make([]byte, 8)
		binary.BigEndian.PutUint64(tmp, length)
		buf.Write(tmp)
	}
}
*/
