# Simple Fix

[![Generic badge](https://img.shields.io/github/v/release/b2broker/simplefix-go.svg?style=for-the-badge)](https://github.com/b2broker/simplefix-go/releases/latest)
[![Generic badge](https://goreportcard.com/badge/github.com/b2broker/simplefix-go?style=for-the-badge)](https://goreportcard.com/report/github.com/b2broker/simplefix-go)
[![Generic badge](https://img.shields.io/github/stars/b2broker/simplefix-go?style=for-the-badge&logo=GitHub)](https://github.com/b2broker/simplefix-go/stargazers)
[![Generic badge](https://img.shields.io/badge/Go->=1.16-blue.svg?style=for-the-badge&logo=go)](https://golang.org/doc/go1.16)
[![Generic badge](https://img.shields.io/badge/semver-semantic_release-blue.svg?style=for-the-badge&logo=semantic-release)](https://github.com/go-semantic-release/semantic-release)

:warning: This is a beta-version of SimpleFix Go. This library is under development and is subject to future changes.

<details>
<summary>Table of contents</summary>

## Table of contents

- [What is FIX?](#what-is-fix)
- [Why SimpleFix Go?](#why-simplefix-go)
- [Installing SimpleFix Go](#installation)
- [Using the Generator](#generator)
- [Getting started with SimpleFix Go](#getting-started)
- [Customizing messages](#customizing-messages)

</details>

## What is FIX?

FIX (a shorthand for Financial Information eXchange) is a widely adopted electronic communications protocol used for real-time exchange of information related to trading and markets. FIX connects the global ecosystem of venues, asset managers, banks/brokers, vendors and regulators by standardizing the communication among participants. FIX is a public-domain specification owned and maintained by [FIX Protocol, Ltd (FPL)](https://www.fixtrading.org/).

## Why SimpleFix Go?

SimpleFix Go is an open-source library that allows you to quickly integrate FIX messaging into your environment. The library is entirely written in Go, making it perfect for solutions powered by this programming language. SimpleFix Go supports any FIX API version. To learn about specifics of various FIX protocol versions, refer to the [OnixS website](https://www.onixs.biz/fix-dictionary.html).

You can provide your own extensions to SimpleFix Go and create a custom FIX dialect. This guide explains the basics of SimpleFix Go installation and provides examples illustrating how to configure and customize the library according to your requirements.

### Main features

- [x] Adding custom fields to the FIX protocol
- [x] Adding custom messages to the FIX protocol
- [ ] Adding custom types to the FIX protocol
- [x] Built-in session pipelines features
- [x] Built-in Acceptor (for the server side)
- [x] Built-in Initiator (for the client side)
- [ ] Validation of incoming messages
- [x] Validation of outgoing messages
- [x] A [demo server](https://docs.marksman.b2broker.com/en/fix-api.html#demo-mode) complete with mock data
- [x] Anything missing? Let us know!

## Installing SimpleFix Go

To install SimpleFix Go, download the library by executing the following command:

```sh
$ go get -u github.com/b2broker/simplefix-go
```

2. Install the *Generator* if you want to use your own XML schema providing a custom set of FIX messaging options:

```sh
$ cd $GOPATH/src/github.com/b2broker/simplefix-go && go install ./...
```

## Using the Generator

The *Generator* is used to define the structure of FIX messages, as well as specify their tags and define message type constants and methods required to support any FIX API version.

Examples of code produced by the *Generator* can be found in the [./tests](https://github.com/b2broker/simplefix-go/tree/master/tests/fix44) directory containing an automatically generated Go library based on a stripped-down FIX version 4.4. The library code is generated according to a scheme located in the [./source](https://github.com/b2broker/simplefix-go/tree/master/source) directory.

### Generating a basic FIX library

The following code generates a FIX library based on an [XML schema](https://github.com/b2broker/simplefix-go/blob/master/source/fix44.xml) defining the library structure:

```sh
fixgen -o=./fix44 -s=./source/fix44.xml -t=./source/types.xml
```

After executing this command, the generated library code will be located in the [./fix44](https://github.com/b2broker/simplefix-go/tree/master/tests/fix44) directory.

### Specifying Generator parameters

To create a custom FIX messaging library, prepare two XML files and specify the following parameters for the `fixgen` command:

`-o` — the output directory

`-s` — the path to the main XML schema

`-t` — the path to an XML file specifying value type mapping and informing the *Generator* about proper type casting (although the original FIX protocol features a lot of different value types, Go uses a smaller set of types that should be mapped to the FIX API)

Sample XML files are located in the [./source](https://github.com/b2broker/simplefix-go/blob/master/source/) directory. You can use the existing files or modify them as required.

## Getting started with SimpleFix Go

In this section, you will learn how to specify the session options and start a new FIX session as a client or as a server.

### Specifying session options

The following sample code illustrates how to use a message builder to create various standard messages, as well as define fields and message tags required for FIX session pipelines. The `fixgen` command will generate the required structure in almost no time.

```
// In this code, fixgen is your generated FIX package:

var sessionOpts = session.Opts{
	MessageBuilders: session.MessageBuilders{
		HeaderBuilder:        fixgen.Header{}.New(),
		TrailerBuilder:       fixgen.Trailer{}.New(),
		LogonBuilder:         fixgen.Logon{}.New(),
		LogoutBuilder:        fixgen.Logout{}.New(),
		RejectBuilder:        fixgen.Reject{}.New(),
		HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
		TestRequestBuilder:   fixgen.TestRequest{}.New(),
		ResendRequestBuilder: fixgen.ResendRequest{}.New(),
	},
	Tags: &messages.Tags{
		MsgType:         mustConvToInt(fixgen.FieldMsgType),
		MsgSeqNum:       mustConvToInt(fixgen.FieldMsgSeqNum),
		HeartBtInt:      mustConvToInt(fixgen.FieldHeartBtInt),
		EncryptedMethod: mustConvToInt(fixgen.FieldEncryptMethod),
	},
	AllowedEncryptedMethods: map[string]struct{}{
		fixgen.EnumEncryptMethodNoneother: {},
	},
	SessionErrorCodes: &messages.SessionErrorCodes{
		InvalidTagNumber:            mustConvToInt(fixgen.EnumSessionRejectReasonInvalidtagnumber),
		RequiredTagMissing:          mustConvToInt(fixgen.EnumSessionRejectReasonRequiredtagmissing),
		TagNotDefinedForMessageType: mustConvToInt(fixgen.EnumSessionRejectReasonTagNotDefinedForThisMessageType),
		UndefinedTag:                mustConvToInt(fixgen.EnumSessionRejectReasonUndefinedtag),
		TagSpecialWithoutValue:      mustConvToInt(fixgen.EnumSessionRejectReasonTagspecifiedwithoutavalue),
		IncorrectValue:              mustConvToInt(fixgen.EnumSessionRejectReasonValueisincorrectoutofrangeforthistag),
		IncorrectDataFormatValue:    mustConvToInt(fixgen.EnumSessionRejectReasonIncorrectdataformatforvalue),
		DecryptionProblem:           mustConvToInt(fixgen.EnumSessionRejectReasonDecryptionproblem),
		SignatureProblem:            mustConvToInt(fixgen.EnumSessionRejectReasonSignatureproblem),
		CompIDProblem:               mustConvToInt(fixgen.EnumSessionRejectReasonCompidproblem),
		Other:                       mustConvToInt(fixgen.EnumSessionRejectReasonOther),
	},
}
```

### Starting as a client

The *Initiator* is a FIX API client that connects to an existing server.

The default *Initiator* implementation can be found in the [./initiator/main.go](https://github.com/b2broker/simplefix-go/blob/master/examples/initiator/main.go) file.

### Starting as a server

The *Acceptor* is a listener that accepts and handles client connection requests. According to the FIX protocol, the *Acceptor* can be both a provider and receiver of data, meaning that it can send requests to the clients as well as read data streams received from them.

The default *Acceptor* implementation can be found in the [./acceptor/main.go](https://github.com/b2broker/simplefix-go/blob/master/examples/acceptor/main.go) file.


## Customizing messages

### Adding custom message fields

The SimpleFix library features a default session package which provides all the necessary functionality for the standard FIX API pipelines, such as authentication, logging out, heartbeats, message rejects and handling of typical errors. For this reason, if you want to customize the default messages, such as `Logon <A>` or `Heartbeat <0>`, you should configure the `Session` structure in one of the following ways:

- by integrating the existing structure into your custom procedure by way of composition, or simply copying the existing structure and modifying it in your client code

- by modifying the message builder to account for the messages that you want to customize


#### Customizing the message builder

The standard `Session` structure accepts each `MessageBuilder` instance as an auto-generated `Opts` field. Each message builder is an interface, which means that you can create a custom builder, and the library will use it as the default one.

For example, if you want to add a new `CounterPartyID` field (tag number 22000) to you `Logon` message, you should modify your XML schema by adding a new field to the `fields` section and to your custom `Logon` message:

```xml

<fix>
    . . .
    <messages>
        <message name='Logon' msgcat='admin' msgtype='A'>
            . . .
            <field name='CounterPartyID' required='Y'/>
            ...
        </message>
    </messages>
    . . .
    <fields>
        <field number="22000" name="CounterPartyID" type="STRING"/>
    </fields>
    . . .
</fix>
```

While the rest of the code is generated by `fixgen`, you should specify this field manually using a custom builder:

```
// Your FIX package is generated by fixgen:

type CustomLogon struct {
    *fixgen.Logon
}

func (cl *CustomLogon) New() messages.LogonBuilder {
    l := cl.New()
    l.SetFieldCounterParty(os.Getenv("COUNTER_PARTY_ID"))
    return l
}
```

After this, you can use your `CustomLogon` structure (with a new field added to it) as a `LogonBuilder` in the default FIX API pipelines.
