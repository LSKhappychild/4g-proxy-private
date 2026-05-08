package s1ap

import (
	"4g-proxy/lib/s1ap/s1apType"
	"encoding/binary"
	"fmt"
)

// S1AP Message info extracted from PDU
type S1APMessage struct {
	ProcedureCode int64
	ProcedureName string
	MessageType   string // "initiating", "successful", "unsuccessful"
	MMEUEID       *int64
	ENBUEID       *int64
	NASPDU        []byte
	Cause         string
}

// Decode decodes S1AP PDU from raw bytes and extracts message information
func Decode(data []byte) (*S1APMessage, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("S1AP message too short: %d bytes", len(data))
	}

	msg := &S1APMessage{}

	// S1AP PDU is encoded in APER format
	// First byte: bit 7 = extension, bits 6-5 = PDU type (CHOICE index)
	// 0x00-0x1F = initiatingMessage, 0x20-0x3F = successfulOutcome, 0x40-0x5F = unsuccessfulOutcome
	pduType := (data[0] >> 5) & 0x03

	switch pduType {
	case 0: // initiatingMessage
		msg.MessageType = "initiating"
		if err := decodeInitiatingMessage(data, msg); err != nil {
			return msg, err
		}
	case 1: // successfulOutcome
		msg.MessageType = "successful"
		if err := decodeSuccessfulOutcome(data, msg); err != nil {
			return msg, err
		}
	case 2: // unsuccessfulOutcome
		msg.MessageType = "unsuccessful"
		if err := decodeUnsuccessfulOutcome(data, msg); err != nil {
			return msg, err
		}
	default:
		return nil, fmt.Errorf("unknown S1AP PDU type: %d", pduType)
	}

	// Set procedure name from code
	if name, ok := s1apType.ProcedureCodeNames[msg.ProcedureCode]; ok {
		msg.ProcedureName = name
	} else {
		msg.ProcedureName = fmt.Sprintf("Unknown(%d)", msg.ProcedureCode)
	}

	return msg, nil
}

func decodeInitiatingMessage(data []byte, msg *S1APMessage) error {
	if len(data) < 4 {
		return fmt.Errorf("initiating message too short")
	}

	// Extract procedure code (after PDU type indicator)
	// Procedure code is typically in the second byte for small values
	msg.ProcedureCode = int64(data[1])

	// Parse the Protocol IEs to find NAS-PDU and UE IDs
	return parseProtocolIEs(data, msg)
}

func decodeSuccessfulOutcome(data []byte, msg *S1APMessage) error {
	if len(data) < 4 {
		return fmt.Errorf("successful outcome too short")
	}

	msg.ProcedureCode = int64(data[1])
	return parseProtocolIEs(data, msg)
}

func decodeUnsuccessfulOutcome(data []byte, msg *S1APMessage) error {
	if len(data) < 4 {
		return fmt.Errorf("unsuccessful outcome too short")
	}

	msg.ProcedureCode = int64(data[1])
	return parseProtocolIEs(data, msg)
}

// parseProtocolIEs extracts key IEs from S1AP message
func parseProtocolIEs(data []byte, msg *S1APMessage) error {
	if len(data) < 6 {
		return nil
	}

	// Skip PDU header (procedure code, criticality, value length)
	// Find Protocol IE container
	offset := 4

	// Look for length encoding
	if offset >= len(data) {
		return nil
	}

	// Variable length encoding - check for extended length
	firstLenByte := data[offset]
	var containerLen int
	if firstLenByte&0x80 == 0 {
		// Short form: length in lower 7 bits
		containerLen = int(firstLenByte & 0x7f)
		offset++
	} else {
		// Long form or fragmented
		numOctets := int(firstLenByte & 0x7f)
		if numOctets == 0 {
			// Indefinite length not supported
			return nil
		}
		offset++
		if offset+numOctets > len(data) {
			return nil
		}
		for i := 0; i < numOctets; i++ {
			containerLen = (containerLen << 8) | int(data[offset])
			offset++
		}
	}

	// Parse number of IEs
	if offset >= len(data) {
		return nil
	}

	// The container starts with count of IEs
	// This is sequence-of encoding

	// Scan through the data looking for specific IE IDs
	return scanForIEs(data[offset:], msg)
}

// scanForIEs scans through the Protocol IE list to find key IEs
func scanForIEs(data []byte, msg *S1APMessage) error {
	if len(data) < 2 {
		return nil
	}

	// Number of IEs (constrained whole number encoding)
	numIEs := int(binary.BigEndian.Uint16(data[0:2]) >> 8)
	if numIEs == 0 || len(data) < 2 {
		// Try alternative interpretation
		numIEs = int(data[0])
	}

	offset := 1
	if data[0] > 127 {
		offset = 2
	}

	for i := 0; i < numIEs && offset < len(data)-4; i++ {
		if offset+4 > len(data) {
			break
		}

		// IE ID (constrained whole number, 0..65535)
		ieID := int64(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2

		// Criticality (enumerated, 2 bits, but byte-aligned)
		offset++ // skip criticality

		// Value length
		if offset >= len(data) {
			break
		}

		valueLen := 0
		lenByte := data[offset]
		if lenByte&0x80 == 0 {
			valueLen = int(lenByte)
			offset++
		} else {
			numOctets := int(lenByte & 0x7f)
			offset++
			if offset+numOctets > len(data) {
				break
			}
			for j := 0; j < numOctets; j++ {
				valueLen = (valueLen << 8) | int(data[offset])
				offset++
			}
		}

		if offset+valueLen > len(data) {
			valueLen = len(data) - offset
		}

		ieValue := data[offset : offset+valueLen]

		// Process known IEs
		switch ieID {
		case s1apType.ProtocolIEIDMMEUES1APID:
			if len(ieValue) >= 4 {
				v := int64(binary.BigEndian.Uint32(ieValue[0:4]))
				msg.MMEUEID = &v
			} else if len(ieValue) >= 2 {
				v := int64(binary.BigEndian.Uint16(ieValue[0:2]))
				msg.MMEUEID = &v
			}
		case s1apType.ProtocolIEIDENBUES1APID:
			if len(ieValue) >= 3 {
				v := int64(ieValue[0])<<16 | int64(ieValue[1])<<8 | int64(ieValue[2])
				msg.ENBUEID = &v
			} else if len(ieValue) >= 2 {
				v := int64(binary.BigEndian.Uint16(ieValue[0:2]))
				msg.ENBUEID = &v
			}
		case s1apType.ProtocolIEIDNASPDU:
			msg.NASPDU = make([]byte, len(ieValue))
			copy(msg.NASPDU, ieValue)
		case s1apType.ProtocolIEIDCause:
			msg.Cause = decodeCause(ieValue)
		}

		offset += valueLen
	}

	return nil
}

func decodeCause(data []byte) string {
	if len(data) < 1 {
		return "unknown"
	}

	// Cause is a CHOICE
	choiceIdx := (data[0] >> 5) & 0x07

	causeValue := int64(0)
	if len(data) >= 2 {
		causeValue = int64(data[1])
	}

	switch choiceIdx {
	case 0: // radioNetwork
		return fmt.Sprintf("radioNetwork(%d)", causeValue)
	case 1: // transport
		return fmt.Sprintf("transport(%d)", causeValue)
	case 2: // nas
		return fmt.Sprintf("nas(%d)", causeValue)
	case 3: // protocol
		return fmt.Sprintf("protocol(%d)", causeValue)
	case 4: // misc
		return fmt.Sprintf("misc(%d)", causeValue)
	default:
		return fmt.Sprintf("unknown(%d:%d)", choiceIdx, causeValue)
	}
}

// GetProcedureName returns the name of an S1AP procedure code
func GetProcedureName(code int64) string {
	if name, ok := s1apType.ProcedureCodeNames[code]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", code)
}

// IsNASCarryingMessage returns true if the procedure code typically carries NAS PDU
func IsNASCarryingMessage(code int64) bool {
	switch code {
	case s1apType.ProcedureCodeInitialUEMessage,
		s1apType.ProcedureCodeUplinkNASTransport,
		s1apType.ProcedureCodeDownlinkNASTransport,
		s1apType.ProcedureCodeInitialContextSetup:
		return true
	default:
		return false
	}
}
