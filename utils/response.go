package utils

import "github.com/gin-gonic/gin"

func Success(data interface{}, message string) gin.H {
	return gin.H{
		"success": true,
		"data":    data,
		"message": message,
	}
}

func Error(message string) gin.H {
	return gin.H{
		"success": false,
		"data":    nil,
		"message": message,
	}
}