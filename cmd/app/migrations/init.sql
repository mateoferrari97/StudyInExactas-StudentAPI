CREATE DATABASE IF NOT EXISTS university;

USE university;

CREATE TABLE IF NOT EXISTS faculty
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(50)                        NOT NULL,
    uri        VARCHAR(128),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS career
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    faculty_id BIGINT                             NOT NULL,
    name       VARCHAR(50)                        NOT NULL,
    uri        VARCHAR(128),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (faculty_id) REFERENCES faculty (id)
);

CREATE TABLE IF NOT EXISTS subject
(
    id   BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    uri  VARCHAR(128),
    meet VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS career_subject
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    career_id      BIGINT NOT NULL,
    subject_id     BIGINT NOT NULL,
    correlative_id BIGINT,
    hours          BIGINT,
    type           VARCHAR(64),
    points         BIGINT,
    FOREIGN KEY (career_id) REFERENCES career (id),
    FOREIGN KEY (subject_id) REFERENCES subject (id),
    FOREIGN KEY (correlative_id) REFERENCES subject (id)
);

CREATE TABLE IF NOT EXISTS professorship
(
    id                BIGINT AUTO_INCREMENT PRIMARY KEY,
    career_subject_id BIGINT                             NOT NULL,
    name              VARCHAR(50)                        NOT NULL,
    created_at        DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (career_subject_id) REFERENCES career_subject (id)
);

CREATE TABLE IF NOT EXISTS schedule
(
    professorship_id BIGINT NOT NULL,
    day              BIGINT NOT NULL,
    start            TIME   NOT NULL,
    end              TIME   NOT NULL,
    FOREIGN KEY (professorship_id) REFERENCES professorship (id)
);

CREATE TABLE IF NOT EXISTS material
(
    id               BIGINT AUTO_INCREMENT PRIMARY KEY,
    professorship_id BIGINT                             NOT NULL,
    uri              VARCHAR(128)                       NOT NULL,
    description      VARCHAR(128)                       NOT NULL,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS professor
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(50)                        NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS professorship_professor
(
    professorship_id BIGINT       NOT NULL,
    professor_id     BIGINT       NOT NULL,
    role             VARCHAR(128) NOT NULL,
    FOREIGN KEY (professorship_id) REFERENCES professorship (id),
    FOREIGN KEY (professor_id) REFERENCES professor (id)
);

CREATE TABLE IF NOT EXISTS student
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(50)                        NOT NULL,
    email      VARCHAR(128)                       NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS student_career
(
    student_id BIGINT,
    career_id  BIGINT,
    FOREIGN KEY (student_id) REFERENCES student (id),
    FOREIGN KEY (career_id) REFERENCES career (id)
);

CREATE TABLE IF NOT EXISTS student_career_subject
(
    student_id        BIGINT      NOT NULL,
    career_subject_id BIGINT      NOT NULL,
    status            VARCHAR(50) NOT NULL,
    description       VARCHAR(128),
    FOREIGN KEY (student_id) REFERENCES student (id),
    FOREIGN KEY (career_subject_id) REFERENCES career_subject (id)
);