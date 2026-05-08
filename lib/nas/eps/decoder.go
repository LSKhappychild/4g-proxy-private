package eps

import (
	"fmt"
)

// NASMessage represents a decoded NAS EPS message
type NASMessage struct {
	SecurityHeaderType uint8
	ProtocolDiscr      uint8
	MessageType        uint8
	MessageTypeName    string

	// For EMM messages
	EMMMessage *EMMMessage

	// For ESM messages
	ESMMessage *ESMMessage

	// Raw payload (for security protected messages, this is the inner message)
	InnerMessage *NASMessage
}

// EMMMessage contains EMM-specific fields
type EMMMessage struct {
	// Attach Request specific
	AttachType    *uint8
	NASKeySetId   *uint8
	EPSMobileId   []byte
	EPSMobileIdType *uint8

	// Detach Request specific
	DetachType    *uint8
	SwitchOff     *bool

	// TAU Request specific
	TAUType       *uint8

	// Service Request specific
	ServiceType   *uint8
	KSIAndSeqNum  *uint8

	// Extended Service Request specific
	ExtServiceType *uint8

	// Authentication
	RAND          []byte
	AUTN          []byte
	RES           []byte

	// Identity
	IdentityType  *uint8
	MobileId      []byte
	MobileIdType  *uint8

	// Security Mode
	NASSecAlg     *uint8
	ReplayedUESecCap []byte
}

// ESMMessage contains ESM-specific fields
type ESMMessage struct {
	EPSBearerId   uint8
	PTI           uint8

	// PDN Connectivity specific
	PDNType       *uint8
	RequestType   *uint8
	APN           []byte

	// Bearer Context specific
	QCI           *uint8
	LinkedEPSBearerId *uint8
}

// Decode decodes a NAS EPS message from raw bytes
func Decode(data []byte) (*NASMessage, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("NAS message too short: %d bytes", len(data))
	}

	msg := &NASMessage{}

	// First byte: Security header type (4 bits) + Protocol discriminator (4 bits)
	msg.SecurityHeaderType = (data[0] >> 4) & 0x0f
	msg.ProtocolDiscr = data[0] & 0x0f

	// Handle security protected messages
	if msg.ProtocolDiscr == ProtocolDiscriminatorEPSMobilityManagement &&
		msg.SecurityHeaderType != SecurityHeaderTypePlainNAS &&
		msg.SecurityHeaderType != SecurityHeaderTypeServiceRequestHeader {
		// This is a security protected NAS message
		// Format: Security header (1) + MAC (4) + Sequence number (1) + Plain NAS message
		if len(data) < 7 {
			return msg, fmt.Errorf("security protected message too short")
		}

		// Skip security header (MAC + SQN) and decode inner message
		innerData := data[6:]
		innerMsg, err := Decode(innerData)
		if err != nil {
			return msg, err
		}
		msg.InnerMessage = innerMsg
		msg.MessageType = innerMsg.MessageType
		msg.MessageTypeName = innerMsg.MessageTypeName
		msg.EMMMessage = innerMsg.EMMMessage
		msg.ESMMessage = innerMsg.ESMMessage
		return msg, nil
	}

	// Handle Service Request (special security header type 12)
	if msg.SecurityHeaderType == SecurityHeaderTypeServiceRequestHeader {
		return decodeServiceRequest(data, msg)
	}

	// Plain NAS message
	if msg.ProtocolDiscr == ProtocolDiscriminatorEPSMobilityManagement {
		return decodeEMMMessage(data, msg)
	} else if msg.ProtocolDiscr == ProtocolDiscriminatorEPSSessionManagement {
		return decodeESMMessage(data, msg)
	}

	return msg, fmt.Errorf("unknown protocol discriminator: 0x%02x", msg.ProtocolDiscr)
}

// decodeServiceRequest decodes a Service Request message (special format)
func decodeServiceRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	msg.MessageType = 0 // Service Request has no standard message type
	msg.MessageTypeName = "ServiceRequest"
	msg.EMMMessage = &EMMMessage{}

	if len(data) >= 2 {
		ksi := (data[0] >> 4) & 0x07
		seqNum := data[1]
		combined := (ksi << 5) | (seqNum & 0x1f)
		msg.EMMMessage.KSIAndSeqNum = &combined
		msg.EMMMessage.ServiceType = &ksi
	}

	return msg, nil
}

// decodeEMMMessage decodes an EMM message
func decodeEMMMessage(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 2 {
		return msg, fmt.Errorf("EMM message too short")
	}

	msg.MessageType = data[1]

	if name, ok := EMMMessageTypeNames[msg.MessageType]; ok {
		msg.MessageTypeName = name
	} else {
		msg.MessageTypeName = fmt.Sprintf("EMMUnknown(0x%02x)", msg.MessageType)
	}

	msg.EMMMessage = &EMMMessage{}

	// Decode specific message types
	switch msg.MessageType {
	case EMMAttachRequest:
		return decodeAttachRequest(data[2:], msg)
	case EMMDetachRequest:
		return decodeDetachRequest(data[2:], msg)
	case EMMTAURequest:
		return decodeTAURequest(data[2:], msg)
	case EMMExtServiceReq:
		return decodeExtServiceRequest(data[2:], msg)
	case EMMAuthRequest:
		return decodeAuthRequest(data[2:], msg)
	case EMMAuthResponse:
		return decodeAuthResponse(data[2:], msg)
	case EMMIdentityRequest:
		return decodeIdentityRequest(data[2:], msg)
	case EMMIdentityResponse:
		return decodeIdentityResponse(data[2:], msg)
	case EMMSecModeCommand:
		return decodeSecModeCommand(data[2:], msg)
	}

	return msg, nil
}

// decodeESMMessage decodes an ESM message
func decodeESMMessage(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 3 {
		return msg, fmt.Errorf("ESM message too short")
	}

	msg.ESMMessage = &ESMMessage{
		EPSBearerId: (data[0] >> 4) & 0x0f,
		PTI:         data[1],
	}

	msg.MessageType = data[2]

	if name, ok := ESMMessageTypeNames[msg.MessageType]; ok {
		msg.MessageTypeName = name
	} else {
		msg.MessageTypeName = fmt.Sprintf("ESMUnknown(0x%02x)", msg.MessageType)
	}

	// Decode specific message types
	switch msg.MessageType {
	case ESMPDNConnectivityReq:
		return decodePDNConnectivityRequest(data[3:], msg)
	case ESMActivateDefaultEPSBearerCtxReq:
		return decodeActivateDefaultBearerReq(data[3:], msg)
	}

	return msg, nil
}

// Decode Attach Request
func decodeAttachRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	// EPS attach type + NAS key set identifier
	attachType := data[0] & 0x07
	nasKeySetId := (data[0] >> 4) & 0x07
	msg.EMMMessage.AttachType = &attachType
	msg.EMMMessage.NASKeySetId = &nasKeySetId

	// Old GUTI or IMSI (EPS mobile identity)
	if len(data) > 1 {
		idLen := int(data[1])
		if len(data) > 2+idLen {
			msg.EMMMessage.EPSMobileId = data[2 : 2+idLen]
			if idLen > 0 {
				idType := data[2] & 0x07
				msg.EMMMessage.EPSMobileIdType = &idType
			}
		}
	}

	return msg, nil
}

// Decode Detach Request
func decodeDetachRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	detachType := data[0] & 0x07
	switchOff := (data[0] & 0x08) != 0
	msg.EMMMessage.DetachType = &detachType
	msg.EMMMessage.SwitchOff = &switchOff

	return msg, nil
}

// Decode TAU Request
func decodeTAURequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	tauType := data[0] & 0x07
	nasKeySetId := (data[0] >> 4) & 0x07
	msg.EMMMessage.TAUType = &tauType
	msg.EMMMessage.NASKeySetId = &nasKeySetId

	return msg, nil
}

// Decode Extended Service Request
func decodeExtServiceRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	serviceType := data[0] & 0x0f
	nasKeySetId := (data[0] >> 4) & 0x07
	msg.EMMMessage.ExtServiceType = &serviceType
	msg.EMMMessage.NASKeySetId = &nasKeySetId

	return msg, nil
}

// Decode Authentication Request
func decodeAuthRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	nasKeySetId := data[0] & 0x07
	msg.EMMMessage.NASKeySetId = &nasKeySetId

	// RAND (16 bytes)
	if len(data) > 17 {
		msg.EMMMessage.RAND = data[1:17]
	}

	// AUTN (16 bytes after length indicator)
	if len(data) > 18 {
		autnLen := int(data[17])
		if len(data) > 18+autnLen {
			msg.EMMMessage.AUTN = data[18 : 18+autnLen]
		}
	}

	return msg, nil
}

// Decode Authentication Response
func decodeAuthResponse(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	// RES length
	resLen := int(data[0])
	if len(data) > 1+resLen {
		msg.EMMMessage.RES = data[1 : 1+resLen]
	}

	return msg, nil
}

// Decode Identity Request
func decodeIdentityRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	idType := data[0] & 0x07
	msg.EMMMessage.IdentityType = &idType

	return msg, nil
}

// Decode Identity Response
func decodeIdentityResponse(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 2 {
		return msg, nil
	}

	idLen := int(data[0])
	if len(data) > 1+idLen {
		msg.EMMMessage.MobileId = data[1 : 1+idLen]
		if idLen > 0 {
			idType := data[1] & 0x07
			msg.EMMMessage.MobileIdType = &idType
		}
	}

	return msg, nil
}

// Decode Security Mode Command
func decodeSecModeCommand(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	nasSecAlg := data[0]
	msg.EMMMessage.NASSecAlg = &nasSecAlg

	// Skip to UE security capabilities (after NAS key set identifier)
	if len(data) > 2 {
		ueSecCapLen := int(data[2])
		if len(data) > 3+ueSecCapLen {
			msg.EMMMessage.ReplayedUESecCap = data[3 : 3+ueSecCapLen]
		}
	}

	return msg, nil
}

// Decode PDN Connectivity Request
func decodePDNConnectivityRequest(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 1 {
		return msg, nil
	}

	pdnType := data[0] & 0x07
	reqType := (data[0] >> 4) & 0x07
	msg.ESMMessage.PDNType = &pdnType
	msg.ESMMessage.RequestType = &reqType

	// Look for APN in optional IEs
	offset := 1
	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}
		iei := data[offset]
		if iei == 0x28 { // APN IEI
			offset++
			apnLen := int(data[offset])
			offset++
			if offset+apnLen <= len(data) {
				msg.ESMMessage.APN = data[offset : offset+apnLen]
			}
			break
		}
		// Skip other IEs
		offset++
		if offset < len(data) {
			ieLen := int(data[offset])
			offset += 1 + ieLen
		}
	}

	return msg, nil
}

// Decode Activate Default EPS Bearer Context Request
func decodeActivateDefaultBearerReq(data []byte, msg *NASMessage) (*NASMessage, error) {
	if len(data) < 2 {
		return msg, nil
	}

	// EPS QoS
	qosLen := int(data[0])
	if qosLen > 0 && len(data) > 1 {
		qci := data[1]
		msg.ESMMessage.QCI = &qci
	}

	return msg, nil
}

// GetMessageTypeName returns the human-readable name of a NAS message
func GetMessageTypeName(protocolDiscr, messageType uint8) string {
	if protocolDiscr == ProtocolDiscriminatorEPSMobilityManagement {
		if name, ok := EMMMessageTypeNames[messageType]; ok {
			return name
		}
		return fmt.Sprintf("EMMUnknown(0x%02x)", messageType)
	} else if protocolDiscr == ProtocolDiscriminatorEPSSessionManagement {
		if name, ok := ESMMessageTypeNames[messageType]; ok {
			return name
		}
		return fmt.Sprintf("ESMUnknown(0x%02x)", messageType)
	}
	return fmt.Sprintf("Unknown(0x%02x)", messageType)
}

// IsAttachRelated returns true if the message is related to Attach procedure
func IsAttachRelated(messageType uint8) bool {
	switch messageType {
	case EMMAttachRequest, EMMAttachAccept, EMMAttachComplete, EMMAttachReject:
		return true
	}
	return false
}

// IsDetachRelated returns true if the message is related to Detach procedure
func IsDetachRelated(messageType uint8) bool {
	switch messageType {
	case EMMDetachRequest, EMMDetachAccept:
		return true
	}
	return false
}

// IsTAURelated returns true if the message is related to TAU procedure
func IsTAURelated(messageType uint8) bool {
	switch messageType {
	case EMMTAURequest, EMMTAUAccept, EMMTAUComplete, EMMTAUReject:
		return true
	}
	return false
}

// IsServiceRelated returns true if the message is related to Service procedure
func IsServiceRelated(messageType uint8) bool {
	switch messageType {
	case EMMExtServiceReq, EMMCPServiceReq, EMMServiceReject:
		return true
	}
	return false
}
