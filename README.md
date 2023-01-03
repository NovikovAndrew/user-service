# REST API

/users - list if users --- 200, 404, 500
/users/:id - user by id --- 200, 404, 500
POST user/:id - create user --- 204, 4xx, Header Location: url 
PUT user/:id - update all users --- 204/200, 400, 404, 500
PATCH user/:id - partial update user --- 204/200, 400, 404, 500
DELETE user/:id - delete user --- 204, 404, 400