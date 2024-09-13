package main

import (
	"fmt"
	"n0ctRnull/blog-api/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Post struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// TODO: Get a single blog post

func getPost(c *gin.Context) {

	post := Post{}

	Id := c.Param("postId")

	query := "SELECT * FROM post WHERE id=$1"

	err := database.Db.QueryRow(query, Id).Scan(&post.Id, &post.Title, &post.Content, &post.Category, pq.Array(&post.Tags))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return

	}

	c.IndentedJSON(http.StatusOK, post)

}

// TODO: Update a blog post
func updatePost(c *gin.Context) {
	var updatedPost Post

	if err := c.BindJSON(&updatedPost); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	Id := c.Param("postId")
	query := "UPDATE post SET title=$1, content=$2, category=$3, tags=$4 WHERE id=$5"

	_, err := database.Db.Exec(query, updatedPost.Title, updatedPost.Content, updatedPost.Category, pq.Array(updatedPost.Tags), Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated"})

}

func deletePost(c *gin.Context) {
	Id := c.Param("postId")

	query := "DELETE FROM post WHERE id=$1"

	_, err := database.Db.Exec(query, Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.IndentedJSON(http.StatusOK, "post deleted")
}

func getPosts(c *gin.Context) {

	posts := []Post{}
	query := "SELECT * FROM post ORDER BY id ASC "
	rows, err := database.Db.Query(query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Couldn't get posts")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post

		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Category, pq.Array(&post.Tags))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't parse"})
			return
		}

		posts = append(posts, post)

	}

	if err = rows.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error iterating rows"})
		return
	}

	c.IndentedJSON(200, posts)
}

func addPost(c *gin.Context) {
	newPost := Post{}

	if err := c.BindJSON(&newPost); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})

	}

	query := "INSERT into post(title,content,category,tags) values ($1,$2,$3,$4)"

	_, err := database.Db.Exec(query, newPost.Title, newPost.Content, newPost.Category, pq.Array(newPost.Tags))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, "Couldn't create post")
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Post created"})
}

func main() {

	router := gin.Default()
	database.ConnectDatabase()
	router.GET("/posts", getPosts)
	router.GET("/posts/:postId", getPost)
	router.DELETE("/posts/:postId", deletePost)
	router.PUT("/posts/:postId", updatePost)
	router.POST("/posts", addPost)

	router.Run("localhost:8080")

}
