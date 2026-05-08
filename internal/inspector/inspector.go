package inspector

import (
	"4g-proxy/lib/nas/eps"
	"4g-proxy/lib/s1ap"
	"4g-proxy/lib/s1ap/s1apType"
	"fmt"
	"log"
	"strings"
)

// Direction indicates message direction
type Direction int

const (
	DirectionUplink   Direction = iota // eNB -> MME
	DirectionDownlink                  // MME -> eNB
)

func (d Direction) String() string {
	if d == DirectionUplink {
		return "UL"
	}
	return "DL"
}

// InspectionResult contains the analysis of an S1AP message
type InspectionResult struct {
	Direction     Direction
	S1APProcedure string
	ProcedureCode int64
	MessageType   string // "initiating", "successful", "unsuccessful"
	MMEUEID       *int64
	ENBUEID       *int64

	// NAS message info (if present)
	HasNAS          bool
	NASProtocol     string // "EMM" or "ESM"
	NASMessageType  string
	NASRawType      uint8

	// Specific signal types for drop control
	SignalType SignalType

	// Handover-specific details
	HandoverSubType HandoverSubType
	HandoverDetails *HandoverDetails

	// Human-readable summary
	Summary string
}

// HandoverDetails contains detailed information for handover messages
type HandoverDetails struct {
	// Handover type (0=intraLTE, 1=LTEtoUTRAN, 2=LTEtoGERAN, etc.)
	HandoverType *int
	// Cause category and value
	Cause string
	// Target info (for HandoverRequired/Request)
	TargetENBID   string
	TargetTAC     string
	TargetPLMN    string
	// Source info (for PathSwitch)
	SourceMMEUEID *int64
	// New cell info (for HandoverNotify)
	NewCellID string
	NewTAI    string
}

// SignalType represents different signal types that can be dropped
type SignalType int

const (
	SignalTypeUnknown SignalType = iota
	SignalTypeAttach
	SignalTypeDetach
	SignalTypeTAU
	SignalTypeServiceRequest
	SignalTypeUEContextRelease
	SignalTypePDNConnectivity
	SignalTypeHandover
	SignalTypeReset
	SignalTypePaging
)

func (s SignalType) String() string {
	switch s {
	case SignalTypeAttach:
		return "Attach"
	case SignalTypeDetach:
		return "Detach"
	case SignalTypeTAU:
		return "TAU"
	case SignalTypeServiceRequest:
		return "ServiceRequest"
	case SignalTypeUEContextRelease:
		return "UEContextRelease"
	case SignalTypePDNConnectivity:
		return "PDNConnectivity"
	case SignalTypeHandover:
		return "Handover"
	case SignalTypeReset:
		return "Reset"
	case SignalTypePaging:
		return "Paging"
	default:
		return "Unknown"
	}
}

// HandoverSubType provides detailed handover message classification
type HandoverSubType int

const (
	HandoverSubTypeNone HandoverSubType = iota
	HandoverSubTypeRequired           // Source eNB -> MME: Initiate HO
	HandoverSubTypeRequest            // MME -> Target eNB: Prepare target
	HandoverSubTypeRequestAck         // Target eNB -> MME: Target ready
	HandoverSubTypeCommand            // MME -> Source eNB: Execute HO
	HandoverSubTypePreparationFailure // MME -> Source eNB: Target rejected
	HandoverSubTypeFailure            // MME -> Source eNB: HO failed
	HandoverSubTypeCancel             // Source eNB -> MME: Cancel HO
	HandoverSubTypeCancelAck          // MME -> Source eNB: Cancel confirmed
	HandoverSubTypeNotify             // Target eNB -> MME: UE arrived
	HandoverSubTypeENBStatusTransfer  // Source eNB -> MME: PDCP status
	HandoverSubTypeMMEStatusTransfer  // MME -> Target eNB: PDCP status
	HandoverSubTypePathSwitchReq      // Target eNB -> MME: X2 HO path switch
	HandoverSubTypePathSwitchAck      // MME -> Target eNB: Path switch OK
	HandoverSubTypePathSwitchFailure  // MME -> Target eNB: Path switch failed
)

func (h HandoverSubType) String() string {
	switch h {
	case HandoverSubTypeRequired:
		return "HandoverRequired"
	case HandoverSubTypeRequest:
		return "HandoverRequest"
	case HandoverSubTypeRequestAck:
		return "HandoverRequestAck"
	case HandoverSubTypeCommand:
		return "HandoverCommand"
	case HandoverSubTypePreparationFailure:
		return "HandoverPreparationFailure"
	case HandoverSubTypeFailure:
		return "HandoverFailure"
	case HandoverSubTypeCancel:
		return "HandoverCancel"
	case HandoverSubTypeCancelAck:
		return "HandoverCancelAck"
	case HandoverSubTypeNotify:
		return "HandoverNotify"
	case HandoverSubTypeENBStatusTransfer:
		return "eNBStatusTransfer"
	case HandoverSubTypeMMEStatusTransfer:
		return "MMEStatusTransfer"
	case HandoverSubTypePathSwitchReq:
		return "PathSwitchRequest"
	case HandoverSubTypePathSwitchAck:
		return "PathSwitchRequestAck"
	case HandoverSubTypePathSwitchFailure:
		return "PathSwitchRequestFailure"
	default:
		return ""
	}
}

// Inspect analyzes an S1AP message and returns inspection results
func Inspect(data []byte, dir Direction) (*InspectionResult, error) {
	result := &InspectionResult{
		Direction:  dir,
		SignalType: SignalTypeUnknown,
	}

	// Decode S1AP message
	s1apMsg, err := s1ap.Decode(data)
	if err != nil {
		return result, fmt.Errorf("S1AP decode error: %w", err)
	}

	result.S1APProcedure = s1apMsg.ProcedureName
	result.ProcedureCode = s1apMsg.ProcedureCode
	result.MessageType = s1apMsg.MessageType
	result.MMEUEID = s1apMsg.MMEUEID
	result.ENBUEID = s1apMsg.ENBUEID

	// Determine signal type from S1AP procedure
	result.SignalType = getSignalTypeFromS1AP(s1apMsg.ProcedureCode)

	// For handover messages, get the specific subtype
	if result.SignalType == SignalTypeHandover {
		result.HandoverSubType = getHandoverSubType(s1apMsg.ProcedureCode, s1apMsg.MessageType)
		result.HandoverDetails = &HandoverDetails{}
		if s1apMsg.Cause != "" {
			result.HandoverDetails.Cause = s1apMsg.Cause
		}
	}

	// Decode NAS if present
	if len(s1apMsg.NASPDU) > 0 {
		nasMsg, nasErr := eps.Decode(s1apMsg.NASPDU)
		if nasErr == nil {
			result.HasNAS = true
			result.NASMessageType = nasMsg.MessageTypeName
			result.NASRawType = nasMsg.MessageType

			if nasMsg.ProtocolDiscr == eps.ProtocolDiscriminatorEPSMobilityManagement {
				result.NASProtocol = "EMM"
				// Refine signal type from NAS message
				result.SignalType = getSignalTypeFromEMM(nasMsg.MessageType)
			} else if nasMsg.ProtocolDiscr == eps.ProtocolDiscriminatorEPSSessionManagement {
				result.NASProtocol = "ESM"
				result.SignalType = getSignalTypeFromESM(nasMsg.MessageType)
			}
		}
	}

	// Build summary
	result.Summary = buildSummary(result)

	return result, nil
}

func getSignalTypeFromS1AP(procCode int64) SignalType {
	switch procCode {
	case s1apType.ProcedureCodeUEContextRelease, s1apType.ProcedureCodeUEContextReleaseReq:
		return SignalTypeUEContextRelease
	case s1apType.ProcedureCodeHandoverPreparation, s1apType.ProcedureCodeHandoverResourceAlloc,
		s1apType.ProcedureCodeHandoverNotification, s1apType.ProcedureCodeHandoverCancel,
		s1apType.ProcedureCodeENBStatusTransfer, s1apType.ProcedureCodeMMEStatusTransfer:
		return SignalTypeHandover
	case s1apType.ProcedureCodePathSwitchRequest:
		return SignalTypeHandover
	case s1apType.ProcedureCodeReset:
		return SignalTypeReset
	case s1apType.ProcedureCodePaging:
		return SignalTypePaging
	default:
		return SignalTypeUnknown
	}
}

// getHandoverSubType determines the specific handover message type
func getHandoverSubType(procCode int64, msgType string) HandoverSubType {
	switch procCode {
	case s1apType.ProcedureCodeHandoverPreparation: // Code 0
		switch msgType {
		case "initiating":
			return HandoverSubTypeRequired
		case "successful":
			return HandoverSubTypeCommand
		case "unsuccessful":
			return HandoverSubTypePreparationFailure
		}
	case s1apType.ProcedureCodeHandoverResourceAlloc: // Code 1
		switch msgType {
		case "initiating":
			return HandoverSubTypeRequest
		case "successful":
			return HandoverSubTypeRequestAck
		case "unsuccessful":
			return HandoverSubTypeFailure
		}
	case s1apType.ProcedureCodeHandoverNotification: // Code 2
		return HandoverSubTypeNotify
	case s1apType.ProcedureCodePathSwitchRequest: // Code 3
		switch msgType {
		case "initiating":
			return HandoverSubTypePathSwitchReq
		case "successful":
			return HandoverSubTypePathSwitchAck
		case "unsuccessful":
			return HandoverSubTypePathSwitchFailure
		}
	case s1apType.ProcedureCodeHandoverCancel: // Code 4
		switch msgType {
		case "initiating":
			return HandoverSubTypeCancel
		case "successful":
			return HandoverSubTypeCancelAck
		}
	case s1apType.ProcedureCodeENBStatusTransfer: // Code 24
		return HandoverSubTypeENBStatusTransfer
	case s1apType.ProcedureCodeMMEStatusTransfer: // Code 25
		return HandoverSubTypeMMEStatusTransfer
	}
	return HandoverSubTypeNone
}

func getSignalTypeFromEMM(msgType uint8) SignalType {
	switch {
	case eps.IsAttachRelated(msgType):
		return SignalTypeAttach
	case eps.IsDetachRelated(msgType):
		return SignalTypeDetach
	case eps.IsTAURelated(msgType):
		return SignalTypeTAU
	case eps.IsServiceRelated(msgType):
		return SignalTypeServiceRequest
	case msgType == 0: // Service Request (special format)
		return SignalTypeServiceRequest
	default:
		return SignalTypeUnknown
	}
}

func getSignalTypeFromESM(msgType uint8) SignalType {
	switch msgType {
	case eps.ESMPDNConnectivityReq, eps.ESMPDNConnectivityRej,
		eps.ESMPDNDisconnectReq, eps.ESMPDNDisconnectRej:
		return SignalTypePDNConnectivity
	default:
		return SignalTypeUnknown
	}
}

func buildSummary(r *InspectionResult) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("[%s]", r.Direction))

	// For handover messages, use the detailed subtype name
	if r.SignalType == SignalTypeHandover && r.HandoverSubType != HandoverSubTypeNone {
		parts = append(parts, r.HandoverSubType.String())
	} else {
		parts = append(parts, r.S1APProcedure)
	}

	if r.HasNAS {
		parts = append(parts, fmt.Sprintf("(%s: %s)", r.NASProtocol, r.NASMessageType))
	}

	if r.MMEUEID != nil || r.ENBUEID != nil {
		ueInfo := "UE:"
		if r.MMEUEID != nil {
			ueInfo += fmt.Sprintf(" MME-UE-ID=%d", *r.MMEUEID)
		}
		if r.ENBUEID != nil {
			ueInfo += fmt.Sprintf(" eNB-UE-ID=%d", *r.ENBUEID)
		}
		parts = append(parts, ueInfo)
	}

	// Add cause for handover messages if present
	if r.HandoverDetails != nil && r.HandoverDetails.Cause != "" {
		parts = append(parts, fmt.Sprintf("Cause: %s", r.HandoverDetails.Cause))
	}

	if r.SignalType != SignalTypeUnknown {
		parts = append(parts, fmt.Sprintf("[Signal: %s]", r.SignalType))
	}

	return strings.Join(parts, " ")
}

// LogInspection logs the inspection result
func LogInspection(result *InspectionResult) {
	log.Printf("S1AP %s", result.Summary)
}

// Inspector provides continuous message inspection with drop control
type Inspector struct {
	verbose bool
}

// NewInspector creates a new Inspector
func NewInspector(verbose bool) *Inspector {
	return &Inspector{
		verbose: verbose,
	}
}

// InspectAndLog inspects a message and logs the result
func (i *Inspector) InspectAndLog(data []byte, dir Direction) (*InspectionResult, error) {
	result, err := Inspect(data, dir)
	if err != nil {
		if i.verbose {
			log.Printf("Inspection error: %v", err)
		}
		return result, err
	}

	if i.verbose {
		LogInspection(result)
	}

	return result, nil
}
