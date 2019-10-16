package policyinsights

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// PolicyEventsClient is the client for the PolicyEvents methods of the Policyinsights service.
type PolicyEventsClient struct {
	BaseClient
}

// NewPolicyEventsClient creates an instance of the PolicyEventsClient client.
func NewPolicyEventsClient() PolicyEventsClient {
	return NewPolicyEventsClientWithBaseURI(DefaultBaseURI)
}

// NewPolicyEventsClientWithBaseURI creates an instance of the PolicyEventsClient client.
func NewPolicyEventsClientWithBaseURI(baseURI string) PolicyEventsClient {
	return PolicyEventsClient{NewWithBaseURI(baseURI)}
}

// ListQueryResultsForManagementGroup queries policy events for the resources under the management group.
// Parameters:
// managementGroupName - management group name.
// top - maximum number of records to return.
// orderBy - ordering expression using OData notation. One or more comma-separated column names with an
// optional "desc" (the default) or "asc", e.g. "$orderby=PolicyAssignmentId, ResourceId asc".
// selectParameter - select expression using OData notation. Limits the columns on each record to just those
// requested, e.g. "$select=PolicyAssignmentId, ResourceId".
// from - ISO 8601 formatted timestamp specifying the start time of the interval to query. When not specified,
// the service uses ($to - 1-day).
// toParameter - ISO 8601 formatted timestamp specifying the end time of the interval to query. When not
// specified, the service uses request time.
// filter - oData filter expression.
// apply - oData apply expression for aggregations.
func (client PolicyEventsClient) ListQueryResultsForManagementGroup(ctx context.Context, managementGroupName string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (result PolicyEventsQueryResults, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PolicyEventsClient.ListQueryResultsForManagementGroup")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("policyinsights.PolicyEventsClient", "ListQueryResultsForManagementGroup", err.Error())
	}

	req, err := client.ListQueryResultsForManagementGroupPreparer(ctx, managementGroupName, top, orderBy, selectParameter, from, toParameter, filter, apply)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForManagementGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListQueryResultsForManagementGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForManagementGroup", resp, "Failure sending request")
		return
	}

	result, err = client.ListQueryResultsForManagementGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForManagementGroup", resp, "Failure responding to request")
	}

	return
}

// ListQueryResultsForManagementGroupPreparer prepares the ListQueryResultsForManagementGroup request.
func (client PolicyEventsClient) ListQueryResultsForManagementGroupPreparer(ctx context.Context, managementGroupName string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"managementGroupName":       autorest.Encode("path", managementGroupName),
		"managementGroupsNamespace": autorest.Encode("path", "Microsoft.Management"),
		"policyEventsResource":      autorest.Encode("path", "default"),
	}

	const APIVersion = "2017-08-09-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(orderBy) > 0 {
		queryParameters["$orderby"] = autorest.Encode("query", orderBy)
	}
	if len(selectParameter) > 0 {
		queryParameters["$select"] = autorest.Encode("query", selectParameter)
	}
	if from != nil {
		queryParameters["$from"] = autorest.Encode("query", *from)
	}
	if toParameter != nil {
		queryParameters["$to"] = autorest.Encode("query", *toParameter)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if len(apply) > 0 {
		queryParameters["$apply"] = autorest.Encode("query", apply)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/{managementGroupsNamespace}/managementGroups/{managementGroupName}/providers/Microsoft.PolicyInsights/policyEvents/{policyEventsResource}/queryResults", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListQueryResultsForManagementGroupSender sends the ListQueryResultsForManagementGroup request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyEventsClient) ListQueryResultsForManagementGroupSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListQueryResultsForManagementGroupResponder handles the response to the ListQueryResultsForManagementGroup request. The method always
// closes the http.Response Body.
func (client PolicyEventsClient) ListQueryResultsForManagementGroupResponder(resp *http.Response) (result PolicyEventsQueryResults, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListQueryResultsForResource queries policy events for the resource.
// Parameters:
// resourceID - resource ID.
// top - maximum number of records to return.
// orderBy - ordering expression using OData notation. One or more comma-separated column names with an
// optional "desc" (the default) or "asc", e.g. "$orderby=PolicyAssignmentId, ResourceId asc".
// selectParameter - select expression using OData notation. Limits the columns on each record to just those
// requested, e.g. "$select=PolicyAssignmentId, ResourceId".
// from - ISO 8601 formatted timestamp specifying the start time of the interval to query. When not specified,
// the service uses ($to - 1-day).
// toParameter - ISO 8601 formatted timestamp specifying the end time of the interval to query. When not
// specified, the service uses request time.
// filter - oData filter expression.
// apply - oData apply expression for aggregations.
func (client PolicyEventsClient) ListQueryResultsForResource(ctx context.Context, resourceID string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (result PolicyEventsQueryResults, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PolicyEventsClient.ListQueryResultsForResource")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("policyinsights.PolicyEventsClient", "ListQueryResultsForResource", err.Error())
	}

	req, err := client.ListQueryResultsForResourcePreparer(ctx, resourceID, top, orderBy, selectParameter, from, toParameter, filter, apply)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForResource", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListQueryResultsForResourceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForResource", resp, "Failure sending request")
		return
	}

	result, err = client.ListQueryResultsForResourceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForResource", resp, "Failure responding to request")
	}

	return
}

// ListQueryResultsForResourcePreparer prepares the ListQueryResultsForResource request.
func (client PolicyEventsClient) ListQueryResultsForResourcePreparer(ctx context.Context, resourceID string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"policyEventsResource": autorest.Encode("path", "default"),
		"resourceId":           resourceID,
	}

	const APIVersion = "2017-08-09-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(orderBy) > 0 {
		queryParameters["$orderby"] = autorest.Encode("query", orderBy)
	}
	if len(selectParameter) > 0 {
		queryParameters["$select"] = autorest.Encode("query", selectParameter)
	}
	if from != nil {
		queryParameters["$from"] = autorest.Encode("query", *from)
	}
	if toParameter != nil {
		queryParameters["$to"] = autorest.Encode("query", *toParameter)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if len(apply) > 0 {
		queryParameters["$apply"] = autorest.Encode("query", apply)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{resourceId}/providers/Microsoft.PolicyInsights/policyEvents/{policyEventsResource}/queryResults", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListQueryResultsForResourceSender sends the ListQueryResultsForResource request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyEventsClient) ListQueryResultsForResourceSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListQueryResultsForResourceResponder handles the response to the ListQueryResultsForResource request. The method always
// closes the http.Response Body.
func (client PolicyEventsClient) ListQueryResultsForResourceResponder(resp *http.Response) (result PolicyEventsQueryResults, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListQueryResultsForResourceGroup queries policy events for the resources under the resource group.
// Parameters:
// subscriptionID - microsoft Azure subscription ID.
// resourceGroupName - resource group name.
// top - maximum number of records to return.
// orderBy - ordering expression using OData notation. One or more comma-separated column names with an
// optional "desc" (the default) or "asc", e.g. "$orderby=PolicyAssignmentId, ResourceId asc".
// selectParameter - select expression using OData notation. Limits the columns on each record to just those
// requested, e.g. "$select=PolicyAssignmentId, ResourceId".
// from - ISO 8601 formatted timestamp specifying the start time of the interval to query. When not specified,
// the service uses ($to - 1-day).
// toParameter - ISO 8601 formatted timestamp specifying the end time of the interval to query. When not
// specified, the service uses request time.
// filter - oData filter expression.
// apply - oData apply expression for aggregations.
func (client PolicyEventsClient) ListQueryResultsForResourceGroup(ctx context.Context, subscriptionID string, resourceGroupName string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (result PolicyEventsQueryResults, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PolicyEventsClient.ListQueryResultsForResourceGroup")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("policyinsights.PolicyEventsClient", "ListQueryResultsForResourceGroup", err.Error())
	}

	req, err := client.ListQueryResultsForResourceGroupPreparer(ctx, subscriptionID, resourceGroupName, top, orderBy, selectParameter, from, toParameter, filter, apply)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListQueryResultsForResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForResourceGroup", resp, "Failure sending request")
		return
	}

	result, err = client.ListQueryResultsForResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForResourceGroup", resp, "Failure responding to request")
	}

	return
}

// ListQueryResultsForResourceGroupPreparer prepares the ListQueryResultsForResourceGroup request.
func (client PolicyEventsClient) ListQueryResultsForResourceGroupPreparer(ctx context.Context, subscriptionID string, resourceGroupName string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"policyEventsResource": autorest.Encode("path", "default"),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"subscriptionId":       autorest.Encode("path", subscriptionID),
	}

	const APIVersion = "2017-08-09-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(orderBy) > 0 {
		queryParameters["$orderby"] = autorest.Encode("query", orderBy)
	}
	if len(selectParameter) > 0 {
		queryParameters["$select"] = autorest.Encode("query", selectParameter)
	}
	if from != nil {
		queryParameters["$from"] = autorest.Encode("query", *from)
	}
	if toParameter != nil {
		queryParameters["$to"] = autorest.Encode("query", *toParameter)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if len(apply) > 0 {
		queryParameters["$apply"] = autorest.Encode("query", apply)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PolicyInsights/policyEvents/{policyEventsResource}/queryResults", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListQueryResultsForResourceGroupSender sends the ListQueryResultsForResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyEventsClient) ListQueryResultsForResourceGroupSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListQueryResultsForResourceGroupResponder handles the response to the ListQueryResultsForResourceGroup request. The method always
// closes the http.Response Body.
func (client PolicyEventsClient) ListQueryResultsForResourceGroupResponder(resp *http.Response) (result PolicyEventsQueryResults, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListQueryResultsForSubscription queries policy events for the resources under the subscription.
// Parameters:
// subscriptionID - microsoft Azure subscription ID.
// top - maximum number of records to return.
// orderBy - ordering expression using OData notation. One or more comma-separated column names with an
// optional "desc" (the default) or "asc", e.g. "$orderby=PolicyAssignmentId, ResourceId asc".
// selectParameter - select expression using OData notation. Limits the columns on each record to just those
// requested, e.g. "$select=PolicyAssignmentId, ResourceId".
// from - ISO 8601 formatted timestamp specifying the start time of the interval to query. When not specified,
// the service uses ($to - 1-day).
// toParameter - ISO 8601 formatted timestamp specifying the end time of the interval to query. When not
// specified, the service uses request time.
// filter - oData filter expression.
// apply - oData apply expression for aggregations.
func (client PolicyEventsClient) ListQueryResultsForSubscription(ctx context.Context, subscriptionID string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (result PolicyEventsQueryResults, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PolicyEventsClient.ListQueryResultsForSubscription")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("policyinsights.PolicyEventsClient", "ListQueryResultsForSubscription", err.Error())
	}

	req, err := client.ListQueryResultsForSubscriptionPreparer(ctx, subscriptionID, top, orderBy, selectParameter, from, toParameter, filter, apply)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForSubscription", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListQueryResultsForSubscriptionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForSubscription", resp, "Failure sending request")
		return
	}

	result, err = client.ListQueryResultsForSubscriptionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyEventsClient", "ListQueryResultsForSubscription", resp, "Failure responding to request")
	}

	return
}

// ListQueryResultsForSubscriptionPreparer prepares the ListQueryResultsForSubscription request.
func (client PolicyEventsClient) ListQueryResultsForSubscriptionPreparer(ctx context.Context, subscriptionID string, top *int32, orderBy string, selectParameter string, from *date.Time, toParameter *date.Time, filter string, apply string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"policyEventsResource": autorest.Encode("path", "default"),
		"subscriptionId":       autorest.Encode("path", subscriptionID),
	}

	const APIVersion = "2017-08-09-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(orderBy) > 0 {
		queryParameters["$orderby"] = autorest.Encode("query", orderBy)
	}
	if len(selectParameter) > 0 {
		queryParameters["$select"] = autorest.Encode("query", selectParameter)
	}
	if from != nil {
		queryParameters["$from"] = autorest.Encode("query", *from)
	}
	if toParameter != nil {
		queryParameters["$to"] = autorest.Encode("query", *toParameter)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if len(apply) > 0 {
		queryParameters["$apply"] = autorest.Encode("query", apply)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.PolicyInsights/policyEvents/{policyEventsResource}/queryResults", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListQueryResultsForSubscriptionSender sends the ListQueryResultsForSubscription request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyEventsClient) ListQueryResultsForSubscriptionSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListQueryResultsForSubscriptionResponder handles the response to the ListQueryResultsForSubscription request. The method always
// closes the http.Response Body.
func (client PolicyEventsClient) ListQueryResultsForSubscriptionResponder(resp *http.Response) (result PolicyEventsQueryResults, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
