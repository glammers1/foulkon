package api

import (
	"testing"

	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPI_AddPolicy(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		policyName  string
		path        string
		statements  []Statement

		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getUserByExternalIDResult *User

		addPolicyMethodResult       *Policy
		getPolicyByNameMethodResult *Policy
		wantError                   error

		getPolicyByNameMethodErr error
		addPolicyMethodErr       error
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			addPolicyMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCasePolicyAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			wantError: &Error{
				Code:    POLICY_ALREADY_EXIST,
				Message: "Unable to create policy, policy with org 123 and name test already exist",
			},
		},
		"ErrorCaseBadName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "**!^#~",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name **!^#~",
			},
		},
		"ErrorCaseEmptyActions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
				{
					Effect:  "allow",
					Actions: []string{},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Empty actions",
			},
		},
		"ErrorCaseEmptyResources": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Empty resources",
			},
		},
		"ErrorCaseBadOrgName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "**!^#~",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **!^#~",
			},
		},
		"ErrorCaseBadPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/**!^#~path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**!^#~path/",
			},
		},
		"ErrorCaseBadStatement": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "idufhefmfcasfluhf",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid effect: idufhefmfcasfluhf - Only 'allow' and 'deny' accepted",
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/1/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:policy/path/test",
			},
		},
		"ErrorCaseDenyResource": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			path:       "/path/",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/1/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policy",
						Org:  "example",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_CREATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_CREATE_POLICY,
								},
								Resources: []string{
									CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
		},
		"ErrorCaseAddPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			addPolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseGetPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[AddPolicyMethod][0] = testcase.addPolicyMethodResult
		testRepo.ArgsOut[AddPolicyMethod][1] = testcase.addPolicyMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.AddPolicy(testcase.requestInfo, testcase.policyName, testcase.path, testcase.org, testcase.statements)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.addPolicyMethodResult, policy)
	}
}

func TestAuthAPI_GetPolicyByName(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		policyName  string

		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getUserByExternalIDResult *User

		getPolicyByNameMethodResult *Policy
		wantError                   error

		getPolicyByNameMethodErr error
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCaseInternalError": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseBadPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "~#**!",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name ~#**!",
			},
		},
		"ErrorCaseBadOrgName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "~#**!",
			policyName: "p1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org ~#**!",
			},
		},
		"ErrorCasePolicyNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "policyUser",
			getPolicyByNameMethodResult: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
				Statements: &[]Statement{
					{
						Effect: "deny",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_POLICY, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/1/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:example:policy/path/policyUser",
			},
		},
		"ErrorCaseDenyResourceErr": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "policyUser",
			getPolicyByNameMethodResult: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
				Statements: &[]Statement{
					{
						Effect: "deny",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_POLICY, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/1/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "example",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:example:policy/path/policyUser",
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.GetPolicyByName(testcase.requestInfo, testcase.org, testcase.policyName)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.getPolicyByNameMethodResult, policy)
	}
}

func TestAuthAPI_ListPolicies(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		filter      *Filter
		// Expected result
		expectedPolicies []PolicyIdentity
		totalResult      int
		wantError        error
		// Manager Results
		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getUserByExternalIDResult *User
		getUserByExternalIDErr    error
		// Manager Errors
		getPoliciesFilteredMethodResult []Policy
		getPoliciesFilteredMethodErr    error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org: "example",
			},
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
				{
					Org:  "example",
					Name: "policyDenied",
				},
			},
			totalResult: 2,
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"OkCaseAdminNoOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &testFilter,
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			totalResult: 1,
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"OkCaseUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			filter: &Filter{
				Org: "example",
			},
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			totalResult: 1,
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/1/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "example",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_LIST_POLICIES,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_LIST_POLICIES,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path2/"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseMaxLimitSize": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:   "123",
				Limit: 10000,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: limit 10000, max limit allowed: 1000",
			},
		},
		"ErrorCaseInvalidPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				PathPrefix: "/path*/ /*",
				Org:        "123",
				Limit:      0,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: pathPrefix /path*/ /*",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org: "!#$$%**^",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !#$$%**^",
			},
		},
		"ErrorCaseInternalErrorGetPoliciesFiltered": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				PathPrefix: "/path/",
			},
			getPoliciesFilteredMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			filter: &Filter{
				PathPrefix: "/path/",
				Org:        "123",
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPoliciesFilteredMethod][0] = testcase.getPoliciesFilteredMethodResult
		testRepo.ArgsOut[GetPoliciesFilteredMethod][1] = testcase.totalResult
		testRepo.ArgsOut[GetPoliciesFilteredMethod][2] = testcase.getPoliciesFilteredMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policies, total, err := testAPI.ListPolicies(testcase.requestInfo, testcase.filter)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedPolicies, policies)
		assert.Equal(t, testcase.totalResult, total, "Error in test case %v", x)
	}
}

func TestAuthAPI_UpdatePolicy(t *testing.T) {
	testcases := map[string]struct {
		requestInfo   RequestInfo
		org           string
		policyName    string
		path          string
		newPolicyName string
		newPath       string
		statements    []Statement
		newStatements []Statement

		getPolicyByNameMethodResult *Policy
		getGroupsByUserIDResult     []TestUserGroupRelation
		getAttachedPoliciesResult   []TestPolicyGroupRelation
		getUserByExternalIDResult   *User
		updatePolicyMethodResult    *Policy

		wantError error

		getPolicyByNameMethodErr error
		getUserByExternalIDErr   error
		updatePolicyMethodErr    error

		getPolicyByNameMethodSpecialFunc func(string, string) (*Policy, error)
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						},
					},
				},
			},
		},
		"ErrorCaseInvalidPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "**!~#",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name **!~#",
			},
		},
		"ErrorCaseInvalidOrgName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "**!~#",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **!~#",
			},
		},
		"ErrorCaseInvalidNewPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "**!~#",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new name **!~#",
			},
		},
		"ErrorCaseInvalidNewPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/**~#!/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new path /**~#!/",
			},
		},
		"ErrorCaseInvalidNewStatements": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "jblkasdjgp",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid effect: jblkasdjgp - Only 'allow' and 'deny' accepted",
			},
		},
		"ErrorCaseGetPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCasePolicyNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseAuthUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						},
					},
				},
			},
			getUserByExternalIDErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
		},
		"ErrorCaseDenyResource": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path/test",
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path/test",
			},
		},
		"ErrorCaseNewPolicyAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path/"),
								},
							},
						},
					}, nil
				}
				return &Policy{
					ID:   "test2",
					Name: "test2",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				}, nil
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    POLICY_ALREADY_EXIST,
				Message: "Policy name: test2 already exists",
			},
		},
		"ErrorCaseNoPermissionsToRetrieveTarget": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path/"),
								},
							},
						},
					}, nil
				}
				return &Policy{
					ID:   "test2",
					Name: "test2",
					Org:  "123",
					Path: "/path2/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path2/"),
							},
						},
					},
				}, nil
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path2/test",
			},
		},
		"ErrorCaseNoPermissionsToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path/"),
								},
							},
						},
					}, nil
				}
				return nil, &database.Error{
					Code: database.POLICY_NOT_FOUND,
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path2/test2",
			},
		},
		"ErrorCaseExplicitDenyPermissionsToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path/"),
								},
							},
						},
					}, nil
				}
				return nil, &database.Error{
					Code: database.POLICY_NOT_FOUND,
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									CreateUrn("123", RESOURCE_POLICY, "/path2/", "test2"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path2/test2",
			},
		},
		"ErrorCaseErrorUpdatingPolicy": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path/"),
								},
							},
						},
					}, nil
				}
				return nil, &database.Error{
					Code: database.POLICY_NOT_FOUND,
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
								},
							},
						},
					},
				},
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_UPDATE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
								},
							},
						},
					},
				},
			},
			updatePolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[UpdatePolicyMethod][0] = testcase.updatePolicyMethodResult
		testRepo.ArgsOut[UpdatePolicyMethod][1] = testcase.updatePolicyMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.SpecialFuncs[GetPolicyByNameMethod] = testcase.getPolicyByNameMethodSpecialFunc
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.UpdatePolicy(testcase.requestInfo, testcase.org, testcase.policyName, testcase.newPolicyName, testcase.newPath, testcase.newStatements)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.updatePolicyMethodResult, policy)
	}
}

func TestAuthAPI_RemovePolicy(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		name        string

		getPolicyByNameMethodResult *Policy
		getPolicyByNameMethodErr    error
		getGroupsByUserIDResult     []TestUserGroupRelation
		getAttachedPoliciesResult   []TestPolicyGroupRelation
		getUserByExternalIDResult   *User
		getUserByExternalIDErr      error
		deletePolicyErr             error

		wantError error
	}{
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "example",
			name: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCaseInvalidName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "123",
			name: "invalid*",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "**!^#$%",
			name: "invalid",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **!^#$%",
			},
		},
		"ErrorCasePolicyNotExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "123",
			name: "policy",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseNotEnoughPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_DELETE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_DELETE_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseRemoveFail": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			deletePolicyErr: &database.Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[RemovePolicyMethod][0] = testcase.deletePolicyErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		err := testAPI.RemovePolicy(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_ListAttachedGroups(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		filter      *Filter
		// Expected result
		expectedGroups []PolicyGroups
		totalResult    int
		wantError      error
		// Manager Results
		getGroupsByUserIDResult     []TestUserGroupRelation
		getAttachedPoliciesResult   []TestPolicyGroupRelation
		getUserByExternalIDResult   *User
		getAttachedGroupsResult     []TestPolicyGroupRelation
		getPolicyByNameMethodResult *Policy
		// Manager Errors
		getAttachedGroupsErr     error
		getPolicyByNameMethodErr error
	}{
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:        "example",
				PolicyName: "test",
				Limit:      0,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getAttachedGroupsResult: []TestPolicyGroupRelation{
				{
					Group: &Group{
						ID:   "Group1",
						Org:  "org1",
						Name: "group1",
					},
				},
				{
					Group: &Group{
						ID:   "Group2",
						Org:  "org2",
						Name: "group2",
					},
				},
			},
			expectedGroups: []PolicyGroups{
				{
					Group: "group1",
				},
				{
					Group: "group2",
				},
			},
		},
		"ErrorCaseMaxLimitSize": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:        "123",
				PolicyName: "p1",
				Limit:      10000,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: limit 10000, max limit allowed: 1000",
			},
		},
		"ErrorCaseInvalidName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:        "123",
				PolicyName: "invalid*",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: policy invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:        "!*^**~$%",
				PolicyName: "p1",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !*^**~$%",
			},
		},
		"ErrorCasePolicyNotExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:        "123",
				PolicyName: "policy",
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			filter: &Filter{
				Org:        "example",
				PolicyName: "test",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseNotEnoughPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			filter: &Filter{
				Org:        "example",
				PolicyName: "test",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
									POLICY_ACTION_LIST_ATTACHED_GROUPS,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_LIST_ATTACHED_GROUPS,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseGetAttachedPoliciesFail": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org:        "example",
				PolicyName: "test",
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getAttachedGroupsErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][0] = testcase.getAttachedGroupsResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][1] = testcase.totalResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][2] = testcase.getAttachedGroupsErr
		groups, total, err := testAPI.ListAttachedGroups(testcase.requestInfo, testcase.filter)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedGroups, groups)
		assert.Equal(t, testcase.totalResult, total, "Error in test case %v", x)
	}
}
