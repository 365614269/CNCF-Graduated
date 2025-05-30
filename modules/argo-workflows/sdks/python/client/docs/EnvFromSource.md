# EnvFromSource

EnvFromSource represents the source of a set of ConfigMaps or Secrets

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**config_map_ref** | [**ConfigMapEnvSource**](ConfigMapEnvSource.md) |  | [optional] 
**prefix** | **str** | Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER. | [optional] 
**secret_ref** | [**SecretEnvSource**](SecretEnvSource.md) |  | [optional] 
**any string name** | **bool, date, datetime, dict, float, int, list, str, none_type** | any string name can be used but the value must be the correct type | [optional]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


