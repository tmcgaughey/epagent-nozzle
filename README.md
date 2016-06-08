# epagent-nozzle (1.0)


# Description
Creates a Nozzle for the [Cloud Foundry Firehose](https://github.com/cloudfoundry/loggregator) and reports metrics to the CA APM EPAgent REST interface. 

## APM version
This has been tested with APM 10.  An EPAgent 9.7.1 or greater is required for the REST interface.

## Supported third party versions
Tested with Cloud Foundry 215

## Limitations
Handles the [dropsonde-protocol events](https://github.com/cloudfoundry/dropsonde-protocol/tree/master/events) except for Log.  It also keeps the older Heartbeat event.
By nature, the Firehose may generate a lot of metrics.  CA suggestes monitoring the load at the Enterprise Manager.
If the nozzle is unable to keep up with the Firehose, a message is posted to the log and the metric `Firehose|slowConsumerAlert` will be non-zero.  Try scaling up the resources for the nozzle.

## License
[Apache License, Version 2.0, January 2004](http://www.apache.org/licenses/). See [Licensing](https://communities.ca.com/docs/DOC-231150910#license) on the CA APM Developer Community.


# Installation Instructions

## Prerequisites
* A [RESTful CA APM EPAgent](https://wiki.ca.com/display/APMDEVOPS97/EPAgent+Overview) must be installed and the [HTTP port must be set](https://wiki.ca.com/display/APMDEVOPS97/Configure+the+EPAgent+RESTful+Interface).  The host and port running the EPAgent should be reachable from the epagent-monitor.

* You must have a user with access to the Cloud Foundry Firehose configured in your manifest

```
properties:
  uaa:
    clients:
      epagent-nozzle:
        access-token-validity: 1209600
        authorized-grant-types: authorization_code,client_credentials,refresh_token
        override: true
        secret: <password>
        scope: openid,oauth.approvals,doppler.firehose
        authorities: oauth.login,doppler.firehose

```

* Golang installed and configured (see [here](https://golang.org/doc/install) for a tutorial on how to do this).
* godep (see [here](https://github.com/tools/godep) for installation instructions).

## Dependencies
* APM EPAgent 9.7.1+
* Golang 1.4.2 and 1.5 were tested
* Godep

## Installation
Once you've met all the prerequisites, you'll need to download the library and install the dependencies:

```
go get github.com/tmcgaughey/epagent-nozzle
cd $GOPATH/src/github.com/tmcgaughey/epagent-nozzle
godep restore
godep go build
```

## Configuration
Edit the **config/caapm-firehose-nozzle.json** file to designate:

 - The [UAA API URL](https://github.com/cloudfoundry/uaa).  The status URL specified in the prerequisites 
 - Credentials matching those in the manifest
 - A doppler endpoint and subscription ID   
 - The path to the EPAgent's REST interface
 - A frequency to flush metrics to the EPAgent
 

# Usage Instructions
From the installation location, execute the fieldpack with:

```
./epagent-nozzle
```

Output will appear on the console or log, and hopefully the Introscope Investigator as well!


## Metric description

Metrics are published under the path `Firehose`.

Following is a brief overview of the metrics that epagent-nozzle will extract from the Firehose and send off to CA APM.

The EPAgent Data Types referenced below can be found in the [EPAgent documentation](https://wiki.ca.com/display/APMDEVOPS98/Configure+the+EPAgent+RESTful+Interface#ConfiguretheEPAgentRESTfulInterface-Types)  

### CounterEvent

CounterEvents represent the increment of a counter. epagent-nozzle will send these as LongCounter metric. These metrics appear in the Investigator under `Firehose|ops|<counterName>`.

### ContainerMetric

CPU, RAM and disk usage metrics for app containers will be sent through as LongCounter metrics. If generated, these metrics appear in the Investigator under `Firehose|apps|<appID>|<containerMetric>|<instanceIndex>`.

### Heartbeat

Heartbeat Events indicate liveness of the emitter and provide counts of the number of Events processed by the emitter. These metrics are processed as Perinterval metrics. These metrics appear in the Investigator under `Firehose|ops|<Origin>|heartbeats`.

### HTTPStartStop

HTTP requests passing through the Cloud Foundry routers get recorded as HTTPStartStop Events. The following table gives an overview of the HTTP metrics:

| Name | Description | APM Metric Type |
| ---- | ----------- | ------------------ |
| HttpStartStopResponseTime | HTTP response times in milliseconds | IntCounter |
| HttpStartStopStatusCodeCount | A count of each HTTP status code | PerintervalCounter |
| HttpStartStopErrorCount | A count of each HTTP error code | PerintervalCounter |
| HttpStartStopHttpRequestCount | A count of each HTTP request | PerintervalCounter |

HTTP metrics are found under `Firehose|http`.  For all HTTPStartStop Events, the hostname is extracted from the URI and used in the Metric name, but not the path.  This is to prevent metric explosion.

### ValueMetric

Any ValueMetric Event that appears on the Firehose will be sent as a LongCounter metric. This includes metrics such as numCPUS, numGoRoutines, memoryStats, etc. These metrics appear in the Investigator under `Firehose|ops|<Origin>`.



## Custom Management Modules
None provided.

## Custom type viewers
None provided.

## Name Formatter Replacements
None provided.

## Debugging and Troubleshooting
The log output will indicate if the nozzle is unable to connect to the Firehose, or send to the EPAgent.  Metric issues at the EPAgent are also logged.  

## Support
This document and associated tools are made available from CA Technologies as examples and provided at no charge as a courtesy to the CA APM Community at large. This resource may require modification for use in your environment. However, please note that this resource is not supported by CA Technologies, and inclusion in this site should not be construed to be an endorsement or recommendation by CA Technologies. These utilities are not covered by the CA Technologies software license agreement and there is no explicit or implied warranty from CA Technologies. They can be used and distributed freely amongst the CA APM Community, but not sold. As such, they are unsupported software, provided as is without warranty of any kind, express or implied, including but not limited to warranties of merchantability and fitness for a particular purpose. CA Technologies does not warrant that this resource will meet your requirements or that the operation of the resource will be uninterrupted or error free or that any defects will be corrected. The use of this resource implies that you understand and agree to the terms listed herein.

Although these utilities are unsupported, please let us know if you have any problems or questions by adding a comment to the CA APM Community Site area where the resource is located, so that the Author(s) may attempt to address the issue or question.

Unless explicitly stated otherwise this field pack is only supported on the same platforms as the APM core agent. See [APM Compatibility Guide](http://www.ca.com/us/support/ca-support-online/product-content/status/compatibility-matrix/application-performance-management-compatibility-guide.aspx).


# Contributing
The [CA APM Community](https://communities.ca.com/community/ca-apm) is the primary means of interfacing with other users and with the CA APM product team.  The [developer subcommunity](https://communities.ca.com/community/ca-apm/ca-developer-apm) is where you can learn more about building APM-based assets, find code examples, and ask questions of other developers and the CA APM product team.

If you wish to contribute to this or any other project, please refer to [easy instructions](https://communities.ca.com/docs/DOC-231150910) available on the CA APM Developer Community.


# Change log
Changes for each version of the field pack.

Version | Author | Comment
--------|--------|--------
1.0 | Tim McGaughey | First version of the field pack.

## Support URL
https://github.com/tmcgaughey/epagent-nozzle/issues

## Short Description
Monitors Cloud Foundry

## Categories
Cloud
