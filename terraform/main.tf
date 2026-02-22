locals {
	region = "us-central1"
	project_id = "cloudlab5and6"

	db_admin = "whoami"
	db_password = "aspuoighbvaws3"
}

# Identity for the app
resource "google_service_account" "app_sa" {
  account_id   = "jdm-app-sa"
  display_name = "Service Account for JDM Registry App"
}

# Add role for the app's identity
resource "google_project_iam_member" "sql_client" {
	project = locals.project_id
	role = "roles/cloudsql.client"
	member = "serviceAccount:${google_service_account.app_sa.email}"
}

resource "google_sql_database_instance" "dbi" {
	name = "jdm-db-instance"
	database_version = "POSTGRES_14"
	region = locals.region 
	settings {
		tier = "db-f1-micro"
		ip_configuration {
			ipv4_enabled = true
		}
	}
	deletion_protection = false
}

resource "google_sql_database" "database" {
	name = "jdm_parts_db"
	instance = google_sql_database_instance.dbi.name
}

resource "google_sql_user" "admin" {
	name = locals.db_admin 
	instance = google_sql_database_instance.dbi.name
	password = locals.db_password
}

resource "google_cloud_run_service" "app" {
	name = "jdm-registry-service"
	location  = locals.region

	template {
		spec {
			service_account_name = google_service_account.app_sa.email

			containers {
				image = "gcr.io/${locals.project_id}/jdm-app:latest"	
				env {
					name = "DATABASE_URL"
					value = "host=/cloudsql/${google_sql_database_instance.dbi.connection_name} user=${locals.db_admin} password=${locals.db_password} dbname=jdm_parts_db sslmode=disable"
				}
			}
		}
		metadata {
            annotations = {
                "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.dbi.connection_name
            }
        }
	}
}

resource "google_cloud_run_service_iam_member" "public_access" {
	location = google_cloud_run_service.app.location	
	service = google_cloud_run_service.app.name
	role = "roles/run.invoker"
	member = "allUsers"
}

