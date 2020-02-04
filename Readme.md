# Config

Environment variables are used to pass security-critical settings to deploy the app.
Process (shell) environment has higher priority than the `.env` file.

Example `.env` file:

```env
DB_NAME=wall-dev
DB_HOST=127.0.0.1
COOKIE_DOMAIN=
SESSION_KEY=keepsecret
SITE_NAME=localhost
SITE_URL=http://localhost/
IMAGE_URL=http://127.0.0.1:8080/uploads/
```

**For production, enable CORS, setting up domain and URLs.**

For development, do not set `COOKIE_DOMAIN`. This suppresses sending emails and turns "developer mode"
