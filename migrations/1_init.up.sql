DROP USER IF EXISTS 'mysqladmin'@'%';

CREATE USER 'mysqladmin'@'%' IDENTIFIED BY 'mysqladmin';
GRANT ALL PRIVILEGES ON logging.* TO 'mysqladmin'@'%';
FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS logging.audio_file (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    link VARCHAR(1000),
    date_created DATETIME
);


CREATE TABLE IF NOT EXISTS logging.text_file (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
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