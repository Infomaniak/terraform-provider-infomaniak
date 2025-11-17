package dbaas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func getDbaasBackupScheduleResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud project",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"dbaas_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the dbaas",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "BackupSchedule identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the backup",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scheduled_at": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Use the given time as the time to create the scheduled backup (24 hour, UTC)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"retention": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The number of backups to keep for the schedule",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"is_pitr_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Enable/Disable point in time recovery",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
