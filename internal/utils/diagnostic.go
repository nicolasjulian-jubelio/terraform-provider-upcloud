package utils

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func diagWarningFromUpcloudProblem(err *upcloud.Problem, details string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  err.Title,
		Detail:   details,
	}
}

func diagBindingRemovedWarningFromUpcloudProblem(err *upcloud.Problem, name string) diag.Diagnostic {
	return diagWarningFromUpcloudProblem(err,
		fmt.Sprintf("Binding to an existing remote object '%s' will be removed from the state. Next plan will include action to re-create the object if you choose to keep it in config.", name),
	)
}
