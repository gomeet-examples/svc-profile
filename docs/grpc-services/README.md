# Protocol Documentation
<a name="top"/>

## Table of Contents

- [profile.proto](#profile.proto)
    - [EmptyMessage](#grpc.gomeetexamples.profile.EmptyMessage)
    - [ProfileCreationRequest](#grpc.gomeetexamples.profile.ProfileCreationRequest)
    - [ProfileInfo](#grpc.gomeetexamples.profile.ProfileInfo)
    - [ProfileList](#grpc.gomeetexamples.profile.ProfileList)
    - [ProfileListRequest](#grpc.gomeetexamples.profile.ProfileListRequest)
    - [ProfileRequest](#grpc.gomeetexamples.profile.ProfileRequest)
    - [ProfileResponse](#grpc.gomeetexamples.profile.ProfileResponse)
    - [ProfileResponseLight](#grpc.gomeetexamples.profile.ProfileResponseLight)
    - [ServiceStatus](#grpc.gomeetexamples.profile.ServiceStatus)
    - [ServicesStatusList](#grpc.gomeetexamples.profile.ServicesStatusList)
    - [VersionResponse](#grpc.gomeetexamples.profile.VersionResponse)
  
    - [Genders](#grpc.gomeetexamples.profile.Genders)
    - [ServiceStatus.Status](#grpc.gomeetexamples.profile.ServiceStatus.Status)
  
  
    - [Profile](#grpc.gomeetexamples.profile.Profile)
  

- [Scalar Value Types](#scalar-value-types)



<a name="profile.proto"/>
<p align="right"><a href="#top">Top</a></p>

## profile.proto



<a name="grpc.gomeetexamples.profile.EmptyMessage"/>

### EmptyMessage







<a name="grpc.gomeetexamples.profile.ProfileCreationRequest"/>

### ProfileCreationRequest
ProfileCreationRequest encodes a profile creation request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gender | [Genders](#grpc.gomeetexamples.profile.Genders) |  | profile role |
| email | [string](#string) |  | profile email |
| name | [string](#string) |  | profile name |
| birthday | [string](#string) |  | profile birthday |






<a name="grpc.gomeetexamples.profile.ProfileInfo"/>

### ProfileInfo
ProfileInfo encodes information about a profile.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | internal profile ID |
| gender | [Genders](#grpc.gomeetexamples.profile.Genders) |  | profile role |
| email | [string](#string) |  | profile email |
| name | [string](#string) |  | profile name |
| birthday | [string](#string) |  | profile birthday |
| created_at | [string](#string) |  | creation time (UTC - RFC 3339 format) |
| updated_at | [string](#string) |  | modification time (UTC - RFC 3339 format) |
| deleted_at | [string](#string) |  | deletion time (UTC - RFC 3339 format if the profile was logically deleted, empty otherwise) |






<a name="grpc.gomeetexamples.profile.ProfileList"/>

### ProfileList
ProfileList encodes the result of a ProfileListRequest.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result_set_size | [uint32](#uint32) |  | total number of results |
| has_more | [bool](#bool) |  | true if there are more results for the ProfileListRequest |
| profiles | [ProfileInfo](#grpc.gomeetexamples.profile.ProfileInfo) | repeated | list of ProfileInfo messages |






<a name="grpc.gomeetexamples.profile.ProfileListRequest"/>

### ProfileListRequest
ProfileListRequest encodes a set of criteria for the retrieval of a list of profiles.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_number | [uint32](#uint32) |  | page number (starting from 1) |
| page_size | [uint32](#uint32) |  | number of results in a page |
| order | [string](#string) |  | result ordering specification (default &#34;created_at asc&#34;) |
| exclude_soft_deleted | [bool](#bool) |  | if true, excludes logically-deleted profiles from the result set |
| soft_deleted_only | [bool](#bool) |  | if true, restricts the result set to logically-deleted profiles |
| gender | [Genders](#grpc.gomeetexamples.profile.Genders) |  | role to search for |






<a name="grpc.gomeetexamples.profile.ProfileRequest"/>

### ProfileRequest
ProfileRequest encodes a profile identifier.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | profile ID |






<a name="grpc.gomeetexamples.profile.ProfileResponse"/>

### ProfileResponse
ProfileResponse encodes the result of a profile operation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | indicates whether the operation (authentication, creation, update or delete) was successful |
| info | [ProfileInfo](#grpc.gomeetexamples.profile.ProfileInfo) |  | profile information (unreliable if the operation failed) |






<a name="grpc.gomeetexamples.profile.ProfileResponseLight"/>

### ProfileResponseLight
ProfileResponseLight encodes the result of a profile operation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | indicates whether the operation was successful |






<a name="grpc.gomeetexamples.profile.ServiceStatus"/>

### ServiceStatus
SeviceStatus represents a sub services status message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of service |
| version | [string](#string) |  | version of service |
| status | [ServiceStatus.Status](#grpc.gomeetexamples.profile.ServiceStatus.Status) |  | status of service see enum Status |
| e_msg | [string](#string) |  |  |






<a name="grpc.gomeetexamples.profile.ServicesStatusList"/>

### ServicesStatusList
ServicesStatusList is the sub services status list


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| services | [ServiceStatus](#grpc.gomeetexamples.profile.ServiceStatus) | repeated |  |






<a name="grpc.gomeetexamples.profile.VersionResponse"/>

### VersionResponse
VersionMessage represents a version message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Id represents the message identifier. |
| version | [string](#string) |  |  |





 


<a name="grpc.gomeetexamples.profile.Genders"/>

### Genders


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOW | 0 | normaly never |
| MALE | 1 | male gender |
| FEMALE | 2 | female gender |



<a name="grpc.gomeetexamples.profile.ServiceStatus.Status"/>

### ServiceStatus.Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| OK | 0 |  |
| UNAVAILABLE | 1 |  |


 

 


<a name="grpc.gomeetexamples.profile.Profile"/>

### Profile


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Version | [EmptyMessage](#grpc.gomeetexamples.profile.EmptyMessage) | [VersionResponse](#grpc.gomeetexamples.profile.EmptyMessage) | Version method receives no paramaters and returns a version message. |
| ServicesStatus | [EmptyMessage](#grpc.gomeetexamples.profile.EmptyMessage) | [ServicesStatusList](#grpc.gomeetexamples.profile.EmptyMessage) | ServicesStatus method receives no paramaters and returns all services status message |
| Create | [ProfileCreationRequest](#grpc.gomeetexamples.profile.ProfileCreationRequest) | [ProfileResponse](#grpc.gomeetexamples.profile.ProfileCreationRequest) | Create attempts to create a new profile. |
| Read | [ProfileRequest](#grpc.gomeetexamples.profile.ProfileRequest) | [ProfileInfo](#grpc.gomeetexamples.profile.ProfileRequest) | Read returns information about an existing profile. |
| List | [ProfileListRequest](#grpc.gomeetexamples.profile.ProfileListRequest) | [ProfileList](#grpc.gomeetexamples.profile.ProfileListRequest) | List returns a list of profiles matching a set of criteria. |
| Update | [ProfileInfo](#grpc.gomeetexamples.profile.ProfileInfo) | [ProfileResponse](#grpc.gomeetexamples.profile.ProfileInfo) | Update attempts to update an existing profile. |
| SoftDelete | [ProfileRequest](#grpc.gomeetexamples.profile.ProfileRequest) | [ProfileResponse](#grpc.gomeetexamples.profile.ProfileRequest) | SoftDelete attempts to delete an existing profile logically. |
| HardDelete | [ProfileRequest](#grpc.gomeetexamples.profile.ProfileRequest) | [ProfileResponseLight](#grpc.gomeetexamples.profile.ProfileRequest) | HardDelete attempts to delete an existing profile physically. |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

