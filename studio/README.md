# Supabase Studio
For Supa-Manager to keep the UI everyone is used to we have a patched version of the Supabase studio.
This folder has the script we use to build that automatically applies all the patches. To add or update a patch please refer to the `patches` folder.

> [!Note]
> Currently the Supabase Studio image needs certain variables set at build time. As such you must run the command yourself to produce a "patched" image for your environment.
> Please refer to the .env.example file for the required variables

## Usage
`./build.sh [branch] [docker tag] [.env file (optional)]` - This will build the patched version of the Supabase Studio and tag it with the specified tag

`./patch.sh [branch]` - This will apply the patches to the specified branch

## Creating new Patches
1. Download a version of the Studio and patch it using `./patch.sh`
2. `git add . && git commit -m "checkpoint"` - Create a checkpoint so patches are applied incrementally
3. Make your changes to the studio codebase
4. `git diff > patches/00-yourpatch.patch` - Create a patch file and increment the number based on the last patch
