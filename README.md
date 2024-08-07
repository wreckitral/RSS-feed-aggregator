# RSS Feed Aggregator
The RSS Feed Aggregator API is a high-performance service built using Go,
using sqlc for type-safe SQL queries and goose for database migrations.
This API allows you to manage RSS feeds and articles.

## How to Run
This API can be run on your local development system using two methods.

### Directly
if you have golang and postgresql installed
- set your env
```
PORT=<port>
DBCONN=postgres://<username>:<password>@<hostname>:5432/<dbname>?sslmode=disable
```
- install [goose](https://github.com/pressly/goose)
- `cd sql/schema`
- `goose postgres "postgres://<username>:<password>@<hostname>:5432/<dbname>?sslmode=disable" up`
- `make run`

### Using Docker
to run the project,
- `docker-compose up -d`
- `docker-compose down`
to turn down the container

## Endpoints
<table>
  <tr>
    <td>Endpoint</td><td>Method</td><td>Status</td><td>Request</td><td>Response</td>
  </tr>
  <tr>
  <td> /users </td>
  <td> POST </td>
  <td> 201 </td>
  <td>

  ```json
  {
    "name": "udin"
  }
  ```

  </td>
  <td>

  ```json
  {
    "id": "b8239f5e-0227-43d8-b60c-5a184ea80e95",
    "createdAt": "2024-06-07T19:05:09.815201Z",
    "updatedAt": "2024-06-07T19:05:09.815201Z",
    "name": "udin",
    "apiKey": "<API KEY>"
  }
  ```

  </td>

  </tr>
  <tr>
  <td>  /users </td>
  <td> GET </td>
  <td> 200 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  </td>
  <td>

  ```json
  {
    "id": "b8239f5e-0227-43d8-b60c-5a184ea80e95",
    "createdAt": "2024-06-07T19:05:09.815201Z",
    "updatedAt": "2024-06-07T19:05:09.815201Z",
    "name": "udin",
    "apiKey": "<API KEY>"
  }
  ```

  </td>
  </tr>
  <tr>
  <td> /feeds </td>
  <td> POST </td>
  <td> 201 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  ```json
  {
    "name": "The Verge Tech",
    "url": "http://www.theverge.com/tech/rss/index.xml"
  }
  ```

  </td>
  <td>

  ```json
  {
    "feed": {
        "id": "b87fc085-431a-4606-92f6-bd12f8fd741c",
        "createdAt": "2024-06-14T08:04:42.546561Z",
        "updatedAt": "2024-06-14T08:04:42.546561Z",
        "name": "The Verge Tech",
        "url": "http://www.theverge.com/tech/rss/index.xml",
        "userId": "082536ec-3598-4f6d-93d4-0fe23cfc623d"
    },
    "feed_follow": {
        "id": "5d146edf-7715-42f5-a324-94b2771e3e2a",
        "feedId": "b87fc085-431a-4606-92f6-bd12f8fd741c",
        "userId": "082536ec-3598-4f6d-93d4-0fe23cfc623d",
        "createdAt": "2024-06-14T08:04:42.547676Z",
        "updatedAt": "2024-06-14T08:04:42.547676Z"
    }
  }
  ```

  </td>
  </tr>
    </tr>
  <tr>
  <td> /feeds </td>
  <td> GET </td>
  <td> 200 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  </td>
  <td>

  ```json
  [
    {
        "id": "ad1416d3-89ec-48c9-b0eb-5621f8d4ba5c",
        "createdAt": "2024-06-14T07:49:22.232084Z",
        "updatedAt": "2024-06-14T07:53:05.541544Z",
        "name": "CNN Tech",
        "url": "http://www.buzzfeed.com/tvandmovies.xml",
        "userId": "082536ec-3598-4f6d-93d4-0fe23cfc623d"
    },
    {
        "id": "e0a2b609-ab42-460b-8f08-4fb44887790c",
        "createdAt": "2024-06-10T06:57:58.735041Z",
        "updatedAt": "2024-06-14T07:53:05.545006Z",
        "name": "CNN Tech",
        "url": "http://rss.cnn.com/rss/cnn_tech.rss",
        "userId": "082536ec-3598-4f6d-93d4-0fe23cfc623d"
    },
    ...
  ]
  ```

  </td>
  </tr>

  <tr>
  <td> /feed_follows </td>
  <td> POST </td>
  <td> 201 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  ```json
  {
    "feedId": "ad1416d3-89ec-48c9-b0eb-5621f8d4ba5c"
  }
  ```

  </td>
  <td>

  ```json
  {
    "id": "e4dc4679-b99f-4273-8055-d8db4ed6b7d3",
    "feedId": "ad1416d3-89ec-48c9-b0eb-5621f8d4ba5c",
    "userId": "b8239f5e-0227-43d8-b60c-5a184ea80e95",
    "createdAt": "2024-07-14T23:47:21.23965Z",
    "updatedAt": "2024-07-14T23:47:21.23965Z"
  }
  ```

  </td>
  </tr>

  <tr>
  <td>  /feed_follows </td>
  <td> GET </td>
  <td> 200 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  </td>
  <td>

  ```json
  [
    {
        "id": "a700f056-9833-4171-9463-1acdc9e57cc2",
        "feedId": "91b64c0a-b549-49a0-8c10-149db699289a",
        "userId": "b8239f5e-0227-43d8-b60c-5a184ea80e95",
        "createdAt": "2024-06-09T11:35:55.716246Z",
        "updatedAt": "2024-06-09T11:35:55.716246Z"
    },
    {
        "id": "e4dc4679-b99f-4273-8055-d8db4ed6b7d3",
        "feedId": "ad1416d3-89ec-48c9-b0eb-5621f8d4ba5c",
        "userId": "b8239f5e-0227-43d8-b60c-5a184ea80e95",
        "createdAt": "2024-07-14T23:47:21.23965Z",
        "updatedAt": "2024-07-14T23:47:21.23965Z"
    }
    ...
  ]
  ```

  </td>
  </tr>

  <tr>
  <td>  /feed_follows/{feedFollowId} </td>
  <td> DELETE </td>
  <td> 200 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  </td>
  <td>

  ```json
  {
    "statusCode": 200,
    "msg": "feed with feed follow id: 5d146edf-7715-42f5-a324-94b2771e3e2a successfully unfollowed"
  }
  ```

  </td>
  </tr>

  <tr>
  <td>  /posts </td>
  <td> GET </td>
  <td> 200 </td>
  <td>

  ```
  'Authorization': ApiKey <API KEY>
  ```

  </td>
  <td>

  ```json
  [
    {
        "id": "59364a5d-7a32-48a0-bc5d-eedae4559576",
        "createdAt": "2024-06-10T06:58:02.636745Z",
        "updatedAt": "2024-06-10T06:58:02.636745Z",
        "title": "",
        "url": "https://www.cnn.com/videos/tech/2024/04/05/ne-yo-ai-impact-music-contd-lcl-vpx.cnn",
        "description": "Grammy award-winning artist Ne-Yo joins CNN's Laura Coates to discuss the impact of artificial intelligence on the music industry.",
        "publishedAt": "0001-01-01T00:00:00Z",
        "feedId": "e0a2b609-ab42-460b-8f08-4fb44887790c"
    },
    {
        "id": "19175151-2db2-4ce7-b965-8a2be24b7c52",
        "createdAt": "2024-06-10T06:58:02.638624Z",
        "updatedAt": "2024-06-10T06:58:02.638624Z",
        "title": "",
        "url": "https://www.cnn.com/videos/tech/2024/03/17/gps-0317-tiktok-ban-in-the-us.cnn",
        "description": "Fareed hosts a spirited debate on the House bill that could lead to a US ban on TikTok.",
        "publishedAt": "0001-01-01T00:00:00Z",
        "feedId": "e0a2b609-ab42-460b-8f08-4fb44887790c"
    },
    ...
  ]
  ```

  </td>
  </tr>

</table>

