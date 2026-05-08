package s1apType

import "github.com/free5gc/aper"

// S1AP-PDU
type S1APPDU struct {
	Present       int
	InitiatingMsg *InitiatingMessage
	SuccessfulMsg *SuccessfulOutcome
	UnsuccessMsg  *UnsuccessfulOutcome
}

const (
	S1APPDUPresentNothing = iota
	S1APPDUPresentInitiatingMessage
	S1APPDUPresentSuccessfulOutcome
	S1APPDUPresentUnsuccessfulOutcome
)

// InitiatingMessage
type InitiatingMessage struct {
	ProcedureCode ProcedureCode
	Criticality   Criticality
	Value         InitiatingMessageValue
}

type ProcedureCode struct {
	Value int64 `aper:"valueLB:0,valueUB:255"`
}

type InitiatingMessageValue struct {
	Present                     int
	InitialUEMessage            *InitialUEMessage
	UplinkNASTransport          *UplinkNASTransport
	DownlinkNASTransport        *DownlinkNASTransport
	InitialContextSetupRequest  *InitialContextSetupRequest
	UEContextReleaseCommand     *UEContextReleaseCommand
	UEContextReleaseRequest     *UEContextReleaseRequest
	ERABSetupRequest            *ERABSetupRequest
	S1SetupRequest              *S1SetupRequest
	HandoverRequired            *HandoverRequired
	HandoverRequest             *HandoverRequest
	PathSwitchRequest           *PathSwitchRequest
	Reset                       *Reset
	Paging                      *Paging
}

const (
	InitiatingMessagePresentNothing = iota
	InitiatingMessagePresentInitialUEMessage
	InitiatingMessagePresentUplinkNASTransport
	InitiatingMessagePresentDownlinkNASTransport
	InitiatingMessagePresentInitialContextSetupRequest
	InitiatingMessagePresentUEContextReleaseCommand
	InitiatingMessagePresentUEContextReleaseRequest
	InitiatingMessagePresentERABSetupRequest
	InitiatingMessagePresentS1SetupRequest
	InitiatingMessagePresentHandoverRequired
	InitiatingMessagePresentHandoverRequest
	InitiatingMessagePresentPathSwitchRequest
	InitiatingMessagePresentReset
	InitiatingMessagePresentPaging
)

// SuccessfulOutcome
type SuccessfulOutcome struct {
	ProcedureCode ProcedureCode
	Criticality   Criticality
	Value         SuccessfulOutcomeValue
}

type SuccessfulOutcomeValue struct {
	Present                      int
	InitialContextSetupResponse  *InitialContextSetupResponse
	UEContextReleaseComplete     *UEContextReleaseComplete
	ERABSetupResponse            *ERABSetupResponse
	S1SetupResponse              *S1SetupResponse
	HandoverCommand              *HandoverCommand
	HandoverRequestAcknowledge   *HandoverRequestAcknowledge
	PathSwitchRequestAcknowledge *PathSwitchRequestAcknowledge
	ResetAcknowledge             *ResetAcknowledge
}

const (
	SuccessfulOutcomePresentNothing = iota
	SuccessfulOutcomePresentInitialContextSetupResponse
	SuccessfulOutcomePresentUEContextReleaseComplete
	SuccessfulOutcomePresentERABSetupResponse
	SuccessfulOutcomePresentS1SetupResponse
	SuccessfulOutcomePresentHandoverCommand
	SuccessfulOutcomePresentHandoverRequestAcknowledge
	SuccessfulOutcomePresentPathSwitchRequestAcknowledge
	SuccessfulOutcomePresentResetAcknowledge
)

// UnsuccessfulOutcome
type UnsuccessfulOutcome struct {
	ProcedureCode ProcedureCode
	Criticality   Criticality
	Value         UnsuccessfulOutcomeValue
}

type UnsuccessfulOutcomeValue struct {
	Present                     int
	InitialContextSetupFailure  *InitialContextSetupFailure
	S1SetupFailure              *S1SetupFailure
	HandoverPreparationFailure  *HandoverPreparationFailure
	HandoverFailure             *HandoverFailure
	PathSwitchRequestFailure    *PathSwitchRequestFailure
}

const (
	UnsuccessfulOutcomePresentNothing = iota
	UnsuccessfulOutcomePresentInitialContextSetupFailure
	UnsuccessfulOutcomePresentS1SetupFailure
	UnsuccessfulOutcomePresentHandoverPreparationFailure
	UnsuccessfulOutcomePresentHandoverFailure
	UnsuccessfulOutcomePresentPathSwitchRequestFailure
)

// Common IE types
type MMEUEID struct {
	Value int64 `aper:"valueLB:0,valueUB:4294967295"`
}

type ENBUEID struct {
	Value int64 `aper:"valueLB:0,valueUB:16777215"`
}

type NASPDU struct {
	Value []byte
}

type TAI struct {
	PLMNIdentity PLMNIdentity
	TAC          TAC
}

type PLMNIdentity struct {
	Value []byte `aper:"sizeLB:3,sizeUB:3"`
}

type TAC struct {
	Value []byte `aper:"sizeLB:2,sizeUB:2"`
}

type EUTRANGCI struct {
	PLMNIdentity PLMNIdentity
	CellIdentity CellIdentity
}

type CellIdentity struct {
	Value aper.BitString `aper:"sizeLB:28,sizeUB:28"`
}

type GUMMEI struct {
	PLMNIdentity PLMNIdentity
	MMEGroupID   MMEGroupID
	MMEC         MMEC
}

type MMEGroupID struct {
	Value []byte `aper:"sizeLB:2,sizeUB:2"`
}

type MMEC struct {
	Value []byte `aper:"sizeLB:1,sizeUB:1"`
}

type GlobalENBID struct {
	PLMNIdentity PLMNIdentity
	ENBID        ENBID
}

type ENBID struct {
	Present  int
	MacroENB *aper.BitString `aper:"sizeLB:20,sizeUB:20"`
	HomeENB  *aper.BitString `aper:"sizeLB:28,sizeUB:28"`
}

const (
	ENBIDPresentNothing = iota
	ENBIDPresentMacroENBID
	ENBIDPresentHomeENBID
)

type RRCEstablishmentCause int64

const (
	RRCEstablishmentCauseEmergency           RRCEstablishmentCause = 0
	RRCEstablishmentCauseHighPriorityAccess  RRCEstablishmentCause = 1
	RRCEstablishmentCauseMtAccess            RRCEstablishmentCause = 2
	RRCEstablishmentCauseMoSignalling        RRCEstablishmentCause = 3
	RRCEstablishmentCauseMoData              RRCEstablishmentCause = 4
	RRCEstablishmentCauseDelayTolerantAccess RRCEstablishmentCause = 5
)

// Cause
type Cause struct {
	Present         int
	RadioNetwork    *CauseRadioNetwork
	Transport       *CauseTransport
	NAS             *CauseNas
	Protocol        *CauseProtocol
	Misc            *CauseMisc
}

const (
	CausePresentNothing = iota
	CausePresentRadioNetwork
	CausePresentTransport
	CausePresentNAS
	CausePresentProtocol
	CausePresentMisc
)

type CauseRadioNetwork int64

const (
	CauseRadioNetworkUnspecified                          CauseRadioNetwork = 0
	CauseRadioNetworkTxResourcesUnavailable               CauseRadioNetwork = 1
	CauseRadioNetworkHandoverCancelled                    CauseRadioNetwork = 2
	CauseRadioNetworkSuccessfulHandover                   CauseRadioNetwork = 3
	CauseRadioNetworkReleaseDueToEutranGeneratedReason    CauseRadioNetwork = 4
	CauseRadioNetworkUserInactivity                       CauseRadioNetwork = 5
	CauseRadioNetworkRadioConnectionWithUeLost            CauseRadioNetwork = 6
	CauseRadioNetworkFailureInRadioInterfaceProcedure     CauseRadioNetwork = 7
	CauseRadioNetworkCSFallbackTriggered                  CauseRadioNetwork = 8
	CauseRadioNetworkUENotAvailableForPSService           CauseRadioNetwork = 9
	CauseRadioNetworkRadioResourcesNotAvailable           CauseRadioNetwork = 10
)

type CauseTransport int64

const (
	CauseTransportTransportResourceUnavailable CauseTransport = 0
	CauseTransportUnspecified                  CauseTransport = 1
)

type CauseNas int64

const (
	CauseNasNormalRelease     CauseNas = 0
	CauseNasAuthFailure       CauseNas = 1
	CauseNasDetach            CauseNas = 2
	CauseNasUnspecified       CauseNas = 3
	CauseNasCSGSubExpired     CauseNas = 4
)

type CauseProtocol int64

const (
	CauseProtocolTransferSyntaxError                        CauseProtocol = 0
	CauseProtocolAbstractSyntaxErrorReject                  CauseProtocol = 1
	CauseProtocolAbstractSyntaxErrorIgnoreAndNotify         CauseProtocol = 2
	CauseProtocolMessageNotCompatibleWithReceiverState      CauseProtocol = 3
	CauseProtocolSemanticError                              CauseProtocol = 4
	CauseProtocolAbstractSyntaxErrorFalselyConstructedMsg   CauseProtocol = 5
	CauseProtocolUnspecified                                CauseProtocol = 6
)

type CauseMisc int64

const (
	CauseMiscControlProcessingOverload          CauseMisc = 0
	CauseMiscNotEnoughUserPlaneProcessingRes    CauseMisc = 1
	CauseMiscHardwareFailure                    CauseMisc = 2
	CauseMiscOMIntervention                     CauseMisc = 3
	CauseMiscUnspecified                        CauseMisc = 4
	CauseMiscUnknownPLMN                        CauseMisc = 5
)

// ResetType
type ResetType struct {
	Present int
	S1InterfaceReset *ResetAll
	PartOfS1Interface *UEAssocLogicalS1ConnectionList
}

const (
	ResetTypePresentNothing = iota
	ResetTypePresentS1Interface
	ResetTypePresentPartOfS1Interface
)

type ResetAll int64

const (
	ResetAllResetAll ResetAll = 0
)

type UEAssocLogicalS1ConnectionList struct {
	List []UEAssocLogicalS1ConnectionItem
}

type UEAssocLogicalS1ConnectionItem struct {
	MMEUEID *MMEUEID
	ENBUEID *ENBUEID
}

// UE Security Capabilities
type UESecurityCapabilities struct {
	EncryptionAlgorithms aper.BitString `aper:"sizeLB:16,sizeUB:16"`
	IntegrityProtection  aper.BitString `aper:"sizeLB:16,sizeUB:16"`
}

// UE Aggregate Maximum Bit Rate
type UEAggregateMaximumBitrate struct {
	DL int64 `aper:"valueLB:0,valueUB:10000000000"`
	UL int64 `aper:"valueLB:0,valueUB:10000000000"`
}

// Security Key
type SecurityKey struct {
	Value aper.BitString `aper:"sizeLB:256,sizeUB:256"`
}

// Paging DRX
type PagingDRX int64

const (
	PagingDRXv32  PagingDRX = 0
	PagingDRXv64  PagingDRX = 1
	PagingDRXv128 PagingDRX = 2
	PagingDRXv256 PagingDRX = 3
)

// CN Domain
type CNDomain int64

const (
	CNDomainPS CNDomain = 0
	CNDomainCS CNDomain = 1
)

// UE Paging ID
type UEPagingID struct {
	Present int
	SIMSI   *SIMSI
	IMSI    *IMSI
}

const (
	UEPagingIDPresentNothing = iota
	UEPagingIDPresentSIMSI
	UEPagingIDPresentIMSI
)

type SIMSI struct {
	Value []byte
}

type IMSI struct {
	Value []byte
}

// Supported TAs
type SupportedTAs struct {
	List []SupportedTAItem
}

type SupportedTAItem struct {
	TAC            TAC
	BroadcastPLMNs BroadcastPLMNs
}

type BroadcastPLMNs struct {
	List []PLMNIdentity
}

// Served GUMMEIs
type ServedGUMMEIs struct {
	List []ServedGUMMEIsItem
}

type ServedGUMMEIsItem struct {
	ServedPLMNs    ServedPLMNs
	ServedGroupIDs ServedGroupIDs
	ServedMMECs    ServedMMECs
}

type ServedPLMNs struct {
	List []PLMNIdentity
}

type ServedGroupIDs struct {
	List []MMEGroupID
}

type ServedMMECs struct {
	List []MMEC
}
