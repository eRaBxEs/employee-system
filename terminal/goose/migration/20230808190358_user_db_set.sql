-- +goose Up
-- +goose StatementBegin
IF NOT EXIST CREATE TABLE User (
    ID INT PRIMARY KEY IDENTITY(1,1),
    Username NVARCHAR(50),
    Password NVARCHAR(50),
);

-- +goose StatementEnd


-- Drop the User table
IF OBJECT_ID('User', 'U') IS NOT NULL
    DROP TABLE User;