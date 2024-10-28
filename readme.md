# Image Feed API

## Overview

This is a simple RESTful API for uploading images, adding captions, and commenting on posts. The API allows users to manage posts and comments.

## Features

- **Upload Images**: Create posts with images and captions.
- **Commenting**: Users can comment on posts and delete their comments.
- **Image Handling**: Images are stored in it's original format and also converted to `.jpg` format and resized to 600x600 for display.

## Technologies Used

- **Programming Language**: Go
- **Web Framework**: Echo
- **Database**: PostgreSQL
- **Containerization**: Docker

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your machine.
- PostgreSQL database setup (or use a Dockerized PostgreSQL).

### How To Run
  - Please visit [How To Run](https://github.com/mfsyahrz/image-feed-api/blob/master/docs/how_to_run.md) which located in docs folder.

## API 

This section outlines the available API endpoints for posts and comments

| Method | Endpoint       | Description               |
|--------|-----------------|---------------------------|
| `POST` | `/posts`      | Create new Post with caption |
| `GET` | `/posts`      | Get Posts with Pagination |
| `POST` | `/posts/[postID]/comments`      | Commenting on a post |
| `DELETE` | `/posts/[postID]/comments/[commentID]`      | Delete a comment |

- For more detail, please refer to [postman collection](https://github.com/mfsyahrz/image-feed-api/blob/master/docs/how_to_run.md) which located in docs folder.

## Productionizing Guide 
- For more detail, please refer to [productionizing guide](https://github.com/mfsyahrz/image-feed-api/blob/main/docs/productionizing_guide.md) which located in docs folder.