package core

import "github.com/gin-gonic/gin"

func APINoMethodFound() gin.H {
	return APIBadResponse("invalid_action")
}

func APIBadResponse(error string) gin.H {
	return gin.H{"error": error}
}

func APISuccessResp() gin.H {
	return gin.H{"status": true}
}

