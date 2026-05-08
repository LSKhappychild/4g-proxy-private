package s1apType

// S1AP Procedure Codes (from 3GPP TS 36.413)
const (
	ProcedureCodeHandoverPreparation    int64 = 0
	ProcedureCodeHandoverResourceAlloc  int64 = 1
	ProcedureCodeHandoverNotification   int64 = 2
	ProcedureCodePathSwitchRequest      int64 = 3
	ProcedureCodeHandoverCancel         int64 = 4
	ProcedureCodeERABSetup              int64 = 5
	ProcedureCodeERABModify             int64 = 6
	ProcedureCodeERABRelease            int64 = 7
	ProcedureCodeERABReleaseIndication  int64 = 8
	ProcedureCodeInitialContextSetup    int64 = 9
	ProcedureCodePaging                 int64 = 10
	ProcedureCodeDownlinkNASTransport   int64 = 11
	ProcedureCodeInitialUEMessage       int64 = 12
	ProcedureCodeUplinkNASTransport     int64 = 13
	ProcedureCodeReset                  int64 = 14
	ProcedureCodeErrorIndication        int64 = 15
	ProcedureCodeNASNonDeliveryInd      int64 = 16
	ProcedureCodeS1Setup                int64 = 17
	ProcedureCodeUEContextReleaseReq    int64 = 18
	ProcedureCodeDownlinkS1cdma2000     int64 = 19
	ProcedureCodeUplinkS1cdma2000       int64 = 20
	ProcedureCodeUEContextModification  int64 = 21
	ProcedureCodeUECapabilityInfoInd    int64 = 22
	ProcedureCodeUEContextRelease       int64 = 23
	ProcedureCodeENBStatusTransfer      int64 = 24
	ProcedureCodeMMEStatusTransfer      int64 = 25
	ProcedureCodeDeactivateTrace        int64 = 26
	ProcedureCodeTraceStart             int64 = 27
	ProcedureCodeTraceFailureIndication int64 = 28
	ProcedureCodeENBConfigUpdate        int64 = 29
	ProcedureCodeMMEConfigUpdate        int64 = 30
	ProcedureCodeLocationReportingCtrl  int64 = 31
	ProcedureCodeLocationReportFailInd  int64 = 32
	ProcedureCodeLocationReport         int64 = 33
	ProcedureCodeOverloadStart          int64 = 34
	ProcedureCodeOverloadStop           int64 = 35
	ProcedureCodeWriteReplaceWarning    int64 = 36
	ProcedureCodeENBDirectInfoTransfer  int64 = 37
	ProcedureCodeMMEDirectInfoTransfer  int64 = 38
	ProcedureCodePrivateMessage         int64 = 39
	ProcedureCodeENBConfigTransfer      int64 = 40
	ProcedureCodeMMEConfigTransfer      int64 = 41
	ProcedureCodeCellTrafficTrace       int64 = 42
	ProcedureCodeKill                   int64 = 43
	ProcedureCodeDownlinkUEAssocLPPa    int64 = 44
	ProcedureCodeUplinkUEAssocLPPa      int64 = 45
	ProcedureCodeDownlinkNonUEAssocLPPa int64 = 46
	ProcedureCodeUplinkNonUEAssocLPPa   int64 = 47
	ProcedureCodeUERadioCapMatch        int64 = 48
	ProcedureCodePWSRestartIndication   int64 = 49
	ProcedureCodeE_RABModifyIndication  int64 = 50
)

// Protocol IE IDs
const (
	ProtocolIEIDMMEUES1APID               int64 = 0
	ProtocolIEIDENBUES1APID               int64 = 8
	ProtocolIEIDNASPDU                    int64 = 26
	ProtocolIEIDTAI                       int64 = 67
	ProtocolIEIDEUTRANCGI                 int64 = 100
	ProtocolIEIDRRCEstablishmentCause     int64 = 134
	ProtocolIEIDGUMMEI                    int64 = 75
	ProtocolIEIDUESecurityCapabilities    int64 = 107
	ProtocolIEIDAggregateMaxBitrate       int64 = 66
	ProtocolIEIDERABToBeSetupListCtxtReq  int64 = 24
	ProtocolIEIDERABSetupListCtxtSURes    int64 = 51
	ProtocolIEIDSecurityKey               int64 = 73
	ProtocolIEIDCause                     int64 = 2
	ProtocolIEIDResetType                 int64 = 92
	ProtocolIEIDGlobalENBID               int64 = 59
	ProtocolIEIDSupportedTAs              int64 = 64
	ProtocolIEIDPagingDRX                 int64 = 137
	ProtocolIEIDCSGIdList                 int64 = 128
	ProtocolIEIDERABToBeSetupListBearerSU int64 = 16
	ProtocolIEIDERABSetupListBearerSURes  int64 = 28
	ProtocolIEIDERABToBeModifyListBearerMod int64 = 30
	ProtocolIEIDERABModifyListBearerMod   int64 = 31
	ProtocolIEIDERABList                  int64 = 17
	ProtocolIEIDERABItem                  int64 = 18
	ProtocolIEIDUEAggregateMaxBitrate     int64 = 66
	ProtocolIEIDERABToBeReleasedList      int64 = 33
	ProtocolIEIDUEIdentityIndexValue      int64 = 80
	ProtocolIEIDUEPagingID                int64 = 43
	ProtocolIEIDCNDomain                  int64 = 109
	ProtocolIEIDTAIList                   int64 = 46
	ProtocolIEIDSourceMMEUES1APID         int64 = 88
	ProtocolIEIDRelativeMMECapacity       int64 = 87
	ProtocolIEIDMMEName                   int64 = 61
	ProtocolIEIDServedGUMMEIs             int64 = 105
	ProtocolIEIDTimeToWait                int64 = 65
	ProtocolIEIDCriticalityDiagnostics    int64 = 58
	ProtocolIEIDENBName                   int64 = 60
	ProtocolIEIDDefaultPagingDRX          int64 = 137
)

// Criticality
const (
	CriticalityReject Criticality = 0
	CriticalityIgnore Criticality = 1
	CriticalityNotify Criticality = 2
)

type Criticality int64

// Presence
const (
	PresenceOptional    Presence = 0
	PresenceConditional Presence = 1
	PresenceMandatory   Presence = 2
)

type Presence int64

// Procedure Code to Name mapping
var ProcedureCodeNames = map[int64]string{
	ProcedureCodeHandoverPreparation:    "HandoverPreparation",
	ProcedureCodeHandoverResourceAlloc:  "HandoverResourceAllocation",
	ProcedureCodeHandoverNotification:   "HandoverNotification",
	ProcedureCodePathSwitchRequest:      "PathSwitchRequest",
	ProcedureCodeHandoverCancel:         "HandoverCancel",
	ProcedureCodeERABSetup:              "E-RABSetup",
	ProcedureCodeERABModify:             "E-RABModify",
	ProcedureCodeERABRelease:            "E-RABRelease",
	ProcedureCodeERABReleaseIndication:  "E-RABReleaseIndication",
	ProcedureCodeInitialContextSetup:    "InitialContextSetup",
	ProcedureCodePaging:                 "Paging",
	ProcedureCodeDownlinkNASTransport:   "DownlinkNASTransport",
	ProcedureCodeInitialUEMessage:       "InitialUEMessage",
	ProcedureCodeUplinkNASTransport:     "UplinkNASTransport",
	ProcedureCodeReset:                  "Reset",
	ProcedureCodeErrorIndication:        "ErrorIndication",
	ProcedureCodeNASNonDeliveryInd:      "NASNonDeliveryIndication",
	ProcedureCodeS1Setup:                "S1Setup",
	ProcedureCodeUEContextReleaseReq:    "UEContextReleaseRequest",
	ProcedureCodeDownlinkS1cdma2000:     "DownlinkS1cdma2000tunneling",
	ProcedureCodeUplinkS1cdma2000:       "UplinkS1cdma2000tunneling",
	ProcedureCodeUEContextModification:  "UEContextModification",
	ProcedureCodeUECapabilityInfoInd:    "UECapabilityInfoIndication",
	ProcedureCodeUEContextRelease:       "UEContextRelease",
	ProcedureCodeENBStatusTransfer:      "eNBStatusTransfer",
	ProcedureCodeMMEStatusTransfer:      "MMEStatusTransfer",
	ProcedureCodeDeactivateTrace:        "DeactivateTrace",
	ProcedureCodeTraceStart:             "TraceStart",
	ProcedureCodeTraceFailureIndication: "TraceFailureIndication",
	ProcedureCodeENBConfigUpdate:        "ENBConfigurationUpdate",
	ProcedureCodeMMEConfigUpdate:        "MMEConfigurationUpdate",
	ProcedureCodeLocationReportingCtrl:  "LocationReportingControl",
	ProcedureCodeLocationReportFailInd:  "LocationReportingFailureIndication",
	ProcedureCodeLocationReport:         "LocationReport",
	ProcedureCodeOverloadStart:          "OverloadStart",
	ProcedureCodeOverloadStop:           "OverloadStop",
	ProcedureCodeWriteReplaceWarning:    "WriteReplaceWarning",
	ProcedureCodeENBDirectInfoTransfer:  "eNBDirectInformationTransfer",
	ProcedureCodeMMEDirectInfoTransfer:  "MMEDirectInformationTransfer",
	ProcedureCodePrivateMessage:         "PrivateMessage",
	ProcedureCodeENBConfigTransfer:      "eNBConfigurationTransfer",
	ProcedureCodeMMEConfigTransfer:      "MMEConfigurationTransfer",
	ProcedureCodeCellTrafficTrace:       "CellTrafficTrace",
	ProcedureCodeKill:                   "Kill",
	ProcedureCodeDownlinkUEAssocLPPa:    "DownlinkUEAssociatedLPPaTransport",
	ProcedureCodeUplinkUEAssocLPPa:      "UplinkUEAssociatedLPPaTransport",
	ProcedureCodeDownlinkNonUEAssocLPPa: "DownlinkNonUEAssociatedLPPaTransport",
	ProcedureCodeUplinkNonUEAssocLPPa:   "UplinkNonUEAssociatedLPPaTransport",
	ProcedureCodeUERadioCapMatch:        "UERadioCapabilityMatch",
	ProcedureCodePWSRestartIndication:   "PWSRestartIndication",
	ProcedureCodeE_RABModifyIndication:  "E-RABModificationIndication",
}
