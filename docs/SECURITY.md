# Security & Production Notes

## API Security
- By default, the API is open and does not require authentication.
- For production deployments, consider adding authentication (e.g., API keys, JWT, or OAuth2).
- Use HTTPS in production to encrypt traffic.

## Database Security
- Use strong, unique passwords for your PostgreSQL user.
- Restrict database access to only trusted hosts/networks.
- Regularly back up your database.

## Rate Limiting
- Consider adding rate limiting middleware to prevent abuse.

## CORS
- If you plan to call the API from browsers, configure CORS headers appropriately.

## Updates & Patching
- Regularly update Go, dependencies, and Docker images to the latest stable versions.

## Error Handling
- Avoid exposing sensitive error details in API responses.

## Environment Variables
- Store secrets (DB credentials, API keys) in environment variables or a secrets manager, not in source code.

## Logging
- Log errors and important events, but avoid logging sensitive data.

## Monitoring
- Use monitoring tools to track uptime, errors, and performance in production.

---

For more advanced security, consider using a reverse proxy (e.g., Nginx) and Web Application Firewall (WAF).
