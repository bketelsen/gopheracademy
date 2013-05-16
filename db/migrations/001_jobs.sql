-- +goose Up
CREATE TABLE jobs (
    id int NOT NULL AUTO_INCREMENT,
    title varchar(255),
    location varchar(255),
    jobtype	int,
    description text,
    companyname	varchar(255),
    companywebsite varchar(255),
    companylogourl varchar(255),
    applyurl varchar(255),
    applyemail varchar(255),
    additionalinstructions text,
    purchaseremail varchar(255),
    approved bool,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE jobs;