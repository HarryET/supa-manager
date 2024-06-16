# Version Management Service
This service provides the latest version of the individual containers/services used by supa-manager.

We have a hosted version available at https://supamanager.io/versions, this is updated as regularly as possible.

If you don't want any external dependencies you are welcome to host it yourself and either mirror or manually update.

## Configuration
- `DATABASE_URL` - PostgresSQL Connection String
- `PUSHING_ACCOUNTS` - Comma separated list of github accounts that are allowed to push new versions (needs a SSH key)
- `LISTEN_ADDRESS` - Address to listen on (default: `0.0.0.0:8081`)

## Usage
### Push new service version
1. Find the new version's image and tag i.e. `supabase/gotrue` & `0.1.0`
2. Sign the string `supabase/gotrue:0.1.0` with an allowed private key
3. Send a POST request to `/:service/update-version` with the following body:
```json
{
  "service": "supabase/gotrue",
  "version": "0.1.0"
}
```
and the header `signature` with the signed string

## Getting the latest versions for all services
Send a GET request to `/` and you will receive a JSON object with all the services and their latest versions.

## Getting all versions for a specific service
Send a GET request to `/:service` and you will receive a JSON object with all the versions for that service.
