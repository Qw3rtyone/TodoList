CREATE TABLE IF NOT EXISTS Users(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Title VARCHAR (127) NOT NULL,
    Note VARCHAR (127),
    Completed NUMBER (1),
    PRIMARY KEY (ID)
)