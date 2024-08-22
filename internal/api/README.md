# GeoSpatial Article Management System API Documentation

This document outlines the API endpoints for the GeoSpatial Article Management System, which includes authentication, article management, spatial data handling, layer management, layer group management, and MVT (Mapbox Vector Tiles) functionality.


## Table of Contents
1. [Authentication API](#authentication-api)
2. [Article API](#article-api)
3. [Spatial Data API](#spatial-data-api)
4. [Layer API](#layer-api)
5. [Layer Group API](#layer-group-api)
6. [MVT API](#mvt-api)

## Authentication API

### POST /signup
Create a new user account.

**Request Body:**
```json
{
    "username": "newuser",
    "password": "securepassword123"
}
```

**Response:**
```json
{
    "message": "User created successfully"
}
```

### POST /signin
Authenticate an existing user.

**Request Body:**
```json
{
    "username": "existinguser",
    "password": "userpassword123"
}
```

**Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "username": "existinguser",
        "created_at": "2023-05-01T10:00:00Z",
        "updated_at": "2023-05-01T10:00:00Z"
    }
}
```

### POST /logout
Log out the current user.

**Response:**
```json
{
    "message": "Logged out successfully"
}
```

## Article API

### GET /articles
Retrieve all articles.

**Response:**
```json
[
    {
        "id": 1,
        "title": "First Article",
        "author": "John Doe",
        "content": "This is the content of the first article.",
        "image_url": "https://example.com/image1.jpg",
        "created_by": 1,
        "created_at": "2023-05-01T12:00:00Z"
    },
    {
        "id": 2,
        "title": "Second Article",
        "author": "Jane Smith",
        "content": "This is the content of the second article.",
        "image_url": "https://example.com/image2.jpg",
        "created_by": 2,
        "created_at": "2023-05-02T14:30:00Z"
    }
]
```

### GET /articles/:id
Retrieve a specific article by ID.

**Example:** `GET /articles/1`

**Response:**
```json
{
    "id": 1,
    "title": "First Article",
    "author": "John Doe",
    "content": "This is the content of the first article.",
    "image_url": "https://example.com/image1.jpg",
    "created_by": 1,
    "created_at": "2023-05-01T12:00:00Z"
}
```

### POST /articles
Create a new article.

**Request Body:**
```json
{
    "title": "New Article Title",
    "content": "Article content goes here.",
    "image_url": "https://example.com/image.jpg"
}
```

**Response:**
```json
{
    "id": 3,
    "title": "New Article Title",
    "author": "Current User",
    "content": "Article content goes here.",
    "image_url": "https://example.com/image.jpg",
    "created_by": 1,
    "created_at": "2023-05-03T09:15:00Z"
}
```

### PUT /articles/:id
Update an existing article.

**Example:** `PUT /articles/1`

**Request Body:**
```json
{
    "title": "Updated Article Title",
    "content": "Updated article content.",
    "image_url": "https://example.com/updated-image.jpg"
}
```

**Response:**
```json
{
    "id": 1,
    "title": "Updated Article Title",
    "author": "John Doe",
    "content": "Updated article content.",
    "image_url": "https://example.com/updated-image.jpg",
    "created_by": 1,
    "created_at": "2023-05-01T12:00:00Z"
}
```

### DELETE /articles/:id
Delete a specific article.

**Example:** `DELETE /articles/1`

**Response:**
```json
{
    "message": "Article deleted successfully"
}
```

## Spatial Data API

### POST /spatial-data
Upload new spatial data.

**Request Body:** `multipart/form-data`
- `table_name`: "new_spatial_data"
- `type`: "point"
- `file`: [GeoJSON file]

**Response:**
```json
{
    "message": "Spatial data created successfully"
}
```

### PUT /spatial-data/:table_name
Update existing spatial data.

**Example:** `PUT /spatial-data/existing_table`

**Request Body:** `multipart/form-data`
- `table_name`: "updated_table_name" (optional)
- `file`: [Updated GeoJSON file] (optional)

**Response:**
```json
{
    "message": "Spatial data updated successfully"
}
```

### DELETE /spatial-data/:table_name
Delete a spatial data table.

**Example:** `DELETE /spatial-data/existing_table`

**Response:**
```json
{
    "message": "Spatial data deleted successfully"
}
```

### GET /spatial-data
Get all spatial data

**Response:**
```json
[
    {
        "id": 1,
        "table_name": "cities",
        "type": "point",
        "created_at": "2023-05-01T10:00:00Z",
        "updated_at": "2023-05-01T10:00:00Z",
        "created_by": 1,
        "updated_by": 1
    },
    {
        "id": 2,
        "table_name": "rivers",
        "type": "linestring",
        "created_at": "2023-05-02T11:30:00Z",
        "updated_at": "2023-05-02T11:30:00Z",
        "created_by": 2,
        "updated_by": 2
    }
]
```

## Layer API

### GET /layers
Retrieve layers.

**Query Parameters:**
- `id`: Comma-separated list of layer IDs or "*" for all layers

**Examples:**
- `GET /layers?id=1,2,3` (retrieve specific layers)
- `GET /layers?id=*` (retrieve all layers)

**Response:**
```json
[
    {
        "id": 1,
        "layer_name": "City Layer",
        "coordinate": [0, 0],
        "layer": {
            "id": "cities",
            "source": {
                "type": "vector",
                "tiles": "http://localhost:8080/mvt/cities/{z}/{x}/{y}"
            },
            "source-layer": "cities",
            "type": "circle",
            "paint": {
                "circle-color": "#FF5733",
                "circle-radius": 5
            }
        }
    },
    {
        "id": 2,
        "layer_name": "River Layer",
        "coordinate": [0, 0],
        "layer": {
            "id": "rivers",
            "source": {
                "type": "vector",
                "tiles": "http://localhost:8080/mvt/rivers/{z}/{x}/{y}"
            },
            "source-layer": "rivers",
            "type": "line",
            "paint": {
                "line-color": "#3366FF",
                "line-width": 2
            }
        }
    }
]
```

### POST /layers
Create a new layer.

**Request Body:**
```json
{
    "spatial_data_id": 1,
    "layer_name": "New Layer",
    "coordinate": [0, 0],
    "color": "#FF5733"
}
```

**Response:**
```json
{
    "id": 3,
    "spatial_data_id": 1,
    "layer_name": "New Layer",
    "coordinate": [0, 0],
    "color": "#FF5733",
    "created_at": "2023-05-03T14:00:00Z",
    "updated_at": "2023-05-03T14:00:00Z",
    "created_by": 1,
    "updated_by": 1
}
```

### PUT /layers/:id
Update an existing layer.

**Example:** `PUT /layers/1`

**Request Body:**
```json
{
    "layer_name": "Updated Layer Name",
    "coordinate": [1, 1],
    "color": "#33FF57"
}
```

**Response:**
```json
{
    "id": 1,
    "spatial_data_id": 1,
    "layer_name": "Updated Layer Name",
    "coordinate": [1, 1],
    "color": "#33FF57",
    "created_at": "2023-05-01T10:00:00Z",
    "updated_at": "2023-05-03T15:30:00Z",
    "created_by": 1,
    "updated_by": 1
}
```

### DELETE /layers/:id
Delete a specific layer.

**Example:** `DELETE /layers/1`

**Response:**
```json
{
    "message": "Layer deleted successfully"
}
```

## Layer Group API

### GET /layer-groups
Retrieve all layer groups.

**Response:**
```json
[
  {
    "group_id": 1,
    "group_name": "City Group",
    "layers": [
      {
        "layer_id": 1,
        "layer_name": "Building Layer"
      },
      {
        "layer_id": 2,
        "layer_name": "Road Layer"
      }
    ]
  },
  {
    "group_id": 2,
    "group_name": "Water Group",
    "layers": [
      {
        "layer_id": 3,
        "layer_name": "River Layer"
      },
      {
        "layer_id": 4,
        "layer_name": "Lake Layer"
      }
    ]
  }
]
```

### POST /layer-groups
Create a new layer group.

**Request Body:**
```json
{
    "group_name": "New Layer Group"
}
```

**Response:**
```json
{
    "message": "Group created successfully"
}
```

### POST /layer-groups/add-layer
Add a layer to a group.

**Request Body:**
```json
{
    "layer_id": 1,
    "group_id": 2
}
```

**Response:**
```json
{
    "message": "Layer added to group successfully"
}
```

### DELETE /layer-groups/remove-layer
Remove a layer from a group.

**Query Parameters:**
- `layer_id`: 1
- `group_id`: 2

**Response:**
```json
{
    "message": "Layer removed from group successfully"
}
```

### DELETE /layer-groups/:id
Delete a group

**Response:**
```json
{
    "message": "Group deleted successfully"
}
```

## MVT API

### GET /mvt/:table_name/:z/:x/:y
Retrieve a vector tile for a specific table and tile coordinates.

**Example:** `GET /mvt/my_spatial_data/12/1234/5678`

**Response:**
Binary data (application/x-protobuf)

Note: The response for this endpoint is binary data representing the vector tile, not JSON.