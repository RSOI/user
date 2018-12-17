package main

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/user/controller"
	"github.com/RSOI/user/ui"
	"github.com/RSOI/user/utils"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func sendResponse(ctx *fasthttp.RequestCtx, r ui.Response, nolog ...bool) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)
	utils.LOG(fmt.Sprintf("Sending response. Status: %d", r.Status))

	doLog := true
	if len(nolog) > 0 {
		doLog = !nolog[0]
	}

	if doLog {
		controller.LogStat(ctx.Path(), r.Status, r.Error)
	}

	content, _ := json.Marshal(r)
	ctx.Write(content)
}

func indexGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Get service stats (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.IndexGET(ctx.Host())
	r.Status, r.Error = ui.ErrToResponse(err)

	nolog := true
	sendResponse(ctx, r, nolog)
}

func userPUT(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Creating new user (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.UserPUT(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	if r.Status == 200 {
		r.Status = 201 // REST :)
	}
	sendResponse(ctx, r)
}

func userGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Get user (%s)", ctx.Path()))
	var err error
	var r ui.Response

	id := ctx.UserValue("id").(string)
	r.Data, err = controller.UserGET(id)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func userPATCH(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Update user (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.UserPATCH(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func initRoutes() *fasthttprouter.Router {
	utils.LOG("Setup router...")
	router := fasthttprouter.New()
	router.GET("/", indexGET)
	router.PUT("/signup", userPUT)
	router.GET("/user/id:id", userGET)
	router.PATCH("/update", userPATCH)

	return router
}
