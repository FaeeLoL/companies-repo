package handlers

var createCompaniesSchema = []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "id": {
      "type": "string",
      "format": "uuid",
      "description": "A unique identifier for the company (UUID format)"
    },
    "name": {
      "type": "string",
      "maxLength": 15,
      "description": "The name of the company, must be unique and not exceed 15 characters"
    },
    "description": {
      "type": "string",
      "maxLength": 3000,
      "description": "An optional description of the company, up to 3000 characters",
      "nullable": true
    },
    "employees_count": {
      "type": "integer",
      "minimum": 1,
      "description": "The number of employees in the company, must be at least 1"
    },
    "registered": {
      "type": "boolean",
      "description": "Indicates whether the company is officially registered"
    },
    "type": {
      "type": "string",
      "enum": [
        "Corporations",
        "NonProfit",
        "Cooperative",
        "Sole Proprietorship"
      ],
      "description": "The type of the company, must be one of the predefined values"
    }
  },
  "required": ["id", "name", "employees_count", "registered", "type"],
  "additionalProperties": false
}`)

var patchCompaniesSchema = []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "id": {
      "type": "string",
      "format": "uuid",
      "description": "A unique identifier for the company (UUID format)"
    },
    "name": {
      "type": "string",
      "maxLength": 15,
      "description": "The name of the company, must be unique and not exceed 15 characters"
    },
    "description": {
      "type": "string",
      "maxLength": 3000,
      "description": "An optional description of the company, up to 3000 characters",
      "nullable": true
    },
    "employees_count": {
      "type": "integer",
      "minimum": 1,
      "description": "The number of employees in the company, must be at least 1"
    },
    "registered": {
      "type": "boolean",
      "description": "Indicates whether the company is officially registered"
    },
    "type": {
      "type": "string",
      "enum": [
        "Corporations",
        "NonProfit",
        "Cooperative",
        "Sole Proprietorship"
      ],
      "description": "The type of the company, must be one of the predefined values"
    }
  },
  "required": [],
  "oneOf": [
    { "required": ["id"] },
    { "required": ["name"] }
  ],
  "additionalProperties": false
}`)
