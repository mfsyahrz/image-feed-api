# Prototype Design Details

This document outlines the design considerations and assumptions made in this prototype. Our focus is to create a flexible and scalable foundation for handling image uploads and high-throughput comment retrieval, with adaptability for production requirements.

---

## 1. Assumptions

- **Centralized Storage**: Images are uploaded to a single, centralized storage solution without variability in storage hosts, simplifying initial setup and ensuring consistency. In production, a solution like AWS S3 could be used.

- **Read-Heavy Workload**: The application is assumed to be read-heavy, with frequent access to comments.

---

## 2. Considerations

- **Flexible Storage Configuration**: The base URL for storage is assigned at runtime based on configuration, allowing seamless switching between local and cloud storage environments (e.g., AWS S3). This flexibility supports future migration to scalable storage solutions.

- **Scalable Architecture**: The design is modular to facilitate the introduction of background processing, caching, and distributed storage, enabling easy scaling as the application grows.

- **Efficient Retrieval**: Cursor-based pagination and caching strategies are anticipated for high-volume comment retrieval, ensuring fast access even under heavy load.

---

## 3. Design

- **File Upload Handling**: Images are stored locally for simplicity, with a structure that allows moving to cloud storage. Processing logic (e.g., resizing, format conversion) can be offloaded to background workers in production for efficiency.

- **Modular API**: Stateless and RESTful API endpoints support horizontal scaling, with middleware handling cross-cutting concerns like logging, validation, and authorization later when we need to introduce it.

- **Efficient and Consistent Writes**: Although optimized for reads, the system maintains efficient write performance by using transactions for comment creation.  PostgreSQL’s transactional and atomic capabilities ensure data consistency, so the comment count is always accurate, even under concurrent writes. This way, we can minimize race conditions and maintains high performance for both read and write operations. This design can be adapted to eventual consistency model when greater write scalability is required.
