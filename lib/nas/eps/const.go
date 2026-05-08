package eps

// Security Header Types
const (
	SecurityHeaderTypePlainNAS                                          uint8 = 0
	SecurityHeaderTypeIntegrityProtected                                uint8 = 1
	SecurityHeaderTypeIntegrityProtectedAndCiphered                     uint8 = 2
	SecurityHeaderTypeIntegrityProtectedWithNewEPSSecurityContext       uint8 = 3
	SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSSecCtx     uint8 = 4
	SecurityHeaderTypeServiceRequestHeader                              uint8 = 12
)

// Protocol Discriminator
const (
	ProtocolDiscriminatorEPSMobilityManagement uint8 = 0x07
	ProtocolDiscriminatorEPSSessionManagement  uint8 = 0x02
)

// EMM Message Types (from 3GPP TS 24.301)
const (
	// Attach
	EMMAttachRequest    uint8 = 65  // 0x41
	EMMAttachAccept     uint8 = 66  // 0x42
	EMMAttachComplete   uint8 = 67  // 0x43
	EMMAttachReject     uint8 = 68  // 0x44

	// Detach
	EMMDetachRequest    uint8 = 69  // 0x45
	EMMDetachAccept     uint8 = 70  // 0x46

	// Tracking Area Update
	EMMTAURequest       uint8 = 72  // 0x48
	EMMTAUAccept        uint8 = 73  // 0x49
	EMMTAUComplete      uint8 = 74  // 0x4a
	EMMTAUReject        uint8 = 75  // 0x4b

	// Extended Service Request
	EMMExtServiceReq    uint8 = 76  // 0x4c

	// Control Plane Service Request
	EMMCPServiceReq     uint8 = 77  // 0x4d

	// Service Reject
	EMMServiceReject    uint8 = 78  // 0x4e

	// GUTI Reallocation
	EMMGUTIReallocCmd   uint8 = 80  // 0x50
	EMMGUTIReallocCmpl  uint8 = 81  // 0x51

	// Authentication
	EMMAuthRequest      uint8 = 82  // 0x52
	EMMAuthResponse     uint8 = 83  // 0x53
	EMMAuthReject       uint8 = 84  // 0x54
	EMMAuthFailure      uint8 = 92  // 0x5c

	// Identity
	EMMIdentityRequest  uint8 = 85  // 0x55
	EMMIdentityResponse uint8 = 86  // 0x56

	// Security Mode
	EMMSecModeCommand   uint8 = 93  // 0x5d
	EMMSecModeComplete  uint8 = 94  // 0x5e
	EMMSecModeReject    uint8 = 95  // 0x5f

	// EMM Status
	EMMStatus           uint8 = 96  // 0x60

	// EMM Information
	EMMInformation      uint8 = 97  // 0x61

	// Downlink/Uplink NAS Transport
	EMMDownlinkNASTransport uint8 = 98  // 0x62
	EMMUplinkNASTransport   uint8 = 99  // 0x63

	// CS Service Notification
	EMMCSServiceNotif   uint8 = 100 // 0x64

	// Downlink Generic NAS Transport
	EMMDownlinkGenericNASTransport uint8 = 104 // 0x68
	EMMUplinkGenericNASTransport   uint8 = 105 // 0x69
)

// ESM Message Types (from 3GPP TS 24.301)
const (
	// Default EPS Bearer Context
	ESMActivateDefaultEPSBearerCtxReq uint8 = 193 // 0xc1
	ESMActivateDefaultEPSBearerCtxAcc uint8 = 194 // 0xc2
	ESMActivateDefaultEPSBearerCtxRej uint8 = 195 // 0xc3

	// Dedicated EPS Bearer Context
	ESMActivateDedicatedEPSBearerCtxReq uint8 = 197 // 0xc5
	ESMActivateDedicatedEPSBearerCtxAcc uint8 = 198 // 0xc6
	ESMActivateDedicatedEPSBearerCtxRej uint8 = 199 // 0xc7

	// Modify EPS Bearer Context
	ESMModifyEPSBearerCtxReq uint8 = 201 // 0xc9
	ESMModifyEPSBearerCtxAcc uint8 = 202 // 0xca
	ESMModifyEPSBearerCtxRej uint8 = 203 // 0xcb

	// Deactivate EPS Bearer Context
	ESMDeactivateEPSBearerCtxReq uint8 = 205 // 0xcd
	ESMDeactivateEPSBearerCtxAcc uint8 = 206 // 0xce

	// PDN Connectivity
	ESMPDNConnectivityReq uint8 = 208 // 0xd0
	ESMPDNConnectivityRej uint8 = 209 // 0xd1

	// PDN Disconnect
	ESMPDNDisconnectReq uint8 = 210 // 0xd2
	ESMPDNDisconnectRej uint8 = 211 // 0xd3

	// Bearer Resource Allocation
	ESMBearerResourceAllocReq uint8 = 212 // 0xd4
	ESMBearerResourceAllocRej uint8 = 213 // 0xd5

	// Bearer Resource Modification
	ESMBearerResourceModReq uint8 = 214 // 0xd6
	ESMBearerResourceModRej uint8 = 215 // 0xd7

	// ESM Information
	ESMInformationReq uint8 = 217 // 0xd9
	ESMInformationRsp uint8 = 218 // 0xda

	// Notification
	ESMNotification uint8 = 219 // 0xdb

	// ESM Status
	ESMStatus uint8 = 232 // 0xe8

	// Remote UE Report
	ESMRemoteUEReport    uint8 = 233 // 0xe9
	ESMRemoteUERspRsp    uint8 = 234 // 0xea

	// ESM Data Transport
	ESMDataTransport uint8 = 235 // 0xeb
)

// EMM Message Type Names
var EMMMessageTypeNames = map[uint8]string{
	EMMAttachRequest:    "AttachRequest",
	EMMAttachAccept:     "AttachAccept",
	EMMAttachComplete:   "AttachComplete",
	EMMAttachReject:     "AttachReject",
	EMMDetachRequest:    "DetachRequest",
	EMMDetachAccept:     "DetachAccept",
	EMMTAURequest:       "TrackingAreaUpdateRequest",
	EMMTAUAccept:        "TrackingAreaUpdateAccept",
	EMMTAUComplete:      "TrackingAreaUpdateComplete",
	EMMTAUReject:        "TrackingAreaUpdateReject",
	EMMExtServiceReq:    "ExtendedServiceRequest",
	EMMCPServiceReq:     "ControlPlaneServiceRequest",
	EMMServiceReject:    "ServiceReject",
	EMMGUTIReallocCmd:   "GUTIReallocationCommand",
	EMMGUTIReallocCmpl:  "GUTIReallocationComplete",
	EMMAuthRequest:      "AuthenticationRequest",
	EMMAuthResponse:     "AuthenticationResponse",
	EMMAuthReject:       "AuthenticationReject",
	EMMAuthFailure:      "AuthenticationFailure",
	EMMIdentityRequest:  "IdentityRequest",
	EMMIdentityResponse: "IdentityResponse",
	EMMSecModeCommand:   "SecurityModeCommand",
	EMMSecModeComplete:  "SecurityModeComplete",
	EMMSecModeReject:    "SecurityModeReject",
	EMMStatus:           "EMMStatus",
	EMMInformation:      "EMMInformation",
	EMMDownlinkNASTransport: "DownlinkNASTransport",
	EMMUplinkNASTransport:   "UplinkNASTransport",
	EMMCSServiceNotif:   "CSServiceNotification",
	EMMDownlinkGenericNASTransport: "DownlinkGenericNASTransport",
	EMMUplinkGenericNASTransport:   "UplinkGenericNASTransport",
}

// ESM Message Type Names
var ESMMessageTypeNames = map[uint8]string{
	ESMActivateDefaultEPSBearerCtxReq: "ActivateDefaultEPSBearerContextRequest",
	ESMActivateDefaultEPSBearerCtxAcc: "ActivateDefaultEPSBearerContextAccept",
	ESMActivateDefaultEPSBearerCtxRej: "ActivateDefaultEPSBearerContextReject",
	ESMActivateDedicatedEPSBearerCtxReq: "ActivateDedicatedEPSBearerContextRequest",
	ESMActivateDedicatedEPSBearerCtxAcc: "ActivateDedicatedEPSBearerContextAccept",
	ESMActivateDedicatedEPSBearerCtxRej: "ActivateDedicatedEPSBearerContextReject",
	ESMModifyEPSBearerCtxReq: "ModifyEPSBearerContextRequest",
	ESMModifyEPSBearerCtxAcc: "ModifyEPSBearerContextAccept",
	ESMModifyEPSBearerCtxRej: "ModifyEPSBearerContextReject",
	ESMDeactivateEPSBearerCtxReq: "DeactivateEPSBearerContextRequest",
	ESMDeactivateEPSBearerCtxAcc: "DeactivateEPSBearerContextAccept",
	ESMPDNConnectivityReq: "PDNConnectivityRequest",
	ESMPDNConnectivityRej: "PDNConnectivityReject",
	ESMPDNDisconnectReq:   "PDNDisconnectRequest",
	ESMPDNDisconnectRej:   "PDNDisconnectReject",
	ESMBearerResourceAllocReq: "BearerResourceAllocationRequest",
	ESMBearerResourceAllocRej: "BearerResourceAllocationReject",
	ESMBearerResourceModReq:   "BearerResourceModificationRequest",
	ESMBearerResourceModRej:   "BearerResourceModificationReject",
	ESMInformationReq:    "ESMInformationRequest",
	ESMInformationRsp:    "ESMInformationResponse",
	ESMNotification:      "ESMNotification",
	ESMStatus:            "ESMStatus",
	ESMRemoteUEReport:    "RemoteUEReport",
	ESMRemoteUERspRsp:    "RemoteUEReportResponse",
	ESMDataTransport:     "ESMDataTransport",
}

// Attach Type values
const (
	AttachTypeEPSAttach         uint8 = 1
	AttachTypeCombinedAttach    uint8 = 2
	AttachTypeEPSEmergency      uint8 = 6
)

// Detach Type values
const (
	DetachTypeEPSDetach                uint8 = 1
	DetachTypeIMSIDetach               uint8 = 2
	DetachTypeCombinedDetach           uint8 = 3
)

// TAU Request Type values
const (
	TAUTypeTA                    uint8 = 0
	TAUTypeCombinedTA            uint8 = 1
	TAUTypeCombinedTAwithIMSI    uint8 = 2
	TAUTypePeriodic              uint8 = 3
)

// Service Type values
const (
	ServiceTypeMOCSFB            uint8 = 0
	ServiceTypeMTCSFB            uint8 = 1
	ServiceTypeMOCSFBEmergency   uint8 = 2
	ServiceTypePacketServices    uint8 = 8
)
