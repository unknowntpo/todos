# {{classname}}

All URIs are relative to *https://todos.unknowntpo.net*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1HealthcheckGet**](DefaultApi.md#V1HealthcheckGet) | **Get** /v1/healthcheck | Returns status of service.
[**V1TasksGet**](DefaultApi.md#V1TasksGet) | **Get** /v1/tasks | Returns all tasks for user identified by userId.
[**V1TasksPost**](DefaultApi.md#V1TasksPost) | **Post** /v1/tasks | Create a new task for user &#x27;user_id&#x27;
[**V1TasksTaskIdDelete**](DefaultApi.md#V1TasksTaskIdDelete) | **Delete** /v1/tasks/{taskId} | Delete task by id for specific user.
[**V1TasksTaskIdGet**](DefaultApi.md#V1TasksTaskIdGet) | **Get** /v1/tasks/{taskId} | Returns all tasks for user identified by userId.
[**V1TasksTaskIdPatch**](DefaultApi.md#V1TasksTaskIdPatch) | **Patch** /v1/tasks/{taskId} | Update task for specific user.
[**V1TokensAuthenticationPost**](DefaultApi.md#V1TokensAuthenticationPost) | **Post** /v1/tokens/authentication | Authenticate the user based on given token.
[**V1UsersActivationPut**](DefaultApi.md#V1UsersActivationPut) | **Put** /v1/users/activation | Activate the user by the given token.
[**V1UsersRegistrationPost**](DefaultApi.md#V1UsersRegistrationPost) | **Post** /v1/users/registration | Register user based on given information.

# **V1HealthcheckGet**
> HealthcheckResponse V1HealthcheckGet(ctx, )
Returns status of service.

None.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**HealthcheckResponse**](HealthcheckResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1TasksGet**
> GetAllTasksResponse V1TasksGet(ctx, userId)
Returns all tasks for user identified by userId.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 

### Return type

[**GetAllTasksResponse**](GetAllTasksResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1TasksPost**
> CreateTaskResponse V1TasksPost(ctx, userId, optional)
Create a new task for user 'user_id'

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
 **optional** | ***DefaultApiV1TasksPostOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1TasksPostOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateTaskRequest**](CreateTaskRequest.md)|  | 

### Return type

[**CreateTaskResponse**](CreateTaskResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1TasksTaskIdDelete**
> DeleteTaskByIdResponse V1TasksTaskIdDelete(ctx, userId, taskId)
Delete task by id for specific user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
  **taskId** | **int64**| tasks&#x27;s id. | 

### Return type

[**DeleteTaskByIdResponse**](DeleteTaskByIdResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1TasksTaskIdGet**
> GetTaskByIdResponse V1TasksTaskIdGet(ctx, userId, taskId)
Returns all tasks for user identified by userId.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
  **taskId** | **int64**| tasks&#x27;s id. | 

### Return type

[**GetTaskByIdResponse**](GetTaskByIdResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1TasksTaskIdPatch**
> UpdateTaskByIdResponse V1TasksTaskIdPatch(ctx, userId, taskId, optional)
Update task for specific user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
  **taskId** | **int64**| tasks&#x27;s id. | 
 **optional** | ***DefaultApiV1TasksTaskIdPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1TasksTaskIdPatchOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**optional.Interface of UpdateTaskByIdRequest**](UpdateTaskByIdRequest.md)|  | 

### Return type

[**UpdateTaskByIdResponse**](UpdateTaskByIdResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1TokensAuthenticationPost**
> AuthenticationResponse V1TokensAuthenticationPost(ctx, optional)
Authenticate the user based on given token.

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DefaultApiV1TokensAuthenticationPostOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1TokensAuthenticationPostOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of AuthenticationRequest**](AuthenticationRequest.md)|  | 

### Return type

[**AuthenticationResponse**](AuthenticationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1UsersActivationPut**
> UserActivationResponse V1UsersActivationPut(ctx, token)
Activate the user by the given token.

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **token** | **string**| token that represents the user who want to be activated. | 

### Return type

[**UserActivationResponse**](UserActivationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1UsersRegistrationPost**
> UserRegistrationResponse V1UsersRegistrationPost(ctx, optional)
Register user based on given information.

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DefaultApiV1UsersRegistrationPostOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1UsersRegistrationPostOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of UserRegistrationRequest**](UserRegistrationRequest.md)|  | 

### Return type

[**UserRegistrationResponse**](UserRegistrationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

