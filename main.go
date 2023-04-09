package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set up the Textract client
	// session will automatically reads your AWS credentials from ~/.aws/credentials
	// to manually set the creds, use: Credentials: credentials.NewStaticCredentials
	// inside the aws.Config{} struct
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-southeast-1"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	textractClient := textract.New(sess)

	r := gin.Default()
	r.GET("/analyze", func(c *gin.Context) {

		imageFile, _, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer imageFile.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, imageFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Analyze the image with Textract
		resp, err := AnalyzeImage(textractClient, buf.Bytes())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	r.Run("127.0.0.1:8080")
}
