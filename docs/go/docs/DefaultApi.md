# {{classname}}

All URIs are relative to *https://todos.unknowntpo.net*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1HealthcheckGet**](DefaultApi.md#V1HealthcheckGet) | **Get** /v1/healthcheck | Returns status of service.
[**V1TokensAuthenticationPost**](DefaultApi.md#V1TokensAuthenticationPost) | **Post** /v1/tokens/authentication | Authenticate the user based on given token.
[**V1UserIdTasksGet**](DefaultApi.md#V1UserIdTasksGet) | **Get** /v1/{userId}/tasks | Returns all tasks for user identified by userId.
[**V1UserIdTasksPost**](DefaultApi.md#V1UserIdTasksPost) | **Post** /v1/{userId}/tasks | Create a new task for user &#x27;user_id&#x27;
[**V1UserIdTasksTaskIdDelete**](DefaultApi.md#V1UserIdTasksTaskIdDelete) | **Delete** /v1/{userId}/tasks/{taskId} | Delete task by id for specific user.
[**V1UserIdTasksTaskIdGet**](DefaultApi.md#V1UserIdTasksTaskIdGet) | **Get** /v1/{userId}/tasks/{taskId} | Returns all tasks for user identified by userId.
[**V1UserIdTasksTaskIdPatch**](DefaultApi.md#V1UserIdTasksTaskIdPatch) | **Patch** /v1/{userId}/tasks/{taskId} | Update task for specific user.
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

# **V1UserIdTasksGet**
> GetAllTasksResponse V1UserIdTasksGet(ctx, userId)
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

# **V1UserIdTasksPost**
> CreateTaskResponse V1UserIdTasksPost(ctx, userId, optional)
Create a new task for user 'user_id'

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
 **optional** | ***DefaultApiV1UserIdTasksPostOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1UserIdTasksPostOpts struct
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

# **V1UserIdTasksTaskIdDelete**
> DeleteTaskByIdResponse V1UserIdTasksTaskIdDelete(ctx, userId, taskId)
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

# **V1UserIdTasksTaskIdGet**
> GetTaskByIdResponse V1UserIdTasksTaskIdGet(ctx, userId, taskId)
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

# **V1UserIdTasksTaskIdPatch**
> UpdateTaskByIdResponse V1UserIdTasksTaskIdPatch(ctx, userId, taskId, optional)
Update task for specific user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
  **taskId** | **int64**| tasks&#x27;s id. | 
 **optional** | ***DefaultApiV1UserIdTasksTaskIdPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1UserIdTasksTaskIdPatchOpts struct
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

