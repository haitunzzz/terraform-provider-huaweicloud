// Generated by PMS #339
package gaussdb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceGaussdbMysqlPtModifyRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbMysqlPtModifyRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the parameter template ID.`,
			},
			"histories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the parameter modify records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter name.`,
						},
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the old parameter value.`,
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the new parameter value.`,
						},
						"update_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the change status.`,
						},
						"is_applied": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the parameter has been applied.`,
						},
						"updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the modification time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
						"applied": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
					},
				},
			},
		},
	}
}

type MysqlPtModifyRecordsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newMysqlPtModifyRecordsDSWrapper(d *schema.ResourceData, meta interface{}) *MysqlPtModifyRecordsDSWrapper {
	return &MysqlPtModifyRecordsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGaussdbMysqlPtModifyRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newMysqlPtModifyRecordsDSWrapper(d, meta)
	listModifyHistoryRst, err := wrapper.ListModifyHistory()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listModifyHistoryToSchema(listModifyHistoryRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API GaussDBforMySQL GET /v3/{project_id}/configurations/{configuration_id}/modify-history
func (w *MysqlPtModifyRecordsDSWrapper) ListModifyHistory() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "gaussdb")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/configurations/{configuration_id}/modify-history"
	uri = strings.ReplaceAll(uri, "{configuration_id}", w.Get("configuration_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("histories", "offset", "limit", 0).
		Request().
		Result()
}

func (w *MysqlPtModifyRecordsDSWrapper) listModifyHistoryToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("histories", schemas.SliceToList(body.Get("histories"),
			func(histories gjson.Result) any {
				return map[string]any{
					"parameter_name": histories.Get("parameter_name").Value(),
					"old_value":      histories.Get("old_value").Value(),
					"new_value":      histories.Get("new_value").Value(),
					"update_result":  histories.Get("update_result").Value(),
					"is_applied":     histories.Get("is_applied").Value(),
					"updated":        histories.Get("updated").Value(),
					"applied":        histories.Get("applied").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
