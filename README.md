# Building X Operations API Library

> Be sure to review the CHANGELOG.md document to understand the current state of this project.

## What is this For?
In spring of 2022 Siemens Building Technologies released their Operations API. This is a REST-style API that can be used to remotely access various Siemens Building Technologies products. As this API is designed to support a wide range of different products and use cases, the JSON syntax can be relatively complex and difficult to understand. The goal of this project, then, is to abstract the complexities of the Operations API syntax from a developer and provide a friendly, easy-to-use library to enable rapid application development.

[Cloudline Solutions, LLC](https://cloudline-solutions.com) has maintained an intimate relationship with Siemens Building Technologies, and specifically the cloud program, since 2011. We have participated in the development of the core cloud infrastructure (currently known as Building X) and we are now a value-added partner for the Building X ecosystem. Feel free to [contact us](mailto:info@cloudline-solutions.com) if you would like assistance with your project.

## Why Go (aka Golang)?
Cloudline Solutions specializes in serverless cloud architectures on AWS. Within this environment, the Lambda is the central resource for code execution and Go is an excellent language for this use case. Go is compiled, strong-typed and its deployment artifact is self-contained with all dependencies and its runtime. This makes Go fast, secure and predictable within the Lambda (and other) runtime environment. 

## Data Model
The library exposes a data model that is a simplification of the native Building X API data model but is, in general, a one-to-one with it. The following tables capture the models and their properties.

### Session
The session object holds key information need by every call to the Building X API.
| Name  | Type | Description |
| ---   | ---   | --- |
| Partition | String | The partition ID (provided by Siemens) that segregates user data. |
| JWT | String | The authentication token provided after a successful credential exchange with the Building X OAuth provider.
| IsInitialized | Boolean | Indicates whether or not a successful authentication has occurred. |

### Location
The location object represents a physical location where one or more Building X compatible devices are installed.
| Name  | Type | Description |
| ---   | ---   | --- |
| ID | String | The unique identifier for the location
| Name | String | The name of the location
| Description | String | A description for the location
| TimeZone | String | The timezone of the location

### Device
The device object represents either a logical or physical device installed at a location.
| Name  | Type | Description |
| ---   | ---   | --- |
| ID | String | The unique identifier for the device
| Name | String | The name of the device
| Description | String | A description of the device
| Model | String | The model number, if any, of the device
| Serial | String | The serial number, if any, of the device
| OnlineStatus | String | The online status of the device. Possible values are "online", "offline" or "unknown"

### Point
This point object represents a logical or physical point residing on a device.
| Name  | Type | Description |
| ---   | ---   | --- |
| ID | String | The unique identifier for the point
| Name | String | The name of the point
| Description | String | A description of the point
| DataType | String | The data type of the point. Possible values are "boolean", "string" or "number".
| Status | String | Indicates the status of the point. Possible values are "ok" or "fail". 
| StringValue | String | The value of the point, represented as a string.
| Timestamp | Time | The date and time the point was last updated.


### PointHistory

## Required Environment Variables

| Name  | Description |
| ---   | --- |


## Example Usage

## Things to Know
- only point value is settable

## Integration Tests


## Known Issues

- Certain data is missing when retrieving a single Device object from the Building X Operations API. Specifically, the Name, Description and OnlineStatus properties of the Device object returned from GetSingleDevice() will be missing until the underlying API issue is resolved.