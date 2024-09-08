CREATE TABLE tasks IF not Exists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    note TEXT,
    status INT,
    due_on TIMESTAMP NULL DEFAULT NULL
);
