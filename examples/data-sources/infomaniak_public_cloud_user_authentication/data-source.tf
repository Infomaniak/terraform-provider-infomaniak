data "infomaniak_public_cloud_user_authentication" "clouds_yaml" {
  public_cloud_id         = 1234
  public_cloud_project_id = 5678
  public_cloud_user_id    = 42
  type                    = "clouds.yaml"
  region                  = "pub1"
}
