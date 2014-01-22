package lastfm

import "encoding/xml"

//album.addTags (empty)

//album.getBuylinks
type AlbumGetBuylinks struct {
	XMLName   xml.Name `xml:"affiliations"`
	Physicals []struct {
		SupplierName string `xml:"supplierName"`
		SupplierIcon string `xml:"supplierIcon"`
		Price        struct {
			Currency  string `xml:"currency"`
			Amount    string `xml:"amount"`
			Formatted string `xml:"formatted"`
		} `xml:"price"`
		Buylink  string `xml:"buylink"`
		IsSearch string `xml:"isSearch"`
	} `xml:"physicals>affiliation"`
	Downloads []struct {
		SupplierName string `xml:"supplierName"`
		SupplierIcon string `xml:"supplierIcon"`
		Price        struct {
			Currency  string `xml:"currency"`
			Amount    string `xml:"amount"`
			Formatted string `xml:"formatted"`
		} `xml:"price"`
		Buylink  string `xml:"buylink"`
		IsSearch string `xml:"isSearch"`
	} `xml:"downloads>affiliation"`
}

//album.getInfo
type AlbumGetInfo struct {
	XMLName     xml.Name `xml:"album"`
	Name        string   `xml:"name"`
	Artist      string   `xml:"artist"`
	Id          string   `xml:"id"`
	Mbid        string   `xml:"mbid"`
	Url         string   `xml:"url"`
	ReleaseDate string   `xml:"releasedate"`
	Images      []struct {
		Size string `xml:"size,attr"`
		Url  string `xml:",chardata"`
	} `xml:"image"`
	Listeners string `xml:"listeners"`
	PlayCount string `xml:"plyacount"`
	Tracks    []struct {
		Rank       string `xml:"rank,attr"`
		Name       string `xml:"name"`
		Duration   string `xml:"duration"`
		Mbid       string `xml:"Mbid"`
		Url        string `xml:"url"`
		Streamable struct {
			FullTrack  string `xml:"fulltrack,attr"`
			Streamable string `xml:"streamable"`
		} `xml:"streamable"`
		Artist struct {
			Name string `xml:"name"`
			Mbid string `xml:"mbid"`
			Url  string `xml:"url"`
		} `xml:"artist"`
	} `xml:"tracks>track"`
	TopTags []struct {
		Name string `xml:"name"`
		Url  string `xml:"url"`
	} `xml:"toptags>tag"`
	Wiki struct {
		Published string `xml:"published"`
		Summary   string `xml:"summary"`
		Content   string `xml:"content"`
	} `xml:"wiki"`
}

//album.getShouts
type AlbumGetShouts struct {
	XMLName    xml.Name `xml:"shouts"`
	Artist     string   `xml:"artist,attr"`
	Album      string   `xml:"album,attr"`
	Total      int      `xml:"total,attr"`
	Page       int      `xml:"page,attr"`
	PerPage    int      `xml:"perPage,attr"`
	TotalPages int      `xml:"totalPages,attr"`
	Shouts     []struct {
		Body   string `xml:"body"`
		Author string `xml:"author"`
		Date   string `xml:"date"`
	} `xml:"shout"`
}

//album.getTags
type AlbumGetTags struct {
	XMLName xml.Name `xml:"tags"`
	Artist  string   `xml:"artist,attr"`
	Album   string   `xml:"album,attr"`
	Tags    []struct {
		Name string `xml:"name"`
		Url  string `xml:"url"`
	} `xml:"tag"`
}

//album.getTopTags
type AlbumGetTopTags struct {
	XMLName xml.Name `xml:"toptags"`
	Aritist string   `xml:"artist,attr"`
	Album   string   `xml:"album,attr"`
	Tags    []struct {
		Name  string `xml:"name"`
		Count string `xml:"count"`
		Url   string `xml:"url"`
	} `xml:"tag"`
}

//album.removeTag (empty)

//album.search
type AlbumSearch struct {
	XMLName    xml.Name `xml:"results"`
	OpenSearch string   `xml:"opensearch,attr"`
	For        string   `xml:"for,attr"`
	Query      struct {
		Role        string `xml:"role,attr"`
		SearchTerms string `xml:"searchTrems,attr"`
		StartPage   string `xml:"startPage,attr"`
	} `xml:"Query"`
	TotalResults string `xml:"totalResults"`
	StartIndex   string `xml:"startIndex"`
	ItemsPerPage string `xml:"itemsPerPage"`
	AlbumMatches []struct {
		Name   string `xml:"name"`
		Artist string `xml:"artist"`
		Id     string `xml:"id"`
		Url    string `xml:"url"`
		Images []struct {
			Size string `xml:"size,attr"`
			Url  string `xml:",chardata"`
		} `xml:"image"`
		Streamable string `xml:"streamable"`
		Mbid       string `xml:"mbid"`
	} `xml:"albummatches>album"`
}

//album.share (empty)