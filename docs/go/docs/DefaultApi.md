# {{classname}}

All URIs are relative to *https://todos.unknowntpo.net*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1HealthcheckGet**](DefaultApi.md#V1HealthcheckGet) | **Get** /v1/healthcheck | Returns status of service.
[**V1TasksUserIdGet**](DefaultApi.md#V1TasksUserIdGet) | **Get** /v1/tasks/{userId} | Returns all tasks for user identified by userId.
[**V1TasksUserIdPost**](DefaultApi.md#V1TasksUserIdPost) | **Post** /v1/tasks/{userId} | Create a new task for user &#x27;user_id&#x27;
[**V1TasksUserIdTaskIdDelete**](DefaultApi.md#V1TasksUserIdTaskIdDelete) | **Delete** /v1/tasks/{userId}/{taskId} | Delete task by id for specific user.
[**V1TasksUserIdTaskIdGet**](DefaultApi.md#V1TasksUserIdTaskIdGet) | **Get** /v1/tasks/{userId}/{taskId} | Returns all tasks for user identified by userId.
[**V1TasksUserIdTaskIdPatch**](DefaultApi.md#V1TasksUserIdTaskIdPatch) | **Patch** /v1/tasks/{userId}/{taskId} | Update task for specific user.
[**V1TokensAuthenticationPost**](DefaultApi.md#V1TokensAuthenticationPost) | **Post** /v1/tokens/authentication | Authenticate the user based on given token.
[**V1UsersActivationPost**](DefaultApi.md#V1UsersActivationPost) | **Post** /v1/users/activation | Activate the user by the given token.
[**V1UsersRegistrationPost**](DefaultApi.md#V1UsersRegistrationPost) | **Post** /v1/users/registration | Returns registered user information.

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

# **V1TasksUserIdGet**
> GetAllTasksResponse V1TasksUserIdGet(ctx, userId)
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

# **V1TasksUserIdPost**
> CreateTaskResponse V1TasksUserIdPost(ctx, userId, optional)
Create a new task for user 'user_id'

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
 **optional** | ***DefaultApiV1TasksUserIdPostOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1TasksUserIdPostOpts struct
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

# **V1TasksUserIdTaskIdDelete**
> DeleteTaskByIdResponse V1TasksUserIdTaskIdDelete(ctx, userId, taskId)
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

# **V1TasksUserIdTaskIdGet**
> GetTaskByIdResponse V1TasksUserIdTaskIdGet(ctx, userId, taskId)
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

# **V1TasksUserIdTaskIdPatch**
> UpdateTaskByIdResponse V1TasksUserIdTaskIdPatch(ctx, userId, taskId, optional)
Update task for specific user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **int64**| user&#x27;s id. | 
  **taskId** | **int64**| tasks&#x27;s id. | 
 **optional** | ***DefaultApiV1TasksUserIdTaskIdPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiV1TasksUserIdTaskIdPatchOpts struct
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

# **V1UsersActivationPost**
> V1UsersActivationPost(ctx, token)
Activate the user by the given token.

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **token** | **string**| token that represents the user who want to be activated. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1UsersRegistrationPost**
> User V1UsersRegistrationPost(ctx, token)
Returns registered user information.

None.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **token** | **string**| token that represents the user who want to be activated. | 

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

