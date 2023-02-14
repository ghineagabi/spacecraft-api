package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddRoutes(r *gin.RouterGroup) {

	r.POST("/spacecraft", insertSpacecraftHandler)
	r.DELETE("/spacecraft", deleteSpacecraftHandler)
	r.PATCH("/spacecraft", updateSpacecraftHandler)

	r.GET("/spacecraft", getSingleSpacecraftHandler)
	r.GET("/filter-spacecrafts", filterSpacecraftsHandler)

}

func insertSpacecraftHandler(ctx *gin.Context) {
	var sc CreateSpacecraft
	var u LoginFromHeader

	if err := ctx.BindJSON(&sc); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	uc, err := DecodeAuth(u.Auth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = CheckCredentials(&uc); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = InsertSpacecraft(&sc); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
	return
}

func deleteSpacecraftHandler(ctx *gin.Context) {
	var Spacecraft SpacecraftId
	var u LoginFromHeader

	if err := ctx.BindJSON(&Spacecraft); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	uc, err := DecodeAuth(u.Auth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = CheckCredentials(&uc); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = DeleteSpacecraftFromDB(Spacecraft.Id); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
	return
}

func updateSpacecraftHandler(ctx *gin.Context) {
	var Spacecraft UpdateSpacecraft
	var u LoginFromHeader

	if err := ctx.BindJSON(&Spacecraft); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	uc, err := DecodeAuth(u.Auth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = CheckCredentials(&uc); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = UpdateSpacecraftFromDB(&Spacecraft); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
	return
}

func getSingleSpacecraftHandler(ctx *gin.Context) {
	var spacecraft DetailedSpacecraft

	_id := ctx.Query("ID")
	if _id == "" {
		ctx.JSON(http.StatusBadRequest, "Id not present in query")
		return
	}
	id, err := strconv.Atoi(_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := GetSpacecraftsFromDB(&spacecraft, id); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, spacecraft)
}

func filterSpacecraftsHandler(ctx *gin.Context) {
	var spacecrafts []FilteredSpacecraft

	name := ctx.Query("name")
	class := ctx.Query("class")
	status := ctx.Query("status")

	if err := GetFilteredSpacecrafts(&spacecrafts, name, class, status); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	ctx.JSON(http.StatusOK, spacecrafts)
}
