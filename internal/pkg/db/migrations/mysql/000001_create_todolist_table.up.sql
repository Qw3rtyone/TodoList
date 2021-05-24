CREATE TABLE IF NOT EXISTS Todolist(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Title VARCHAR (127) NOT NULL,
    Note VARCHAR (127),
    Completed BOOLEAN,
    PRIMARY KEY (ID)
)