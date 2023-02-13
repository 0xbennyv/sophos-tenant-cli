# Sophos Tenant API's
A multiplatform CLI application for working with Sophos Enterprise Dashbaord and Partner Dashboard API's

This applications initial intent is to summarise the health report API for every Sophos tenant in a partner or enterprise tenant.

On initial run if config.json isn't within the directory of the binary then a prompt will appear to set the enterprise or partner tenants API credentials.

Once the credentials have been set the application must be run as a cli the "sophos_tenant_cli.exe healthsummary" will run the health summary function to check every availible tenant and build a CSV of troubled tenants.