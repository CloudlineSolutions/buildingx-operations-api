# Building X Operations API Library

> Be sure to review the CHANGELOG.md document to understand the current state of this project.

## What is this For?
In spring of 2022 Siemens Building Technologies released their Operations API. This is a REST-style API that can be used to remotely access various Siemens Building Technologies products. As this API is designed to support a wide range of different products and use cases, the JSON syntax can be relatively complex. The goal of this project is to abstract the complexities of the Operations API JSON syntax from the developer and provide a friendly, easy-to-use library to enable rapid application development.

[Cloudline Solutions, LLC](https://cloudline-solutions.com) has maintained an intimate relationship with Siemens Building Technologies, and specifically the cloud program, since 2011. We have participated in the development of the core cloud infrastructure (currently known as Building X) and we are now a value-added partner for the Building X ecosystem. Feel free to [contact us](mailto:info@cloudline-solutions.com) if you would like assistance with your project.

## Why Go (aka Golang)?
Cloudline Solutions specializes in serverless cloud architectures on AWS. Within this environment, the AWS Lambda is the primary resource for code execution and Go is an excellent language for this use case. Go is compiled, strongly-typed and its deployment artifact is self-contained with all dependencies and its runtime. This makes Go fast, secure and predictable within the Lambda (and other) runtime environment. 

## Data Model
The library exposes a data model that is a simplification of the native Building X API data model. The following tables capture the models and their properties.

### Session
The session object holds key information needed by every call to the Building X API.
| Name  | Type | Description |
| ---   | ---   | --- |
| Partition | String | The partition ID (provided by Siemens) that segregates user data. |
| JWT | String | The authentication token provided after a successful credential exchange with the Building X OAuth provider.
| IsInitialized | Boolean | Indicates whether or not a successful authentication has occurred. |

### Location
The location object represents a physical location where one or more Building X compatible devices are installed.
| Name  | Type | Description |
| ---   | ---   | --- |
| ID | String | The unique identifier for the location |
| Name | String | The name of the location |
| Description | String | A description for the location |
| Street | String | The street address of the location |
| City | String | The city name of the location |
| PostalCode | String | The location postal code |
| Country | String | The country code for the location |
| TimeZone | String | The timezone of the location |

### Device
The device object represents either a logical or physical device installed at a location.
| Name  | Type | Description |
| ---   | ---   | --- |
| ID | String | The unique identifier for the device |
| Name | String | The name of the device
| Description | String | A description of the device |
| Model | String | The model number, if any, of the device |
| Serial | String | The serial number, if any, of the device |
| OnlineStatus | String | The online status of the device. Possible values are "online", "offline" or "unknown" |

### Point
The point object represents a logical or physical point residing on a device.
| Name  | Type | Description |
| ---   | ---   | --- |
| ID | String | The unique identifier for the point |
| Name | String | The name of the point |
| Description | String | A description of the point |
| DataType | String | The data type of the point. Possible values are "boolean", "string" or "number".|
| Writable | Boolean | Indicates whether or not the point can be commanded. |
| Status | String | Indicates the status of the point. Possible values are "ok" or "fail". |
| StringValue | String | The value of the point, represented as a string. |
| Timestamp | Time | The date and time the point was last updated. |


### PointHistory
The point history object represents a single record (typically a COV) in the history of a point.
| Name  | Type | Description |
| ---   | ---   | --- |
| Timestamp | Time | A timestamp for when the record was created |
| Value | String | The value for the record |

## Required Environment Variables
The library requires certain environment variables to be present at runtime. These are listed in the following table.

| Name  | Description |
| ---   | --- |
| BUILDINGX_CLIENT_ID | The client ID credential for the Operations API.  |
| BUILDINGX_CLIENT_SECRET | The client secret credential for the Operations API. |
| BUILDINGX_AUDIENCE | An authentication service property. Normally this has a value of https://horizon.siemens.com |
| BUILDINGX_ENDPOINT | An authentication service property. Normally this has a value of https://api.bpcloud.siemens.com |
| BUILDINGX_AUTH_URL | An authentication service property. Normally this has a value of https://siemens-bt-015.eu.auth0.com/oauth/token |
| BUILDINGX_PARTITION_ID | This environment variable holds a partition ID and is used for running the integration tests but is not required for the use of the library in a project. |


## Example Usage
The following example code demonstrates a complete set of typical operations, ending in setting a point value. It does not include all available functions, but it gives you a good idea for how to use the library.

The sequence assumes that the first location associated with the partition has a gateway with at least one device under it, with the device having at least one writable point (a boolean in this example).

```

    // initialize the session (uses credentials and other properties stored in environment variables)
    // your partition ID is provided when you set up your account in the SBT cloud portal
    partitionID := "{your partition id goes here}"
	session := Session{}
	err := session.Initialize(partitionID)
	if err != nil {
		// handle the error
	}

    // get all locations for the partition
    locations, err := GetLocations(&session)
	if err != nil {
		// handle the error
	}
    
    // Get all devices associated with the first location 
	devices, err := GetDevicesByLocation(&session, &locations[0])
	if err != nil {
		// handle the error
	}

    // find the gateway (X300 or X200).
	gatewayID := ""
	for _, device := range devices {
		if strings.ToLower(device.Model) == "x300" || strings.ToLower(device.Model) == "x200" {
			gatewayID = device.ID
			break
		}
	}

    // get the devices under a gateway
	gatewayDevices, err := GetDevicesByGateway(&session, gatewayID)
	if err != nil {
		// handle the error
	}

    // get the points associated with the first device under the gateway
	points, err := GetPointsByDevice(&session, &gatewayDevices[0])
	if err != nil {
		// handle the error
	}

    // find a writable point
	writablePoint := Point{}
	for _, p := range points {
		if p.Writable {
			writablePoint = p
		}
	}

    // command the writable point - this example assumes a boolean type
    err := CommandPointValue(&session, &writablePoint, "true")
	if err != nil {
		// handle the error
	}


```

## Things to Know

- only point value is settable. All other object properties are read only.

## Integration Tests
The Go test files in this project represent integration, not unit, tests. This means that the tests expect a working Building X account and API credentials. The tests also assume that an X300 (or X200) gateway is installed with at least one device (ex: PXC4) connected to the gateway.


## Known Issues & Limitations

- Certain data is missing when retrieving a single Device object from the Building X Operations API. Specifically, the Name, Description and OnlineStatus properties of the Device object returned from GetSingleDevice() will be missing until the underlying API issue is resolved.