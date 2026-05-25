package httpapi

import "github.com/go-go-golems/go-go-host/internal/store"

type orgDTO struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type siteDTO struct {
	ID                 string `json:"id"`
	OrgID              string `json:"orgId"`
	Slug               string `json:"slug"`
	Name               string `json:"name"`
	PrimaryHost        string `json:"primaryHost"`
	Status             string `json:"status"`
	ActiveDeploymentID string `json:"activeDeploymentId"`
}

func orgToDTO(org *store.Org) orgDTO {
	return orgDTO{ID: org.ID, Slug: org.Slug, Name: org.Name}
}

func siteToDTO(site store.Site) siteDTO {
	return siteDTO{ID: site.ID, OrgID: site.OrgID, Slug: site.Slug, Name: site.Name, PrimaryHost: site.PrimaryHost, Status: site.Status, ActiveDeploymentID: site.ActiveDeploymentID}
}
