package utils

import "github.com/gin-gonic/gin"

func ApiNoMethodFound() gin.H {
	return gin.H{"error": "invalid_action"}
}

func ApiBadResponse(error string) gin.H {
	return gin.H{"error": error}
}

func ApiSuccessResp() gin.H {
	return gin.H{"status": true}
}
