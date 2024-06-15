# Supabase Studio
For Supa-Manager to keep the UI everyone is used to we have a patched version of the Supabase studio.
This folder has the script we use to build that automatically applies all the patches. To add or update a patch please refer to the `patches` folder.

> [!Note]
> Currently the Supabase Studio image needs certain variables set at build time. As such you must run the command yourself to produce a "patched" image for your environment.
> Please refer to the .env.example file for the required variables

## Usage
`./build [branch] [docker tag] [.env file (optional)]`