package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ejsuarez793/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm Client")
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Kike",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreatBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Unexpected Error: %v\n", err)
	}

	fmt.Printf("Blog has been created: %v", res)
	blogID := res.GetBlog().GetId()

	blog2 := &blogpb.Blog{
		AuthorId: "Kike",
		Title:    "My Second Blog",
		Content:  "Content of the second blog",
	}

	res2, err2 := c.CreateBlog(context.Background(), &blogpb.CreatBlogRequest{Blog: blog2})

	if err2 != nil {
		log.Fatalf("Unexpected Error: %v\n", err2)
	}

	fmt.Printf("Blog has been created: %v", res2)
	// blogID2 := res2.GetBlog().GetId()

	// read Blog
	fmt.Println("\n\nReading the blog")

	_, err3 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "622fda42d4a61483aef0b6c1"})
	if err3 != nil {
		fmt.Printf("Error happened while reading blog: %v\n", err3)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading blog: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update Blog

	fmt.Println("\n\nUpdating the blog")

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Kike 2.0",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions",
	}

	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Error happened while updating blog: %v\n", err)
	}

	fmt.Printf("Blog (updated) was read: %v\n", updateRes)

	// delete Blog
	fmt.Println("\n\nDeleting the blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting blog: %v\n", deleteErr)
	}

	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list Blogs
	fmt.Println("\n\n Listing the blog")
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalf("error while calling ListBlog RPC %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			//we've reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}

		fmt.Println(res.GetBlog())
	}

}
