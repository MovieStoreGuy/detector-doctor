# Detector Doctor
_Ensure your detectors never skip a beat_  



Help ensure the health of the detector and if there is any required actions that need to come of it.
This should helpfully resolve the issue of understand where the issue lays and what needs to be done 
as part of it.

## Development
 
The follow code has been developed against Go version 1.14 and to ensure a transparent build environment docker is also employed.

## Usage

```bash
$ detdoc [--flags] --token ${SFX_API_TOKEN} detectorID, ...
```  
_To get the most up to date flags, run the application locally_


This will query the SignalFx API to read the detectorIDs provided to it establish if there is any obvious errors with the detector,
the result of this will produce a report for each detector request.

## Intent / How it works

The application will query the detector in order to establish what is the current issue with the detector:

| Testing                     | Issue  | Implemented?                                                |
|-----------------------------|--------|-------------------------------------------------------------|
| OverMTS Limt                | User   | Yes, Queries the result of the detector                     |
| Has Notificaiton            | User   | Yes, checks to see if any notification rules are set        |
| Has Program Text            | User   | Yes, checks to see if there is any program text defined     |
| Is locked                   | System | Yes, though not sure what can cause a detector to be locked |
| Valid SignalFlow            | User   | No, need an easy to validate SignalFlow                     |
| Handling Sparse Metrics     | User   | No, need to implement handling SSE on signalflow endpoint   |
| Disabled rules              | User   | Yes, queries the rules set with a detector                  |
| Depends on properties       | User   | No, currently looking at way to handle this check           | 
| Metrics are delayed         | ??     | No, need to be able to query the underlying MTS             |
| Missed alert                | System | No, need to be able to take the signalflow and compare when it should of fired when looking at last modified time |
| Missed fired detector       | System | No, need to figure out how best to do this                  |
| Drop in data points         | ??     | No, relies on querying MTS which isn't implemented          |

