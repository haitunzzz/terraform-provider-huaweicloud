package gaussdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDBMysqlTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getParameterTemplate: Query the GaussDB MySQL parameter Template
	var (
		getParameterTemplateHttpUrl = "v3/{project_id}/configurations/{configuration_id}"
		getParameterTemplateProduct = "gaussdb"
	)
	getParameterTemplateClient, err := cfg.NewServiceClient(getParameterTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB Client: %s", err)
	}

	getParameterTemplatePath := getParameterTemplateClient.Endpoint + getParameterTemplateHttpUrl
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{project_id}",
		getParameterTemplateClient.ProjectID)
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{configuration_id}",
		state.Primary.ID)

	getParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getParameterTemplateResp, err := getParameterTemplateClient.Request("GET",
		getParameterTemplatePath, &getParameterTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL Parameter Template: %s", err)
	}
	return utils.FlattenResponse(getParameterTemplateResp)
}

func TestAccGaussDBMysqlTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_mysql_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBMysqlTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description",
						"test gaussdb mysql parameter template"),
					resource.TestCheckResourceAttr(rName, "datastore_engine", "gaussdb-mysql"),
					resource.TestCheckResourceAttr(rName, "datastore_version", "8.0"),
					resource.TestCheckResourceAttr(rName, "parameter_values.auto_increment_increment", "4"),
					resource.TestCheckResourceAttr(rName, "parameter_values.auto_increment_offset", "5"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testParameterTemplate_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description",
						"test gaussdb mysql parameter template update"),
					resource.TestCheckResourceAttr(rName, "datastore_engine", "gaussdb-mysql"),
					resource.TestCheckResourceAttr(rName, "datastore_version", "8.0"),
					resource.TestCheckResourceAttr(rName, "parameter_values.auto_increment_increment", "6"),
					resource.TestCheckResourceAttr(rName, "parameter_values.auto_increment_offset", "8"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter_values"},
			},
		},
	})
}

func testParameterTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name              = "%s"
  description       = "test gaussdb mysql parameter template"
  datastore_engine  = "gaussdb-mysql"
  datastore_version = "8.0"

  parameter_values = {
    auto_increment_increment = "4"
    auto_increment_offset    = "5"
  }
}
`, name)
}

func testParameterTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name              = "%s"
  description       = "test gaussdb mysql parameter template update"
  datastore_engine  = "gaussdb-mysql"
  datastore_version = "8.0"

  parameter_values = {
    auto_increment_increment = "6"
    auto_increment_offset    = "8"
  }
}
`, name)
}
