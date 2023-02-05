management_org = %SupaManager.Models.Organization{
  slug: "mng",
  name: "Managed Projects"
}
SupaManager.Repo.insert!(management_org)

open_org = %SupaManager.Models.Organization{
  slug: "prj",
  name: "Projects"
}
SupaManager.Repo.insert!(open_org)
