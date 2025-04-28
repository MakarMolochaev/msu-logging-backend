CREATE DATABASE IF NOT EXISTS logging;

CREATE TABLE IF NOT EXISTS logging.audio_file (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    link VARCHAR(1000),
    date_created DATETIME
);


CREATE TABLE IF NOT EXISTS logging.protocols (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    task_id INT UNSIGNED,
    text_full VARCHAR(1000),
    text_short VARCHAR(1000),
    date_created DATETIME
);

CREATE TABLE IF NOT EXISTS logging.valuation (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    usability INT,
    processing_speed INT,
    processing_quality INT,
    reuse_service BOOLEAN,
    comment VARCHAR(1000)
);

CREATE TABLE IF NOT EXISTS logging.metrics (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    image_count INT,
    guest_count INT,
    av_audio_time FLOAT,
    av_process_time FLOAT,
    satisfy_user_count INT
);

CREATE TABLE IF NOT EXISTS logging.tasks (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    task_status VARCHAR(24)
);