# Backend

This project is a RESTful API service based on Go, mainly used for managing and verifying folder structures, images, annotations, and related information. The backend leverages various modern frameworks and tools to enhance development efficiency and maintainability.

## Technology Stack and Feature Overview

- **gin**: A high-performance RESTful API framework responsible for routing, request handling, and middleware management.
- **cors**: Middleware for Cross-Origin Resource Sharing (CORS), allowing frontend applications from different origins to securely access the API.
- **flag**: Used for managing startup parameters, such as `--env` to specify different runtime environments (development/production).
- **viper**: A powerful configuration management tool that automatically loads the corresponding YAML config file based on the specified environment and supports environment variable overrides.
- **fsnotify**: Directory monitoring tool, mainly used to watch for real-time changes (add, delete, modify) in the jobs directory and automatically refresh the data structure.
- **Swagger**: Automatically generates API documentation, making it easy for developers to browse and test the API.

## Configuration Management (flag + viper)

You can specify the environment when starting the service using the `--env` parameter, for example:

```sh
# Start in development environment
./your-app --env=development

# Start in production environment
./your-app --env=production
```

The system will automatically load the corresponding config file at `config/yaml/config.<env>.yaml`, such as:

- `config/yaml/config.development.yaml`
- `config/yaml/config.production.yaml`

Example (development):
```yaml
server:
  port: "8080"
  host: "localhost"
static:
  folder: "../example_root"
database:
  host: ""
  port: ""
  username: ""
  password: ""
  database: ""
logging:
  level: "debug"
  format: "console"
cors:
  allowed_origins: ["http://localhost:3000", "http://localhost:3001"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["*"]
```

## Directory Monitoring (fsnotify)

The system automatically monitors changes (add, delete, modify) in the jobs directory and synchronizes the data structure in real time, without needing to restart the service.

## Swagger API Documentation

After starting the service, you can open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser to view the complete API documentation and testing interface.

## Main API Endpoints

- `GET /api/folder-structure`: Get the folder structure.
- `GET /api/getJobs`: Get a list of all job names.
- `GET /api/getDatasets?job=JOB_NAME`: Get all dataset names under the specified job.
- `GET /api/getImages?job=JOB_NAME&dataset=DATASET_NAME`: Get all images under the specified job/dataset.
- `GET /api/getBase64Images?job=JOB_NAME&dataset=DATASET_NAME&pageIndex=0&pageNumber=10`: Get base64-encoded images by page.
- `GET /api/getAllPages?job=JOB_NAME`: Get all page information.
- `POST /api/savePendingReview`: Save pending review data.
- `GET /api/getPendingReview`: Get all pending review data.
- `POST /api/approvedRemove`: Remove approved data.
- `POST /api/unapprovedRemove`: Clear unapproved data.

For more details, please refer to the Swagger documentation.
