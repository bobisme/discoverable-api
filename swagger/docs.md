/accounts/{id}
==============

GET /api/v1/accounts/{id}
-------------------------

Show an account

get string by ID


### Parameters

```yaml
- description: Account ID
  in: path
  name: id
  required: true
  type: integer
```

### Responses

#### 404: Not Found

```yaml
description: Not Found
schema:
  properties:
    code:
      example: 400
      type: integer
    message:
      example: status bad request
      type: string
  type: object
```

#### 500: Internal Server Error

```yaml
description: Internal Server Error
schema:
  properties:
    code:
      example: 400
      type: integer
    message:
      example: status bad request
      type: string
  type: object
```

#### 200: OK

```yaml
description: OK
schema:
  properties:
    id:
      type: integer
  type: object
```

#### 400: Bad Request

```yaml
description: Bad Request
schema:
  properties:
    code:
      example: 400
      type: integer
    message:
      example: status bad request
      type: string
  type: object
```

/paramexamples
==============

GET /api/v1/paramexamples
-------------------------

Show an account

get string by ID


### Parameters

```yaml
- description: A
  in: body
  name: data
  required: true
  schema:
    properties:
      a:
        description: A is the first parameter I ahve listed
        type: string
      b:
        type: string
      c:
        minimum: 10
        type: integer
    required:
    - a
    type: object
```

### Responses

#### 200: OK

```yaml
description: OK
schema:
  properties:
    id:
      type: integer
  type: object
```

