-- +goose Up
-- +goose StatementBegin
-- Create the Department table
USE CompanyDB;

IF NOT EXIST CREATE TABLE Department (
    ID INT PRIMARY KEY,
    DepartmentName NVARCHAR(50) NOT NULL
);

-- Create the Employee table
IF NOT EXIST CREATE TABLE Employee (
    ID INT PRIMARY KEY IDENTITY(1,1),
    FirstName NVARCHAR(50),
    LastName NVARCHAR(50),
    Email NVARCHAR(100),
    DOB DATE,
    DepartmentID INT,
    Position NVARCHAR(50),
    FOREIGN KEY (DepartmentID) REFERENCES Department(DepartmentID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
USE CompanyDB;

-- Drop the Employee table
IF OBJECT_ID('Employee', 'U') IS NOT NULL
    DROP TABLE Employee;

-- Drop the Department table
IF OBJECT_ID('Department', 'U') IS NOT NULL
    DROP TABLE Department;
-- +goose StatementEnd
