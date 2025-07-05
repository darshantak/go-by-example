package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func decodeJWT(tokenString string) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Error parsing JWT:", err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("JWT Claims:", claims)
		fmt.Println("Email:", claims["email"])

		// Access email: claims["email"]
	}
}
func main() {
	r := gin.Default()

	r.GET("/api/me", func(c *gin.Context) {
		email := c.GetHeader("X-Pomerium-Claim-Email")
		name := c.GetHeader("X-Pomerium-Claim-Name")
		for k, v := range c.Request.Header {
			fmt.Printf("%s: %v\n", k, v)
		}
		jwtAssertion := c.GetHeader("X-Pomerium-Jwt-Assertion")
		decodeJWT(jwtAssertion)
		// if email == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{
			"email": email,
			"name":  name,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Backend running behind Pomerium")
	})

	r.Run(":8080")
}
