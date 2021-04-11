package ytapi

import (
	"net/url"
	"strings"
)

func (c *channelInfoYTInitialData) GetLinks() (res []channelInfoLink) {
	for _, l := range c.Header.Renderer.HeaderLinks.Renderer.PrimaryLinks {
		res = append(res, l)
	}
	for _, l := range c.Header.Renderer.HeaderLinks.Renderer.SecondaryLinks {
		res = append(res, l)
	}
	return
}

func (l *channelInfoLink) Extract() string {
	u := l.NavigationEndpoint.UrlEndpoint.URL
	idx := strings.Index(u, "&q=")
	if idx < 0 {
		return ""
	}
	s, _ := url.QueryUnescape(u[idx+3:])
	return s
}

type channelInfoLink struct {
	NavigationEndpoint struct {
		ClickTrackingParams string `json:"clickTrackingParams"`
		CommandMetadata     struct {
			WebCommandMetadata struct {
				URL         string `json:"url"`
				WebpageType string `json:"webPageType"`
				RootVe      int    `json:"rootVe"`
			} `json:"webCommandMetadata"`
		} `json:"commandMetadata"`
		UrlEndpoint struct {
			URL      string `json:"url"`
			Target   string `json:"target"`
			Nofollow bool   `json:"nofollow"`
		} `json:"urlEndpoint"`
	} `json:"navigationEndpoint"`
	Icon struct {
		Thumbnails []struct {
			URL string `json:"url"`
		} `json:"thumbnails"`
	} `json:"icon"`
	Title struct {
		Simpletext string `json:"simpleText"`
	} `json:"title"`
}
type channelInfoYTInitialData struct {
	Contents struct {
		Renderer struct {
			Tabs []struct {
				TabRenderer struct {
					Endpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetaData     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
								APIURL      string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
					} `json:"endpoint"`
					Title    string `json:"title"`
					Selected bool   `json:"selected"`
					Content  struct {
						Renderer struct {
							Contents []struct {
								Renderer struct {
									Contents []struct {
										Renderer struct {
											Description struct {
												SimpleText string `json:"simpleText"`
											} `json:"description"`
											ViewCountText struct {
												SimpleText string `json:"simpleText"`
											} `json:"viewCountText"`
											JoinedDateText struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"joinedDateText"`
										} `json:"channelAboutFullMetadataRenderer"`
									} `json:"contents"`
									TrackingParams string `json:"trackingParams"`
								} `json:"itemSectionRenderer"`
							} `json:"contents"`
							TrackingParams string `json:"trackingParams"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
					TrackingParams string `json:"trackingParams"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
	Metadata struct {
		Renderer struct {
			Title                string   `json:"title"`
			Description          string   `json:"description"`
			RSSURL               string   `json:"rssUrl"`
			ChannelConversionURL string   `json:"channelConversionUrl"`
			ExternalID           string   `json:"externalId"`
			Keywords             string   `json:"keywords"`
			OwnerURLs            []string `json:"ownerUrls"`
			Avatar               struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"avatar"`
			ChannelURL       string `json:"channelUrl"`
			IsFamilySafe     bool   `json:"isFamilySafe"`
			VanityChannelURL string `json:"vanityChannelUrl"`
		} `json:"channelMetadataRenderer"`
	} `json:"metadata"`
	Header struct {
		Renderer struct {
			ChannelID          string `json:"channelId"`
			Title              string `json:"title"`
			NavigationEndpoint struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				CommandMetadata     struct {
					WebCommandMetadata struct {
						URL         string `json:"url"`
						Webpagetype string `json:"webPageType"`
						Rootve      int    `json:"rootVe"`
						Apiurl      string `json:"apiUrl"`
					} `json:"webCommandMetadata"`
				} `json:"commandMetadata"`
				BrowseEndpoint struct {
					BrowseID         string `json:"browseId"`
					CanonicalBaseURL string `json:"canonicalBaseUrl"`
				} `json:"browseEndpoint"`
			} `json:"navigationEndpoint"`
			Avatar struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"avatar"`
			Banner struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"banner"`
			HeaderLinks struct {
				Renderer struct {
					PrimaryLinks   []channelInfoLink `json:"primaryLinks"`
					SecondaryLinks []channelInfoLink `json:"secondaryLinks"`
				} `json:"channelHeaderLinksRenderer"`
			} `json:"headerLinks"`
			SubscriberCountText struct {
				Accessibility struct {
					AccessibilityData struct {
						Label string `json:"label"`
					} `json:"accessibilityData"`
				} `json:"accessibility"`
				SimpleText string `json:"simpleText"`
			} `json:"subscriberCountText"`
		} `json:"c4TabbedHeaderRenderer"`
	} `json:"header"`
}
