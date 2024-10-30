package feed

import "github.com/BariVakhidov/rssaggregator/internal/model"

func (i *Implementation) validateFeedInfo(feedInfo model.FeedInfo) map[string]string {
	errors := make(map[string]string)

	if len(feedInfo.Name) == 0 {
		errors["name"] = ErrFeedNameRequired.Error()
	}

	if len(feedInfo.Url) == 0 {
		errors["url"] = ErrFeedUrlRequired.Error()
	}

	return errors
}
