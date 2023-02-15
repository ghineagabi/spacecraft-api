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

/*
	Steps:

1) Validates the fields from the provided JSON (in the body request)
2) Checks the Basic Auth, decodes and verifies if the credentials are met
3) Attempts to insert the Spacecraft in the DB. If successful, returns 202. Else, returns 400
*/
func insertSpacecraftHandler(ctx *gin.Context) {
	var sc CreateSpacecraft
	var u LoginFromHeader

	if err := ctx.BindJSON(&sc); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
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

// Similar to insertSpacecraftHandler: Validating (only SpacecraftID in body), credential checking,
// DB Query and status returning
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

// Validation (SpacecraftID in body is required, other fields are optional), credential checking,
// and updating the fields (At least one extra field is required to finish the query)
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

// Gets JSON response as (checks the ID from the query params)
/*
{
    "id": 1,
    "name": "Devastator",
    "class": "Star Destroyer",
    "crew": 35000,
    "image": "image",
    "value": 1999.99,
    "status": "operational",
    "armament": [
        {
            "title": "Turbo Laser",
            "qty": "60"
        },
        {
            "title": "Ion Cannons",
            "qty": "60"
        }
    ]
}*/
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

	if err = GetSpacecraftsFromDB(&spacecraft, id); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, spacecraft)
}

// Filters a Spacecrafts based on the <name, class, status> from the params in the request. Returns:
/*
{
    "data": [
        {
            "id": 1,
            "name": "Devastator",
            "status": "operational"
        },
        {
            "id": 2,
            "name": "Red Five",
            "status": "damaged"
        }
    ]
}
*/
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
