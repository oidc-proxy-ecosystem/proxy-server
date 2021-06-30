# プラグインAPI仕様書
<a name="top"></a>

## インデックス
- [API仕様](#API仕様)

  - [response.proto](#response.proto)
      - [Input](#ncs.protobuf.Input)
      - [Input.HeaderEntry](#ncs.protobuf.Input.HeaderEntry)
      - [Output](#ncs.protobuf.Output)
      - [Output.HeaderEntry](#ncs.protobuf.Output.HeaderEntry)
      - [Response](#ncs.protobuf.Response)
  

  - [transport.proto](#transport.proto)
      - [Interface](#ncs.protobuf.Interface)
      - [Interface.HeaderEntry](#ncs.protobuf.Interface.HeaderEntry)
      - [Reply](#ncs.protobuf.Reply)
      - [Reply.HeaderEntry](#ncs.protobuf.Reply.HeaderEntry)
      - [config](#ncs.protobuf.config)
      - [Transport](#ncs.protobuf.Transport)
  

  - [types.proto](#types.proto)
      - [values](#ncs.protobuf.values)
  

- [スカラー値型](#スカラー値型)

## API仕様


<a name="response.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### response.proto



<a name="ncs.protobuf.Input"></a>

#### Input



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| URL | [string](#string) |  |  |
| method | [string](#string) |  |  |
| header | [Input.HeaderEntry](#ncs.protobuf.Input.HeaderEntry) | repeated |  |
| body | [bytes](#bytes) |  |  |






<a name="ncs.protobuf.Input.HeaderEntry"></a>

#### Input.HeaderEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [values](#ncs.protobuf.values) |  |  |






<a name="ncs.protobuf.Output"></a>

#### Output



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [Output.HeaderEntry](#ncs.protobuf.Output.HeaderEntry) | repeated |  |
| body | [bytes](#bytes) |  |  |






<a name="ncs.protobuf.Output.HeaderEntry"></a>

#### Output.HeaderEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [values](#ncs.protobuf.values) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.Response"></a>

#### Response


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Modify | [Input](#ncs.protobuf.Input) | [Output](#ncs.protobuf.Output) |  |

 <!-- end services -->



<a name="transport.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### transport.proto



<a name="ncs.protobuf.Interface"></a>

#### Interface



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| URL | [string](#string) |  |  |
| header | [Interface.HeaderEntry](#ncs.protobuf.Interface.HeaderEntry) | repeated |  |
| config | [config](#ncs.protobuf.config) |  |  |






<a name="ncs.protobuf.Interface.HeaderEntry"></a>

#### Interface.HeaderEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [values](#ncs.protobuf.values) |  |  |






<a name="ncs.protobuf.Reply"></a>

#### Reply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| URL | [string](#string) |  |  |
| header | [Reply.HeaderEntry](#ncs.protobuf.Reply.HeaderEntry) | repeated |  |
| status | [int32](#int32) |  |  |
| errorMessage | [string](#string) |  |  |






<a name="ncs.protobuf.Reply.HeaderEntry"></a>

#### Reply.HeaderEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [values](#ncs.protobuf.values) |  |  |






<a name="ncs.protobuf.config"></a>

#### config



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| directory | [string](#string) |  |  |
| loadbalancer | [string](#string) |  |  |
| config | [string](#string) |  |  |
| oidc | [string](#string) |  |  |
| saml | [string](#string) |  |  |
| auth | [string](#string) |  |  |
| menu | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.Transport"></a>

#### Transport


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Transport | [Interface](#ncs.protobuf.Interface) | [Reply](#ncs.protobuf.Reply) |  |

 <!-- end services -->



<a name="types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### types.proto



<a name="ncs.protobuf.values"></a>

#### values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## スカラー値型

| .proto Type | Notes | Go Type | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | -------- | --------- | ----------- |
| <a name="double" /> double |  | float64 | double | double | float |
| <a name="float" /> float |  | float32 | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | []byte | string | ByteString | str |
