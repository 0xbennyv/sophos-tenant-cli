# Sophos Dashboard CLI Tool
A multiplatform command-line interface designed for both Sophos Enterprise and Partner Dashboard APIs.

**Primary Purpose**: 
- Streamline and present a concise health report for every tenant within a Sophos partner or enterprise environment.

## Key Features:

1. **Automated Configuration**: 
   - Upon first launch, if `config.json` is absent in the application's directory, users will be prompted to input either enterprise or partner API credentials.
   
2. **CLI Execution**: 
   - After setting up the credentials, execute the tool via the command line. Using the command `sophos_tenant_cli.exe healthsummary` will trigger a comprehensive health check across all accessible tenants, compiling a CSV list of those that require attention.
