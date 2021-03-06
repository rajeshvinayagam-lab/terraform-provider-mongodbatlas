package integration_testing

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func SkipTestExtCred(t *testing.T) {
	if strings.EqualFold(os.Getenv("SKIP_TEST_EXTERNAL_CREDENTIALS"), "true") {
		t.SkipNow()
	}
}
func TestTerraformResourceMongoDBAtlasEncryptionAtRestWithRole_basicAWS(t *testing.T) {
	SkipTestExtCred(t)
	t.Parallel()

	var (
		projectID   = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		accessKey   = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey   = os.Getenv("AWS_SECRET_ACCESS_KEY")
		customerKey = os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")
		awsRegion   = os.Getenv("AWS_REGION")
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-encryptionAtRest-roles",
		Vars: map[string]interface{}{
			"access_key":          accessKey,
			"secret_key":          secretKey,
			"customer_master_key": customerKey,
			"atlas_region":        awsRegion,
			"project_id":          projectID,
			"public_key":          publicKey,
			"private_key":         privateKey,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the IP of the instance
	awsRoleARN := terraform.Output(t, terraformOptions, "aws_iam_role_arn")
	cpaRoleID := terraform.Output(t, terraformOptions, "cpa_role_id")

	fmt.Println(fmt.Sprintf("awsRoleARN : %s", awsRoleARN))
	fmt.Println(fmt.Sprintf("cpaRoleID : %s", cpaRoleID))

	terraformOptionsUpdated := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-encryptionAtRest-roles",
		Vars: map[string]interface{}{
			"access_key":          accessKey,
			"secret_key":          secretKey,
			"customer_master_key": customerKey,
			"atlas_region":        awsRegion,
			"project_id":          projectID,
			"public_key":          publicKey,
			"private_key":         privateKey,
			"aws_iam_role_arn":    awsRoleARN,
		},
	})

	terraform.Apply(t, terraformOptionsUpdated)

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-encryptionAtRest-roles/second_step",
		Vars: map[string]interface{}{
			"customer_master_key": customerKey,
			"atlas_region":        awsRegion,
			"project_id":          projectID,
			"public_key":          publicKey,
			"private_key":         privateKey,
			"cpa_role_id":         cpaRoleID,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsSecond)

	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptionsSecond)

}
