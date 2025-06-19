# FHIR Patient Examples for HAPI FHIR Server

## Quick Examples

### Basic Patient (Minimal Required Fields)
```json
{
  "resourceType": "Patient",
  "active": true,
  "name": [
    {
      "family": "Doe",
      "given": ["John"]
    }
  ],
  "gender": "male",
  "birthDate": "1990-01-01"
}
```

### Patient with Contact Information
```json
{
  "resourceType": "Patient",
  "active": true,
  "name": [
    {
      "use": "official",
      "family": "Smith",
      "given": ["Jane", "Marie"]
    }
  ],
  "telecom": [
    {
      "system": "phone",
      "value": "555-123-4567",
      "use": "mobile"
    },
    {
      "system": "email",
      "value": "jane.smith@email.com"
    }
  ],
  "gender": "female",
  "birthDate": "1985-06-15",
  "address": [
    {
      "use": "home",
      "line": ["123 Main St"],
      "city": "Anytown",
      "state": "NY",
      "postalCode": "12345",
      "country": "US"
    }
  ]
}
```

## Usage with curl

```bash
# Post to HAPI FHIR server
curl -X POST 'http://localhost:8080/fhir/Patient' \
  -H 'Content-Type: application/fhir+json' \
  -H 'Accept: application/fhir+json' \
  -d '{
    "resourceType": "Patient",
    "active": true,
    "name": [{"family": "Test", "given": ["Patient"]}],
    "gender": "unknown"
  }'
```

## Testing with Your Go Application

You can use these JSON payloads to test your external patient endpoints:

```bash
# Test your local API that forwards to HAPI FHIR
curl -X POST 'http://localhost:8080/api/v1/external-patients' \
  -H 'Content-Type: application/json' \
  -d @examples/simple_patient.json
```
