package s1apType

// Protocol IE container generic type
type ProtocolIEField struct {
	Id          int64
	Criticality Criticality
	Value       interface{}
}

// InitialUEMessage
type InitialUEMessage struct {
	ProtocolIEs InitialUEMessageIEs
}

type InitialUEMessageIEs struct {
	ENBUEID               *ENBUEID
	NASPDU                *NASPDU
	TAI                   *TAI
	EUTRANGCI             *EUTRANGCI
	RRCEstablishmentCause *RRCEstablishmentCause
	STMSI                 *STMSI
	CSGId                 *CSGId
	GUMMEI                *GUMMEI
	CellAccessMode        *CellAccessMode
	GWTransportLayerAddr  *TransportLayerAddress
	RelayNodeIndicator    *RelayNodeIndicator
}

type STMSI struct {
	MMEC   MMEC
	MTMSI  MTMSI
}

type MTMSI struct {
	Value []byte `aper:"sizeLB:4,sizeUB:4"`
}

type CSGId struct {
	Value []byte `aper:"sizeLB:27,sizeUB:27"`
}

type CellAccessMode int64

const (
	CellAccessModeHybrid CellAccessMode = 0
)

type TransportLayerAddress struct {
	Value []byte
}

type RelayNodeIndicator int64

const (
	RelayNodeIndicatorTrue RelayNodeIndicator = 0
)

// UplinkNASTransport
type UplinkNASTransport struct {
	ProtocolIEs UplinkNASTransportIEs
}

type UplinkNASTransportIEs struct {
	MMEUEID   *MMEUEID
	ENBUEID   *ENBUEID
	NASPDU    *NASPDU
	EUTRANGCI *EUTRANGCI
	TAI       *TAI
	// Additional optional IEs
	GWTransportLayerAddr  *TransportLayerAddress
}

// DownlinkNASTransport
type DownlinkNASTransport struct {
	ProtocolIEs DownlinkNASTransportIEs
}

type DownlinkNASTransportIEs struct {
	MMEUEID    *MMEUEID
	ENBUEID    *ENBUEID
	NASPDU     *NASPDU
	HandoverRestrictionList *HandoverRestrictionList
	SubscriberProfileIDforRFP *SubscriberProfileIDforRFP
}

type HandoverRestrictionList struct {
	ServingPLMN       PLMNIdentity
	EquivalentPLMNs   *EquivalentPLMNs
	ForbiddenTAs      *ForbiddenTAs
	ForbiddenLAs      *ForbiddenLAs
	ForbiddenInterRAT *ForbiddenInterRAT
}

type EquivalentPLMNs struct {
	List []PLMNIdentity
}

type ForbiddenTAs struct {
	List []ForbiddenTAsItem
}

type ForbiddenTAsItem struct {
	PLMNIdentity PLMNIdentity
	ForbiddenTACs ForbiddenTACs
}

type ForbiddenTACs struct {
	List []TAC
}

type ForbiddenLAs struct {
	List []ForbiddenLAsItem
}

type ForbiddenLAsItem struct {
	PLMNIdentity PLMNIdentity
	ForbiddenLACs ForbiddenLACs
}

type ForbiddenLACs struct {
	List []LAC
}

type LAC struct {
	Value []byte `aper:"sizeLB:2,sizeUB:2"`
}

type ForbiddenInterRAT int64

const (
	ForbiddenInterRATAll     ForbiddenInterRAT = 0
	ForbiddenInterRATGERAN   ForbiddenInterRAT = 1
	ForbiddenInterRATUTRANWB ForbiddenInterRAT = 2
	ForbiddenInterRATCDMA2000 ForbiddenInterRAT = 3
)

type SubscriberProfileIDforRFP struct {
	Value int64 `aper:"valueLB:1,valueUB:256"`
}

// InitialContextSetupRequest
type InitialContextSetupRequest struct {
	ProtocolIEs InitialContextSetupRequestIEs
}

type InitialContextSetupRequestIEs struct {
	MMEUEID                     *MMEUEID
	ENBUEID                     *ENBUEID
	UEAggregateMaximumBitrate   *UEAggregateMaximumBitrate
	ERABToBeSetupListCtxtSUReq  *ERABToBeSetupListCtxtSUReq
	UESecurityCapabilities      *UESecurityCapabilities
	SecurityKey                 *SecurityKey
	TraceActivation             *TraceActivation
	HandoverRestrictionList     *HandoverRestrictionList
	UERadioCapability           *UERadioCapability
	SubscriberProfileIDforRFP   *SubscriberProfileIDforRFP
	CSFallbackIndicator         *CSFallbackIndicator
	SRVCCOperationPossible      *SRVCCOperationPossible
	CSGMembershipStatus         *CSGMembershipStatus
	LAI                         *LAI
	GUMMEI                      *GUMMEI
	MMEUEID2                    *MMEUEID
	ManagementBasedMDTAllowed   *ManagementBasedMDTAllowed
}

type ERABToBeSetupListCtxtSUReq struct {
	List []ERABToBeSetupItemCtxtSUReq
}

type ERABToBeSetupItemCtxtSUReq struct {
	ERABID                     ERABID
	ERABLevelQoSParameters     ERABLevelQoSParameters
	TransportLayerAddress      TransportLayerAddress
	GTPTEID                    GTPTEID
	NASPDU                     *NASPDU
	CorrelationID              *CorrelationID
}

type ERABID struct {
	Value int64 `aper:"valueLB:0,valueUB:15"`
}

type ERABLevelQoSParameters struct {
	QCI                             QCI
	AllocationRetentionPriority     AllocationRetentionPriority
	GBRQoSInformation               *GBRQoSInformation
}

type QCI struct {
	Value int64 `aper:"valueLB:0,valueUB:255"`
}

type AllocationRetentionPriority struct {
	PriorityLevel            PriorityLevel
	PreEmptionCapability     PreEmptionCapability
	PreEmptionVulnerability  PreEmptionVulnerability
}

type PriorityLevel struct {
	Value int64 `aper:"valueLB:0,valueUB:15"`
}

type PreEmptionCapability int64

const (
	PreEmptionCapabilityShallNotTrigger PreEmptionCapability = 0
	PreEmptionCapabilityMayTrigger      PreEmptionCapability = 1
)

type PreEmptionVulnerability int64

const (
	PreEmptionVulnerabilityNotPreemptable PreEmptionVulnerability = 0
	PreEmptionVulnerabilityPreemptable    PreEmptionVulnerability = 1
)

type GBRQoSInformation struct {
	ERABMaxBitrateDL     BitRate
	ERABMaxBitrateUL     BitRate
	ERABGuaranteedBitrateDL BitRate
	ERABGuaranteedBitrateUL BitRate
}

type BitRate struct {
	Value int64 `aper:"valueLB:0,valueUB:10000000000"`
}

type GTPTEID struct {
	Value []byte `aper:"sizeLB:4,sizeUB:4"`
}

type CorrelationID struct {
	Value []byte `aper:"sizeLB:4,sizeUB:4"`
}

type TraceActivation struct {
	EUTRANTraceID              EUTRANTraceID
	InterfacesToTrace          InterfacesToTrace
	TraceDepth                 TraceDepth
	TraceCollectionEntityIPAddr TransportLayerAddress
}

type EUTRANTraceID struct {
	Value []byte `aper:"sizeLB:8,sizeUB:8"`
}

type InterfacesToTrace struct {
	Value []byte `aper:"sizeLB:8,sizeUB:8"`
}

type TraceDepth int64

const (
	TraceDepthMinimum                       TraceDepth = 0
	TraceDepthMedium                        TraceDepth = 1
	TraceDepthMaximum                       TraceDepth = 2
	TraceDepthMinimumWithoutVendorExtension TraceDepth = 3
	TraceDepthMediumWithoutVendorExtension  TraceDepth = 4
	TraceDepthMaximumWithoutVendorExtension TraceDepth = 5
)

type UERadioCapability struct {
	Value []byte
}

type CSFallbackIndicator int64

const (
	CSFallbackIndicatorRequired    CSFallbackIndicator = 0
	CSFallbackIndicatorHighPriority CSFallbackIndicator = 1
)

type SRVCCOperationPossible int64

const (
	SRVCCOperationPossiblePossible SRVCCOperationPossible = 0
)

type CSGMembershipStatus int64

const (
	CSGMembershipStatusMember    CSGMembershipStatus = 0
	CSGMembershipStatusNotMember CSGMembershipStatus = 1
)

type LAI struct {
	PLMNIdentity PLMNIdentity
	LAC          LAC
}

type ManagementBasedMDTAllowed int64

const (
	ManagementBasedMDTAllowedAllowed ManagementBasedMDTAllowed = 0
)

// InitialContextSetupResponse
type InitialContextSetupResponse struct {
	ProtocolIEs InitialContextSetupResponseIEs
}

type InitialContextSetupResponseIEs struct {
	MMEUEID                    *MMEUEID
	ENBUEID                    *ENBUEID
	ERABSetupListCtxtSURes     *ERABSetupListCtxtSURes
	ERABFailedToSetupListCtxtSURes *ERABList
	CriticalityDiagnostics     *CriticalityDiagnostics
}

type ERABSetupListCtxtSURes struct {
	List []ERABSetupItemCtxtSURes
}

type ERABSetupItemCtxtSURes struct {
	ERABID                ERABID
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
}

type ERABList struct {
	List []ERABItem
}

type ERABItem struct {
	ERABID ERABID
	Cause  Cause
}

type CriticalityDiagnostics struct {
	ProcedureCode          *ProcedureCode
	TriggeringMessage      *TriggeringMessage
	ProcedureCriticality   *Criticality
	IEsCriticalityDiagnostics *CriticalityDiagnosticsIEList
}

type TriggeringMessage int64

const (
	TriggeringMessageInitiating   TriggeringMessage = 0
	TriggeringMessageSuccessful   TriggeringMessage = 1
	TriggeringMessageUnsuccessful TriggeringMessage = 2
)

type CriticalityDiagnosticsIEList struct {
	List []CriticalityDiagnosticsIEItem
}

type CriticalityDiagnosticsIEItem struct {
	IECriticality Criticality
	IEID          int64
	TypeOfError   TypeOfError
}

type TypeOfError int64

const (
	TypeOfErrorNotUnderstood TypeOfError = 0
	TypeOfErrorMissing       TypeOfError = 1
)

// InitialContextSetupFailure
type InitialContextSetupFailure struct {
	ProtocolIEs InitialContextSetupFailureIEs
}

type InitialContextSetupFailureIEs struct {
	MMEUEID                *MMEUEID
	ENBUEID                *ENBUEID
	Cause                  *Cause
	CriticalityDiagnostics *CriticalityDiagnostics
}

// UEContextReleaseCommand
type UEContextReleaseCommand struct {
	ProtocolIEs UEContextReleaseCommandIEs
}

type UEContextReleaseCommandIEs struct {
	UEID  *UEID
	Cause *Cause
}

type UEID struct {
	Present int
	MMEUEID *MMEUEID
	UEIDPair *UEIDPair
}

const (
	UEIDPresentNothing = iota
	UEIDPresentMMEUEID
	UEIDPresentUEIDPair
)

type UEIDPair struct {
	MMEUEID MMEUEID
	ENBUEID ENBUEID
}

// UEContextReleaseRequest
type UEContextReleaseRequest struct {
	ProtocolIEs UEContextReleaseRequestIEs
}

type UEContextReleaseRequestIEs struct {
	MMEUEID     *MMEUEID
	ENBUEID     *ENBUEID
	Cause       *Cause
	GWContextReleaseIndication *GWContextReleaseIndication
}

type GWContextReleaseIndication int64

const (
	GWContextReleaseIndicationTrue GWContextReleaseIndication = 0
)

// UEContextReleaseComplete
type UEContextReleaseComplete struct {
	ProtocolIEs UEContextReleaseCompleteIEs
}

type UEContextReleaseCompleteIEs struct {
	MMEUEID                    *MMEUEID
	ENBUEID                    *ENBUEID
	CriticalityDiagnostics     *CriticalityDiagnostics
	UserLocationInformation    *UserLocationInformation
}

type UserLocationInformation struct {
	EUTRANGCI EUTRANGCI
	TAI       TAI
}

// E-RABSetupRequest
type ERABSetupRequest struct {
	ProtocolIEs ERABSetupRequestIEs
}

type ERABSetupRequestIEs struct {
	MMEUEID                       *MMEUEID
	ENBUEID                       *ENBUEID
	UEAggregateMaximumBitrate     *UEAggregateMaximumBitrate
	ERABToBeSetupListBearerSUReq  *ERABToBeSetupListBearerSUReq
}

type ERABToBeSetupListBearerSUReq struct {
	List []ERABToBeSetupItemBearerSUReq
}

type ERABToBeSetupItemBearerSUReq struct {
	ERABID                     ERABID
	ERABLevelQoSParameters     ERABLevelQoSParameters
	TransportLayerAddress      TransportLayerAddress
	GTPTEID                    GTPTEID
	NASPDU                     NASPDU
	CorrelationID              *CorrelationID
}

// E-RABSetupResponse
type ERABSetupResponse struct {
	ProtocolIEs ERABSetupResponseIEs
}

type ERABSetupResponseIEs struct {
	MMEUEID                    *MMEUEID
	ENBUEID                    *ENBUEID
	ERABSetupListBearerSURes   *ERABSetupListBearerSURes
	ERABFailedToSetupListBearerSURes *ERABList
	CriticalityDiagnostics     *CriticalityDiagnostics
}

type ERABSetupListBearerSURes struct {
	List []ERABSetupItemBearerSURes
}

type ERABSetupItemBearerSURes struct {
	ERABID                ERABID
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
}

// S1SetupRequest
type S1SetupRequest struct {
	ProtocolIEs S1SetupRequestIEs
}

type S1SetupRequestIEs struct {
	GlobalENBID         *GlobalENBID
	ENBName             *ENBName
	SupportedTAs        *SupportedTAs
	DefaultPagingDRX    *PagingDRX
	CSGIdList           *CSGIdList
}

type ENBName struct {
	Value string
}

type CSGIdList struct {
	List []CSGId
}

// S1SetupResponse
type S1SetupResponse struct {
	ProtocolIEs S1SetupResponseIEs
}

type S1SetupResponseIEs struct {
	MMEName             *MMEName
	ServedGUMMEIs       *ServedGUMMEIs
	RelativeMMECapacity *RelativeMMECapacity
	MMERelaySupportIndicator *MMERelaySupportIndicator
	CriticalityDiagnostics *CriticalityDiagnostics
}

type MMEName struct {
	Value string
}

type RelativeMMECapacity struct {
	Value int64 `aper:"valueLB:0,valueUB:255"`
}

type MMERelaySupportIndicator int64

const (
	MMERelaySupportIndicatorTrue MMERelaySupportIndicator = 0
)

// S1SetupFailure
type S1SetupFailure struct {
	ProtocolIEs S1SetupFailureIEs
}

type S1SetupFailureIEs struct {
	Cause                  *Cause
	TimeToWait             *TimeToWait
	CriticalityDiagnostics *CriticalityDiagnostics
}

type TimeToWait int64

const (
	TimeToWait1s   TimeToWait = 0
	TimeToWait2s   TimeToWait = 1
	TimeToWait5s   TimeToWait = 2
	TimeToWait10s  TimeToWait = 3
	TimeToWait20s  TimeToWait = 4
	TimeToWait60s  TimeToWait = 5
)

// HandoverRequired
type HandoverRequired struct {
	ProtocolIEs HandoverRequiredIEs
}

type HandoverRequiredIEs struct {
	MMEUEID                *MMEUEID
	ENBUEID                *ENBUEID
	HandoverType           *HandoverType
	Cause                  *Cause
	TargetID               *TargetID
	DirectForwardingPathAvailability *DirectForwardingPathAvailability
	SRVCCHOIndication      *SRVCCHOIndication
	SourceToTargetTransparentContainer *SourceToTargetTransparentContainer
	SourceToTargetTransparentContainerSecondary *SourceToTargetTransparentContainer
	MSClassmark2           *MSClassmark2
	MSClassmark3           *MSClassmark3
	CSGId                  *CSGId
	CellAccessMode         *CellAccessMode
	PSServiceNotAvailable  *PSServiceNotAvailable
}

type HandoverType int64

const (
	HandoverTypeIntraLTE        HandoverType = 0
	HandoverTypeLTEtoUTRAN      HandoverType = 1
	HandoverTypeLTEtoGERAN      HandoverType = 2
	HandoverTypeUTRANtoLTE      HandoverType = 3
	HandoverTypeGERANtoLTE      HandoverType = 4
)

type TargetID struct {
	Present     int
	TargeteNBID *TargeteNBID
	TargetRNCID *TargetRNCID
	CGI         *CGI
}

const (
	TargetIDPresentNothing = iota
	TargetIDPresentTargeteNBID
	TargetIDPresentTargetRNCID
	TargetIDPresentCGI
)

type TargeteNBID struct {
	GlobalENBID GlobalENBID
	TAI         TAI
}

type TargetRNCID struct {
	LAI    LAI
	RNCID  int64
	ExtendedRNCID *int64
}

type CGI struct {
	PLMNIdentity PLMNIdentity
	LAC          LAC
	CI           CI
	RAC          *RAC
}

type CI struct {
	Value []byte `aper:"sizeLB:2,sizeUB:2"`
}

type RAC struct {
	Value []byte `aper:"sizeLB:1,sizeUB:1"`
}

type DirectForwardingPathAvailability int64

const (
	DirectForwardingPathAvailabilityAvailable DirectForwardingPathAvailability = 0
)

type SRVCCHOIndication int64

const (
	SRVCCHOIndicationPSandCS SRVCCHOIndication = 0
	SRVCCHOIndicationCSonly  SRVCCHOIndication = 1
)

type SourceToTargetTransparentContainer struct {
	Value []byte
}

type MSClassmark2 struct {
	Value []byte
}

type MSClassmark3 struct {
	Value []byte
}

type PSServiceNotAvailable int64

const (
	PSServiceNotAvailablePSServiceNotAvailable PSServiceNotAvailable = 0
)

// HandoverCommand
type HandoverCommand struct {
	ProtocolIEs HandoverCommandIEs
}

type HandoverCommandIEs struct {
	MMEUEID                    *MMEUEID
	ENBUEID                    *ENBUEID
	HandoverType               *HandoverType
	NASSecurityParamsFromEUTRAN *NASSecurityParamsFromEUTRAN
	ERABSubjectToForwardingList *ERABSubjectToForwardingList
	ERABToReleaseListHOCmd     *ERABList
	TargetToSourceTransparentContainer *TargetToSourceTransparentContainer
	TargetToSourceTransparentContainerSecondary *TargetToSourceTransparentContainer
	CriticalityDiagnostics     *CriticalityDiagnostics
}

type NASSecurityParamsFromEUTRAN struct {
	Value []byte
}

type ERABSubjectToForwardingList struct {
	List []ERABDataForwardingItem
}

type ERABDataForwardingItem struct {
	ERABID                ERABID
	DLTransportLayerAddr  *TransportLayerAddress
	DLGTPTEID             *GTPTEID
	ULTransportLayerAddr  *TransportLayerAddress
	ULGTPTEID             *GTPTEID
}

type TargetToSourceTransparentContainer struct {
	Value []byte
}

// HandoverRequest
type HandoverRequest struct {
	ProtocolIEs HandoverRequestIEs
}

type HandoverRequestIEs struct {
	MMEUEID                    *MMEUEID
	HandoverType               *HandoverType
	Cause                      *Cause
	UEAggregateMaximumBitrate  *UEAggregateMaximumBitrate
	ERABToBeSetupListHOReq     *ERABToBeSetupListHOReq
	SourceToTargetTransparentContainer *SourceToTargetTransparentContainer
	UESecurityCapabilities     *UESecurityCapabilities
	HandoverRestrictionList    *HandoverRestrictionList
	TraceActivation            *TraceActivation
	RequestType                *RequestType
	SRVCCOperationPossible     *SRVCCOperationPossible
	SecurityContext            *SecurityContext
	NASSecurityParamsToEUTRAN  *NASSecurityParamsToEUTRAN
	CSGId                      *CSGId
	CSGMembershipStatus        *CSGMembershipStatus
	GUMMEI                     *GUMMEI
	MMEUEID2                   *MMEUEID
}

type ERABToBeSetupListHOReq struct {
	List []ERABToBeSetupItemHOReq
}

type ERABToBeSetupItemHOReq struct {
	ERABID                     ERABID
	TransportLayerAddress      TransportLayerAddress
	GTPTEID                    GTPTEID
	ERABLevelQoSParameters     ERABLevelQoSParameters
}

type RequestType struct {
	EventType        EventType
	ReportArea       ReportArea
}

type EventType int64

const (
	EventTypeDirect                EventType = 0
	EventTypeChangeOfServeCell     EventType = 1
	EventTypeStopChangeOfServeCell EventType = 2
)

type ReportArea int64

const (
	ReportAreaECGI ReportArea = 0
)

type SecurityContext struct {
	NextHopChainingCount int64
	NextHopParameter     SecurityKey
}

type NASSecurityParamsToEUTRAN struct {
	Value []byte
}

// HandoverRequestAcknowledge
type HandoverRequestAcknowledge struct {
	ProtocolIEs HandoverRequestAcknowledgeIEs
}

type HandoverRequestAcknowledgeIEs struct {
	MMEUEID                            *MMEUEID
	ENBUEID                            *ENBUEID
	ERABAdmittedList                   *ERABAdmittedList
	ERABFailedToSetupListHOReqAck      *ERABList
	TargetToSourceTransparentContainer *TargetToSourceTransparentContainer
	CSGId                              *CSGId
	CriticalityDiagnostics             *CriticalityDiagnostics
}

type ERABAdmittedList struct {
	List []ERABAdmittedItem
}

type ERABAdmittedItem struct {
	ERABID                ERABID
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
	DLTransportLayerAddr  *TransportLayerAddress
	DLGTPTEID             *GTPTEID
	ULTransportLayerAddr  *TransportLayerAddress
	ULGTPTEID             *GTPTEID
}

// HandoverPreparationFailure
type HandoverPreparationFailure struct {
	ProtocolIEs HandoverPreparationFailureIEs
}

type HandoverPreparationFailureIEs struct {
	MMEUEID                *MMEUEID
	ENBUEID                *ENBUEID
	Cause                  *Cause
	CriticalityDiagnostics *CriticalityDiagnostics
}

// HandoverFailure
type HandoverFailure struct {
	ProtocolIEs HandoverFailureIEs
}

type HandoverFailureIEs struct {
	MMEUEID                *MMEUEID
	Cause                  *Cause
	CriticalityDiagnostics *CriticalityDiagnostics
}

// PathSwitchRequest
type PathSwitchRequest struct {
	ProtocolIEs PathSwitchRequestIEs
}

type PathSwitchRequestIEs struct {
	ENBUEID                    *ENBUEID
	ERABToBeSwitchedDLList     *ERABToBeSwitchedDLList
	SourceMMEUEID              *MMEUEID
	EUTRANGCI                  *EUTRANGCI
	TAI                        *TAI
	UESecurityCapabilities     *UESecurityCapabilities
	CSGId                      *CSGId
	CellAccessMode             *CellAccessMode
	SourceMMEGUMMEI            *GUMMEI
	CSGMembershipStatus        *CSGMembershipStatus
}

type ERABToBeSwitchedDLList struct {
	List []ERABToBeSwitchedDLItem
}

type ERABToBeSwitchedDLItem struct {
	ERABID                ERABID
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
}

// PathSwitchRequestAcknowledge
type PathSwitchRequestAcknowledge struct {
	ProtocolIEs PathSwitchRequestAcknowledgeIEs
}

type PathSwitchRequestAcknowledgeIEs struct {
	MMEUEID                    *MMEUEID
	ENBUEID                    *ENBUEID
	UEAggregateMaximumBitrate  *UEAggregateMaximumBitrate
	ERABToBeSwitchedULList     *ERABToBeSwitchedULList
	ERABToBeReleasedList       *ERABList
	SecurityContext            *SecurityContext
	CriticalityDiagnostics     *CriticalityDiagnostics
	MMEUEID2                   *MMEUEID
}

type ERABToBeSwitchedULList struct {
	List []ERABToBeSwitchedULItem
}

type ERABToBeSwitchedULItem struct {
	ERABID                ERABID
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
}

// PathSwitchRequestFailure
type PathSwitchRequestFailure struct {
	ProtocolIEs PathSwitchRequestFailureIEs
}

type PathSwitchRequestFailureIEs struct {
	MMEUEID                *MMEUEID
	ENBUEID                *ENBUEID
	Cause                  *Cause
	CriticalityDiagnostics *CriticalityDiagnostics
}

// Reset
type Reset struct {
	ProtocolIEs ResetIEs
}

type ResetIEs struct {
	Cause     *Cause
	ResetType *ResetType
}

// ResetAcknowledge
type ResetAcknowledge struct {
	ProtocolIEs ResetAcknowledgeIEs
}

type ResetAcknowledgeIEs struct {
	UEAssocLogicalS1ConnectionList *UEAssocLogicalS1ConnectionList
	CriticalityDiagnostics         *CriticalityDiagnostics
}

// Paging
type Paging struct {
	ProtocolIEs PagingIEs
}

type PagingIEs struct {
	UEIdentityIndexValue *UEIdentityIndexValue
	UEPagingID           *UEPagingID
	PagingDRX            *PagingDRX
	CNDomain             *CNDomain
	TAIList              *TAIList
	CSGIdList            *CSGIdList
	PagingPriority       *PagingPriority
}

type UEIdentityIndexValue struct {
	Value []byte `aper:"sizeLB:10,sizeUB:10"`
}

type TAIList struct {
	List []TAIItem
}

type TAIItem struct {
	TAI TAI
}

type PagingPriority int64

const (
	PagingPriorityLevel1 PagingPriority = 0
	PagingPriorityLevel2 PagingPriority = 1
	PagingPriorityLevel3 PagingPriority = 2
	PagingPriorityLevel4 PagingPriority = 3
	PagingPriorityLevel5 PagingPriority = 4
	PagingPriorityLevel6 PagingPriority = 5
	PagingPriorityLevel7 PagingPriority = 6
	PagingPriorityLevel8 PagingPriority = 7
)
