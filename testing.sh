# Get Posts
curl localhost:9000/api/post

# Get Single Post
curl localhost:9000/api/post/{id}

# Create New Post
curl -X POST -d '{"title":"my new post", "content": "such a good writer"}' localhost:9000/api/post

# Modify Post
curl -X PUT -d '{"content": "this is better"}' localhost:9000/api/post/{id}

# Delete Post
curl -X DELETE localhost:9000/api/post/{id}
