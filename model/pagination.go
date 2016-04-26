package models

import (
	"math"
	"net/url"

	"github.com/jinzhu/gorm"
)

type Pagination struct {
	Query        *gorm.DB
	TotalEntites int        `json:"total_entites" `
	PerPage      int        `json:"per_page" `
	Path         string     `json:"path" `
	Page         int        `json:"page" `
	UrlQuery     url.Values `json:"url_query" `
	TotalPages   int        `json:"total_pages" `
}

func (p *Pagination) Paginate(page int) *gorm.DB {
	p.Page = page
	// TODO Paginate bug
	p.Query.Count(&p.TotalEntites)
	if p.TotalEntites == 0 {
		return p.Query
	}

	p.TotalPages = int(math.Ceil(float64(p.TotalEntites) / float64(p.PerPage)))

	if !(p.Page > 0 && p.Page <= p.TotalPages) {
		p.Page = 1
	}

	query := p.Query.Offset((p.Page - 1) * p.PerPage).Limit(p.PerPage)

	return query

}

func (p *Pagination) CanShowPre() bool {
	if p.Page <= 1 {
		return false
	} else {
		return true
	}
}

func (p *Pagination) CanShowNext() bool {
	return p.Page < p.TotalPages
}
