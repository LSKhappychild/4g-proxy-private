# S1AP Messages Reference

This document describes all S1AP messages that the proxy can interpret, with special focus on the S1 Handover procedure.

## Supported S1AP Procedures

The proxy recognizes all 51 S1AP procedure codes defined in 3GPP TS 36.413:

### UE-Associated Procedures

| Code | Procedure | Direction | Signal Type |
|------|-----------|-----------|-------------|
| 12 | InitialUEMessage | eNB→MME | (from NAS) |
| 13 | UplinkNASTransport | eNB→MME | (from NAS) |
| 11 | DownlinkNASTransport | MME→eNB | (from NAS) |
| 9 | InitialContextSetup | MME→eNB | (from NAS) |
| 18 | UEContextReleaseRequest | eNB→MME | UEContextRelease |
| 23 | UEContextRelease | MME→eNB | UEContextRelease |
| 21 | UEContextModification | MME→eNB | - |
| 22 | UECapabilityInfoIndication | eNB→MME | - |

### E-RAB Procedures

| Code | Procedure | Direction | Signal Type |
|------|-----------|-----------|-------------|
| 5 | E-RABSetup | MME→eNB | - |
| 6 | E-RABModify | MME→eNB | - |
| 7 | E-RABRelease | MME→eNB | - |
| 8 | E-RABReleaseIndication | eNB→MME | - |
| 50 | E-RABModificationIndication | eNB→MME | - |

### Handover Procedures

| Code | Procedure | Direction | Signal Type |
|------|-----------|-----------|-------------|
| 0 | HandoverPreparation | eNB→MME | **Handover** |
| 1 | HandoverResourceAllocation | MME→eNB | **Handover** |
| 2 | HandoverNotification | eNB→MME | **Handover** |
| 4 | HandoverCancel | eNB→MME | **Handover** |
| 3 | PathSwitchRequest | eNB→MME | **Handover** |
| 24 | eNBStatusTransfer | eNB→MME | **Handover** |
| 25 | MMEStatusTransfer | MME→eNB | **Handover** |

### Interface Management

| Code | Procedure | Direction | Signal Type |
|------|-----------|-----------|-------------|
| 17 | S1Setup | eNB→MME | - |
| 14 | Reset | Both | Reset |
| 29 | ENBConfigurationUpdate | eNB→MME | - |
| 30 | MMEConfigurationUpdate | MME→eNB | - |
| 34 | OverloadStart | MME→eNB | - |
| 35 | OverloadStop | MME→eNB | - |
| 15 | ErrorIndication | Both | - |

### Paging

| Code | Procedure | Direction | Signal Type |
|------|-----------|-----------|-------------|
| 10 | Paging | MME→eNB | **Paging** |

---

## S1 Handover Procedure (Detailed)

### Intra-LTE S1-based Handover (X2 not available)

This is the full S1 handover when X2 interface is not available between source and target eNBs.

```
    Source eNB              MME                Target eNB
        |                    |                      |
        |   (1) HANDOVER REQUIRED                   |
        |------------------->|                      |
        |                    |                      |
        |                    |   (2) HANDOVER REQUEST
        |                    |--------------------->|
        |                    |                      |
        |                    |   (3) HANDOVER REQUEST ACK
        |                    |<---------------------|
        |                    |                      |
        |   (4) HANDOVER COMMAND                    |
        |<-------------------|                      |
        |                    |                      |
        |   (5) eNB STATUS TRANSFER                 |
        |------------------->|                      |
        |                    |                      |
        |                    |   (6) MME STATUS TRANSFER
        |                    |--------------------->|
        |                    |                      |
        |                    |   (7) HANDOVER NOTIFY
        |                    |<---------------------|
        |                    |                      |
        |   (8) UE CONTEXT RELEASE COMMAND          |
        |<-------------------|                      |
        |                    |                      |
        |   (9) UE CONTEXT RELEASE COMPLETE         |
        |------------------->|                      |
        |                    |                      |
```

### Messages in S1 Handover

#### 1. HANDOVER REQUIRED (Source eNB → MME)
- **Procedure Code**: 0 (HandoverPreparation)
- **Message Type**: Initiating
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID
  - Handover Type (intraLTE, LTEtoUTRAN, etc.)
  - Cause
  - Target ID (Target eNB ID + TAI)
  - Source-to-Target Transparent Container (RRC container)

#### 2. HANDOVER REQUEST (MME → Target eNB)
- **Procedure Code**: 1 (HandoverResourceAllocation)
- **Message Type**: Initiating
- **Key IEs**:
  - MME-UE-S1AP-ID
  - Handover Type
  - Cause
  - UE Aggregate Maximum Bit Rate
  - E-RAB to Be Setup List
  - Source-to-Target Transparent Container
  - UE Security Capabilities
  - Security Context

#### 3. HANDOVER REQUEST ACKNOWLEDGE (Target eNB → MME)
- **Procedure Code**: 1 (HandoverResourceAllocation)
- **Message Type**: Successful
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID (new, from target)
  - E-RAB Admitted List
  - E-RAB Failed to Setup List
  - Target-to-Source Transparent Container

#### 4. HANDOVER COMMAND (MME → Source eNB)
- **Procedure Code**: 0 (HandoverPreparation)
- **Message Type**: Successful
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID
  - Handover Type
  - Target-to-Source Transparent Container
  - E-RAB Subject to Forwarding List

#### 5. eNB STATUS TRANSFER (Source eNB → MME)
- **Procedure Code**: 24
- **Message Type**: Initiating
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID
  - eNB Status Transfer Transparent Container (PDCP SN status)

#### 6. MME STATUS TRANSFER (MME → Target eNB)
- **Procedure Code**: 25
- **Message Type**: Initiating
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID
  - eNB Status Transfer Transparent Container

#### 7. HANDOVER NOTIFY (Target eNB → MME)
- **Procedure Code**: 2 (HandoverNotification)
- **Message Type**: Initiating
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID
  - E-UTRAN CGI (new cell)
  - TAI (new tracking area)

#### 8. UE CONTEXT RELEASE COMMAND (MME → Source eNB)
- **Procedure Code**: 23
- **Message Type**: Initiating
- **Key IEs**:
  - UE S1AP IDs (pair)
  - Cause

#### 9. UE CONTEXT RELEASE COMPLETE (Source eNB → MME)
- **Procedure Code**: 23
- **Message Type**: Successful
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID

---

### X2 Handover with Path Switch

When X2 handover is used, S1AP is only involved for path switch:

```
    Source eNB              MME                Target eNB
        |                    |                      |
        |   (X2 Handover between eNBs)              |
        |<==========================================|
        |                    |                      |
        |                    |   (1) PATH SWITCH REQUEST
        |                    |<---------------------|
        |                    |                      |
        |                    |   (2) PATH SWITCH REQUEST ACK
        |                    |--------------------->|
        |                    |                      |
        |   (3) UE CONTEXT RELEASE                  |
        |<-------------------|                      |
```

#### PATH SWITCH REQUEST (Target eNB → MME)
- **Procedure Code**: 3
- **Message Type**: Initiating
- **Key IEs**:
  - eNB-UE-S1AP-ID (new)
  - Source MME-UE-S1AP-ID
  - E-UTRAN CGI
  - TAI
  - UE Security Capabilities
  - E-RAB to Be Switched in Downlink List

#### PATH SWITCH REQUEST ACKNOWLEDGE (MME → Target eNB)
- **Procedure Code**: 3
- **Message Type**: Successful
- **Key IEs**:
  - MME-UE-S1AP-ID
  - eNB-UE-S1AP-ID
  - E-RAB to Be Switched in Uplink List
  - Security Context

---

## Handover Failure Scenarios

### HANDOVER PREPARATION FAILURE (MME → Source eNB)
- **Procedure Code**: 0
- **Message Type**: Unsuccessful
- Sent when target eNB rejects the handover

### HANDOVER FAILURE (MME → Source eNB)
- **Procedure Code**: 1
- **Message Type**: Unsuccessful
- Sent when HandoverRequest fails at target

### HANDOVER CANCEL (Source eNB → MME)
- **Procedure Code**: 4
- **Message Type**: Initiating
- Sent by source eNB to cancel ongoing handover

### PATH SWITCH REQUEST FAILURE (MME → Target eNB)
- **Procedure Code**: 3
- **Message Type**: Unsuccessful
- Sent when path switch fails

---

## Using the Proxy for Handover Testing

### Drop All Handover Messages

```bash
./dropctl drop handover
```

This will drop:
- HandoverRequired
- HandoverRequest
- HandoverRequestAcknowledge
- HandoverCommand
- HandoverNotify
- eNBStatusTransfer
- MMEStatusTransfer
- PathSwitchRequest

### Delay Handover Messages

```bash
# Add 2 second delay to handover
./dropctl delay handover 2000
```

### Monitor Handover Flow

Run the proxy with verbose logging:

```bash
./4g-proxy --mme <MME_IP> --mme-port 36412 -v
```

Example output during S1 handover:
```
S1AP [UL] HandoverPreparation UE: MME-UE-ID=1 eNB-UE-ID=1 [Signal: Handover]
S1AP [DL] HandoverResourceAllocation UE: MME-UE-ID=1 [Signal: Handover]
S1AP [UL] HandoverResourceAllocation UE: MME-UE-ID=1 eNB-UE-ID=2 [Signal: Handover]
S1AP [DL] HandoverPreparation UE: MME-UE-ID=1 eNB-UE-ID=1 [Signal: Handover]
S1AP [UL] eNBStatusTransfer UE: MME-UE-ID=1 eNB-UE-ID=1 [Signal: Handover]
S1AP [DL] MMEStatusTransfer UE: MME-UE-ID=1 eNB-UE-ID=2 [Signal: Handover]
S1AP [UL] HandoverNotification UE: MME-UE-ID=1 eNB-UE-ID=2 [Signal: Handover]
S1AP [DL] UEContextRelease UE: MME-UE-ID=1 eNB-UE-ID=1 [Signal: UEContextRelease]
S1AP [UL] UEContextRelease UE: MME-UE-ID=1 eNB-UE-ID=1 [Signal: UEContextRelease]
```

---

## NAS Messages (EMM/ESM)

The proxy also decodes NAS messages carried in S1AP. See [NAS_MESSAGES.md](NAS_MESSAGES.md) for details.
