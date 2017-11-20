// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"github.com/miniflux/miniflux2/server/api/payload"
	"github.com/miniflux/miniflux2/server/core"
)

// CreateFeed is the API handler to create a new feed.
func (c *Controller) CreateFeed(ctx *core.Context, request *core.Request, response *core.Response) {
	userID := ctx.GetUserID()
	feedURL, categoryID, err := payload.DecodeFeedCreationPayload(request.GetBody())
	if err != nil {
		response.Json().BadRequest(err)
		return
	}

	feed, err := c.feedHandler.CreateFeed(userID, categoryID, feedURL)
	if err != nil {
		response.Json().ServerError(errors.New("Unable to create this feed"))
		return
	}

	response.Json().Created(feed)
}

// RefreshFeed is the API handler to refresh a feed.
func (c *Controller) RefreshFeed(ctx *core.Context, request *core.Request, response *core.Response) {
	userID := ctx.GetUserID()
	feedID, err := request.GetIntegerParam("feedID")
	if err != nil {
		response.Json().BadRequest(err)
		return
	}

	err = c.feedHandler.RefreshFeed(userID, feedID)
	if err != nil {
		response.Json().ServerError(errors.New("Unable to refresh this feed"))
		return
	}

	response.Json().NoContent()
}

// UpdateFeed is the API handler that is used to update a feed.
func (c *Controller) UpdateFeed(ctx *core.Context, request *core.Request, response *core.Response) {
	userID := ctx.GetUserID()
	feedID, err := request.GetIntegerParam("feedID")
	if err != nil {
		response.Json().BadRequest(err)
		return
	}

	newFeed, err := payload.DecodeFeedModificationPayload(request.GetBody())
	if err != nil {
		response.Json().BadRequest(err)
		return
	}

	originalFeed, err := c.store.GetFeedById(userID, feedID)
	if err != nil {
		response.Json().NotFound(errors.New("Unable to find this feed"))
		return
	}

	if originalFeed == nil {
		response.Json().NotFound(errors.New("Feed not found"))
		return
	}

	originalFeed.Merge(newFeed)
	if err := c.store.UpdateFeed(originalFeed); err != nil {
		response.Json().ServerError(errors.New("Unable to update this feed"))
		return
	}

	response.Json().Created(originalFeed)
}

// GetFeeds is the API handler that get all feeds that belongs to the given user.
func (c *Controller) GetFeeds(ctx *core.Context, request *core.Request, response *core.Response) {
	feeds, err := c.store.GetFeeds(ctx.GetUserID())
	if err != nil {
		response.Json().ServerError(errors.New("Unable to fetch feeds from the database"))
		return
	}

	response.Json().Standard(feeds)
}

// GetFeed is the API handler to get a feed.
func (c *Controller) GetFeed(ctx *core.Context, request *core.Request, response *core.Response) {
	userID := ctx.GetUserID()
	feedID, err := request.GetIntegerParam("feedID")
	if err != nil {
		response.Json().BadRequest(err)
		return
	}

	feed, err := c.store.GetFeedById(userID, feedID)
	if err != nil {
		response.Json().ServerError(errors.New("Unable to fetch this feed"))
		return
	}

	if feed == nil {
		response.Json().NotFound(errors.New("Feed not found"))
		return
	}

	response.Json().Standard(feed)
}

// RemoveFeed is the API handler to remove a feed.
func (c *Controller) RemoveFeed(ctx *core.Context, request *core.Request, response *core.Response) {
	userID := ctx.GetUserID()
	feedID, err := request.GetIntegerParam("feedID")
	if err != nil {
		response.Json().BadRequest(err)
		return
	}

	if !c.store.FeedExists(userID, feedID) {
		response.Json().NotFound(errors.New("Feed not found"))
		return
	}

	if err := c.store.RemoveFeed(userID, feedID); err != nil {
		response.Json().ServerError(errors.New("Unable to remove this feed"))
		return
	}

	response.Json().NoContent()
}