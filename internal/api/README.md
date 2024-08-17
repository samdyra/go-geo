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

Handles user authentication and session management.

### POST /signup
Create a new user account.

**Request Body:**
```json
{
    "username": "newuser",
    "password": "securepassword123"
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

### POST /logout
Log out the current user.

## Article API

Manages articles within the system.

### GET /articles
Retrieve all articles.

### GET /articles/:id
Retrieve a specific article by ID.

**Example:** `GET /articles/1`

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

### DELETE /articles/:id
Delete a specific article.

**Example:** `DELETE /articles/1`

## Spatial Data API

Handles geospatial data operations.

### POST /spatial-data
Upload new spatial data.

**Request Body:** `multipart/form-data`
- `table_name`: "new_spatial_data"
- `type`: "point"
- `file`: [GeoJSON file]

### PUT /spatial-data/:table_name
Update existing spatial data.

**Example:** `PUT /spatial-data/existing_table`

**Request Body:** `multipart/form-data`
- `table_name`: "updated_table_name" (optional)
- `file`: [Updated GeoJSON file] (optional)

### DELETE /spatial-data/:table_name
Delete a spatial data table.

**Example:** `DELETE /spatial-data/existing_table`

## Layer API

Manages layers for visualizing spatial data.

### GET /layers
Retrieve layers.

**Query Parameters:**
- `id`: Comma-separated list of layer IDs or "*" for all layers

**Examples:**
- `GET /layers?id=1,2,3` (retrieve specific layers)
- `GET /layers?id=*` (retrieve all layers)

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

### DELETE /layers/:id
Delete a specific layer.

**Example:** `DELETE /layers/1`

## Layer Group API

Manages layer groups for organizing spatial data.

### GET /layer-groups
Retrieve all layer groups.

### POST /layer-groups
Create a new layer group.

**Request Body:**
```json
{
    "group_name": "New Layer Group"
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

### DELETE /layer-groups/remove-layer
Remove a layer from a group.

**Query Parameters:**
- `layer_id`: 1
- `group_id`: 2

## MVT API

Provides Mapbox Vector Tiles for efficient map rendering.

### GET /mvt/:table_name/:z/:x/:y
Retrieve a vector tile for a specific table and tile coordinates.

**Example:** `GET /mvt/my_spatial_data/12/1234/5678`

- `:table_name` - Name of the spatial data table
- `:z` - Zoom level
- `:x` - X coordinate of the tile
- `:y` - Y coordinate of the tile