package postman

import (
	"encoding/json"
	"io"
	"sort"
)

var (
	DefaultCollectionName = "Root"
)

type (
	// Collection describes the full API documentation
	Collection struct {
		Info      Info       `json:"info"`
		Items     []Item     `json:"item"`
		Variables []Variable `json:"variable"`
	}

	// Info describes the postman info section
	Info struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Schema      string `json:"schema"`
	}

	// Field describes query, header, form-data and urlencoded field of request
	Variable struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Disabled    bool   `json:"disabled"`
		Enabled     bool   `json:"enabled"` // this field is used by env
	}

	// URL describes URL of the request
	URL struct {
		Raw         string     `json:"raw"`
		Host        []string   `json:"host"`
		Path        []string   `json:"path"`
		Description string     `json:"description"`
		Query       []Variable `json:"query"`
		Variables   []Variable `json:"variable"`
	}

	// Body describes a request body
	Body struct {
		Mode       string     `json:"mode"`
		FormData   []Variable `json:"formdata"`
		URLEncoded []Variable `json:"urlencoded"`
		Raw        string     `json:"raw"`
	}

	// Request describes a request
	Request struct {
		Method      string     `json:"method"`
		Headers     []Variable `json:"header"`
		Body        Body       `json:"body"`
		URL         URL        `json:"url"`
		Description string     `json:"description"`
	}

	// Response describes a request resposne
	Response struct {
		ID              string     `json:"id"`
		Name            string     `json:"name"`
		OriginalRequest Request    `json:"originalRequest"`
		Status          string     `json:"status"`
		Code            int        `json:"code"`
		Headers         []Variable `json:"header"`
		Body            string     `json:"body"`
		PreviewLanguage string     `json:"_postman_previewlanguage"`
	}

	// RequestItem describes a request item
	RequestItem struct {
		Name        string        `json:"name"`
		Items       []RequestItem `json:"item"`
		Request     Request       `json:"request"`
		Responses   []Response    `json:"response"`
		IsSubFolder bool          `json:"_postman_isSubFolder"`
	}

	// Item describes a request collection/item
	Item struct {
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Items       []RequestItem `json:"item"`
		RequestItem
	}
)

// ParseFrom create new collection from specified reader
func (d *Collection) ParseFrom(rdr io.Reader) error {
	dcr := json.NewDecoder(rdr)
	if err := dcr.Decode(&d); err != nil {
		return err
	}

	// build base payloads
	d.buildItems()
	d.removeEmptyItems()

	// after building all the collections, remove disabled fields
	d.removeItemRequestDisabledField()
	d.removeItemResponseRequestDisabledField()

	// sort the collections in lexical order
	d.sortItems()

	return nil
}

func (d *Collection) buildItems() {
	c := Item{}
	c.Name = DefaultCollectionName
	for i := len(d.Items) - 1; i >= 0; i-- {
		if len(d.Items[i].Items) <= 0 {
			if d.Items[i].Request.Method == "" { //a collection with no sub-items and request method is empty
				continue
			}
			c.Items = append(c.Items, RequestItem{
				Name:      d.Items[i].Name,
				Request:   d.Items[i].Request,
				Responses: d.Items[i].Responses,
			})
			d.Items = append(d.Items[:i], d.Items[i+1:]...)
		} else {
			collection := Item{}
			for j := len(d.Items[i].Items) - 1; j >= 0; j-- {
				built := d.buildSubChildItems(d.Items[i].Items[j], &collection, d.Items[i].Name)
				if built {
					d.Items[i].Items = append(d.Items[i].Items[:j], d.Items[i].Items[j+1:]...) //removing the sub folder from the parent to make it a collection itself
				}
			}
		}
	}
	d.Items = append(d.Items, c)
}

func (d *Collection) buildSubChildItems(itm RequestItem, c *Item, pn string) bool {
	if itm.IsSubFolder {
		collection := Item{}
		collection.Name = pn + "/" + itm.Name
		collection.IsSubFolder = true
		for _, i := range itm.Items {
			d.buildSubChildItems(i, &collection, collection.Name)
		}
		d.Items = append(d.Items, collection)
		return true
	}
	c.Items = append(c.Items, RequestItem{
		Name:      itm.Name,
		Request:   itm.Request,
		Responses: itm.Responses,
	})

	return false
}

func (d *Collection) sortItems() {
	sort.Slice(d.Items, func(i int, j int) bool {
		if d.Items[i].Name == DefaultCollectionName {
			return true
		}
		return d.Items[i].Name < d.Items[j].Name
	})

	for index := range d.Items {
		sort.Slice(d.Items[index].Items, func(i, j int) bool {
			return d.Items[index].Items[i].Name < d.Items[index].Items[j].Name
		})
	}
}

func (d *Collection) removeEmptyItems() {
	for i := 0; i < len(d.Items); i++ {
		if len(d.Items[i].Items) == 0 {
			d.Items = append(d.Items[:i], d.Items[i+1:]...) //popping the certain empty collection
			d.removeEmptyItems()                            //recursion followed by break to ensure proper indexing after a pop
			break
		}
	}
	return
}

func (d *Collection) removeItemRequestDisabledField() {
	for i, c := range d.Items {
		for j := range c.Items {
			// remove disabled headers
			for k := len(d.Items[i].Items[j].Request.Headers) - 1; k >= 0; k-- {
				if d.Items[i].Items[j].Request.Headers[k].Disabled {
					d.Items[i].Items[j].Request.Headers = append(d.Items[i].Items[j].Request.Headers[:k], d.Items[i].Items[j].Request.Headers[k+1:]...)
				}
			}
			// remove disabled queries
			for l := len(d.Items[i].Items[j].Request.URL.Query) - 1; l >= 0; l-- {
				if d.Items[i].Items[j].Request.URL.Query[l].Disabled {
					d.Items[i].Items[j].Request.URL.Query = append(d.Items[i].Items[j].Request.URL.Query[:l], d.Items[i].Items[j].Request.URL.Query[l+1:]...)
				}
			}
			// remove disabled form-data
			for m := len(d.Items[i].Items[j].Request.Body.FormData) - 1; m >= 0; m-- {
				if d.Items[i].Items[j].Request.Body.FormData[m].Disabled {
					d.Items[i].Items[j].Request.Body.FormData = append(d.Items[i].Items[j].Request.Body.FormData[:m], d.Items[i].Items[j].Request.Body.FormData[m+1:]...)
				}
			}
			// remove disabled urlencoded
			for n := len(d.Items[i].Items[j].Request.Body.URLEncoded) - 1; n >= 0; n-- {
				if d.Items[i].Items[j].Request.Body.URLEncoded[n].Disabled {
					d.Items[i].Items[j].Request.Body.URLEncoded = append(d.Items[i].Items[j].Request.Body.URLEncoded[:n], d.Items[i].Items[j].Request.Body.URLEncoded[n+1:]...)
				}
			}
		}
	}
	return
}

func (d *Collection) removeItemResponseRequestDisabledField() {
	for i, c := range d.Items {
		for j, item := range c.Items {
			for o := range item.Responses {
				// remove disabled headers
				for k := len(d.Items[i].Items[j].Responses[o].OriginalRequest.Headers) - 1; k >= 0; k-- {
					if d.Items[i].Items[j].Responses[o].OriginalRequest.Headers[k].Disabled {
						d.Items[i].Items[j].Responses[o].OriginalRequest.Headers = append(d.Items[i].Items[j].Responses[o].OriginalRequest.Headers[:k], d.Items[i].Items[j].Responses[o].OriginalRequest.Headers[k+1:]...)
					}
				}
				// remove disabled queries
				for l := len(d.Items[i].Items[j].Responses[o].OriginalRequest.URL.Query) - 1; l >= 0; l-- {
					if d.Items[i].Items[j].Responses[o].OriginalRequest.URL.Query[l].Disabled {
						d.Items[i].Items[j].Responses[o].OriginalRequest.URL.Query = append(d.Items[i].Items[j].Responses[o].OriginalRequest.URL.Query[:l], d.Items[i].Items[j].Responses[o].OriginalRequest.URL.Query[l+1:]...)
					}
				}
				// remove disabled form-data
				for m := len(d.Items[i].Items[j].Responses[o].OriginalRequest.Body.FormData) - 1; m >= 0; m-- {
					if d.Items[i].Items[j].Responses[o].OriginalRequest.Body.FormData[m].Disabled {
						d.Items[i].Items[j].Responses[o].OriginalRequest.Body.FormData = append(d.Items[i].Items[j].Responses[o].OriginalRequest.Body.FormData[:m], d.Items[i].Items[j].Responses[o].OriginalRequest.Body.FormData[m+1:]...)
					}
				}
				// remove disabled urlencoded
				for n := len(d.Items[i].Items[j].Responses[o].OriginalRequest.Body.URLEncoded) - 1; n >= 0; n-- {
					if d.Items[i].Items[j].Responses[o].OriginalRequest.Body.URLEncoded[n].Disabled {
						d.Items[i].Items[j].Responses[o].OriginalRequest.Body.URLEncoded = append(d.Items[i].Items[j].Responses[o].OriginalRequest.Body.URLEncoded[:n], d.Items[i].Items[j].Responses[o].OriginalRequest.Body.URLEncoded[n+1:]...)
					}
				}
			}
		}
	}
	return
}
