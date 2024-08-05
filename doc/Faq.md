# Q: Why the `tencent_cloud` and `tencent_cos` are separate types in the connectors? (`aliyun` and `aliyun_oss`, etc.)
A: The underlying APIs are designed differently.

The API for the common services of Tencent cloud is designed to be called using POST method (recommended, GET method is an alternative), with parameters encoded in the header or body.

Meanwhile, the API for Tencent COS is designed to be RESTful-like, using HTTP methods (GET, POST, PUT, DELETE, etc.) as actions on the resource, and parts of URI as resource identifier.

This project uses the official SDK to connect to the common services of Tencent cloud, and uses the reflection of SDK to connect to Tencent COS for simplification. Therefore, unique types for the two ways are required.
